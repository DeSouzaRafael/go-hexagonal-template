package util

import (
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
)

func CurrentExecutionEnvironmentProduction() bool {
	return config.AppConfig.Environment == "prd"
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", domain.ErrInvalidPassword
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), nil
}

func ComparePassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
