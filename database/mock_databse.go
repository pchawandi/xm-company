package database

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDatabase is a mock implementation of the Database interface
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDatabase) Where(query interface{}, args ...interface{}) Database {
	args1 := m.Called(query, args)
	return args1.Get(0).(Database)
}

func (m *MockDatabase) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(value, conds)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDatabase) Model(model interface{}) *gorm.DB {
	args := m.Called(model)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDatabase) First(dest interface{}, conds ...interface{}) Database {
	args := m.Called(dest, conds)
	return args.Get(0).(Database)
}

func (m *MockDatabase) Updates(values interface{}) *gorm.DB {
	args := m.Called(values)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDatabase) Error() error {
	args := m.Called()
	return args.Error(0)
}
