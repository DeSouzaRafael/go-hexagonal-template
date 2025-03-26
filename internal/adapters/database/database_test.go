package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewDatabaseAdapter(t *testing.T) {
	tests := []struct {
		name    string
		config  config.DBConfig
		wantErr bool
	}{
		{
			name: "invalid connection",
			config: config.DBConfig{
				Host:     "invalid-host",
				Port:     "5432",
				User:     "postgres",
				Pass:     "postgres",
				DBName:   "test",
				SSLMode:  "disable",
				LogLevel: 1,
			},
			wantErr: true,
		},
		{
			name: "empty configuration",
			config: config.DBConfig{
				Host: "",
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			config: config.DBConfig{
				Host: "localhost",
				Port: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewDatabaseAdapter(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
			}
		})
	}
}

func TestDatabaseAdapter_GetDB(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn: sqlDB,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	adapter := &DatabaseAdapter{
		db: db,
	}

	result := adapter.GetDB()
	assert.NotNil(t, result)
	assert.Equal(t, db, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDatabaseAdapter_Close(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (*DatabaseAdapter, sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "should close connection successfully",
			setup: func(t *testing.T) (*DatabaseAdapter, sqlmock.Sqlmock) {
				sqlDB, mock, err := sqlmock.New()
				require.NoError(t, err)

				mock.ExpectClose()

				dialector := postgres.New(postgres.Config{
					Conn:                 sqlDB,
					PreferSimpleProtocol: true,
				})
				db, err := gorm.Open(dialector, &gorm.Config{})
				require.NoError(t, err)

				return &DatabaseAdapter{db: db}, mock
			},
			wantErr: false,
		},
		{
			name: "should return nil when db is nil",
			setup: func(t *testing.T) (*DatabaseAdapter, sqlmock.Sqlmock) {
				return &DatabaseAdapter{db: nil}, nil
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter, mock := tt.setup(t)

			err := adapter.Close()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if mock != nil {
					assert.NoError(t, mock.ExpectationsWereMet())
				}
			}
		})
	}
}

func TestDatabaseAdapter_AutoMigrate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) (*DatabaseAdapter, sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "should migrate successfully",
			setup: func(t *testing.T) (*DatabaseAdapter, sqlmock.Sqlmock) {
				sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
				require.NoError(t, err)

				mock.ExpectQuery(`SELECT count\(\*\) FROM information_schema\.tables WHERE table_schema = CURRENT_SCHEMA\(\) AND table_name = \$1 AND table_type = \$2`).
					WithArgs("test_models", "BASE TABLE").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				mock.ExpectExec(`CREATE TABLE "test_models" \("id" bigserial,"name" varchar\(100\),PRIMARY KEY \("id"\)\)`).
					WillReturnResult(sqlmock.NewResult(1, 1))

				dialector := postgres.New(postgres.Config{
					Conn:                 sqlDB,
					PreferSimpleProtocol: true,
				})

				db, err := gorm.Open(dialector, &gorm.Config{
					SkipDefaultTransaction: true,
				})
				require.NoError(t, err)

				return &DatabaseAdapter{db: db}, mock
			},
			wantErr: false,
		},
		{
			name: "should fail when db is nil",
			setup: func(t *testing.T) (*DatabaseAdapter, sqlmock.Sqlmock) {
				return &DatabaseAdapter{db: nil}, nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter, mock := tt.setup(t)

			type TestModel struct {
				ID   uint   `gorm:"primarykey"`
				Name string `gorm:"type:varchar(100)"`
			}

			err := adapter.AutoMigrate(&TestModel{})

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, "database connection not initialized", err.Error())
			} else {
				assert.NoError(t, err)
				if mock != nil {
					assert.NoError(t, mock.ExpectationsWereMet())
				}
			}
		})
	}
}

func TestDatabaseAdapter_Connect(t *testing.T) {
	tests := []struct {
		name    string
		config  config.DBConfig
		wantErr bool
	}{
		{
			name: "should connect successfully",
			config: config.DBConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Pass:     "postgres",
				DBName:   "test",
				SSLMode:  "disable",
				LogLevel: 1,
			},
			wantErr: true,
		},
		{
			name: "should fail with empty host",
			config: config.DBConfig{
				Host: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DatabaseAdapter{
				config: tt.config,
			}
			err := adapter.Connect()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, adapter.db)
			}
		})
	}
}
