package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var secretKey = []byte(config.AppConfig.Jwt.Secret)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing Token"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Badly formatted token"})
		}

		tokenString := parts[1]
		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		if claims.ExpiresAt.Before(time.Now()) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Expired token"})
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
