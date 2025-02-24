package repository

import (
	"fmt"
	"log/slog"
	"todo-app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewTodoRepo(db *gorm.DB, logger *slog.Logger) *TodoRepo {
	return &TodoRepo{db: db, logger: logger}
}

func (r *TodoRepo) Create(todo *domain.Todo) error {
	if err := r.db.Create(todo).Error; err != nil {
		r.logger.Error("Failed to create todo", "error", err, "todo_id", todo.ID)
		return fmt.Errorf("%w: %v", domain.ErrDatabaseOperation, err)
	}
	r.logger.Info("Todo created", "todo_id", todo.ID)
	return nil
}

func (r *TodoRepo) Update(todo *domain.Todo) error {
	if err := r.db.Save(todo).Error; err != nil {
		r.logger.Error("Failed to update todo", "error", err, "todo_id", todo.ID)
		return fmt.Errorf("%w: %v", domain.ErrDatabaseOperation, err)
	}
	r.logger.Info("Todo updated", "todo_id", todo.ID)
	return nil
}

func (r *TodoRepo) FindAll() ([]domain.Todo, error) {
	var todos []domain.Todo
	if err := r.db.Find(&todos).Error; err != nil {
		r.logger.Error("Failed to list todos", "error", err)
		return nil, fmt.Errorf("%w: %v", domain.ErrDatabaseOperation, err)
	}
	r.logger.Info("Todos retrieved", "count", len(todos))
	return todos, nil
}

func (r *TodoRepo) FindByID(id uuid.UUID) (*domain.Todo, error) {
	var todo domain.Todo
	err := r.db.First(&todo, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warn("Todo not found", "todo_id", id)
			return nil, domain.ErrNotFound
		}
		r.logger.Error("Failed to find todo", "error", err, "todo_id", id)
		return nil, fmt.Errorf("%w: %v", domain.ErrDatabaseOperation, err)
	}
	r.logger.Info("Todo retrieved", "todo_id", id)
	return &todo, nil
}

func (r *TodoRepo) Delete(id uuid.UUID) error {
	result := r.db.Delete(&domain.Todo{}, "id = ?", id)
	if err := result.Error; err != nil {
		r.logger.Error("Failed to delete todo", "error", err, "todo_id", id)
		return fmt.Errorf("%w: %v", domain.ErrDatabaseOperation, err)
	}
	if result.RowsAffected == 0 {
		r.logger.Warn("Todo not found for deletion", "todo_id", id)
		return domain.ErrNotFound
	}
	r.logger.Info("Todo deleted", "todo_id", id)
	return nil
}
