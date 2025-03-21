package util_test

import (
	"testing"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/domain"
	"github.com/DeSouzaRafael/go-hexagonal-template/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestCurrentExecutionEnvironmentProduction(t *testing.T) {
	config.AppConfig = &config.Config{App: config.App{Name: "app", Environment: "production"}}

	// Success environment prd
	result := util.CurrentExecutionEnvironmentProduction()
	assert.True(t, result)

	// Another environment
	config.AppConfig.App.Environment = "development"
	result = util.CurrentExecutionEnvironmentProduction()
	assert.False(t, result)
}

func TestHashPassword(t *testing.T) {

	// Success
	password := "validpassword"
	hashedPassword, err := util.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	// Failure
	password = ""
	hashedPassword, err = util.HashPassword(password)
	assert.Error(t, err)
	assert.Equal(t, err, domain.ErrInvalidPassword)
	assert.Empty(t, hashedPassword)
}

func TestComparePassword(t *testing.T) {

	// Success
	password := "password123"
	hashedPassword, err := util.HashPassword(password)
	assert.NoError(t, err)
	err = util.ComparePassword(password, hashedPassword)
	assert.NoError(t, err)

	// Failure
	err = util.ComparePassword("wrongpassword", hashedPassword)
	assert.Error(t, err)
}
