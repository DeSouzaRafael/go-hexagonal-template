package middleware

import (
	"net/http"
	"strings"

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

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, domain.ErrMissingToken)
		}

		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return []byte(config.AppConfig.Jwt.Secret), nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, domain.ErrInvalidToken)
		}

		c.Set("user_id", claims.Subject)

		return next(c)
	}
}
