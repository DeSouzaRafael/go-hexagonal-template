package util_test

import (
	"testing"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/DeSouzaRafael/go-hexagonal-template/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestCurrentExecutionEnvironmentProduction(t *testing.T) {
	config.AppConfig = &config.Config{App: config.App{Name: "app", Environment: "production"}}

	result := util.CurrentExecutionEnvironmentProduction()
	assert.True(t, result)

	config.AppConfig.App.Environment = "development"
	result = util.CurrentExecutionEnvironmentProduction()
	assert.False(t, result)
}

func TestHashPassword(t *testing.T) {
	password := "validpassword"
	hashedPassword, err := util.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	other, err := util.HashPassword("anotherpassword")
	assert.NoError(t, err)
	assert.NotEqual(t, hashedPassword, other)
}
