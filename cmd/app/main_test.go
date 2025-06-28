package main

import (
	"errors"
	"testing"

	container "github.com/DeSouzaRafael/go-hexagonal-template/internal"
	"github.com/DeSouzaRafael/go-hexagonal-template/internal/core/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) GetDB() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDatabase) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDatabase) AutoMigrate(models ...interface{}) error {
	args := m.Called(models)
	return args.Error(0)
}

type MockWebServer struct {
	mock.Mock
}

func (m *MockWebServer) Start() error {
	args := m.Called()
	return args.Error(0)
}

type MockContainer struct {
	Handlers container.Handlers
}

func TestRunWithDependencies(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		setupMocks    func(*MockDatabase, *MockWebServer) (port.Database, ContainerFactory, WebServiceFactory)
		expectedError error
	}{
		{
			name: "successful run",
			setupMocks: func(mockDB *MockDatabase, mockWebServer *MockWebServer) (port.Database, ContainerFactory, WebServiceFactory) {
				// Setup database mock
				mockDB.On("AutoMigrate", mock.Anything).Return(nil)

				// Setup web server mock
				mockWebServer.On("Start").Return(nil)

				// Create container factory
				containerFactory := func(db port.Database) *container.Container {
					return &container.Container{
						Handlers: container.Handlers{},
					}
				}

				// Create web service factory
				webServiceFactory := func(h container.Handlers) WebServer {
					return mockWebServer
				}

				return mockDB, containerFactory, webServiceFactory
			},
			expectedError: nil,
		},
		{
			name: "migration error",
			setupMocks: func(mockDB *MockDatabase, mockWebServer *MockWebServer) (port.Database, ContainerFactory, WebServiceFactory) {
				// Setup database mock with error
				mockDB.On("AutoMigrate", mock.Anything).Return(errors.New("migration error"))

				// Create container factory
				containerFactory := func(db port.Database) *container.Container {
					return &container.Container{
						Handlers: container.Handlers{},
					}
				}

				// Create web service factory
				webServiceFactory := func(h container.Handlers) WebServer {
					return mockWebServer
				}

				return mockDB, containerFactory, webServiceFactory
			},
			expectedError: errors.New("migration error"),
		},
		{
			name: "web server error",
			setupMocks: func(mockDB *MockDatabase, mockWebServer *MockWebServer) (port.Database, ContainerFactory, WebServiceFactory) {
				// Setup database mock
				mockDB.On("AutoMigrate", mock.Anything).Return(nil)

				// Setup web server mock with error
				mockWebServer.On("Start").Return(errors.New("web server error"))

				// Create container factory
				containerFactory := func(db port.Database) *container.Container {
					return &container.Container{
						Handlers: container.Handlers{},
					}
				}

				// Create web service factory
				webServiceFactory := func(h container.Handlers) WebServer {
					return mockWebServer
				}

				return mockDB, containerFactory, webServiceFactory
			},
			expectedError: errors.New("web server error"),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockDB := new(MockDatabase)
			mockWebServer := new(MockWebServer)
			db, containerFactory, webServiceFactory := tc.setupMocks(mockDB, mockWebServer)

			// Execute
			err := runWithDependencies(db, containerFactory, webServiceFactory)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify expectations
			mockDB.AssertExpectations(t)
			mockWebServer.AssertExpectations(t)
		})
	}
}
