# Build stage
FROM golang:1.24-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git ca-certificates tzdata gcc musl-dev

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go mod файлы
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

# Создать пользователя и группу
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Создать volume-директорию и назначить владельца
RUN mkdir -p /data && chown -R appuser:appgroup /data

WORKDIR /app

# Копируем бинарник и конфиг с нужным владельцем
COPY --from=builder --chown=appuser:appgroup /app/main /app/main
COPY --from=builder --chown=appuser:appgroup /app/config.yaml /app/config.yaml

USER appuser

# Открываем порты
EXPOSE 8080 9090

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Запускаем приложение
CMD ["./main"] 