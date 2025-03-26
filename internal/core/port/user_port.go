package port

import (
	"context"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
)

type UserRepository interface {
	Repository[domain.User]
	GetUserByName(ctx context.Context, name string) (*domain.User, error)
}

type UserService interface {
	GetUser(ctx context.Context, id string) (*domain.User, error)
	GetUserByName(ctx context.Context, name string) (*domain.User, error)
	ListUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
}
