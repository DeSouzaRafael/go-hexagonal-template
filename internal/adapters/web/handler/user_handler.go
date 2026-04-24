package handler

import (
	"errors"
	"net/http"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/port"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	svc port.UserService
}

func NewUserHandler(svc port.UserService) UserHandler {
	return UserHandler{
		svc,
	}
}

func (uh *UserHandler) GetUser(ctx echo.Context) error {
	id := ctx.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID format"})
	}

	user, err := uh.svc.GetUser(ctx.Request().Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, user)
}

func (uh *UserHandler) ListUsers(ctx echo.Context) error {
	users, err := uh.svc.ListUsers(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, users)
}

type updateUserRequest struct {
	Name     string `json:"name"     validate:"omitempty,min=1"`
	Password string `json:"password" validate:"omitempty,min=8"`
}

func (uh *UserHandler) UpdateUser(ctx echo.Context) error {
	var req updateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	id := ctx.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID format"})
	}

	user := &domain.User{
		ID:       parsedID,
		Name:     req.Name,
		Password: req.Password,
	}

	user, err = uh.svc.UpdateUser(ctx.Request().Context(), user)
	if err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
		} else if errors.Is(err, domain.ErrNoUpdatedData) {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		} else if errors.Is(err, domain.ErrConflictingData) {
			return ctx.JSON(http.StatusConflict, map[string]string{"message": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, user)
}

func (uh *UserHandler) DeleteUser(ctx echo.Context) error {
	id := ctx.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID format"})
	}

	err := uh.svc.DeleteUser(ctx.Request().Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return ctx.NoContent(http.StatusNoContent)
}
