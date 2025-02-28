package repositories

import (
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/ports"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type RepositoryGORM[T ports.Model] struct {
	db *gorm.DB
}

func NewRepositoryGORM[T ports.Model](db *gorm.DB) *RepositoryGORM[T] {
	return &RepositoryGORM[T]{db: db}
}

func (r *RepositoryGORM[T]) GetByID(id uuid.UUID) (T, error) {
	var model T
	result := r.db.First(&model, "id = ?", id)
	return model, result.Error
}

func (r *RepositoryGORM[T]) GetAll() ([]T, error) {
	var model []T
	result := r.db.Find(&model)
	return model, result.Error
}

func (r *RepositoryGORM[T]) Create(model T) (T, error) {
	err := r.db.Create(&model).Error
	return model, err
}

func (r *RepositoryGORM[T]) Update(model T) error {
	return r.db.Save(&model).Error
}

func (r *RepositoryGORM[T]) DeleteByID(id uuid.UUID) error {
	var model T
	return r.db.Delete(&model, id).Error
}
