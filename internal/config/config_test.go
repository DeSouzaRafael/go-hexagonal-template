package config_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
)

const sampleConfig = `database:
  host: "localhost"
  port: "5432"
  user: "testuser"
  pass: "testpass"
  dbname: "testdb"
  sslmode: "disable"
  loglevel: 1
jwt:
  secret: "mysecret"
  expiration: 3600
environment: "test"
webservice:
  port: "8080"
  domain: "yourdomain.com"
`

func getTestConfigPath(t *testing.T) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get current file path")
	}
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "config.yml")
}

func createTempConfigFile(t *testing.T, configPath string) {
	err := os.WriteFile(configPath, []byte(sampleConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
}

func deleteTempConfigFile(t *testing.T, configPath string) {
	err := os.Remove(configPath)
	if err != nil {
		t.Fatalf("Failed to delete temp config file: %v", err)
	}
}

func validateConfigValues(t *testing.T) {
	// Validate environment
	if config.AppConfig.Environment != "test" {
		t.Errorf("Expected environment 'test', got '%s'", config.AppConfig.Environment)
	}

	// Validate database configuration
	validateDBConfig(t)

	// Validate JWT configuration
	validateJWTConfig(t)

	// Validate webservice configuration
	if config.AppConfig.WebService.Port != "8080" {
		t.Errorf("Expected webservice.port '8080', got '%s'", config.AppConfig.WebService.Port)
	}
}

func validateDBConfig(t *testing.T) {
	if config.AppConfig.Database.Host != "localhost" {
		t.Errorf("Expected database.host 'localhost', got '%s'", config.AppConfig.Database.Host)
	}

	if config.AppConfig.Database.Port != "5432" {
		t.Errorf("Expected database.port '5432', got '%s'", config.AppConfig.Database.Port)
	}

	if config.AppConfig.Database.User != "testuser" {
		t.Errorf("Expected database.user 'testuser', got '%s'", config.AppConfig.Database.User)
	}

	if config.AppConfig.Database.Pass != "testpass" {
		t.Errorf("Expected database.pass 'testpass', got '%s'", config.AppConfig.Database.Pass)
	}

	if config.AppConfig.Database.DBName != "testdb" {
		t.Errorf("Expected database.dbname 'testdb', got '%s'", config.AppConfig.Database.DBName)
	}

	if config.AppConfig.Database.SSLMode != "disable" {
		t.Errorf("Expected database.sslmode 'disable', got '%s'", config.AppConfig.Database.SSLMode)
	}

	if config.AppConfig.Database.LogLevel != 1 {
		t.Errorf("Expected database.loglevel 1, got %d", config.AppConfig.Database.LogLevel)
	}
}

func validateJWTConfig(t *testing.T) {
	if config.AppConfig.Jwt.Secret != "mysecret" {
		t.Errorf("Expected jwt.secret 'mysecret', got '%s'", config.AppConfig.Jwt.Secret)
	}

	if config.AppConfig.Jwt.Expiration != 3600 {
		t.Errorf("Expected jwt.expiration 3600, got %d", config.AppConfig.Jwt.Expiration)
	}
}

func TestInitConfig(t *testing.T) {
	configPath := getTestConfigPath(t)

	// Create temp config file
	createTempConfigFile(t, configPath)
	defer deleteTempConfigFile(t, configPath)

	// Initialize configuration
	config.InitConfig()

	// Check if config.AppConfig was populated correctly
	if config.AppConfig == nil {
		t.Fatal("AppConfig is nil after initialization")
	}

	// Validate configuration values
	validateConfigValues(t)
}
