package repository

import (
	"log/slog"
	"testing"
	"todo-app/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&domain.Todo{})
	assert.NoError(t, err)

	return db
}

func TestTodoRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	logger := slog.Default()
	repo := NewTodoRepo(db, logger)

	t.Run("success", func(t *testing.T) {
		todo := &domain.Todo{
			Title:       "Test Todo",
			Description: "Test Description",
			Status:      "IN_PROGRESS",
		}

		err := repo.Create(todo)

		assert.NoError(t, err)
		assert.NotZero(t, todo.ID)

		// Verify the todo was created
		var found domain.Todo
		err = db.First(&found, todo.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, todo.Title, found.Title)
	})
}

func TestTodoRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	logger := slog.Default()
	repo := NewTodoRepo(db, logger)

	t.Run("success", func(t *testing.T) {
		// Create a test todo
		todo := &domain.Todo{
			Title:       "Test Todo",
			Description: "Test Description",
			Status:      "IN_PROGRESS",
		}
		err := db.Create(todo).Error
		assert.NoError(t, err)

		// Find the created todo
		found, err := repo.FindByID(todo.ID)

		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, todo.ID, found.ID)
		assert.Equal(t, todo.Title, found.Title)
	})

	t.Run("not found", func(t *testing.T) {
		nonExistentID := uuid.New()
		found, err := repo.FindByID(nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, found)
	})
}
