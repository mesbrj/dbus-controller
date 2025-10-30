package service

import (
	"encoding/xml"
	"fmt"
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/mesbrj/dbus-controller/internal/model"
)

// DBusService provides D-Bus operations
type DBusService struct {
	systemConn    *dbus.Conn
	sessionConn   *dbus.Conn
	mutex         sync.RWMutex
	subscriptions map[string]*SignalHandler
}

// SignalHandler manages signal subscriptions
type SignalHandler struct {
	channel chan *dbus.Signal
	active  bool
	mu      sync.RWMutex
}

// NewDBusService creates a new D-Bus service instance
func NewDBusService() *DBusService {
	service := &DBusService{
		subscriptions: make(map[string]*SignalHandler),
	}

	// Initialize system bus connection
	if conn, err := dbus.SystemBus(); err == nil {
		service.systemConn = conn
	}

	// Initialize session bus connection
	if conn, err := dbus.SessionBus(); err == nil {
		service.sessionConn = conn
	}

	return service
}

// Close closes all D-Bus connections
func (s *DBusService) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.systemConn != nil {
		s.systemConn.Close()
	}
	if s.sessionConn != nil {
		s.sessionConn.Close()
	}

	// Close all signal subscriptions
	for _, handler := range s.subscriptions {
		handler.mu.Lock()
		handler.active = false
		close(handler.channel)
		handler.mu.Unlock()
	}
}

// getConnection returns the appropriate D-Bus connection
func (s *DBusService) getConnection(busType string) (*dbus.Conn, error) {
	switch busType {
	case "system":
		if s.systemConn == nil {
			return nil, fmt.Errorf("system bus not available")
		}
		return s.systemConn, nil
	case "session":
		if s.sessionConn == nil {
			return nil, fmt.Errorf("session bus not available")
		}
		return s.sessionConn, nil
	default:
		return nil, fmt.Errorf("invalid bus type: %s", busType)
	}
}

// ListServices returns all services on the specified bus
func (s *DBusService) ListServices(busType string) ([]string, error) {
	conn, err := s.getConnection(busType)
	if err != nil {
		return nil, err
	}

	var services []string
	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&services)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	return services, nil
}

// GetServiceInfo returns detailed information about a service
func (s *DBusService) GetServiceInfo(busType, serviceName string) (*model.ServiceInfo, error) {
	conn, err := s.getConnection(busType)
	if err != nil {
		return nil, err
	}

	// Get service owner
	var owner string
	err = conn.BusObject().Call("org.freedesktop.DBus.GetNameOwner", 0, serviceName).Store(&owner)
	if err != nil {
		owner = "unknown"
	}

	// Get introspection data
	introspectionResult, err := s.IntrospectService(busType, serviceName)
	if err != nil {
		return &model.ServiceInfo{
			Name:  serviceName,
			Owner: owner,
		}, nil
	}

	interfaces := make([]string, 0)
	if introspectionResult.ParsedData != nil {
		for _, iface := range introspectionResult.ParsedData.Interfaces {
			interfaces = append(interfaces, iface.Name)
		}
	}

	return &model.ServiceInfo{
		Name:          serviceName,
		Owner:         owner,
		Interfaces:    interfaces,
		ObjectPaths:   []string{"/"}, // Default object path
		Introspection: introspectionResult,
	}, nil
}

// IntrospectService returns introspection data for a service
func (s *DBusService) IntrospectService(busType, serviceName string) (*model.IntrospectionResult, error) {
	conn, err := s.getConnection(busType)
	if err != nil {
		return nil, err
	}

	obj := conn.Object(serviceName, "/")
	var xmlData string

	err = obj.Call("org.freedesktop.DBus.Introspectable.Introspect", 0).Store(&xmlData)
	if err != nil {
		return nil, fmt.Errorf("failed to introspect service %s: %w", serviceName, err)
	}

	result := &model.IntrospectionResult{
		Service:    serviceName,
		ObjectPath: "/",
		XML:        xmlData,
		Timestamp:  time.Now(),
	}

	// Parse the introspection XML
	parsed, err := s.parseIntrospectionXML(xmlData)
	if err == nil {
		result.ParsedData = parsed
	}

	return result, nil
}

