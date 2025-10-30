package api

import (
	"github.com/go-fuego/fuego"
	"github.com/mesbrj/dbus-controller/internal/handler"
	"github.com/mesbrj/dbus-controller/internal/service"
)

// SetupRoutes configures all API routes
func SetupRoutes(s *fuego.Server, dbusService service.DBusServiceInterface) {
	h := handler.NewHandler(dbusService)

	// Bus management routes
	fuego.Get(s, "/buses", h.ListBuses)
	fuego.Get(s, "/buses/{busType}", h.GetBusInfo)

	// Service routes
	fuego.Get(s, "/buses/{busType}/services", h.ListServices)
	fuego.Get(s, "/buses/{busType}/services/{serviceName}", h.GetService)

	// Interface routes
	fuego.Get(s, "/buses/{busType}/services/{serviceName}/interfaces", h.ListInterfaces)
	fuego.Get(s, "/buses/{busType}/services/{serviceName}/interfaces/{interfaceName}", h.GetInterface)

	// Method routes
	fuego.Get(s, "/buses/{busType}/services/{serviceName}/interfaces/{interfaceName}/methods", h.ListMethods)
	fuego.Post(s, "/buses/{busType}/services/{serviceName}/interfaces/{interfaceName}/methods/{methodName}/call", h.CallMethod)

	// Property routes
	fuego.Get(s, "/buses/{busType}/services/{serviceName}/interfaces/{interfaceName}/properties", h.ListProperties)
	fuego.Get(s, "/buses/{busType}/services/{serviceName}/interfaces/{interfaceName}/properties/{propertyName}", h.GetProperty)
	fuego.Put(s, "/buses/{busType}/services/{serviceName}/interfaces/{interfaceName}/properties/{propertyName}", h.SetProperty)

	// Signal routes
	fuego.Get(s, "/buses/{busType}/services/{serviceName}/interfaces/{interfaceName}/signals", h.ListSignals)
	fuego.Post(s, "/buses/{busType}/services/{serviceName}/interfaces/{interfaceName}/signals/{signalName}/subscribe", h.SubscribeToSignal)

	// Introspection routes
	fuego.Get(s, "/buses/{busType}/services/{serviceName}/introspect", h.IntrospectService)
}
