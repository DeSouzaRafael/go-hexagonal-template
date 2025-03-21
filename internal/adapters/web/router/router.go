package router

import (
	"net/http"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/handler"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/middleware"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func InitRouter(
	e *echo.Echo,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
) {
	e.Use(echoMiddleware.CORSWithConfig(middleware.CorsConfig()))
	e.Use(echoMiddleware.Recover(), echoMiddleware.Logger())

	v1 := e.Group("/api/v1")
	authV1 := v1.Group("", middleware.Middleware)
	authV1.GET("/profile", getUserProfile)

	user := authV1.Group("/users")
	user.GET("", userHandler.ListUsers)
	user.GET("/:id", userHandler.GetUser)
	user.PUT("/:id", userHandler.UpdateUser)
	user.DELETE("/:id", userHandler.DeleteUser)

	auth := v1.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
}

func getUserProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)
	return c.JSON(http.StatusOK, map[string]string{"message": "User profile", "user_id": userID})
}
