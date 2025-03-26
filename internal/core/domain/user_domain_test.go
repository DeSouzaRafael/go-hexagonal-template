package domain_test

import (
	"testing"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserBeforeCreate(t *testing.T) {
	// Given
	user := &domain.User{}

	// When
	err := user.BeforeCreate(&gorm.DB{})

	// Then
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.IsType(t, uuid.UUID{}, user.ID)
}
