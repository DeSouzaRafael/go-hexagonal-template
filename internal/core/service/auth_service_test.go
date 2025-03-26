package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService(t *testing.T) {
	config.AppConfig = &config.Config{
		Jwt: config.JWT{
			Secret:     "test-secret",
			Expiration: 3600,
		},
	}

	t.Run("successful registration", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := service.NewAuthService(mockRepo)
		user := &domain.User{Name: "Rafa", Password: "password123"}

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(user, nil)

		result, err := service.Register(context.Background(), user)

		assert.NoError(t, err)
		assert.NotEqual(t, "password123", result.Password)
		mockRepo.AssertExpectations(t)
	})

	t.Run("successful login", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := service.NewAuthService(mockRepo)
		hashedPw, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user := &domain.User{Name: "Rafa", Password: string(hashedPw)}

		mockRepo.On("GetUserByName", mock.Anything, "Rafa").Return(user, nil)

		token, err := service.Login(context.Background(), "Rafa", "password123")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("login with invalid credentials", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := service.NewAuthService(mockRepo)

		mockRepo.On("GetUserByName", mock.Anything, "Rafa").Return(&domain.User{}, errors.New("user not found"))

		token, err := service.Login(context.Background(), "Rafa", "password123")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
