package repositories

import (
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
