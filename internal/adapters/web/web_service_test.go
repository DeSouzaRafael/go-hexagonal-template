package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	container "github.com/DeSouzaRafael/go-hexagonal-template/internal"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/handler"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewWebService(t *testing.T) {
	// Setup
	config.AppConfig = &config.Config{
		WebService: config.WebService{
			Port: "8080",
		},
	}

	// Create mock handlers
	handlers := container.Handlers{
		AuthHandler: handler.AuthHandler{},
		UserHandler: handler.UserHandler{},
	}

	// Create web service
	webService := web.NewWebService(handlers)

	// Assert
	assert.NotNil(t, webService)
	assert.NotNil(t, webService.Echo)

	// Test that routes are registered
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/profile", nil)
	rec := httptest.NewRecorder()
	webService.Echo.ServeHTTP(rec, req)

	// We expect a 401 Unauthorized because the route exists but requires authentication
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// Test that non-existent routes return 404
	req = httptest.NewRequest(http.MethodGet, "/non-existent-route", nil)
	rec = httptest.NewRecorder()
	webService.Echo.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestWebService_Start(t *testing.T) {

	// Setup
	config.AppConfig = &config.Config{
		WebService: config.WebService{
			Port: "0", // Use port 0 to let the OS choose an available port
		},
	}

	// Create mock handlers
	handlers := container.Handlers{
		AuthHandler: handler.AuthHandler{},
		UserHandler: handler.UserHandler{},
	}

	// Create web service
	webService := web.NewWebService(handlers)

	// Start the server in a goroutine
	go func() {
		err := webService.Start()
		// We expect an error when the server is closed
		assert.Error(t, err)
	}()

	// Close the server after a short delay
	go func() {
		webService.Echo.Close()
	}()
}
