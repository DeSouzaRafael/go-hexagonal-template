package database

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/config"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/port"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type DatabaseAdapter struct {
	db     *gorm.DB
	config config.DBConfig
}

func NewDatabaseAdapter(config config.DBConfig) (port.Database, error) {
	adapter := &DatabaseAdapter{
		config: config,
	}

	if err := adapter.Connect(); err != nil {
		return nil, err
	}

	return adapter, nil
}

func (d *DatabaseAdapter) Connect() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		d.config.Host, d.config.User, d.config.Pass, d.config.DBName, d.config.Port, d.config.SSLMode,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.LogLevel(d.config.LogLevel),
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql connection: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	d.db = db
	return nil
}

func (d *DatabaseAdapter) GetDB() *gorm.DB {
	return d.db
}

func (d *DatabaseAdapter) Close() error {
	if d.db == nil {
		return nil
	}

	sqlDB, err := d.db.DB()
	if err != nil {
		if err.Error() == "sql: database is closed" {
			return nil
		}
		return fmt.Errorf("failed to get sql connection: %w", err)
	}
	err = sqlDB.Close()
	if err != nil && err.Error() == "sql: database is closed" {
		return nil
	}
	return err
}

func (d *DatabaseAdapter) Migrate() error {
	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.config.User, d.config.Pass, d.config.Host, d.config.Port, d.config.DBName, d.config.SSLMode,
	)

	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, dsn)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
