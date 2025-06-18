package services

import (
	"context"
	"time"

	"todo_app_go/internal/cache"
	"todo_app_go/internal/events"
	"todo_app_go/internal/logger"
	"todo_app_go/internal/metrics"
	"todo_app_go/internal/models"

	"go.uber.org/zap"
)

type TodoService struct {
	repo     models.TodoRepository
	cache    *cache.RedisCache
	producer *events.KafkaProducer
}

func NewTodoService(repo models.TodoRepository, cache *cache.RedisCache, producer *events.KafkaProducer) *TodoService {
	return &TodoService{
		repo:     repo,
		cache:    cache,
		producer: producer,
	}
}

func (s *TodoService) CreateTodo(ctx context.Context, req models.TodoCreateRequest) (*models.Todo, error) {
	start := time.Now()
	defer func() {
		metrics.TodoOperationsDuration.WithLabelValues("create").Observe(time.Since(start).Seconds())
	}()

	// Создаем todo в базе данных
	todo, err := s.repo.Create(req.Task)
	if err != nil {
		metrics.TodoOperationsTotal.WithLabelValues("create", "error").Inc()
		return nil, err
	}

	// Сохраняем в кэш
	if s.cache != nil {
		if err := s.cache.SetTodo(todo, 30*time.Minute); err != nil {
			logger.Warn("Failed to cache todo", zap.Error(err))
		}
		// Инвалидируем список todos
		if err := s.cache.InvalidateTodos(); err != nil {
			logger.Warn("Failed to invalidate todos cache", zap.Error(err))
		}
	}

	// Публикуем событие
	if s.producer != nil {
		event := events.CreateTodoCreatedEvent(*todo)
		if err := s.producer.PublishTodoEvent(event); err != nil {
			logger.Error("Failed to publish todo created event", zap.Error(err))
		}
	}

	metrics.TodoOperationsTotal.WithLabelValues("create", "success").Inc()
	logger.Info("Todo created successfully", zap.Int64("todo_id", todo.ID))

	return todo, nil
}

func (s *TodoService) GetTodo(ctx context.Context, id int64) (*models.Todo, error) {
	start := time.Now()
	defer func() {
		metrics.TodoOperationsDuration.WithLabelValues("get").Observe(time.Since(start).Seconds())
	}()

	// Пытаемся получить из кэша
	if s.cache != nil {
		if todo, err := s.cache.GetTodo(id); err == nil && todo != nil {
			metrics.TodoOperationsTotal.WithLabelValues("get", "cache_hit").Inc()
			return todo, nil
		}
	}

	// Получаем из базы данных
	todo, err := s.repo.GetByID(id)
	if err != nil {
		metrics.TodoOperationsTotal.WithLabelValues("get", "error").Inc()
		return nil, err
	}

	if todo == nil {
		metrics.TodoOperationsTotal.WithLabelValues("get", "not_found").Inc()
		return nil, nil
	}

	// Сохраняем в кэш
	if s.cache != nil {
		if err := s.cache.SetTodo(todo, 30*time.Minute); err != nil {
			logger.Warn("Failed to cache todo", zap.Error(err))
		}
	}

	metrics.TodoOperationsTotal.WithLabelValues("get", "success").Inc()
	return todo, nil
}

func (s *TodoService) GetAllTodos(ctx context.Context) ([]models.Todo, error) {
	start := time.Now()
	defer func() {
		metrics.TodoOperationsDuration.WithLabelValues("get_all").Observe(time.Since(start).Seconds())
	}()

	// Пытаемся получить из кэша
	if s.cache != nil {
		if todos, err := s.cache.GetTodos(); err == nil && todos != nil {
			metrics.TodoOperationsTotal.WithLabelValues("get_all", "cache_hit").Inc()
			return todos, nil
		}
	}

	// Получаем из базы данных
	todos, err := s.repo.GetAll()
	if err != nil {
		metrics.TodoOperationsTotal.WithLabelValues("get_all", "error").Inc()
		return nil, err
	}

	// Сохраняем в кэш
	if s.cache != nil {
		if err := s.cache.SetTodos(todos, 5*time.Minute); err != nil {
			logger.Warn("Failed to cache todos", zap.Error(err))
		}
	}

	metrics.TodoOperationsTotal.WithLabelValues("get_all", "success").Inc()
	return todos, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, id int64, req models.TodoUpdateRequest) (*models.Todo, error) {
	start := time.Now()
	defer func() {
		metrics.TodoOperationsDuration.WithLabelValues("update").Observe(time.Since(start).Seconds())
	}()

	// Обновляем в базе данных
	todo, err := s.repo.Update(id, req)
	if err != nil {
		metrics.TodoOperationsTotal.WithLabelValues("update", "error").Inc()
		return nil, err
	}

	if todo == nil {
		metrics.TodoOperationsTotal.WithLabelValues("update", "not_found").Inc()
		return nil, nil
	}

	// Обновляем кэш
	if s.cache != nil {
		if err := s.cache.SetTodo(todo, 30*time.Minute); err != nil {
			logger.Warn("Failed to cache updated todo", zap.Error(err))
		}
		// Инвалидируем список todos
		if err := s.cache.InvalidateTodos(); err != nil {
			logger.Warn("Failed to invalidate todos cache", zap.Error(err))
		}
	}

	// Публикуем событие
	if s.producer != nil {
		event := events.CreateTodoUpdatedEvent(*todo)
		if err := s.producer.PublishTodoEvent(event); err != nil {
			logger.Error("Failed to publish todo updated event", zap.Error(err))
		}
	}

	metrics.TodoOperationsTotal.WithLabelValues("update", "success").Inc()
	logger.Info("Todo updated successfully", zap.Int64("todo_id", todo.ID))

	return todo, nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, id int64) error {
	start := time.Now()
	defer func() {
		metrics.TodoOperationsDuration.WithLabelValues("delete").Observe(time.Since(start).Seconds())
	}()

	// Удаляем из базы данных
	err := s.repo.Delete(id)
	if err != nil {
		metrics.TodoOperationsTotal.WithLabelValues("delete", "error").Inc()
		return err
	}

	// Удаляем из кэша
	if s.cache != nil {
		if err := s.cache.DeleteTodo(id); err != nil {
			logger.Warn("Failed to delete todo from cache", zap.Error(err))
		}
		// Инвалидируем список todos
		if err := s.cache.InvalidateTodos(); err != nil {
			logger.Warn("Failed to invalidate todos cache", zap.Error(err))
		}
	}

	// Публикуем событие
	if s.producer != nil {
		event := events.CreateTodoDeletedEvent(id)
		if err := s.producer.PublishTodoEvent(event); err != nil {
			logger.Error("Failed to publish todo deleted event", zap.Error(err))
		}
	}

	metrics.TodoOperationsTotal.WithLabelValues("delete", "success").Inc()
	logger.Info("Todo deleted successfully", zap.Int64("todo_id", id))

	return nil
}
