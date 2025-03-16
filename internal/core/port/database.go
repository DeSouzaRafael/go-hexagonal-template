package port

import (
	"gorm.io/gorm"
)

type Database interface {
	GetDB() *gorm.DB
	Close() error
	AutoMigrate(models ...interface{}) error
}
