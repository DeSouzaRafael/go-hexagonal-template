package token

import (
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/middleware"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/port"
)

// JWTTokenGenerator implements the TokenGenerator interface using JWT
type JWTTokenGenerator struct{}

// NewJWTTokenGenerator creates a new JWTTokenGenerator
func NewJWTTokenGenerator() port.TokenGenerator {
	return &JWTTokenGenerator{}
}

// GenerateToken generates a JWT token for the given user ID
func (g *JWTTokenGenerator) GenerateToken(userID string) (string, error) {
	return middleware.GenerateJWT(userID)
}
