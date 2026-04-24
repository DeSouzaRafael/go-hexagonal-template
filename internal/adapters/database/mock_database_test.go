package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMockDatabaseAdapter(t *testing.T) {
	adapter, err := NewMockDatabaseAdapter()
	assert.NoError(t, err)
	assert.NotNil(t, adapter)
}

func TestMockDatabaseAdapter_GetDB(t *testing.T) {
	adapter := &MockDatabaseAdapter{}
	db := adapter.GetDB()
	assert.Nil(t, db)
}

func TestMockDatabaseAdapter_Close(t *testing.T) {
	adapter := &MockDatabaseAdapter{}
	err := adapter.Close()
	assert.NoError(t, err)
}

func TestMockDatabaseAdapter_Migrate(t *testing.T) {
	adapter := &MockDatabaseAdapter{}
	err := adapter.Migrate()
	assert.NoError(t, err)
}
