package port

import (
	"context"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
)

type UserRepository interface {
	Repository[domain.User]
}

type UserService interface {
	Register(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
}
