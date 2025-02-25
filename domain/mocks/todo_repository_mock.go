package mocks

import (
	"todo-app/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) Create(todo *domain.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoRepository) Update(todo *domain.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTodoRepository) FindByID(id uuid.UUID) (*domain.Todo, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Todo), args.Error(1)
}

func (m *MockTodoRepository) FindAll() ([]domain.Todo, error) {
	args := m.Called()
	return args.Get(0).([]domain.Todo), args.Error(1)
}
