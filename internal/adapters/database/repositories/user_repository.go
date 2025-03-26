package repositories

import (
	"context"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/port"
)

type UserRepo struct {
	*RepositoryGORM[domain.User]
}

func NewUserRepository(db port.Database) *UserRepo {
	return &UserRepo{
		RepositoryGORM: NewRepositoryGORM[domain.User](db),
	}
}

func (ur *UserRepo) GetUserByName(ctx context.Context, name string) (*domain.User, error) {
	var user domain.User
	err := ur.db.GetDB().First(&user, "name = ?", name).Error
	return &user, err
}
