package token

import (
	"time"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/port"
	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenGenerator struct{}

func NewJWTTokenGenerator() port.TokenGenerator {
	return &JWTTokenGenerator{}
}

func (g *JWTTokenGenerator) GenerateToken(userID string) (string, error) {
	expiration := time.Now().Add(time.Second * time.Duration(config.AppConfig.Jwt.Expiration))
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expiration),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(config.AppConfig.Jwt.Secret))
}
