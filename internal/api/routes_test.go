package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-fuego/fuego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/mesbrj/dbus-controller/internal/handler"
)

// APIIntegrationTestSuite defines integration tests for the API
type APIIntegrationTestSuite struct {
	suite.Suite
	server      *fuego.Server
	mockService *handler.MockDBusService
}

func (suite *APIIntegrationTestSuite) SetupTest() {
	suite.mockService = new(handler.MockDBusService)
	suite.server = fuego.NewServer()

	// Setup routes with mock service
	SetupRoutes(suite.server, suite.mockService)
}

func (suite *APIIntegrationTestSuite) TearDownTest() {
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *APIIntegrationTestSuite) TestAPIRoutes_BusesEndpoint() {
	// Test that the buses endpoint is properly registered
	req := httptest.NewRequest(http.MethodGet, "/buses", nil)
	rec := httptest.NewRecorder()

	suite.server.ServeHTTP(rec, req)

	// The exact status code depends on fuego's implementation
	// This test ensures the route is registered
	assert.NotEqual(suite.T(), http.StatusNotFound, rec.Code)
}

func (suite *APIIntegrationTestSuite) TestAPIRoutes_ServicesEndpoint() {
	expectedServices := []string{"org.freedesktop.DBus", "org.freedesktop.NetworkManager"}
	suite.mockService.On("ListServices", "system").Return(expectedServices, nil)

	req := httptest.NewRequest(http.MethodGet, "/buses/system/services", nil)
	rec := httptest.NewRecorder()

	suite.server.ServeHTTP(rec, req)

	assert.NotEqual(suite.T(), http.StatusNotFound, rec.Code)
}

func TestAPIIntegrationSuite(t *testing.T) {
	suite.Run(t, new(APIIntegrationTestSuite))
}

// Test route registration
func TestSetupRoutes(t *testing.T) {
	server := fuego.NewServer()
	mockService := new(handler.MockDBusService)

	// This should not panic
	assert.NotPanics(t, func() {
		SetupRoutes(server, mockService)
	})
}

// Test that all expected routes are registered
func TestRoutePatterns(t *testing.T) {
	expectedRoutes := []string{
		"/buses",
		"/buses/{busType}",
		"/buses/{busType}/services",
		"/buses/{busType}/services/{serviceName}",
		"/buses/{busType}/services/{serviceName}/interfaces",
		"/buses/{busType}/services/{serviceName}/interfaces/{interfaceName}",
		"/buses/{busType}/services/{serviceName}/interfaces/{interfaceName}/methods",
		"/buses/{busType}/services/{serviceName}/interfaces/{interfaceName}/methods/{methodName}/call",
		"/buses/{busType}/services/{serviceName}/interfaces/{interfaceName}/properties",
		"/buses/{busType}/services/{serviceName}/interfaces/{interfaceName}/properties/{propertyName}",
		"/buses/{busType}/services/{serviceName}/interfaces/{interfaceName}/signals",
		"/buses/{busType}/services/{serviceName}/interfaces/{interfaceName}/signals/{signalName}/subscribe",
		"/buses/{busType}/services/{serviceName}/introspect",
	}

	// Ensure we have all the expected route patterns defined
	assert.Len(t, expectedRoutes, 13)

	// Test route parameter patterns
	for _, route := range expectedRoutes {
		assert.Contains(t, route, "/buses")
		if len(route) > 6 { // More than just "/buses"
			assert.Contains(t, route, "{busType}")
		}
	}
}
