package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigWithUnreadableEnvFile(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")
	require.NoError(t, os.WriteFile(envFile, []byte("APP_NAME=test\n"), 0000))

	orig, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmpDir))
	defer func() { _ = os.Chdir(orig) }()

	err = config.LoadConfig()
	assert.Error(t, err)
}

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
