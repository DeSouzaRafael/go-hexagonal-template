package service

import (
	"context"
	"errors"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/middleware"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/port"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo port.UserRepository
}

func NewAuthService(userRepo port.UserRepository) AuthService {
	return AuthService{
		userRepo: userRepo,
	}
}

func (as *AuthService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	return as.userRepo.Create(ctx, user)
}

func (as *AuthService) Login(ctx context.Context, name, password string) (string, error) {
	user, err := as.userRepo.GetUserByName(ctx, name)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := middleware.GenerateJWT(user.ID.String())
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
