package repositories

import (
	"context"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/port"

	"gorm.io/gorm"
)

type RepositoryGORM[T port.Domain] struct {
	db *gorm.DB
}

func NewRepositoryGORM[T port.Domain](db *gorm.DB) *RepositoryGORM[T] {
	return &RepositoryGORM[T]{db: db}
}

func (repo *RepositoryGORM[T]) Get(ctx context.Context, id interface{}) (*T, error) {
	var domain *T
	err := repo.db.WithContext(ctx).First(domain, id).Error
	return domain, err
}

func (r *RepositoryGORM[T]) List(ctx context.Context) ([]T, error) {
	var domain []T
	result := r.db.WithContext(ctx).Find(&domain)
	return domain, result.Error
}

func (repo *RepositoryGORM[T]) Create(ctx context.Context, domain *T) (*T, error) {
	err := repo.db.WithContext(ctx).Create(domain).Error
	return domain, err
}

func (repo *RepositoryGORM[T]) Update(ctx context.Context, domain *T) error {
	return repo.db.WithContext(ctx).Save(domain).Error
}

func (repo *RepositoryGORM[T]) Delete(ctx context.Context, id interface{}) error {
	var domain *T
	return repo.db.WithContext(ctx).Delete(domain, id).Error
}
