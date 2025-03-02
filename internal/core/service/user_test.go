package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/service"
	"github.com/DeSouzaRafael/go-hexagonal-template/pkg/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Get(ctx context.Context, id interface{}) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context) ([]domain.User, error) {
	fmt.Println("List function called")
	args := m.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id interface{}) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestRegister(t *testing.T) {
	repo := new(MockUserRepository)
	service := service.NewUserService(repo)
	ctx := context.Background()

	user := &domain.User{
		ID:       uuid.New(),
		Name:     "Rafa S",
		Password: "",
	}

	// Failure - invalid password
	result, err := service.Register(ctx, user)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrInvalidPassword, err)

	user.Password = "Password123"
	hashedPassword, _ := util.HashPassword(user.Password)
	expectedUser := *user
	expectedUser.Password = hashedPassword

	// Success
	repo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(&expectedUser, nil).Once()

	result, err = service.Register(ctx, user)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)

	// Failure - user exist
	repo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return((*domain.User)(nil), domain.ErrConflictingData).Once()

	result, err = service.Register(ctx, user)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrConflictingData, err)
}

func TestGetUser(t *testing.T) {
	repo := new(MockUserRepository)
	service := service.NewUserService(repo)
	ctx := context.Background()

	userID := uuid.New()
	expectedUser := &domain.User{ID: userID, Name: "Rafa S"}

	// Success
	repo.On("Get", ctx, userID).Return(expectedUser, nil).Once()
	result, err := service.GetUser(ctx, userID.String())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)

	// Failure - user not exist
	repo.On("Get", ctx, userID).Return((*domain.User)(nil), domain.ErrDataNotFound).Once()
	result, err = service.GetUser(ctx, userID.String())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrDataNotFound, err)
}

func TestListUsers(t *testing.T) {
	repo := new(MockUserRepository)
	service := service.NewUserService(repo)
	ctx := context.Background()

	users := []domain.User{
		{ID: uuid.New(), Name: "Rafa S"},
		{ID: uuid.New(), Name: "Joao D"},
	}

	// Success
	repo.On("List", ctx).Return(users, nil).Once()
	result, err := service.ListUsers(ctx)
	assert.NoError(t, err)
	assert.Equal(t, users, result)

	// Failure - error list users
	repo.On("List", ctx).Return([]domain.User{}, errors.New("db error")).Once()
	result, err = service.ListUsers(ctx)
	assert.Error(t, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	repo := new(MockUserRepository)
	service := service.NewUserService(repo)
	ctx := context.Background()
	updatedUser := &domain.User{
		ID:   uuid.New(),
		Name: "Updated Name",
	}

	// Success
	repo.On("Get", ctx, updatedUser.ID).Return(&domain.User{ID: updatedUser.ID, Name: "Rafa S"}, nil).Once()
	repo.On("Update", ctx, updatedUser).Return(updatedUser, nil).Once()
	result, err := service.UpdateUser(ctx, updatedUser)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser, result)

	// Success - without changing name and changing password
	updatedUser.Name = ""
	updatedUser.Password = "new-pass123"
	repo.On("Get", ctx, updatedUser.ID).Return(&domain.User{ID: updatedUser.ID, Name: "Rafa S"}, nil).Once()
	repo.On("Update", ctx, updatedUser).Return(updatedUser, nil).Once()
	result, err = service.UpdateUser(ctx, updatedUser)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser, result)

	// Failure - update user
	repo.On("Get", ctx, updatedUser.ID).Return(&domain.User{ID: updatedUser.ID, Name: "Rafa S"}, nil).Once()
	repo.On("Update", ctx, updatedUser).Return((*domain.User)(nil), errors.New("db error")).Once()
	result, err = service.UpdateUser(ctx, updatedUser)
	assert.Error(t, err)
	assert.Nil(t, result)

	// Failure - user not found
	repo.On("Get", ctx, updatedUser.ID).Return((*domain.User)(nil), domain.ErrDataNotFound).Once()
	repo.On("Update", ctx, updatedUser).Return((*domain.User)(nil), domain.ErrDataNotFound).Once()
	result, err = service.UpdateUser(ctx, updatedUser)
	assert.Error(t, err)
	assert.Nil(t, result)

}

func TestDeleteUser(t *testing.T) {
	repo := new(MockUserRepository)
	service := service.NewUserService(repo)
	ctx := context.Background()
	id := uuid.New()

	// Success
	repo.On("Get", ctx, id).Return(&domain.User{ID: id, Name: "Rafa S"}, nil).Once()
	repo.On("Delete", ctx, id).Return(nil).Once()
	err := service.DeleteUser(ctx, id.String())
	assert.NoError(t, err)

	// Failure - error delete user
	repo.On("Get", ctx, id).Return(&domain.User{ID: id, Name: "Rafa S"}, nil).Once()
	repo.On("Delete", ctx, id).Return(errors.New("db error")).Once()
	err = service.DeleteUser(ctx, id.String())
	assert.Error(t, err)

	// Failure - user not found
	repo.On("Get", ctx, id).Return((*domain.User)(nil), domain.ErrDataNotFound).Once()
	repo.On("Delete", ctx, id).Return((*domain.User)(nil), domain.ErrDataNotFound).Once()
	err = service.DeleteUser(ctx, id.String())
	assert.Error(t, err)
}
