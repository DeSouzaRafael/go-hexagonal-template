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
	Name     string `json:"name"     validate:"required"      example:"John Doe"`
	Password string `json:"password" validate:"required,min=8" example:"12345678"`
}

// Register godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body registerRequest true "Registration payload"
// @Success      201  {object}  domain.User
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /v1/auth/register [post]
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

	created, err := ah.authService.Register(ctx.Request().Context(), &user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, created)
}

type loginRequest struct {
	Name     string `json:"name"     validate:"required" example:"John Doe"`
	Password string `json:"password" validate:"required" example:"12345678"`
}

// Login godoc
// @Summary      Authenticate a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body loginRequest true "Login credentials"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /v1/auth/login [post]
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

// GetProfile godoc
// @Summary      Get authenticated user profile
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Security     BearerAuth
// @Router       /v1/auth/profile [get]
func (ah *AuthHandler) GetProfile(ctx echo.Context) error {
	userID := ctx.Get("user_id").(string)
	return ctx.JSON(http.StatusOK, map[string]string{"user_id": userID})
}
