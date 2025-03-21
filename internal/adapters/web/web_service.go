package web

import (
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/handler"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/router"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/labstack/echo/v4"
)

type WebService struct {
	Echo *echo.Echo
}

func NewWebService(uh *handler.UserHandler, ah *handler.AuthHandler) *WebService {
	e := echo.New()

	router.InitRouter(e, uh, ah)

	return &WebService{Echo: e}
}

func (s *WebService) Start() error {
	return s.Echo.Start(":" + config.AppConfig.WebService.Port)
}
