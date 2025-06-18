package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"todo_app_go/internal/config"
	"todo_app_go/internal/logger"
	"todo_app_go/internal/metrics"
	"todo_app_go/internal/models"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(cfg config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()

	// Проверяем подключение
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis cache initialized successfully")
	return &RedisCache{
		client: client,
		ctx:    ctx,
	}, nil
}

func (c *RedisCache) GetTodo(id int64) (*models.Todo, error) {
	key := fmt.Sprintf("todo:%d", id)

	data, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			metrics.CacheMissesTotal.Inc()
			return nil, nil // Кэш miss
		}
		return nil, err
	}

	metrics.CacheHitsTotal.Inc()

	var todo models.Todo
	if err := json.Unmarshal([]byte(data), &todo); err != nil {
		return nil, err
	}

	return &todo, nil
}

func (c *RedisCache) SetTodo(todo *models.Todo, expiration time.Duration) error {
	key := fmt.Sprintf("todo:%d", todo.ID)

	data, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	return c.client.Set(c.ctx, key, data, expiration).Err()
}

func (c *RedisCache) DeleteTodo(id int64) error {
	key := fmt.Sprintf("todo:%d", id)
	return c.client.Del(c.ctx, key).Err()
}

func (c *RedisCache) GetTodos() ([]models.Todo, error) {
	key := "todos:all"

	data, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			metrics.CacheMissesTotal.Inc()
			return nil, nil // Кэш miss
		}
		return nil, err
	}

	metrics.CacheHitsTotal.Inc()

	var todos []models.Todo
	if err := json.Unmarshal([]byte(data), &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func (c *RedisCache) SetTodos(todos []models.Todo, expiration time.Duration) error {
	key := "todos:all"

	data, err := json.Marshal(todos)
	if err != nil {
		return err
	}

	return c.client.Set(c.ctx, key, data, expiration).Err()
}

func (c *RedisCache) InvalidateTodos() error {
	key := "todos:all"
	return c.client.Del(c.ctx, key).Err()
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}
