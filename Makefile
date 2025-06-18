.PHONY: help build run test clean docker-build docker-run docker-stop k8s-deploy k8s-delete

# Переменные
APP_NAME=todo-app
DOCKER_IMAGE=$(APP_NAME):latest
K8S_NAMESPACE=default

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Собрать приложение
	@echo "Сборка приложения..."
	go build -o bin/$(APP_NAME) ./cmd/api

run: ## Запустить приложение локально
	@echo "Запуск приложения..."
	go run ./cmd/api

test: ## Запустить тесты
	@echo "Запуск тестов..."
	go test -v ./...

test-coverage: ## Запустить тесты с покрытием
	@echo "Запуск тестов с покрытием..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean: ## Очистить артефакты сборки
	@echo "Очистка артефактов..."
	rm -rf bin/
	rm -f coverage.out

docker-build: ## Собрать Docker образ
	@echo "Сборка Docker образа..."
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Запустить с Docker Compose
	@echo "Запуск с Docker Compose..."
	docker-compose up -d

docker-stop: ## Остановить Docker Compose
	@echo "Остановка Docker Compose..."
	docker-compose down

docker-logs: ## Показать логи Docker Compose
	docker-compose logs -f

docker-clean: ## Очистить Docker ресурсы
	@echo "Очистка Docker ресурсов..."
	docker-compose down -v
	docker system prune -f

k8s-deploy: ## Развернуть в Kubernetes
	@echo "Развертывание в Kubernetes..."
	kubectl apply -f k8s/

k8s-delete: ## Удалить из Kubernetes
	@echo "Удаление из Kubernetes..."
	kubectl delete -f k8s/

k8s-logs: ## Показать логи в Kubernetes
	kubectl logs -f deployment/$(APP_NAME)

k8s-status: ## Показать статус в Kubernetes
	kubectl get pods -l app=$(APP_NAME)
	kubectl get services -l app=$(APP_NAME)

lint: ## Запустить линтер
	@echo "Запуск линтера..."
	golangci-lint run

fmt: ## Форматировать код
	@echo "Форматирование кода..."
	go fmt ./...

vet: ## Проверить код
	@echo "Проверка кода..."
	go vet ./...

deps: ## Обновить зависимости
	@echo "Обновление зависимостей..."
	go mod tidy
	go mod download

dev: ## Запустить в режиме разработки
	@echo "Запуск в режиме разработки..."
	@if [ ! -f .env ]; then cp .env.example .env; fi
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

dev-stop: ## Остановить режим разработки
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml down

# Команды для мониторинга
prometheus: ## Открыть Prometheus
	@echo "Открытие Prometheus..."
	open http://localhost:9091

grafana: ## Открыть Grafana
	@echo "Открытие Grafana..."
	open http://localhost:3000

api-docs: ## Открыть API документацию
	@echo "Открытие API документации..."
	open http://localhost:8080/swagger/index.html

# Команды для тестирования API
test-api: ## Тестировать API
	@echo "Тестирование API..."
	@curl -s http://localhost:8080/health | jq .
	@curl -s http://localhost:8080/api/v1/todos | jq .

create-todo: ## Создать тестовую задачу
	@echo "Создание тестовой задачи..."
	@curl -X POST http://localhost:8080/api/v1/todos \
		-H "Content-Type: application/json" \
		-d '{"task": "Тестовая задача"}' | jq .

# Команды для базы данных
db-migrate: ## Применить миграции
	@echo "Применение миграций..."
	@# Добавить команды для миграций

db-reset: ## Сбросить базу данных
	@echo "Сброс базы данных..."
	rm -f todos.db

# Команды для мониторинга производительности
bench: ## Запустить бенчмарки
	@echo "Запуск бенчмарков..."
	go test -bench=. ./...

profile: ## Создать профиль производительности
	@echo "Создание профиля производительности..."
	go test -cpuprofile=cpu.prof -memprofile=mem.prof ./...

# Команды для релиза
release: ## Создать релиз
	@echo "Создание релиза..."
	@read -p "Введите версию (например, v1.0.0): " version; \
	git tag $$version; \
	git push origin $$version

install-tools: ## Установить инструменты разработки
	@echo "Установка инструментов разработки..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/go-delve/delve/cmd/dlv@latest

# Команды для CI/CD
ci-build: ## Сборка для CI
	@echo "CI сборка..."
	CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

ci-test: ## Тесты для CI
	@echo "CI тесты..."
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Команды для безопасности
security-scan: ## Сканирование безопасности
	@echo "Сканирование безопасности..."
	@# Добавить команды для сканирования безопасности

# Команды для документации
docs: ## Генерировать документацию
	@echo "Генерация документации..."
	swag init -g cmd/api/main.go

docs-serve: ## Запустить сервер документации
	@echo "Запуск сервера документации..."
	@# Добавить команды для запуска сервера документации 