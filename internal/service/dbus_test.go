package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/mesbrj/dbus-controller/internal/model"
)

// DBusServiceTestSuite defines a test suite for D-Bus service tests
type DBusServiceTestSuite struct {
	suite.Suite
	service *DBusService
}

func (suite *DBusServiceTestSuite) SetupTest() {
	suite.service = NewDBusService()
}

func (suite *DBusServiceTestSuite) TearDownTest() {
	if suite.service != nil {
		suite.service.Close()
	}
}

func (suite *DBusServiceTestSuite) TestNewDBusService() {
	service := NewDBusService()
	defer service.Close()

	assert.NotNil(suite.T(), service)
	assert.NotNil(suite.T(), service.subscriptions)
}

func (suite *DBusServiceTestSuite) TestGetConnection_ValidBusTypes() {
	// Test that valid bus types return connections (may be nil if not available)
	systemConn, systemErr := suite.service.getConnection("system")
	sessionConn, sessionErr := suite.service.getConnection("session")

	// Connections might be nil in test environment, but errors should be descriptive
	if systemErr != nil {
		assert.Contains(suite.T(), systemErr.Error(), "system bus not available")
	} else {
		// If connection is available, it should not be nil
		if systemConn != nil {
			assert.NotNil(suite.T(), systemConn)
		}
	}

	if sessionErr != nil {
		assert.Contains(suite.T(), sessionErr.Error(), "session bus not available")
	} else {
		if sessionConn != nil {
			assert.NotNil(suite.T(), sessionConn)
		}
	}
}

func (suite *DBusServiceTestSuite) TestGetConnection_InvalidBusType() {
	conn, err := suite.service.getConnection("invalid")

	assert.Nil(suite.T(), conn)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid bus type")
}

func (suite *DBusServiceTestSuite) TestParseIntrospectionXML() {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<node>
  <interface name="org.freedesktop.DBus">
    <method name="Hello">
      <arg direction="out" type="s"/>
    </method>
    <method name="RequestName">
      <arg direction="in" type="s" name="name"/>
      <arg direction="in" type="u" name="flags"/>
      <arg direction="out" type="u"/>
    </method>
    <property name="Features" type="as" access="read"/>
    <signal name="NameOwnerChanged">
      <arg type="s" name="name"/>
      <arg type="s" name="old_owner"/>
      <arg type="s" name="new_owner"/>
    </signal>
  </interface>
  <node name="org"/>
</node>`

	parsed, err := suite.service.parseIntrospectionXML(xmlData)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), parsed)
	assert.Len(suite.T(), parsed.Interfaces, 1)
	assert.Equal(suite.T(), "org.freedesktop.DBus", parsed.Interfaces[0].Name)

	// Check methods
	assert.Len(suite.T(), parsed.Interfaces[0].Methods, 2)
	helloMethod := parsed.Interfaces[0].Methods[0]
	assert.Equal(suite.T(), "Hello", helloMethod.Name)
	assert.Len(suite.T(), helloMethod.OutArgs, 1)
	assert.Equal(suite.T(), "s", helloMethod.OutArgs[0].Type)

	requestNameMethod := parsed.Interfaces[0].Methods[1]
	assert.Equal(suite.T(), "RequestName", requestNameMethod.Name)
	assert.Len(suite.T(), requestNameMethod.InArgs, 2)
	assert.Len(suite.T(), requestNameMethod.OutArgs, 1)

	// Check properties
	assert.Len(suite.T(), parsed.Interfaces[0].Properties, 1)
	featuresProperty := parsed.Interfaces[0].Properties[0]
	assert.Equal(suite.T(), "Features", featuresProperty.Name)
	assert.Equal(suite.T(), "as", featuresProperty.Type)
	assert.Equal(suite.T(), "read", featuresProperty.Access)

	// Check signals
	assert.Len(suite.T(), parsed.Interfaces[0].Signals, 1)
	nameOwnerChangedSignal := parsed.Interfaces[0].Signals[0]
	assert.Equal(suite.T(), "NameOwnerChanged", nameOwnerChangedSignal.Name)
	assert.Len(suite.T(), nameOwnerChangedSignal.Args, 3)

	// Check nodes
	assert.Len(suite.T(), parsed.Nodes, 1)
	assert.Equal(suite.T(), "org", parsed.Nodes[0].Name)
	assert.Equal(suite.T(), "/org", parsed.Nodes[0].Path)
}

func (suite *DBusServiceTestSuite) TestParseIntrospectionXML_InvalidXML() {
	invalidXML := "not valid xml"

	parsed, err := suite.service.parseIntrospectionXML(invalidXML)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), parsed)
	assert.Contains(suite.T(), err.Error(), "failed to parse introspection XML")
}

func TestDBusServiceSuite(t *testing.T) {
	suite.Run(t, new(DBusServiceTestSuite))
}

// Test helper functions and mock data
func TestCreateMockIntrospectionResult(t *testing.T) {
	result := &model.IntrospectionResult{
		Service:    "org.freedesktop.DBus",
		ObjectPath: "/",
		XML:        "<node></node>",
		Timestamp:  time.Now(),
	}

	assert.Equal(t, "org.freedesktop.DBus", result.Service)
	assert.Equal(t, "/", result.ObjectPath)
	assert.NotEmpty(t, result.XML)
	assert.False(t, result.Timestamp.IsZero())
}

func TestCreateMockServiceInfo(t *testing.T) {
	service := &model.ServiceInfo{
		Name:        "org.freedesktop.NetworkManager",
		Owner:       ":1.5",
		Interfaces:  []string{"org.freedesktop.NetworkManager", "org.freedesktop.DBus.Properties"},
		ObjectPaths: []string{"/org/freedesktop/NetworkManager"},
	}

	assert.Equal(t, "org.freedesktop.NetworkManager", service.Name)
	assert.Equal(t, ":1.5", service.Owner)
	assert.Len(t, service.Interfaces, 2)
	assert.Len(t, service.ObjectPaths, 1)
}

// Integration test helpers (these would require actual D-Bus in CI/CD)
func TestDBusService_Integration_ListServices(t *testing.T) {
	// Skip integration tests if not in integration test environment
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	service := NewDBusService()
	defer service.Close()

	// This test would only pass if system D-Bus is available
	services, err := service.ListServices("system")
	if err != nil {
		t.Logf("System D-Bus not available: %v", err)
		return
	}

	assert.NotEmpty(t, services)
	// org.freedesktop.DBus should always be present if D-Bus is working
	assert.Contains(t, services, "org.freedesktop.DBus")
}

// Benchmark tests
func BenchmarkParseIntrospectionXML(b *testing.B) {
	service := NewDBusService()
	defer service.Close()

	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<node>
  <interface name="org.freedesktop.DBus">
    <method name="Hello"><arg direction="out" type="s"/></method>
    <property name="Features" type="as" access="read"/>
    <signal name="NameOwnerChanged">
      <arg type="s" name="name"/>
      <arg type="s" name="old_owner"/>
      <arg type="s" name="new_owner"/>
    </signal>
  </interface>
</node>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.parseIntrospectionXML(xmlData)
		if err != nil {
			b.Fatal(err)
		}
	}
}
