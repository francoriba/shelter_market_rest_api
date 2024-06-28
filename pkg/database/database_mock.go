// pkg/database/database_mock.go

package database

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// DBInterface defines the interface for database operations
type DBInterface interface {
	Create(value interface{}) *gorm.DB
	// Add other database operations used by your controller if needed
}

// MockDB is a mock struct for DBInterface
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return nil // Implement this if you need to mock Where
}
