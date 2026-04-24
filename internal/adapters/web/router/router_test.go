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

func TestRoutes(t *testing.T) {
	e := echo.New()
	g := e.Group("/api")

	authHandler := handler.AuthHandler{}
	userHandler := handler.UserHandler{}

	handlers := container.Handlers{
		AuthHandler: authHandler,
		UserHandler: userHandler,
	}

	router.AuthRouter(g, handlers)
	router.UserRouter(g, handlers)

	routes := e.Routes()

	expectedRoutes := map[string]bool{
		"POST /api/v1/auth/register": false,
		"POST /api/v1/auth/login":    false,
		"GET /api/v1/auth/profile":   false,
		"GET /api/v1/users":          false,
		"GET /api/v1/users/:id":      false,
		"PUT /api/v1/users/:id":      false,
		"DELETE /api/v1/users/:id":   false,
	}

	for _, route := range routes {
		path := route.Method + " " + route.Path
		if _, exists := expectedRoutes[path]; exists {
			expectedRoutes[path] = true
		}
	}

	for route, found := range expectedRoutes {
		assert.True(t, found, "Route %s was not registered", route)
	}
}

func TestGetProfile(t *testing.T) {
	e := echo.New()
	h := handler.AuthHandler{}

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "test-user-id")

	err := h.GetProfile(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "test-user-id")
}
