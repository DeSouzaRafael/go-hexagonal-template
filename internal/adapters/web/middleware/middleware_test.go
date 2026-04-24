package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/token"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/middleware"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTest() {
	config.AppConfig = &config.Config{
		Jwt: config.JWT{
			Secret:     "test-secret-key",
			Expiration: 3600,
		},
	}
}

func TestMiddleware(t *testing.T) {
	setupTest()
	e := echo.New()

	testHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	e.GET("/protected", testHandler, middleware.Middleware)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	userID := uuid.New().String()
	generator := token.NewJWTTokenGenerator()
	validToken, _ := generator.GenerateToken(userID)

	e.GET("/check-context", func(c echo.Context) error {
		userIDFromContext := c.Get("user_id")
		if userIDFromContext == nil {
			return c.String(http.StatusInternalServerError, "user_id not set in context")
		}
		return c.String(http.StatusOK, userIDFromContext.(string))
	}, middleware.Middleware)

	req = httptest.NewRequest(http.MethodGet, "/check-context", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, userID, rec.Body.String())
}
