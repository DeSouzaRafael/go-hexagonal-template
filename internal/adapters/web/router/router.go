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
) {
	e.Use(echoMiddleware.CORSWithConfig(middleware.CorsConfig()))
	e.Use(echoMiddleware.Recover(), echoMiddleware.Logger())

	v1 := e.Group("/api/v1")
	authV1 := v1.Group("", middleware.Middleware)
	authV1.GET("/profile", getUserProfile)

	user := v1.Group("/users")
	authUser := authV1.Group("/users")
	user.POST("/", userHandler.Register)
	authUser.GET("/", userHandler.ListUsers)
	authUser.GET("/:id", userHandler.GetUser)
	authUser.PUT("/:id", userHandler.UpdateUser)
	authUser.DELETE("/:id", userHandler.DeleteUser)
}

func getUserProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)
	return c.JSON(http.StatusOK, map[string]string{"message": "User profile", "user_id": userID})
}
