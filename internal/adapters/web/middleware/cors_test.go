package middleware_test

import (
	"testing"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/middleware"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestCorsConfig(t *testing.T) {
	t.Run("development allows all origins", func(t *testing.T) {
		config.AppConfig = &config.Config{
			App:        config.App{Environment: "development"},
			WebService: config.WebService{Domain: "localhost"},
		}
		cc := middleware.CorsConfig()
		assert.Equal(t, []string{"*"}, cc.AllowOrigins)
	})

	t.Run("production restricts to domain", func(t *testing.T) {
		config.AppConfig = &config.Config{
			App:        config.App{Environment: "production"},
			WebService: config.WebService{Domain: "example.com"},
		}
		cc := middleware.CorsConfig()
		assert.Equal(t, []string{"https://*.example.com"}, cc.AllowOrigins)
	})
}
