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

type User struct {
	ID   int
	Name string
}

type MockDatabase struct {
	DB *gorm.DB
}

func (m *MockDatabase) GetDB() *gorm.DB {
	return m.DB
}

func (m *MockDatabase) Close() error {
	return nil
}

func (m *MockDatabase) AutoMigrate(models ...interface{}) error {
	return nil
}

func setupMockDB(t *testing.T) (*MockDatabase, sqlmock.Sqlmock) {
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

func TestRepositoryGORM_Get(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repositories.NewRepositoryGORM[User](db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Rafael"))

	user, err := repo.Get(context.Background(), 1)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, 1, user.ID)
	require.Equal(t, "Rafael", user.Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestRepositoryGORM_List(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repositories.NewRepositoryGORM[User](db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Rafael S").
			AddRow(2, "Joao M"))

	users, err := repo.List(context.Background())

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "Rafael S", users[0].Name)
	assert.Equal(t, "Joao M", users[1].Name)
}

func TestRepositoryGORM_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repositories.NewRepositoryGORM[User](db)

	user := &User{Name: "Joao M"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("name") VALUES ($1) RETURNING "id"`)).
		WithArgs("Joao M").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	createdUser, err := repo.Create(context.Background(), user)

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
}

func TestRepositoryGORM_Update(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repositories.NewRepositoryGORM[User](db)

	user := &User{ID: 1, Name: "Joao M"}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "users" SET.*WHERE.*`).
		WithArgs(user.Name, user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), user)

	assert.NoError(t, err)
}

func TestRepositoryGORM_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repositories.NewRepositoryGORM[User](db)

	user := &User{ID: 1}

	mock.ExpectBegin()
	mock.ExpectExec(`^DELETE FROM "users" WHERE.*`).
		WithArgs(user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), user.ID)
	assert.NoError(t, err)
}
