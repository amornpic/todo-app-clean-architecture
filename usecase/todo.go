package usecase

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"todo-app/domain"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type todoUsecase struct {
	repo     domain.TodoRepository
	validate *validator.Validate
	logger   *slog.Logger
}

func NewTodoUsecase(repo domain.TodoRepository, logger *slog.Logger) domain.TodoUsecase {
	return &todoUsecase{
		repo:     repo,
		validate: validator.New(),
		logger:   logger,
	}
}

func (u *todoUsecase) Create(todo *domain.Todo) error {
	if todo.CreatedAt.IsZero() {
		todo.CreatedAt = time.Now().UTC()
	}

	if err := u.validate.Struct(todo); err != nil {
		u.logger.Warn("Validation failed for create", "error", err)
		return fmt.Errorf("%w: %v", domain.ErrValidationFailed, err)
	}

	if err := u.repo.Create(todo); err != nil {
		return err // Error already logged in repository
	}
	return nil
}

func (u *todoUsecase) Update(todo *domain.Todo) error {
	if err := u.validate.Struct(todo); err != nil {
		u.logger.Warn("Validation failed for update", "error", err, "todo_id", todo.ID)
		return fmt.Errorf("%w: %v", domain.ErrValidationFailed, err)
	}

	existing, err := u.repo.FindByID(todo.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return err
		}
		return err // Error already logged in repository
	}
	if existing == nil {
		u.logger.Warn("Todo not found for update", "todo_id", todo.ID)
		return domain.ErrNotFound
	}

	if err := u.repo.Update(todo); err != nil {
		return err // Error already logged in repository
	}
	return nil
}

func (u *todoUsecase) List(sortBy, search string) ([]domain.Todo, error) {
	todos, err := u.repo.FindAll()
	if err != nil {
		return nil, err // Error already logged in repository
	}

	if search != "" {
		filtered := []domain.Todo{}
		for _, t := range todos {
			if strings.Contains(strings.ToLower(t.Title), strings.ToLower(search)) ||
				strings.Contains(strings.ToLower(t.Description), strings.ToLower(search)) {
				filtered = append(filtered, t)
			}
		}
		todos = filtered
		u.logger.Info("Todos filtered", "search", search, "count", len(todos))
	}

	switch sortBy {
	case "title", "date", "status":
		u.sortTodos(todos, sortBy)
		u.logger.Info("Todos sorted", "sort_by", sortBy)
	case "":
		// No sorting
	default:
		u.logger.Warn("Invalid sort parameter", "sort_by", sortBy)
		return nil, fmt.Errorf("invalid sort parameter: %s", sortBy)
	}

	return todos, nil
}

func (u *todoUsecase) Delete(id uuid.UUID) error {
	if id == uuid.Nil {
		u.logger.Warn("Invalid ID for deletion", "todo_id", id)
		return fmt.Errorf("%w: ID cannot be empty", domain.ErrValidationFailed)
	}

	if err := u.repo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (u *todoUsecase) sortTodos(todos []domain.Todo, sortBy string) {
	switch sortBy {
	case "title":
		sortTodos(todos, func(i, j int) bool { return todos[i].Title < todos[j].Title })
	case "date":
		sortTodos(todos, func(i, j int) bool { return todos[i].CreatedAt.Before(todos[j].CreatedAt) })
	case "status":
		sortTodos(todos, func(i, j int) bool { return todos[i].Status < todos[j].Status })
	}
}

func sortTodos(todos []domain.Todo, less func(i, j int) bool) {
	for i := 0; i < len(todos)-1; i++ {
		for j := i + 1; j < len(todos); j++ {
			if !less(i, j) {
				todos[i], todos[j] = todos[j], todos[i]
			}
		}
	}
}
