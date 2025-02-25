package route

import (
	"log/slog"
	"todo-app/api/controller"
	"todo-app/repository"
	"todo-app/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(gin *gin.Engine, db *gorm.DB, logger *slog.Logger) {
	NewTodoRoter(gin, db, logger)
}

func NewTodoRoter(gin *gin.Engine, db *gorm.DB, logger *slog.Logger) {
	repo := repository.NewTodoRepo(db, logger)
	usecase := usecase.NewTodoUsecase(repo, logger)
	tc := controller.NewTodoController(usecase, logger)

	gin.POST("/todos", tc.Create)
	gin.PUT("/todos/:id", tc.Update)
	gin.GET("/todos", tc.List)
	gin.DELETE("/todos/:id", tc.Delete)
}
