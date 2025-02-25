package usecase

import (
	"log/slog"
	"testing"
	"todo-app/domain"
	"todo-app/domain/mocks"

	"github.com/stretchr/testify/assert"
)

func TestTodoUsecase_Create(t *testing.T) {
	mockRepo := new(mocks.MockTodoRepository)
	logger := slog.Default()
	usecase := NewTodoUsecase(mockRepo, logger)

	todo := &domain.Todo{
		Title:       "Test Todo",
		Description: "Test Description",
		Status:      "IN_PROGRESS",
	}

	// Success case
	t.Run("success", func(t *testing.T) {
		mockRepo.On("Create", todo).Return(nil).Once()

		err := usecase.Create(todo)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	// Error case
	t.Run("error", func(t *testing.T) {
		mockRepo.On("Create", todo).Return(assert.AnError).Once()

		err := usecase.Create(todo)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
