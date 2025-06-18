package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"todo_app_go/internal/config"
	"todo_app_go/internal/logger"
	"todo_app_go/internal/metrics"
	"todo_app_go/internal/models"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type TodoEvent struct {
	Type      string      `json:"type"` // "created", "updated", "deleted"
	TodoID    int64       `json:"todo_id"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   models.Todo `json:"payload"`
}

type KafkaProducer struct {
	writer *kafka.Writer
	topic  string
}

type KafkaConsumer struct {
	reader *kafka.Reader
	topic  string
}

func NewKafkaProducer(cfg config.KafkaConfig) (*KafkaProducer, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	logger.Info("Kafka producer initialized successfully")
	return &KafkaProducer{
		writer: writer,
		topic:  cfg.Topic,
	}, nil
}

func (p *KafkaProducer) PublishTodoEvent(event TodoEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(fmt.Sprintf("todo-%d", event.TodoID)),
		Value: data,
	})

	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	metrics.KafkaMessagesPublished.Inc()
	logger.Info("Todo event published",
		zap.String("type", event.Type),
		zap.Int64("todo_id", event.TodoID))

	return nil
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

func NewKafkaConsumer(cfg config.KafkaConfig) (*KafkaConsumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topic,
		GroupID:  cfg.GroupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	logger.Info("Kafka consumer initialized successfully")
	return &KafkaConsumer{
		reader: reader,
		topic:  cfg.Topic,
	}, nil
}

func (c *KafkaConsumer) ConsumeMessages(ctx context.Context, handler func(TodoEvent) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				logger.Error("Failed to read message", zap.Error(err))
				continue
			}

			var event TodoEvent
			if err := json.Unmarshal(m.Value, &event); err != nil {
				logger.Error("Failed to unmarshal event", zap.Error(err))
				continue
			}

			if err := handler(event); err != nil {
				logger.Error("Failed to handle event",
					zap.Error(err),
					zap.String("event_type", event.Type))
				continue
			}

			metrics.KafkaMessagesConsumed.Inc()
		}
	}
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}

// Вспомогательные функции для создания событий
func CreateTodoCreatedEvent(todo models.Todo) TodoEvent {
	return TodoEvent{
		Type:      "created",
		TodoID:    todo.ID,
		Timestamp: time.Now(),
		Payload:   todo,
	}
}

func CreateTodoUpdatedEvent(todo models.Todo) TodoEvent {
	return TodoEvent{
		Type:      "updated",
		TodoID:    todo.ID,
		Timestamp: time.Now(),
		Payload:   todo,
	}
}

func CreateTodoDeletedEvent(todoID int64) TodoEvent {
	return TodoEvent{
		Type:      "deleted",
		TodoID:    todoID,
		Timestamp: time.Now(),
		Payload:   models.Todo{ID: todoID},
	}
}
