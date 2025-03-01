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

func NewUserService(repo port.UserRepository) UserService {
	return UserService{
		repo,
	}
}

func (us *UserService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, domain.ErrInternal
	}

	user.Password = hashedPassword

	user, err = us.repo.Create(ctx, user)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (us *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	var user *domain.User

	user, err := us.repo.Get(ctx, uuid.MustParse(id))
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (us *UserService) ListUsers(ctx context.Context) ([]domain.User, error) {
	var users []domain.User

	users, err := us.repo.List(ctx)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return users, nil
}

func (us *UserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := us.repo.Get(ctx, user.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	emptyData := user.Name == "" &&
		user.Password == ""
	sameData := existingUser.Name == user.Name
	if emptyData || sameData {
		return nil, domain.ErrNoUpdatedData
	}

	var hashedPassword string

	if user.Password != "" {
		hashedPassword, err = util.HashPassword(user.Password)
		if err != nil {
			return nil, domain.ErrInternal
		}
	}

	user.Password = hashedPassword

	err = us.repo.Update(ctx, user)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (us *UserService) DeleteUser(ctx context.Context, id string) error {
	_, err := us.repo.Get(ctx, uuid.MustParse(id))
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	return us.repo.Delete(ctx, uuid.MustParse(id))
}
