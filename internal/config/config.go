package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database   DBConfig
	Jwt        JWT
	App        App
	WebService WebService
}

type App struct {
	Name        string
	Environment string
}

type WebService struct {
	Port   string
	Domain string
}

type JWT struct {
	Secret     string
	Expiration int
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Pass     string
	DBName   string
	SSLMode  string
	LogLevel int
}

var AppConfig *Config

func LoadConfig() error {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return err
	}

	logLevel, _ := strconv.Atoi(os.Getenv("DB_LOG_LEVEL"))
	jwtExpiration, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))

	AppConfig = &Config{
		App: App{
			Name:        os.Getenv("APP_NAME"),
			Environment: os.Getenv("APP_ENV"),
		},
		WebService: WebService{
			Port:   os.Getenv("WEB_PORT"),
			Domain: os.Getenv("WEB_DOMAIN"),
		},
		Jwt: JWT{
			Secret:     os.Getenv("JWT_SECRET"),
			Expiration: jwtExpiration,
		},
		Database: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Pass:     os.Getenv("DB_PASS"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSL_MODE"),
			LogLevel: logLevel,
		},
	}

	return nil
}
