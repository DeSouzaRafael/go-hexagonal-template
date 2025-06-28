package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	container "github.com/DeSouzaRafael/go-hexagonal-template/internal"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/handler"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/router"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func testHandler(c echo.Context) error {
	return c.String(http.StatusOK, "test")
}

func TestRoutes(t *testing.T) {
	// Create a new echo instance
	e := echo.New()
	g := e.Group("/api")

	// Create mock handlers
	authHandler := handler.AuthHandler{}
	userHandler := handler.UserHandler{}

	// Create handlers container with mock handlers
	handlers := container.Handlers{
		AuthHandler: authHandler,
		UserHandler: userHandler,
	}

	// Setup the routers
	router.AuthRouter(g, handlers)
	router.UserRouter(g, handlers)

	// Get all routes
	routes := e.Routes()

	// Expected routes
	expectedRoutes := map[string]bool{
		"POST /api/v1/auth/register": false,
		"POST /api/v1/auth/login":    false,
		"GET /api/v1/auth/profile":   false,
		"GET /api/v1/users":          false,
		"GET /api/v1/users/:id":      false,
		"PUT /api/v1/users/:id":      false,
		"DELETE /api/v1/users/:id":   false,
	}

	// Check that all expected routes are registered
	for _, route := range routes {
		path := route.Method + " " + route.Path
		if _, exists := expectedRoutes[path]; exists {
			expectedRoutes[path] = true
		}
	}

	// Verify all expected routes were found
	for route, found := range expectedRoutes {
		assert.True(t, found, "Route %s was not registered", route)
	}
}

func TestGetUserProfile(t *testing.T) {
	// Create a new echo instance
	e := echo.New()

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set user_id in context
	c.Set("user_id", "test-user-id")

	// Call the handler
	err := router.GetUserProfile(c)

	// Assert no error
	assert.NoError(t, err)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "test-user-id")
}
