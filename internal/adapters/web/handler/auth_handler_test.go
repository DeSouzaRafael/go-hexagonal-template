package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/web/handler"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type AuthService struct {
	mock.Mock
}

func (m *AuthService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *AuthService) Login(ctx context.Context, name, password string) (string, error) {
	args := m.Called(ctx, name, password)
	return args.String(0), args.Error(1)
}

func TestRegister(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		mockService := new(AuthService)
		h := handler.NewAuthHandler(mockService)
		e := echo.New()

		reqBody := map[string]string{
			"name":     "Rafael",
			"password": "12345678",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("Register", mock.Anything, mock.AnythingOfType("*domain.User")).Return(&domain.User{}, nil)

		err := h.Register(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("registration error", func(t *testing.T) {
		mockService := new(AuthService)
		h := handler.NewAuthHandler(mockService)
		e := echo.New()

		reqBody := map[string]string{
			"name":     "Rafael",
			"password": "12345678",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("Register", mock.Anything, mock.AnythingOfType("*domain.User")).Return(&domain.User{}, errors.New("registration failed"))

		err := h.Register(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("invalid binding", func(t *testing.T) {
		mockService := new(AuthService)
		h := handler.NewAuthHandler(mockService)
		e := echo.New()

		invalidBody := []byte(`{"name": 123, "password": true}`)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(invalidBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Register(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestLogin(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		mockService := new(AuthService)
		h := handler.NewAuthHandler(mockService)
		e := echo.New()

		reqBody := map[string]string{
			"name":     "Rafael",
			"password": "12345678",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("Login", mock.Anything, "Rafael", "12345678").Return("token123", nil)

		err := h.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("login error", func(t *testing.T) {
		mockService := new(AuthService)
		h := handler.NewAuthHandler(mockService)
		e := echo.New()

		reqBody := map[string]string{
			"name":     "Rafael",
			"password": "12345678",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("Login", mock.Anything, "Rafael", "12345678").Return("", errors.New("invalid credentials"))

		err := h.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid binding", func(t *testing.T) {
		mockService := new(AuthService)
		h := handler.NewAuthHandler(mockService)
		e := echo.New()

		invalidBody := []byte(`{"name": 123, "password": true}`)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(invalidBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
