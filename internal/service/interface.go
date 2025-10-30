package service

import "github.com/mesbrj/dbus-controller/internal/model"

// DBusServiceInterface defines the interface for D-Bus operations
// This interface allows for easy mocking in tests
type DBusServiceInterface interface {
	ListServices(busType string) ([]string, error)
	GetServiceInfo(busType, serviceName string) (*model.ServiceInfo, error)
	ListInterfaces(busType, serviceName string) ([]string, error)
	GetInterfaceInfo(busType, serviceName, interfaceName string) (*model.InterfaceInfo, error)
	ListMethods(busType, serviceName, interfaceName string) ([]model.MethodInfo, error)
	CallMethod(busType, serviceName, interfaceName, methodName string, args []interface{}) (*model.MethodCallResult, error)
	ListProperties(busType, serviceName, interfaceName string) ([]model.PropertyInfo, error)
	GetProperty(busType, serviceName, interfaceName, propertyName string) (*model.PropertyValue, error)
	SetProperty(busType, serviceName, interfaceName, propertyName string, value interface{}) (*model.PropertyValue, error)
	ListSignals(busType, serviceName, interfaceName string) ([]model.SignalInfo, error)
	SubscribeToSignal(busType, serviceName, interfaceName, signalName string) (*model.SignalSubscription, error)
	IntrospectService(busType, serviceName string) (*model.IntrospectionResult, error)
	Close()
}

// Ensure DBusService implements the interface
var _ DBusServiceInterface = (*DBusService)(nil)
