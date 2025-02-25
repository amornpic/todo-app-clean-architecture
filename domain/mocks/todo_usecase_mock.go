package mocks

import (
	"todo-app/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTodoUsecase struct {
	mock.Mock
}

func (m *MockTodoUsecase) Create(todo *domain.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoUsecase) Update(todo *domain.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoUsecase) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTodoUsecase) List(sortBy, search string) ([]domain.Todo, error) {
	args := m.Called(sortBy, search)
	return args.Get(0).([]domain.Todo), args.Error(1)
}
