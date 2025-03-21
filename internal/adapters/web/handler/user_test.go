package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/handler"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetUserByName(ctx context.Context, name string) (*domain.User, error) {
	args := m.Called(ctx, name)

	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) ListUsers(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGetUser(t *testing.T) {
	mockSvc := new(MockUserService)
	h := handler.NewUserHandler(mockSvc)
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		user := &domain.User{
			ID:       uuid.New(),
			Name:     "Test User",
			Password: "password123",
		}

		mockSvc.On("GetUser", mock.Anything, mock.Anything).Return(user, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(user.ID.String())

		err := h.GetUser(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestListUsers(t *testing.T) {
	mockSvc := new(MockUserService)
	h := handler.NewUserHandler(mockSvc)
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		users := []domain.User{
			{ID: uuid.New(), Name: "User 1", Password: "pass1"},
			{ID: uuid.New(), Name: "User 2", Password: "pass2"},
		}

		mockSvc.On("ListUsers", mock.Anything).Return(users, nil)

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.ListUsers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
func TestUpdateUser(t *testing.T) {
	mockSvc := new(MockUserService)
	h := handler.NewUserHandler(mockSvc)
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		user := &domain.User{
			ID:       userID,
			Name:     "Updated User",
			Password: "newpass123",
		}

		mockSvc.On("UpdateUser", mock.Anything, mock.Anything).Return(user, nil)

		reqBody := `{"name":"Updated User","password":"newpass123"}`
		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(userID.String())

		err := h.UpdateUser(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestDeleteUser(t *testing.T) {
	mockSvc := new(MockUserService)
	h := handler.NewUserHandler(mockSvc)
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()

		mockSvc.On("DeleteUser", mock.Anything, mock.Anything).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(userID.String())

		err := h.DeleteUser(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}
