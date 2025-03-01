package repositories

import (
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"gorm.io/gorm"
)

type UserRepo struct {
	*RepositoryGORM[domain.User]
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{
		RepositoryGORM: NewRepositoryGORM[domain.User](db),
	}
}
