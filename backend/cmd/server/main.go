package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "auth-project/docs"
	"auth-project/internal/controller"
	"auth-project/internal/controller/middleware"
	"auth-project/internal/infrastructure/database"
	"auth-project/internal/repository"
	"auth-project/internal/service"
)

// @title           Backend API
// @version         1.0.0
// @description     API for authentication, board, and task management

// @contact.name   API Support
// @contact.url    http://example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Bearer token authentication using JWT

func main() {
	// Load configuration
	cfg := loadConfig()

	// Initialize database connection
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepositoryPostgres(db)
	sessionRepo := repository.NewSessionRepositoryPostgres(db)
	columnRepo := repository.NewColumnRepositoryPostgres(db)
	taskRepo := repository.NewTaskRepositoryPostgres(db)
	assignmentHistoryRepo := repository.NewAssignmentHistoryRepositoryPostgres(db)

	// Create cached worklog repository for better performance
	baseWorklogRepo := repository.NewWorklogRepositoryPostgres(db)
	worklogRepo := repository.NewCachedWorklogRepository(baseWorklogRepo, 5*time.Minute)

	// Initialize services
	authService := service.NewAuthService(
		userRepo,
		sessionRepo,
		7*24*time.Hour, // 7 days session TTL
	)
	boardService := service.NewBoardService(columnRepo, taskRepo)
	taskService := service.NewTaskService(taskRepo, columnRepo, userRepo, assignmentHistoryRepo, db)
	worklogService := service.NewWorklogService(worklogRepo, taskRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	boardController := controller.NewBoardController(boardService)
	taskController := controller.NewTaskController(taskService)
	worklogController := controller.NewWorklogController(worklogService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Setup Gin router
	router := setupRouter(authController, boardController, taskController, worklogController, authMiddleware)

	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

// Config holds application configuration
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	ServerPort string
}

func loadConfig() Config {
	return Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "p@ssw0rd"),
		DBName:     getEnv("DB_NAME", "auth_db"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func initDatabase(cfg Config) (*pgxpool.Pool, error) {
	dbConfig := database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	}

	ctx := context.Background()
	pool, err := database.NewPool(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	log.Println("Database connection pool established")
	return pool, nil
}

func setupRouter(
	authController *controller.AuthController,
	boardController *controller.BoardController,
	taskController *controller.TaskController,
	worklogController *controller.WorklogController,
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `<html><head><title>Backend API</title></head><body><h1>Backend API</h1><ul><li><a href="/swagger/index.html">Swagger UI</a></li><li><a href="/health">Health</a></li></ul></body></html>`)
	})

	// CORS middleware with more flexible configuration for development
	// Allow any localhost origin with any port
	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		MaxAge:           12 * time.Hour,
	}

	// In development, allow any localhost origin
	corsConfig.AllowOriginFunc = func(origin string) bool {
		// Allow localhost with any port
		return origin == "http://localhost:3000" ||
			origin == "http://localhost:3001" ||
			origin == "http://localhost:8080" ||
			(len(origin) >= 17 && origin[:17] == "http://localhost:")
	}

	router.Use(cors.New(corsConfig))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Legacy auth routes (for compatibility with old frontend builds)
	router.POST("/auth/register", authController.Register)
	router.POST("/auth/login", authController.Login)
	router.POST("/auth/logout", authController.Logout)
	router.GET("/auth/validate", authController.ValidateSession)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/logout", authController.Logout)
			auth.GET("/me", authMiddleware.RequireAuth(), authController.GetCurrentUser)
			auth.GET("/validate", authController.ValidateSession)
		}

		// Users routes (protected)
		users := v1.Group("/users")
		users.Use(authMiddleware.RequireAuth())
		{
			users.GET("", authController.ListUsers)
		}

		// Board routes (protected)
		board := v1.Group("/board")
		board.Use(authMiddleware.RequireAuth())
		{
			board.GET("", boardController.GetBoard)
		}

		// Task routes (protected)
		tasks := v1.Group("/tasks")
		tasks.Use(authMiddleware.RequireAuth())
		{
			tasks.POST("", taskController.CreateTask)
			tasks.GET("/:id", taskController.GetTask)
			tasks.PUT("/:id", taskController.UpdateTask)
			tasks.DELETE("/:id", taskController.DeleteTask)
			tasks.PATCH("/:id/move", taskController.MoveTask)
			tasks.PATCH("/:id/assignee", taskController.UpdateAssignee)
			tasks.GET("/:id/history", taskController.GetAssignmentHistory)
			tasks.POST("/:id/worklogs", worklogController.LogWork)
			tasks.GET("/:id/worklogs", worklogController.GetWorklogs)
		}

		// Column routes (protected)
		columns := v1.Group("/columns")
		columns.Use(authMiddleware.RequireAuth())
		{
			columns.GET("/:columnId/tasks", taskController.GetTasksByColumn)
		}

		// Report routes (protected)
		reports := v1.Group("/reports")
		reports.Use(authMiddleware.RequireAuth())
		reports.Use(middleware.CacheMiddleware(3 * time.Minute))
		{
			reports.GET("/time", worklogController.GetTimeReport)
		}

		// Protected routes example
		protected := v1.Group("/protected")
		protected.Use(authMiddleware.RequireAuth())
		{
			protected.GET("/test", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Access granted to protected route"})
			})
		}
	}

	return router
}
