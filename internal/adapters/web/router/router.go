package router

import (
	"net/http"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/middleware"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func InitRouter(e *echo.Echo) {
	e.Use(echoMiddleware.CORSWithConfig(middleware.CorsConfig()))
	e.Use(echoMiddleware.Recover(), echoMiddleware.Logger())

	apiV1 := e.Group("/api/v1")

	apiV1.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	protected := apiV1.Group("", middleware.Middleware)
	protected.GET("/profile", getUserProfile)
}

func getUserProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)
	return c.JSON(http.StatusOK, map[string]string{"message": "User profile", "user_id": userID})
}
