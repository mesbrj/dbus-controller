package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBusInfo(t *testing.T) {
	bus := BusInfo{
		Type:        "system",
		Description: "System D-Bus",
	}

	assert.Equal(t, "system", bus.Type)
	assert.Equal(t, "System D-Bus", bus.Description)
}

func TestServiceInfo(t *testing.T) {
	service := ServiceInfo{
		Name:        "org.freedesktop.DBus",
		Owner:       ":1.0",
		Interfaces:  []string{"org.freedesktop.DBus"},
		ObjectPaths: []string{"/"},
	}

	assert.Equal(t, "org.freedesktop.DBus", service.Name)
	assert.Equal(t, ":1.0", service.Owner)
	assert.Len(t, service.Interfaces, 1)
	assert.Len(t, service.ObjectPaths, 1)
}

func TestInterfaceInfo(t *testing.T) {
	iface := InterfaceInfo{
		Name: "org.freedesktop.DBus",
		Methods: []MethodInfo{
			{
				Name: "Hello",
				OutArgs: []ArgumentInfo{
					{Name: "unique_name", Type: "s", Direction: "out"},
				},
			},
		},
		Properties: []PropertyInfo{
			{Name: "Features", Type: "as", Access: "read"},
		},
		Signals: []SignalInfo{
			{
				Name: "NameOwnerChanged",
				Args: []ArgumentInfo{
					{Name: "name", Type: "s", Direction: "out"},
					{Name: "old_owner", Type: "s", Direction: "out"},
					{Name: "new_owner", Type: "s", Direction: "out"},
				},
			},
		},
	}

	assert.Equal(t, "org.freedesktop.DBus", iface.Name)
	assert.Len(t, iface.Methods, 1)
	assert.Len(t, iface.Properties, 1)
	assert.Len(t, iface.Signals, 1)

	// Test method
	method := iface.Methods[0]
	assert.Equal(t, "Hello", method.Name)
	assert.Len(t, method.OutArgs, 1)
	assert.Equal(t, "s", method.OutArgs[0].Type)

	// Test property
	property := iface.Properties[0]
	assert.Equal(t, "Features", property.Name)
	assert.Equal(t, "read", property.Access)

	// Test signal
	signal := iface.Signals[0]
	assert.Equal(t, "NameOwnerChanged", signal.Name)
	assert.Len(t, signal.Args, 3)
}

func TestMethodCallResult(t *testing.T) {
	now := time.Now()
	result := MethodCallResult{
		Success:      true,
		ReturnValues: []interface{}{"test_value"},
		Timestamp:    now,
	}

	assert.True(t, result.Success)
	assert.Len(t, result.ReturnValues, 1)
	assert.Equal(t, "test_value", result.ReturnValues[0])
	assert.Equal(t, now, result.Timestamp)
	assert.Empty(t, result.Error)
}

func TestMethodCallResult_WithError(t *testing.T) {
	now := time.Now()
	result := MethodCallResult{
		Success:   false,
		Error:     "Method call failed",
		Timestamp: now,
	}

	assert.False(t, result.Success)
	assert.Equal(t, "Method call failed", result.Error)
	assert.Empty(t, result.ReturnValues)
	assert.Equal(t, now, result.Timestamp)
}

func TestPropertyValue(t *testing.T) {
	now := time.Now()
	prop := PropertyValue{
		Name:      "TestProperty",
		Type:      "s",
		Value:     "test_value",
		Timestamp: now,
	}

	assert.Equal(t, "TestProperty", prop.Name)
	assert.Equal(t, "s", prop.Type)
	assert.Equal(t, "test_value", prop.Value)
	assert.Equal(t, now, prop.Timestamp)
}

func TestSignalSubscription(t *testing.T) {
	now := time.Now()
	subscription := SignalSubscription{
		ID:        "sub_123",
		BusType:   "system",
		Service:   "org.freedesktop.DBus",
		Interface: "org.freedesktop.DBus",
		Signal:    "NameOwnerChanged",
		Active:    true,
		CreatedAt: now,
	}

	assert.Equal(t, "sub_123", subscription.ID)
	assert.Equal(t, "system", subscription.BusType)
	assert.Equal(t, "org.freedesktop.DBus", subscription.Service)
	assert.Equal(t, "org.freedesktop.DBus", subscription.Interface)
	assert.Equal(t, "NameOwnerChanged", subscription.Signal)
	assert.True(t, subscription.Active)
	assert.Equal(t, now, subscription.CreatedAt)
}

func TestIntrospectionResult(t *testing.T) {
	now := time.Now()
	xmlData := `<?xml version="1.0"?><node></node>`

	result := IntrospectionResult{
		Service:    "org.freedesktop.DBus",
		ObjectPath: "/",
		XML:        xmlData,
		Timestamp:  now,
	}

	assert.Equal(t, "org.freedesktop.DBus", result.Service)
	assert.Equal(t, "/", result.ObjectPath)
	assert.Equal(t, xmlData, result.XML)
	assert.Equal(t, now, result.Timestamp)
}

func TestParsedIntrospection(t *testing.T) {
	parsed := ParsedIntrospection{
		Interfaces: []InterfaceInfo{
			{Name: "org.freedesktop.DBus"},
		},
		Nodes: []NodeInfo{
			{Name: "org", Path: "/org"},
		},
	}

	assert.Len(t, parsed.Interfaces, 1)
	assert.Equal(t, "org.freedesktop.DBus", parsed.Interfaces[0].Name)
	assert.Len(t, parsed.Nodes, 1)
	assert.Equal(t, "org", parsed.Nodes[0].Name)
	assert.Equal(t, "/org", parsed.Nodes[0].Path)
}

func TestErrorResponse(t *testing.T) {
	err := ErrorResponse{
		Error:   "BadRequest",
		Message: "Invalid bus type",
		Code:    400,
	}

	assert.Equal(t, "BadRequest", err.Error)
	assert.Equal(t, "Invalid bus type", err.Message)
	assert.Equal(t, 400, err.Code)
}

// Test argument info validation
func TestArgumentInfo_Direction(t *testing.T) {
	inArg := ArgumentInfo{
		Name:      "input",
		Type:      "s",
		Direction: "in",
	}

	outArg := ArgumentInfo{
		Name:      "output",
		Type:      "i",
		Direction: "out",
	}

	assert.Equal(t, "in", inArg.Direction)
	assert.Equal(t, "out", outArg.Direction)

	// Test valid directions
	validDirections := []string{"in", "out"}
	assert.Contains(t, validDirections, inArg.Direction)
	assert.Contains(t, validDirections, outArg.Direction)
}

// Test property access validation
func TestPropertyInfo_Access(t *testing.T) {
	readProp := PropertyInfo{
		Name:   "ReadOnly",
		Type:   "s",
		Access: "read",
	}

	writeProp := PropertyInfo{
		Name:   "WriteOnly",
		Type:   "s",
		Access: "write",
	}

	readWriteProp := PropertyInfo{
		Name:   "ReadWrite",
		Type:   "s",
		Access: "readwrite",
	}

	validAccess := []string{"read", "write", "readwrite"}
	assert.Contains(t, validAccess, readProp.Access)
	assert.Contains(t, validAccess, writeProp.Access)
	assert.Contains(t, validAccess, readWriteProp.Access)
}
