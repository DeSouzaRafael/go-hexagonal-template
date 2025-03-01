package util

import (
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func CurrentExecutionEnvironmentProduction() bool {
	return config.AppConfig.Environment == "prd"
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func ComparePassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
