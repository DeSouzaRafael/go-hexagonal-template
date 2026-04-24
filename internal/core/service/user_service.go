package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/port"
	"github.com/DeSouzaRafael/go-hexagonal-template/pkg/util"
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{
		repo,
	}
}

func (us *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, domain.ErrDataNotFound
	}

	user, err := us.repo.Get(ctx, parsedID)
	if err != nil {
		return nil, domain.ErrDataNotFound
	}

	return user, nil
}

func (us *UserService) GetUserByName(ctx context.Context, name string) (*domain.User, error) {
	user, err := us.repo.GetUserByName(ctx, name)
	if err != nil {
		return nil, domain.ErrDataNotFound
	}

	return user, nil
}

func (us *UserService) ListUsers(ctx context.Context) ([]domain.User, error) {
	users, err := us.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us *UserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := us.repo.Get(ctx, user.ID)
	if err != nil {
		return nil, domain.ErrDataNotFound
	}

	if user.Name == "" {
		user.Name = existingUser.Name
	}

	if user.Password != "" {
		hashed, err := util.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashed
	} else {
		user.Password = existingUser.Password
	}

	err = us.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) DeleteUser(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return domain.ErrDataNotFound
	}

	_, err = us.repo.Get(ctx, parsedID)
	if err != nil {
		return domain.ErrDataNotFound
	}

	return us.repo.Delete(ctx, parsedID)
}
