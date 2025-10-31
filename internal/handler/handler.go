package handler

import (
	"github.com/go-fuego/fuego"
	"github.com/mesbrj/dbus-controller/internal/model"
	"github.com/mesbrj/dbus-controller/internal/service"
)

// Handler contains the HTTP handlers for D-Bus operations
type Handler struct {
	dbusService service.DBusServiceInterface
}

// NewHandler creates a new handler instance
func NewHandler(dbusService service.DBusServiceInterface) *Handler {
	return &Handler{
		dbusService: dbusService,
	}
}

// ListBuses returns available D-Bus types (system, session)
func (h *Handler) ListBuses(c fuego.ContextNoBody) ([]model.BusInfo, error) {
	buses := []model.BusInfo{
		{Type: "system", Description: "System D-Bus"},
		{Type: "session", Description: "Session D-Bus"},
	}
	return buses, nil
}

// GetBusInfo returns information about a specific bus
func (h *Handler) GetBusInfo(c fuego.ContextNoBody) (*model.BusInfo, error) {
	busType := c.PathParam("busType")

	if busType != "system" && busType != "session" {
		return nil, fuego.BadRequestError{Title: "Invalid bus type", Detail: "Bus type must be 'system' or 'session'"}
	}

	bus := &model.BusInfo{
		Type:        busType,
		Description: busType + " D-Bus",
	}

	return bus, nil
}

// ListServices returns all services on the specified bus
func (h *Handler) ListServices(c fuego.ContextNoBody) ([]string, error) {
	busType := c.PathParam("busType")
	return h.dbusService.ListServices(busType)
}

// GetService returns detailed information about a service
func (h *Handler) GetService(c fuego.ContextNoBody) (*model.ServiceInfo, error) {
	busType := c.PathParam("busType")
	serviceName := c.PathParam("serviceName")
	return h.dbusService.GetServiceInfo(busType, serviceName)
}

// ListInterfaces returns all interfaces for a service
func (h *Handler) ListInterfaces(c fuego.ContextNoBody) ([]string, error) {
	busType := c.PathParam("busType")
	serviceName := c.PathParam("serviceName")
	return h.dbusService.ListInterfaces(busType, serviceName)
}

// GetInterface returns detailed information about an interface
func (h *Handler) GetInterface(c fuego.ContextNoBody) (*model.InterfaceInfo, error) {
	busType := c.PathParam("busType")
	serviceName := c.PathParam("serviceName")
	interfaceName := c.PathParam("interfaceName")
	return h.dbusService.GetInterfaceInfo(busType, serviceName, interfaceName)
}

// ListMethods returns all methods for an interface
func (h *Handler) ListMethods(c fuego.ContextNoBody) ([]model.MethodInfo, error) {
	busType := c.PathParam("busType")
	serviceName := c.PathParam("serviceName")
	interfaceName := c.PathParam("interfaceName")
	return h.dbusService.ListMethods(busType, serviceName, interfaceName)
}

// CallMethodRequest represents the request body for method calls
type CallMethodRequest struct {
	Args []interface{} `json:"args,omitempty"`
}

// CallMethod executes a D-Bus method call
func (h *Handler) CallMethod(c *fuego.ContextWithBody[CallMethodRequest]) (*model.MethodCallResult, error) {
	busType := c.PathParam("busType")
	serviceName := c.PathParam("serviceName")
	interfaceName := c.PathParam("interfaceName")
	methodName := c.PathParam("methodName")

	body, err := c.Body()
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid request body", Detail: err.Error()}
	}

	return h.dbusService.CallMethod(busType, serviceName, interfaceName, methodName, body.Args)
}

// ListProperties returns all properties for an interface
func (h *Handler) ListProperties(c fuego.ContextNoBody) ([]model.PropertyInfo, error) {
	busType := c.PathParam("busType")
	serviceName := c.PathParam("serviceName")
	interfaceName := c.PathParam("interfaceName")
	return h.dbusService.ListProperties(busType, serviceName, interfaceName)
}

// GetProperty returns the value of a specific property
func (h *Handler) GetProperty(c fuego.ContextNoBody) (*model.PropertyValue, error) {
	busType := c.PathParam("busType")
	serviceName := c.PathParam("serviceName")
	interfaceName := c.PathParam("interfaceName")
	propertyName := c.PathParam("propertyName")
	return h.dbusService.GetProperty(busType, serviceName, interfaceName, propertyName)
}

// SetPropertyRequest represents the request body for setting properties
type SetPropertyRequest struct {
	Value interface{} `json:"value"`
}

// SetProperty sets the value of a specific property
func (h *Handler) SetProperty(c *fuego.ContextWithBody[SetPropertyRequest]) (*model.PropertyValue, error) {
	busType := c.PathParam("busType")
	serviceName := c.PathParam("serviceName")
	interfaceName := c.PathParam("interfaceName")
	propertyName := c.PathParam("propertyName")

	body, err := c.Body()
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid request body", Detail: err.Error()}
	}

	return h.dbusService.SetProperty(busType, serviceName, interfaceName, propertyName, body.Value)
}

// ListSignals returns all signals for an interface
func (h *Handler) ListSignals(c fuego.ContextNoBody) ([]model.SignalInfo, error) {
	busType := c.PathParam("busType")
	serviceName := c.PathParam("serviceName")
	interfaceName := c.PathParam("interfaceName")
	return h.dbusService.ListSignals(busType, serviceName, interfaceName)
}

// SubscribeToSignal subscribes to a D-Bus signal
func (h *Handler) SubscribeToSignal(c fuego.ContextNoBody) (*model.SignalSubscription, error) {
	busType := c.PathParam("busType")
	serviceName := c.PathParam("serviceName")
	interfaceName := c.PathParam("interfaceName")
	signalName := c.PathParam("signalName")
	return h.dbusService.SubscribeToSignal(busType, serviceName, interfaceName, signalName)
}

// IntrospectService returns the introspection XML for a service
func (h *Handler) IntrospectService(c fuego.ContextNoBody) (*model.IntrospectionResult, error) {
	busType := c.PathParam("busType")
	serviceName := c.PathParam("serviceName")
	return h.dbusService.IntrospectService(busType, serviceName)
}
