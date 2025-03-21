package port

import (
	"context"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
)

type AuthService interface {
	Register(ctx context.Context, user *domain.User) (*domain.User, error)
	Login(ctx context.Context, name, password string) (string, error)
}
