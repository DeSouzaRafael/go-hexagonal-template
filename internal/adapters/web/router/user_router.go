package router

import (
	container "github.com/DeSouzaRafael/go-hexagonal-template/internal"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/middleware"
	"github.com/labstack/echo/v4"
)

func UserRouter(group *echo.Group, handlers container.Handlers) {

	// Set up the user routes
	userV1 := group.Group("/v1/users")

	// Set up the user routes that require authentication
	userV1Auth := userV1.Group("", middleware.Middleware)

	// Set up the user routes
	userV1Auth.GET("", handlers.UserHandler.ListUsers)
	userV1Auth.GET("/:id", handlers.UserHandler.GetUser)
	userV1Auth.PUT("/:id", handlers.UserHandler.UpdateUser)
	userV1Auth.DELETE("/:id", handlers.UserHandler.DeleteUser)
}
