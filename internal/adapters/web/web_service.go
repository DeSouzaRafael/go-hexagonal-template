package web

import (
	container "github.com/DeSouzaRafael/go-hexagonal-template/internal"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/middleware"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/router"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type WebService struct {
	Echo *echo.Echo
}

func NewWebService(h container.Handlers) *WebService {
	e := echo.New()

	// Middlewares configuration for the echo instanc
	e.Use(echoMiddleware.CORSWithConfig(middleware.CorsConfig()))
	e.Use(echoMiddleware.Recover(), echoMiddleware.Logger())

	// Grouping routes
	g := e.Group("/api")

	// Set up the routes
	router.AuthRouter(g, h)
	router.UserRouter(g, h)

	return &WebService{Echo: e}
}

func (s *WebService) Start() error {
	return s.Echo.Start(":" + config.AppConfig.WebService.Port)
}
