package handlers

import (
	"net/http"
	"strconv"

	"todo_app_go/internal/logger"
	"todo_app_go/internal/metrics"
	"todo_app_go/internal/models"
	"todo_app_go/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type TodoHandler struct {
	service  *services.TodoService
	validate *validator.Validate
}

func NewTodoHandler(service *services.TodoService) *TodoHandler {
	return &TodoHandler{
		service:  service,
		validate: validator.New(),
	}
}

// CreateTodo godoc
// @Summary Create a new todo
// @Description Create a new todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body models.TodoCreateRequest true "Todo to create"
// @Success 201 {object} models.Todo
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /todos [post]
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req models.TodoCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleValidationError(c, err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.handleValidationError(c, err)
		return
	}

	todo, err := h.service.CreateTodo(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to create todo", err)
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// GetTodo godoc
// @Summary Get a todo by ID
// @Description Get a specific todo item by its ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /todos/{id} [get]
func (h *TodoHandler) GetTodo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.handleError(c, http.StatusBadRequest, "Invalid todo ID", err)
		return
	}

	todo, err := h.service.GetTodo(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to get todo", err)
		return
	}

	if todo == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// GetAllTodos godoc
// @Summary Get all todos
// @Description Get all todo items
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {array} models.Todo
// @Failure 500 {object} ErrorResponse
// @Router /todos [get]
func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	todos, err := h.service.GetAllTodos(c.Request.Context())
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to get todos", err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

// UpdateTodo godoc
// @Summary Update a todo
// @Description Update an existing todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body models.TodoUpdateRequest true "Todo updates"
// @Success 200 {object} models.Todo
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /todos/{id} [put]
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.handleError(c, http.StatusBadRequest, "Invalid todo ID", err)
		return
	}

	var req models.TodoUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleValidationError(c, err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.handleValidationError(c, err)
		return
	}

	todo, err := h.service.UpdateTodo(c.Request.Context(), id, req)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to update todo", err)
		return
	}

	if todo == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// DeleteTodo godoc
// @Summary Delete a todo
// @Description Delete a todo item by its ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.handleError(c, http.StatusBadRequest, "Invalid todo ID", err)
		return
	}

	err = h.service.DeleteTodo(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to delete todo", err)
		return
	}

	c.Status(http.StatusNoContent)
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *TodoHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status: "healthy",
	})
}

// ReadyCheck godoc
// @Summary Ready check
// @Description Check if the service is ready to serve requests
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /ready [get]
func (h *TodoHandler) ReadyCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status: "ready",
	})
}

// Вспомогательные методы
func (h *TodoHandler) handleError(c *gin.Context, statusCode int, message string, err error) {
	logger.Error(message, zap.Error(err))

	// Увеличиваем метрики ошибок
	metrics.HttpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), strconv.Itoa(statusCode)).Inc()

	c.JSON(statusCode, ErrorResponse{
		Error: message,
	})
}

func (h *TodoHandler) handleValidationError(c *gin.Context, err error) {
	logger.Error("Validation error", zap.Error(err))

	// Увеличиваем метрики ошибок валидации
	metrics.HttpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), "400").Inc()

	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error:   "Validation failed",
		Details: err.Error(),
	})
}

// Response структуры
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type HealthResponse struct {
	Status string `json:"status"`
}
