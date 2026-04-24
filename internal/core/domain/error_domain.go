package domain

import (
	"errors"
)

var (
	ErrInternal        = errors.New("internal error")
	ErrDataNotFound    = errors.New("data not found")
	ErrNoUpdatedData   = errors.New("no data to update")
	ErrConflictingData = errors.New("data conflicts with existing data in unique column")
	ErrTokenDuration   = errors.New("invalid token duration format")
	ErrMissingToken    = errors.New("missing token")
	ErrExpiredToken    = errors.New("access token has expired")
	ErrInvalidToken    = errors.New("access token is invalid")
)
