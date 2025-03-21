package handler

import (
	"net/http"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/port"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	svc port.UserService
}

func NewUserHandler(svc port.UserService) *UserHandler {
	return &UserHandler{
		svc,
	}
}

func (uh *UserHandler) GetUser(ctx echo.Context) error {
	id := ctx.Param("id")

	user, err := uh.svc.GetUser(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, map[string]string{"message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, user)
}

func (uh *UserHandler) ListUsers(ctx echo.Context) error {

	users, err := uh.svc.ListUsers(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusForbidden, map[string]string{"message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, users)
}

type updateUserRequest struct {
	Name     string `json:"name" binding:"omitempty,required" example:"Rafa Dev"`
	Password string `json:"password" binding:"omitempty,required,min=8" example:"12345678"`
}

func (uh *UserHandler) UpdateUser(ctx echo.Context) error {
	var req updateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	id := ctx.Param("id")

	user := &domain.User{
		ID:       uuid.MustParse(id),
		Name:     req.Name,
		Password: req.Password,
	}

	user, err := uh.svc.UpdateUser(ctx.Request().Context(), user)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, map[string]string{"message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, user)
}

func (uh *UserHandler) DeleteUser(ctx echo.Context) error {
	id := ctx.Param("id")

	err := uh.svc.DeleteUser(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, map[string]string{"message": err.Error()})
	}

	return ctx.NoContent(http.StatusNoContent)
}
