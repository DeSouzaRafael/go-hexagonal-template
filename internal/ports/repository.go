package ports

import (
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/database/models"
	"github.com/google/uuid"
)

type Model interface{}

type Repository[T Model] interface {
	GetByID(id uuid.UUID) (T, error)
	GetAll() ([]T, error)
	Create(model T) (T, error)
	Update(model T) error
	DeleteByID(id uuid.UUID) error
}

type UserRepository interface {
	Repository[models.User]
}
