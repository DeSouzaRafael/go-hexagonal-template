package config_test

import (
	"os"
	"testing"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigWithMissingEnvFile(t *testing.T) {
	// Clean environment variables
	os.Clearenv()

	// Test loading without .env file
	err := config.LoadConfig()
	assert.NoError(t, err)

	// Config should be loaded with empty values
	assert.NotNil(t, config.AppConfig)
	assert.Empty(t, config.AppConfig.App.Name)
	assert.Empty(t, config.AppConfig.App.Environment)
	assert.Empty(t, config.AppConfig.WebService.Port)
	assert.Empty(t, config.AppConfig.WebService.Domain)
	assert.Empty(t, config.AppConfig.Jwt.Secret)
	assert.Zero(t, config.AppConfig.Jwt.Expiration)
	assert.Empty(t, config.AppConfig.Database.Host)
	assert.Empty(t, config.AppConfig.Database.Port)
	assert.Empty(t, config.AppConfig.Database.User)
	assert.Empty(t, config.AppConfig.Database.Pass)
	assert.Empty(t, config.AppConfig.Database.DBName)
	assert.Empty(t, config.AppConfig.Database.SSLMode)
	assert.Zero(t, config.AppConfig.Database.LogLevel)
}
