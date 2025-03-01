package port

import (
	"context"
)

type Domain interface{}

type Repository[T Domain] interface {
	Get(ctx context.Context, id interface{}) (*T, error)
	List(ctx context.Context) ([]T, error)
	Create(ctx context.Context, domain *T) (*T, error)
	Update(ctx context.Context, domain *T) error
	Delete(ctx context.Context, id interface{}) error
}
