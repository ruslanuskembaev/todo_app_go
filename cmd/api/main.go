package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"todo_app_go/internal/cache"
	"todo_app_go/internal/config"
	"todo_app_go/internal/database"
	"todo_app_go/internal/events"
	"todo_app_go/internal/handlers"
	"todo_app_go/internal/logger"
	"todo_app_go/internal/middleware"
	"todo_app_go/internal/models"
	"todo_app_go/internal/services"

	// Swagger
	_ "todo_app_go/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// @title Todo App API
// @version 1.0
// @description Production-ready Todo API with Redis, Kafka, Prometheus, Grafana, Kubernetes.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Инициализируем логгер
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Get().Sync()

	logger.Info("Starting Todo Application")

	// Инициализируем базу данных
	db, err := database.NewSQLiteDB(database.Config{
		Path: cfg.Database.Path,
	})
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer db.Close()

	// Автоматически создаём таблицу и столбец updated_at, если их нет
	if err := database.EnsureTodosTableAndColumn(db); err != nil {
		logger.Fatal("Failed to ensure todos table and column", zap.Error(err))
	}

	// Инициализируем репозиторий
	repo := models.NewSQLiteTodoRepository(db)

	// Инициализируем Redis кэш (опционально)
	var redisCache *cache.RedisCache
	if cfg.Redis.Host != "" {
		redisCache, err = cache.NewRedisCache(cfg.Redis)
		if err != nil {
			logger.Warn("Failed to initialize Redis cache, continuing without cache", zap.Error(err))
		} else {
			defer redisCache.Close()
		}
	}

	// Инициализируем Kafka producer (опционально)
	var kafkaProducer *events.KafkaProducer
	if len(cfg.Kafka.Brokers) > 0 {
		kafkaProducer, err = events.NewKafkaProducer(cfg.Kafka)
		if err != nil {
			logger.Warn("Failed to initialize Kafka producer, continuing without events", zap.Error(err))
		} else {
			defer kafkaProducer.Close()
		}
	}

	// Инициализируем сервис
	todoService := services.NewTodoService(repo, redisCache, kafkaProducer)

	// Инициализируем хендлеры
	todoHandler := handlers.NewTodoHandler(todoService)

	// Настраиваем Gin
	if cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.CORSGin())
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(middleware.Timeout(cfg.Server.ReadTimeout))

	// Метрики
	if cfg.Metrics.Enabled {
		router.GET(cfg.Metrics.Path, gin.WrapH(promhttp.Handler()))
	}

	// Health checks
	router.GET("/health", todoHandler.HealthCheck)
	router.GET("/ready", todoHandler.ReadyCheck)

	// API routes
	api := router.Group("/api/v1")
	{
		todos := api.Group("/todos")
		{
			todos.GET("", todoHandler.GetAllTodos)
			todos.POST("", todoHandler.CreateTodo)
			todos.GET("/:id", todoHandler.GetTodo)
			todos.PUT("/:id", todoHandler.UpdateTodo)
			todos.DELETE("/:id", todoHandler.DeleteTodo)
		}
	}

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Создаем HTTP сервер
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Запускаем сервер в горутине
	go func() {
		logger.Info("Starting HTTP server", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Ждем сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
