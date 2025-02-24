package tests

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"
	"todo-app/domain"
	"todo-app/repository"
	"todo-app/usecase"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, *slog.Logger) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	_ = godotenv.Load("../.env")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "todo_test"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect test database: " + err.Error())
	}
	db.AutoMigrate(&domain.Todo{})
	return db, logger
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func TestTodoUsecase(t *testing.T) {
	db, logger := setupTestDB()
	repo := repository.NewTodoRepo(db, logger)
	uc := usecase.NewTodoUsecase(repo, logger)

	t.Run("Create valid todo", func(t *testing.T) {
		todo := &domain.Todo{
			Title:     "Test Task",
			Status:    "IN_PROGRESS",
			CreatedAt: time.Now().UTC(),
		}

		err := uc.Create(todo)
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, todo.ID)
	})

	t.Run("Create invalid title - empty", func(t *testing.T) {
		todo := &domain.Todo{
			Status:    "IN_PROGRESS",
			CreatedAt: time.Now().UTC(),
		}

		err := uc.Create(todo)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrValidationFailed))
	})

	t.Run("Update non-existent todo", func(t *testing.T) {
		todo := &domain.Todo{
			ID:        uuid.New(),
			Title:     "Test",
			Status:    "IN_PROGRESS",
			CreatedAt: time.Now().UTC(),
		}

		err := uc.Update(todo)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNotFound))
	})

	t.Run("List with invalid sort", func(t *testing.T) {
		_, err := uc.List("invalid", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid sort parameter")
	})

	t.Run("Delete existing todo", func(t *testing.T) {
		// Create a todo first
		todo := &domain.Todo{
			Title:     "Test Delete",
			Status:    "IN_PROGRESS",
			CreatedAt: time.Now().UTC(),
		}
		err := uc.Create(todo)
		assert.NoError(t, err)

		// Delete it
		err = uc.Delete(todo.ID)
		assert.NoError(t, err)

		// Verify it's gone
		// _, err = uc.repo.FindByID(todo.ID)
		// assert.Error(t, err)
		// assert.True(t, errors.Is(err, domain.ErrNotFound))
	})

	t.Run("Delete non-existent todo", func(t *testing.T) {
		id := uuid.New()
		err := uc.Delete(id)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNotFound))
	})

	t.Run("Delete with invalid ID", func(t *testing.T) {
		err := uc.Delete(uuid.Nil)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrValidationFailed))
	})
}
