package controller

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-app/domain"
	"todo-app/domain/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

func TestTodoController_Create(t *testing.T) {
	mockUsecase := new(mocks.MockTodoUsecase)
	logger := slog.Default()
	controller := NewTodoController(mockUsecase, logger)
	router := setupRouter()

	router.POST("/todos", controller.Create)

	t.Run("success", func(t *testing.T) {
		todo := domain.Todo{
			Title:       "Test Todo",
			Description: "Test Description",
			Status:      "IN_PROGRESS",
		}

		mockUsecase.On("Create", mock.AnythingOfType("*domain.Todo")).Return(nil).Once()

		body, _ := json.Marshal(todo)
		req := httptest.NewRequest("POST", "/todos", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("bad request", func(t *testing.T) {
		invalidJSON := `{"title": 123,` // Invalid JSON

		req := httptest.NewRequest("POST", "/todos", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestTodoController_Update(t *testing.T) {
	mockUsecase := new(mocks.MockTodoUsecase)
	logger := slog.Default()
	controller := NewTodoController(mockUsecase, logger)
	router := setupRouter()

	router.PUT("/todos/:id", controller.Update)

	t.Run("success", func(t *testing.T) {
		todo := domain.Todo{
			Title:       "Updated Todo",
			Description: "Updated Description",
			Status:      "COMPLETED",
		}

		mockUsecase.On("Update", mock.AnythingOfType("*domain.Todo")).Return(nil).Once()

		body, _ := json.Marshal(todo)
		req := httptest.NewRequest("PUT", "/todos/123e4567-e89b-12d3-a456-426614174000", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		todo := domain.Todo{
			Title:       "Updated Todo",
			Description: "Updated Description",
			Status:      "COMPLETED",
		}

		body, _ := json.Marshal(todo)
		req := httptest.NewRequest("PUT", "/todos/invalid-uuid", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestTodoController_List(t *testing.T) {
	mockUsecase := new(mocks.MockTodoUsecase)
	logger := slog.Default()
	controller := NewTodoController(mockUsecase, logger)
	router := setupRouter()

	router.GET("/todos", controller.List)

	t.Run("success", func(t *testing.T) {
		todos := []domain.Todo{
			{
				Title:       "Todo 1",
				Description: "Description 1",
				Status:      "IN_PROGRESS",
			},
			{
				Title:       "Todo 2",
				Description: "Description 2",
				Status:      "COMPLETED",
			},
		}

		mockUsecase.On("List", "", "").Return(todos, nil).Once()

		req := httptest.NewRequest("GET", "/todos", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("with sort and search", func(t *testing.T) {
		todos := []domain.Todo{}
		mockUsecase.On("List", "title", "test").Return(todos, nil).Once()

		req := httptest.NewRequest("GET", "/todos?sort_by=title&search=test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func TestTodoController_Delete(t *testing.T) {
	mockUsecase := new(mocks.MockTodoUsecase)
	logger := slog.Default()
	controller := NewTodoController(mockUsecase, logger)
	router := setupRouter()

	router.DELETE("/todos/:id", controller.Delete)

	t.Run("success", func(t *testing.T) {
		mockUsecase.On("Delete", mock.AnythingOfType("uuid.UUID")).Return(nil).Once()

		req := httptest.NewRequest("DELETE", "/todos/123e4567-e89b-12d3-a456-426614174000", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/todos/invalid-uuid", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
