package database

import (
	"log"

	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/port"
	"gorm.io/gorm"
)

// MockDatabaseAdapter is a mock implementation of the port.Database interface
// It doesn't actually connect to a database, but provides a mock implementation
// that allows the application to start and respond to API requests
type MockDatabaseAdapter struct {
	db *gorm.DB
}

// NewMockDatabaseAdapter creates a new MockDatabaseAdapter
func NewMockDatabaseAdapter() (port.Database, error) {
	log.Println("Using mock database adapter - no actual database connection will be established")
	log.Println("This is for demonstration purposes only and will not persist data")

	adapter := &MockDatabaseAdapter{}

	// Return the adapter without actually connecting to a database
	return adapter, nil
}

// GetDB returns a nil GORM DB instance
// This will cause errors if actual database operations are attempted
func (d *MockDatabaseAdapter) GetDB() *gorm.DB {
	log.Println("WARNING: Attempting to use mock database - no actual database operations will be performed")
	log.Println("This is a demonstration mode only. To use the application properly, please run it with a real PostgreSQL database.")
	log.Println("The mock database is not intended for actual use, but to demonstrate that the application can start without a database connection.")
	return d.db
}

// Close is a no-op for the mock adapter
func (d *MockDatabaseAdapter) Close() error {
	log.Println("Closing mock database adapter (no-op)")
	return nil
}

// Migrate is a no-op for the mock adapter
func (d *MockDatabaseAdapter) Migrate() error {
	log.Println("Skipping migrations for mock database adapter")
	return nil
}
