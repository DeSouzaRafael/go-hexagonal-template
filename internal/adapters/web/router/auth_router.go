package router

import (
	"net/http"

	container "github.com/DeSouzaRafael/go-hexagonal-template/internal"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/middleware"
	"github.com/labstack/echo/v4"
)

func AuthRouter(group *echo.Group, handlers container.Handlers) {

	// Set up the auth routes
	authV1 := group.Group("/v1/auth")

	// Set up the auth routes that require authentication
	protectedAuthV1 := authV1.Group("", middleware.Middleware)

	// Set up the auth routes
	protectedAuthV1.GET("/profile", getUserProfile)
	authV1.POST("/register", handlers.AuthHandler.Register)
	authV1.POST("/login", handlers.AuthHandler.Login)
}

func getUserProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)
	return c.JSON(http.StatusOK, map[string]string{"message": "User profile", "user_id": userID})
}
