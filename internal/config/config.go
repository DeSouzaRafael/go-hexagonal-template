package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database    DBConfig   `yaml:"database"`
	Jwt         JWT        `yaml:"jwt"`
	Environment string     `yaml:"environment"`
	WebService  WebService `yaml:"webservice"`
}

type WebService struct {
	Port string `yaml:"port"`
}

type JWT struct {
	Secret     string `yaml:"secret"`
	Expiration int    `yaml:"expiration"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
	LogLevel int    `yaml:"loglevel"`
}

var AppConfig *Config

func InitConfig() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("Failed to get current filename")
	}
	configDir := filepath.Dir(filename)
	configFilePath := filepath.Join(configDir, "config.yml")

	configFile, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer configFile.Close()

	decoder := yaml.NewDecoder(configFile)
	AppConfig = &Config{}
	err = decoder.Decode(AppConfig)
	if err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}
}
