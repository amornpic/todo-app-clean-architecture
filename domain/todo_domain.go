package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrNotFound          = errors.New("todo not found")
	ErrValidationFailed  = errors.New("validation failed")
	ErrDatabaseOperation = errors.New("database operation failed")
)

type Todo struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Title       string    `json:"title" gorm:"type:varchar(100);not null" validate:"required,max=100"`
	Description string    `json:"description" gorm:"type:text" validate:"omitempty"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Image       string    `json:"image" gorm:"type:text" validate:"omitempty,base64"`
	Status      string    `json:"status" gorm:"type:varchar(20);not null" validate:"required,oneof=IN_PROGRESS COMPLETED"`
}

type TodoRepository interface {
	Create(todo *Todo) error
	Update(todo *Todo) error
	FindAll() ([]Todo, error)
	FindByID(id uuid.UUID) (*Todo, error)
	Delete(id uuid.UUID) error
}

type TodoUsecase interface {
	Create(todo *Todo) error
	Update(todo *Todo) error
	List(sortBy, search string) ([]Todo, error)
	Delete(id uuid.UUID) error
}

func (t *Todo) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}
