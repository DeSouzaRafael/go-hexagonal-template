package repositories

import (
	"context"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/port"
)

type RepositoryGORM[T any] struct {
	db port.Database
}

func NewRepositoryGORM[T any](db port.Database) *RepositoryGORM[T] {
	return &RepositoryGORM[T]{db: db}
}

func (repo *RepositoryGORM[T]) Get(ctx context.Context, id interface{}) (*T, error) {
	domain := new(T)
	err := repo.db.GetDB().WithContext(ctx).First(domain, id).Error
	return domain, err
}

func (r *RepositoryGORM[T]) List(ctx context.Context) ([]T, error) {
	var domain []T
	result := r.db.GetDB().WithContext(ctx).Find(&domain)
	return domain, result.Error
}

func (repo *RepositoryGORM[T]) Create(ctx context.Context, domain *T) (*T, error) {
	err := repo.db.GetDB().WithContext(ctx).Create(domain).Error
	return domain, err
}

func (repo *RepositoryGORM[T]) Update(ctx context.Context, domain *T) error {
	return repo.db.GetDB().WithContext(ctx).Save(domain).Error
}

func (repo *RepositoryGORM[T]) Delete(ctx context.Context, id interface{}) error {
	var domain T
	return repo.db.GetDB().WithContext(ctx).Model(&domain).Delete(&domain, id).Error
}
