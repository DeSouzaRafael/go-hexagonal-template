package container_test

import (
	"testing"

	container "github.com/DeSouzaRafael/go-hexagonal-template/internal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type mockDB struct{}

func (m *mockDB) Connect() error {
	return nil
}

func (m *mockDB) Close() error {
	return nil
}

func (m *mockDB) AutoMigrate(models ...interface{}) error {
	return nil
}

func (m *mockDB) GetDB() *gorm.DB {
	return nil
}

func TestNewContainer(t *testing.T) {
	// Setup
	mockDB := &mockDB{}

	// Execute
	cont := container.NewContainer(mockDB)

	// Assert
	assert.NotNil(t, cont)
	assert.NotNil(t, cont.Repositories.UserRepository)
	assert.NotNil(t, cont.Services.AuthService)
	assert.NotNil(t, cont.Services.UserService)
	assert.NotNil(t, cont.Handlers.AuthHandler)
	assert.NotNil(t, cont.Handlers.UserHandler)
	assert.NotNil(t, cont.Repositories.UserRepository)

	// Verify container structure
	assert.IsType(t, &container.Container{}, cont)
}
