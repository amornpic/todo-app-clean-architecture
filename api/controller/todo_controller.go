package controller

import (
	"errors"
	"log/slog"
	"net/http"
	"todo-app/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TodoController struct {
	usecase domain.TodoUsecase
	logger  *slog.Logger
}

func NewTodoRoter(r *gin.Engine, u domain.TodoUsecase, logger *slog.Logger) {
	h := &TodoController{usecase: u, logger: logger}
	r.POST("/todos", h.Create)
	r.PUT("/todos/:id", h.Update)
	r.GET("/todos", h.List)
	r.DELETE("/todos/:id", h.Delete)
}

// Create creates a new todo
// @Summary Create a new todo
// @Description Create a new todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body domain.Todo true "Todo object"
// @Success 201 {object} domain.Todo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos [post]
func (h *TodoController) Create(c *gin.Context) {
	var todo domain.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		h.logger.Warn("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	if err := h.usecase.Create(&todo); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, todo)
}

// Update updates an existing todo
// @Summary Update a todo
// @Description Update a todo item by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param todo body domain.Todo true "Todo object"
// @Success 200 {object} domain.Todo
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos/{id} [put]
func (h *TodoController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid UUID", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var todo domain.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		h.logger.Warn("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}
	todo.ID = id

	if err := h.usecase.Update(&todo); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, todo)
}

// List returns a list of todos
// @Summary List todos
// @Description Get a list of todos with optional sorting and searching
// @Tags todos
// @Produce json
// @Param sort_by query string false "Sort by field (title, date, status)"
// @Param search query string false "Search in title or description"
// @Success 200 {array} domain.Todo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos [get]
func (h *TodoController) List(c *gin.Context) {
	sortBy := c.Query("sort_by")
	search := c.Query("search")

	todos, err := h.usecase.List(sortBy, search)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, todos)
}

// Delete removes a todo
// @Summary Delete a todo
// @Description Delete a todo item by ID
// @Tags todos
// @Param id path string true "Todo ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos/{id} [delete]
func (h *TodoController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid UUID", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.usecase.Delete(id); err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TodoController) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidationFailed):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
	case errors.Is(err, domain.ErrDatabaseOperation):
		h.logger.Error("Database error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	default:
		h.logger.Error("Unexpected error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
