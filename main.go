package main

import (
	"fmt"
	"log/slog"
	"os"
	"todo-app/api/route"
	"todo-app/docs"
	"todo-app/domain"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Warn("Error loading .env file, using system environment variables", "error", err)
	}

	// Database connection
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		getEnv("DB_HOST", "postgres"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "todo"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		panic("failed to connect database: " + err.Error())
	}
	logger.Info("Database connected successfully")

	if err := db.AutoMigrate(&domain.Todo{}); err != nil {
		logger.Error("Failed to migrate database", "error", err)
		panic("failed to migrate database: " + err.Error())
	}

	gin := gin.Default()
	docs.SwaggerInfo.BasePath = ""

	route.Setup(gin, db, logger)

	gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := getEnv("APP_PORT", "8080")
	logger.Info("Starting server", "port", port)
	if err := gin.Run(":" + port); err != nil {
		logger.Error("Failed to start server", "error", err)
		panic("failed to start server: " + err.Error())
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
