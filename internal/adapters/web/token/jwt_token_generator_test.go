package token_test

import (
	"testing"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/token"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWTTokenGenerator_GenerateToken(t *testing.T) {
	// Setup test config
	config.AppConfig = &config.Config{
		Jwt: config.JWT{
			Secret:     "test-secret-key",
			Expiration: 3600,
		},
	}

	generator := token.NewJWTTokenGenerator()
	userID := uuid.New().String()

	// Test token generation
	tokenString, err := generator.GenerateToken(userID)
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
