package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/port"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockTokenGenerator implements the port.TokenGenerator interface for testing
type MockTokenGenerator struct {
	mock.Mock
}

// GenerateToken implements port.TokenGenerator
func (m *MockTokenGenerator) GenerateToken(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

// Ensure MockTokenGenerator implements port.TokenGenerator
var _ port.TokenGenerator = (*MockTokenGenerator)(nil)

func TestAuthService(t *testing.T) {
	config.AppConfig = &config.Config{
		Jwt: config.JWT{
			Secret:     "test-secret",
			Expiration: 3600,
		},
	}

	t.Run("successful registration", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenGen := new(MockTokenGenerator)
		service := service.NewAuthService(mockRepo, mockTokenGen)
		user := &domain.User{Name: "Rafa", Password: "password123"}

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(user, nil)

		result, err := service.Register(context.Background(), user)

		assert.NoError(t, err)
		assert.NotEqual(t, "password123", result.Password)
		mockRepo.AssertExpectations(t)
		mockTokenGen.AssertExpectations(t)
	})

	t.Run("successful login", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenGen := new(MockTokenGenerator)
		service := service.NewAuthService(mockRepo, mockTokenGen)
		hashedPw, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user := &domain.User{Name: "Rafa", Password: string(hashedPw)}

		mockRepo.On("GetUserByName", mock.Anything, "Rafa").Return(user, nil)
		mockTokenGen.On("GenerateToken", mock.Anything).Return("test-token", nil)

		token, err := service.Login(context.Background(), "Rafa", "password123")

		assert.NoError(t, err)
		assert.Equal(t, "test-token", token)
		mockRepo.AssertExpectations(t)
		mockTokenGen.AssertExpectations(t)
	})

	t.Run("login with invalid credentials - user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenGen := new(MockTokenGenerator)
		service := service.NewAuthService(mockRepo, mockTokenGen)

		mockRepo.On("GetUserByName", mock.Anything, "Rafa").Return(&domain.User{}, errors.New("user not found"))

		token, err := service.Login(context.Background(), "Rafa", "password123")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
		mockTokenGen.AssertExpectations(t)
	})

	t.Run("login with invalid credentials - wrong password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenGen := new(MockTokenGenerator)
		service := service.NewAuthService(mockRepo, mockTokenGen)

		hashedPw, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
		user := &domain.User{Name: "Rafa", Password: string(hashedPw)}

		mockRepo.On("GetUserByName", mock.Anything, "Rafa").Return(user, nil)

		token, err := service.Login(context.Background(), "Rafa", "wrong-password")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
		mockTokenGen.AssertExpectations(t)
	})

	t.Run("login failure - token generation error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenGen := new(MockTokenGenerator)
		service := service.NewAuthService(mockRepo, mockTokenGen)

		hashedPw, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user := &domain.User{Name: "Rafa", Password: string(hashedPw)}

		mockRepo.On("GetUserByName", mock.Anything, "Rafa").Return(user, nil)
		mockTokenGen.On("GenerateToken", mock.Anything).Return("", errors.New("token generation failed"))

		token, err := service.Login(context.Background(), "Rafa", "password123")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, "failed to generate token", err.Error())
		mockRepo.AssertExpectations(t)
		mockTokenGen.AssertExpectations(t)
	})

	t.Run("registration failure", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockTokenGen := new(MockTokenGenerator)
		service := service.NewAuthService(mockRepo, mockTokenGen)
		user := &domain.User{Name: "Rafa", Password: "password123"}

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return((*domain.User)(nil), errors.New("db error"))

		result, err := service.Register(context.Background(), user)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
		mockTokenGen.AssertExpectations(t)
	})
}
