package repositories_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/adapters/database/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestRepositoryGORM_GetUserByName(t *testing.T) {
	db, mock := setupMockUserDB(t)
	repo := repositories.NewUserRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE name = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("Rafael", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow("550e8400-e29b-41d4-a716-446655440000", "Rafael"))

	user, err := repo.GetUserByName(context.Background(), "Rafael")
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, "550e8400-e29b-41d4-a716-446655440000", user.ID.String())
	require.Equal(t, "Rafael", user.Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
