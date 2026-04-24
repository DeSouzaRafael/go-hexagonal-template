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

func (m *MockDatabase) Migrate() error {
	args := m.Called()
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
	testCases := []struct {
		name          string
		setupMocks    func(*MockDatabase, *MockWebServer) (port.Database, ContainerFactory, WebServiceFactory)
		expectedError error
	}{
		{
			name: "successful run",
			setupMocks: func(mockDB *MockDatabase, mockWebServer *MockWebServer) (port.Database, ContainerFactory, WebServiceFactory) {
				mockDB.On("Migrate").Return(nil)
				mockWebServer.On("Start").Return(nil)

				containerFactory := func(db port.Database) *container.Container {
					return &container.Container{Handlers: container.Handlers{}}
				}
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
				mockDB.On("Migrate").Return(errors.New("migration error"))

				containerFactory := func(db port.Database) *container.Container {
					return &container.Container{Handlers: container.Handlers{}}
				}
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
				mockDB.On("Migrate").Return(nil)
				mockWebServer.On("Start").Return(errors.New("web server error"))

				containerFactory := func(db port.Database) *container.Container {
					return &container.Container{Handlers: container.Handlers{}}
				}
				webServiceFactory := func(h container.Handlers) WebServer {
					return mockWebServer
				}

				return mockDB, containerFactory, webServiceFactory
			},
			expectedError: errors.New("web server error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDB := new(MockDatabase)
			mockWebServer := new(MockWebServer)
			db, containerFactory, webServiceFactory := tc.setupMocks(mockDB, mockWebServer)

			err := runWithDependencies(db, containerFactory, webServiceFactory)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
			mockWebServer.AssertExpectations(t)
		})
	}
}
