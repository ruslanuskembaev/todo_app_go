package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP метрики
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Todo метрики
	TodoOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "todo_operations_total",
			Help: "Total number of todo operations",
		},
		[]string{"operation", "status"},
	)

	TodoOperationsDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "todo_operations_duration_seconds",
			Help:    "Todo operations duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// Cache метрики
	CacheHitsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	CacheMissesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	// Kafka метрики
	KafkaMessagesPublished = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_messages_published_total",
			Help: "Total number of Kafka messages published",
		},
	)

	KafkaMessagesConsumed = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_messages_consumed_total",
			Help: "Total number of Kafka messages consumed",
		},
	)

	// Системные метрики
	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)

	DatabaseConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections",
			Help: "Number of database connections",
		},
	)
)
