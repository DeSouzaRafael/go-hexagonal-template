package repositories_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/database/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockUserDB(t *testing.T) (*MockDatabase, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	mockDatabase := &MockDatabase{DB: gormDB}

	return mockDatabase, mock
}

func TestNewUserRepository(t *testing.T) {
	db, _ := setupMockUserDB(t)

	userRepo := repositories.NewUserRepository(db)

	assert.NotNil(t, userRepo)
	assert.NotNil(t, userRepo.RepositoryGORM)
}
