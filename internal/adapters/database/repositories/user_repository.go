package repositories

import (
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/database/models"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/ports"
	"gorm.io/gorm"
)

type UserRepositoryGORM struct {
	*RepositoryGORM[models.User]
}

func NewUserRepositoryGORM(db *gorm.DB) ports.UserRepository {
	return &UserRepositoryGORM{
		RepositoryGORM: NewRepositoryGORM[models.User](db),
	}
}
