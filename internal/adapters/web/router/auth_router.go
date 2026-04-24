package router

import (
	container "github.com/DeSouzaRafael/go-hexagonal-template/internal"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/middleware"
	"github.com/labstack/echo/v4"
)

func AuthRouter(group *echo.Group, handlers container.Handlers) {
	authV1 := group.Group("/v1/auth")

	protectedAuthV1 := authV1.Group("", middleware.Middleware)
	protectedAuthV1.GET("/profile", handlers.AuthHandler.GetProfile)

	authV1.POST("/register", handlers.AuthHandler.Register)
	authV1.POST("/login", handlers.AuthHandler.Login)
}
