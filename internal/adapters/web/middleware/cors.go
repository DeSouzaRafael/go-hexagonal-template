package middleware

import (
	"github.com/DeSouzaRafael/go-hexagonal-template/pkg/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CorsConfig() middleware.CORSConfig {
	cc := middleware.CORSConfig{
		AllowHeaders:     []string{echo.HeaderAccept, echo.HeaderAcceptEncoding, echo.HeaderAuthorization, echo.HeaderContentLength, echo.HeaderContentType, echo.HeaderOrigin, echo.HeaderXCSRFToken},
		AllowCredentials: true,
		ExposeHeaders:    []string{echo.HeaderAccept, echo.HeaderAcceptEncoding, echo.HeaderAuthorization, echo.HeaderContentLength, echo.HeaderContentType, echo.HeaderOrigin, echo.HeaderXCSRFToken},
	}

	if util.CurrentExecutionEnvironmentProduction() {
		cc.AllowOrigins = []string{"https://*.yourdomain.com"}
	} else {
		cc.AllowOrigins = []string{"*"}
	}

	return cc
}
