package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, domain.ErrMissingToken)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, domain.ErrTokenDuration)
		}

		tokenString := parts[1]
		claims := &jwt.RegisteredClaims{}

		secretKey := []byte(config.AppConfig.Jwt.Secret)

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, domain.ErrInvalidToken)
		}

		if claims.ExpiresAt.Before(time.Now()) {
			return c.JSON(http.StatusUnauthorized, domain.ErrExpiredToken)
		}

		c.Set("user_id", claims.Subject)

		return next(c)
	}
}

func GenerateJWT(userID string) (string, error) {
	expirationTime := time.Now().Add(time.Second * time.Duration(config.AppConfig.Jwt.Expiration))
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	secretKey := []byte(config.AppConfig.Jwt.Secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