// parseIntrospectionXML parses introspection XML data
func (s *DBusService) parseIntrospectionXML(xmlData string) (*model.ParsedIntrospection, error) {
	var node introspect.Node
	err := xml.Unmarshal([]byte(xmlData), &node)
	if err != nil {
		return nil, fmt.Errorf("failed to parse introspection XML: %w", err)
	}

	parsed := &model.ParsedIntrospection{
		Interfaces: make([]model.InterfaceInfo, 0, len(node.Interfaces)),
		Nodes:      make([]model.NodeInfo, 0, len(node.Children)),
	}

	// Convert interfaces
	for _, iface := range node.Interfaces {
		interfaceInfo := model.InterfaceInfo{
			Name:       iface.Name,
			Methods:    make([]model.MethodInfo, 0, len(iface.Methods)),
			Properties: make([]model.PropertyInfo, 0, len(iface.Properties)),
			Signals:    make([]model.SignalInfo, 0, len(iface.Signals)),
		}

		// Convert methods
		for _, method := range iface.Methods {
			methodInfo := model.MethodInfo{
				Name:        method.Name,
				InArgs:      make([]model.ArgumentInfo, 0),
				OutArgs:     make([]model.ArgumentInfo, 0),
				Annotations: make(map[string]string),
			}

			for _, arg := range method.Args {
				argInfo := model.ArgumentInfo{
					Name:      arg.Name,
					Type:      arg.Type,
					Direction: arg.Direction,
				}
				if arg.Direction == "in" {
					methodInfo.InArgs = append(methodInfo.InArgs, argInfo)
				} else {
					methodInfo.OutArgs = append(methodInfo.OutArgs, argInfo)
				}
			}

			interfaceInfo.Methods = append(interfaceInfo.Methods, methodInfo)
		}

		// Convert properties
		for _, prop := range iface.Properties {
			propInfo := model.PropertyInfo{
				Name:        prop.Name,
				Type:        prop.Type,
				Access:      prop.Access,
				Annotations: make(map[string]string),
			}
			interfaceInfo.Properties = append(interfaceInfo.Properties, propInfo)
		}

		// Convert signals
		for _, signal := range iface.Signals {
			signalInfo := model.SignalInfo{
				Name:        signal.Name,
				Args:        make([]model.ArgumentInfo, 0, len(signal.Args)),
				Annotations: make(map[string]string),
			}

			for _, arg := range signal.Args {
				argInfo := model.ArgumentInfo{
					Name:      arg.Name,
					Type:      arg.Type,
					Direction: "out", // Signals always output
				}
				signalInfo.Args = append(signalInfo.Args, argInfo)
			}

			interfaceInfo.Signals = append(interfaceInfo.Signals, signalInfo)
		}

		parsed.Interfaces = append(parsed.Interfaces, interfaceInfo)
	}

	// Convert child nodes
	for _, child := range node.Children {
		nodeInfo := model.NodeInfo{
			Name: child.Name,
			Path: "/" + child.Name,
		}
		parsed.Nodes = append(parsed.Nodes, nodeInfo)
	}

	return parsed, nil
}

// ListInterfaces returns all interfaces for a service
func (s *DBusService) ListInterfaces(busType, serviceName string) ([]string, error) {
	serviceInfo, err := s.GetServiceInfo(busType, serviceName)
	if err != nil {
		return nil, err
	}

	return serviceInfo.Interfaces, nil
}

// GetInterfaceInfo returns detailed information about an interface
func (s *DBusService) GetInterfaceInfo(busType, serviceName, interfaceName string) (*model.InterfaceInfo, error) {
	serviceInfo, err := s.GetServiceInfo(busType, serviceName)
	if err != nil {
		return nil, err
	}

	if serviceInfo.Introspection != nil && serviceInfo.Introspection.ParsedData != nil {
		for _, iface := range serviceInfo.Introspection.ParsedData.Interfaces {
			if iface.Name == interfaceName {
				return &iface, nil
			}
		}
	}

	return nil, fmt.Errorf("interface %s not found", interfaceName)
}

// ListMethods returns all methods for an interface
func (s *DBusService) ListMethods(busType, serviceName, interfaceName string) ([]model.MethodInfo, error) {
	interfaceInfo, err := s.GetInterfaceInfo(busType, serviceName, interfaceName)
	if err != nil {
		return nil, err
	}

	return interfaceInfo.Methods, nil
}

