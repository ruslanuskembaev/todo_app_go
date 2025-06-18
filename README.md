# todo_app_go

Простое TODO-приложение на Go с использованием REST API и SQLite.  
Проект создаётся для практики написания бэкенда

---

## 📦 Стек технологий

- Go 1.22+
- SQLite 
=======
# Todo Application - Production Ready

Современное Go приложение для управления задачами с поддержкой Kubernetes, Redis, Kafka и Grafana.

## 🚀 Возможности

- **RESTful API** для управления задачами
- **Redis кэширование** для улучшения производительности
- **Kafka события** для асинхронной обработки
- **Prometheus метрики** для мониторинга
- **Grafana дашборды** для визуализации
- **Kubernetes готовность** с HPA и health checks
- **Graceful shutdown** для корректного завершения
- **Структурированное логирование** с Zap
- **Валидация данных** с validator
- **CORS поддержка** для веб-приложений

## 🏗️ Архитектура

```
┌─────────────────┐    ┌──────────────┐    ┌──────────────┐
│   Web Client    │    │   Mobile     │    │   API Client │
└─────────┬───────┘    └──────┬───────┘    └──────┬───────┘
          │                   │                   │
          └───────────────────┼───────────────────┘
                              │
                    ┌─────────▼─────────┐
                    │   Load Balancer   │
                    └─────────┬─────────┘
                              │
                    ┌─────────▼─────────┐
                    │   Todo API        │
                    │   (Gin Server)    │
                    └─────────┬─────────┘
                              │
          ┌───────────────────┼───────────────────┐
          │                   │                   │
    ┌─────▼─────┐    ┌────────▼────────┐    ┌────▼────┐
    │   Redis   │    │   SQLite DB     │    │  Kafka  │
    │   Cache   │    │                 │    │ Events  │
    └───────────┘    └─────────────────┘    └─────────┘
```

## 📦 Установка и запуск

### Локальная разработка

1. **Клонируйте репозиторий:**
```bash
git clone <repository-url>
cd todo_app_go
```

2. **Установите зависимости:**
```bash
go mod tidy
```

3. **Запустите с Docker Compose:**
```bash
docker-compose up -d
```

4. **Проверьте работу:**
```bash
# Health check
curl http://localhost:8080/health

# Создать задачу
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"task": "Купить молоко"}'

# Получить все задачи
curl http://localhost:8080/api/v1/todos
```

### Kubernetes развертывание

1. **Соберите Docker образ:**
```bash
docker build -t todo-app:latest .
```

2. **Примените Kubernetes манифесты:**
```bash
kubectl apply -f k8s/
```

3. **Проверьте статус:**
```bash
kubectl get pods
kubectl get services
```

## 🔧 Конфигурация

Приложение использует файл `config.yaml` для настройки:

```yaml
server:
  port: 8080
  host: "0.0.0.0"
  shutdown_timeout: "30s"

database:
  type: "sqlite"
  path: "todos.db"

redis:
  host: "localhost"
  port: 6379

kafka:
  brokers: ["localhost:9092"]
  topic: "todo-events"

metrics:
  enabled: true
  path: "/metrics"

log:
  level: "info"
  format: "json"
```

## 📊 API Endpoints

### Задачи

- `GET /api/v1/todos` - Получить все задачи
- `POST /api/v1/todos` - Создать новую задачу
- `GET /api/v1/todos/{id}` - Получить задачу по ID
- `PUT /api/v1/todos/{id}` - Обновить задачу
- `DELETE /api/v1/todos/{id}` - Удалить задачу

### Системные

- `GET /health` - Health check
- `GET /ready` - Ready check
- `GET /metrics` - Prometheus метрики

## 📈 Мониторинг

### Prometheus метрики

- `http_requests_total` - Общее количество HTTP запросов
- `http_request_duration_seconds` - Время выполнения запросов
- `todo_operations_total` - Операции с задачами
- `cache_hits_total` / `cache_misses_total` - Статистика кэша
- `kafka_messages_published` / `kafka_messages_consumed` - Kafka события

### Grafana дашборды

Доступ к Grafana: http://localhost:3000
- Логин: `admin`
- Пароль: `admin`

## 🔍 Логирование

Приложение использует структурированное логирование с Zap:

```json
{
  "level": "info",
  "timestamp": "2024-01-15T10:30:00Z",
  "caller": "main.go:45",
  "msg": "Todo created successfully",
  "todo_id": 123
}
```

## 🚀 Производительность

- **Кэширование:** Redis для часто запрашиваемых данных
- **Асинхронность:** Kafka для обработки событий
- **Масштабирование:** Kubernetes HPA для автоматического масштабирования
- **Мониторинг:** Prometheus + Grafana для отслеживания метрик

## 🔒 Безопасность

- Непривилегированный пользователь в контейнере
- Health checks для Kubernetes
- Graceful shutdown
- Валидация входных данных
- CORS настройки

## 🧪 Тестирование

```bash
# Запуск тестов
go test ./...

# Запуск с покрытием
go test -cover ./...

# Бенчмарки
go test -bench=. ./...
```

## 📝 Разработка

### Структура проекта

```
todo_app_go/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── cache/
│   ├── config/
│   ├── database/
│   ├── events/
│   ├── handlers/
│   ├── logger/
│   ├── metrics/
│   ├── middleware/
│   ├── models/
│   └── services/
├── k8s/
├── grafana/
├── config.yaml
├── docker-compose.yml
├── Dockerfile
└── README.md
```

### Добавление новых функций

1. Создайте модель в `internal/models/`
2. Добавьте методы в репозиторий
3. Реализуйте бизнес-логику в сервисе
4. Создайте хендлер
5. Добавьте маршрут в `main.go`
6. Напишите тесты

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте feature branch
3. Внесите изменения
4. Добавьте тесты
5. Создайте Pull Request

## 📄 Лицензия

MIT License

## 🆘 Поддержка

При возникновении проблем:

1. Проверьте логи: `docker-compose logs todo-app`
2. Проверьте метрики: http://localhost:9091
3. Проверьте Grafana: http://localhost:3000
4. Создайте Issue в репозитории
>>>>>>> develop
