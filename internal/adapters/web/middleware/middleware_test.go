package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/middleware"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/golang-jwt/jwt/v5"
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

func TestGenerateJWT(t *testing.T) {
	setupTest()
	userID := uuid.New().String()

	// Test token generation
	tokenString, err := middleware.GenerateJWT(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Verify the token
	claims := &jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret-key"), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	assert.Equal(t, userID, claims.Subject)
}

func TestMiddleware(t *testing.T) {
	setupTest()
	e := echo.New()

	// Create a test handler that will be protected by the middleware
	testHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	// Register the handler with the middleware
	e.GET("/protected", testHandler, middleware.Middleware)

	// Test with missing token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Check status code only
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// Test with invalid token format
	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Check status code only
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// Test with invalid token
	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Check status code only
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// Test with valid token
	userID := uuid.New().String()
	validToken, _ := middleware.GenerateJWT(userID)
	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	rec = httptest.NewRecorder()

	// Create a custom handler that checks if user_id is set in context
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