// CallMethod executes a D-Bus method call
func (s *DBusService) CallMethod(busType, serviceName, interfaceName, methodName string, args []interface{}) (*model.MethodCallResult, error) {
	conn, err := s.getConnection(busType)
	if err != nil {
		return nil, err
	}

	obj := conn.Object(serviceName, "/")
	call := obj.Call(interfaceName+"."+methodName, 0, args...)

	result := &model.MethodCallResult{
		Timestamp: time.Now(),
	}

	if call.Err != nil {
		result.Success = false
		result.Error = call.Err.Error()
	} else {
		result.Success = true
		result.ReturnValues = call.Body
	}

	return result, nil
}

// ListProperties returns all properties for an interface
func (s *DBusService) ListProperties(busType, serviceName, interfaceName string) ([]model.PropertyInfo, error) {
	interfaceInfo, err := s.GetInterfaceInfo(busType, serviceName, interfaceName)
	if err != nil {
		return nil, err
	}

	// Try to get actual property values
	for i := range interfaceInfo.Properties {
		if value, err := s.GetProperty(busType, serviceName, interfaceName, interfaceInfo.Properties[i].Name); err == nil {
			interfaceInfo.Properties[i].Value = value.Value
		}
	}

	return interfaceInfo.Properties, nil
}

// GetProperty returns the value of a specific property
func (s *DBusService) GetProperty(busType, serviceName, interfaceName, propertyName string) (*model.PropertyValue, error) {
	conn, err := s.getConnection(busType)
	if err != nil {
		return nil, err
	}

	obj := conn.Object(serviceName, "/")
	variant, err := obj.GetProperty(interfaceName + "." + propertyName)
	if err != nil {
		return nil, fmt.Errorf("failed to get property %s: %w", propertyName, err)
	}

	return &model.PropertyValue{
		Name:      propertyName,
		Type:      variant.Signature().String(),
		Value:     variant.Value(),
		Timestamp: time.Now(),
	}, nil
}

// SetProperty sets the value of a specific property
func (s *DBusService) SetProperty(busType, serviceName, interfaceName, propertyName string, value interface{}) (*model.PropertyValue, error) {
	conn, err := s.getConnection(busType)
	if err != nil {
		return nil, err
	}

	obj := conn.Object(serviceName, "/")
	err = obj.SetProperty(interfaceName+"."+propertyName, dbus.MakeVariant(value))
	if err != nil {
		return nil, fmt.Errorf("failed to set property %s: %w", propertyName, err)
	}

	// Return the updated property value
	return s.GetProperty(busType, serviceName, interfaceName, propertyName)
}

// ListSignals returns all signals for an interface
func (s *DBusService) ListSignals(busType, serviceName, interfaceName string) ([]model.SignalInfo, error) {
	interfaceInfo, err := s.GetInterfaceInfo(busType, serviceName, interfaceName)
	if err != nil {
		return nil, err
	}

	return interfaceInfo.Signals, nil
}

// SubscribeToSignal subscribes to a D-Bus signal
func (s *DBusService) SubscribeToSignal(busType, serviceName, interfaceName, signalName string) (*model.SignalSubscription, error) {
	conn, err := s.getConnection(busType)
	if err != nil {
		return nil, err
	}

	// Create a unique subscription ID
	subscriptionID := fmt.Sprintf("%s:%s:%s:%s", busType, serviceName, interfaceName, signalName)

	// Create signal handler
	handler := &SignalHandler{
		channel: make(chan *dbus.Signal, 100),
		active:  true,
	}

	// Add match rule
	matchRule := fmt.Sprintf("type='signal',interface='%s',member='%s'", interfaceName, signalName)
	if serviceName != "" {
		matchRule += fmt.Sprintf(",sender='%s'", serviceName)
	}

	err = conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchRule).Err
	if err != nil {
		return nil, fmt.Errorf("failed to add match rule: %w", err)
	}

	// Register signal handler
	conn.Signal(handler.channel)

	s.mutex.Lock()
	s.subscriptions[subscriptionID] = handler
	s.mutex.Unlock()

	subscription := &model.SignalSubscription{
		ID:        subscriptionID,
		BusType:   busType,
		Service:   serviceName,
		Interface: interfaceName,
		Signal:    signalName,
		Active:    true,
		CreatedAt: time.Now(),
	}

	return subscription, nil
}
