package handler

import (
	"net/http"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/port"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService port.AuthService
}

func NewAuthHandler(authService port.AuthService) AuthHandler {
	return AuthHandler{
		authService: authService,
	}
}

type registerRequest struct {
	Name     string `json:"name"     validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

func (ah *AuthHandler) Register(ctx echo.Context) error {
	var req registerRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	user := domain.User{
		Name:     req.Name,
		Password: req.Password,
	}

	_, err := ah.authService.Register(ctx.Request().Context(), &user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, user)
}

type loginRequest struct {
	Name     string `json:"name"     validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (ah *AuthHandler) Login(ctx echo.Context) error {
	var req loginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	token, err := ah.authService.Login(ctx.Request().Context(), req.Name, req.Password)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"token": token})
}

func (ah *AuthHandler) GetProfile(ctx echo.Context) error {
	userID := ctx.Get("user_id").(string)
	return ctx.JSON(http.StatusOK, map[string]string{"user_id": userID})
}
