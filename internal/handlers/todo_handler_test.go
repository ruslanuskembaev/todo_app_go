package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"todo_app_go/internal/logger"
	"todo_app_go/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	_ = logger.Init("debug", "console")
}

func TestGetAllTodos(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockRepo{}
	service := services.NewTodoService(repo, nil, nil)
	h := NewTodoHandler(service)

	r := gin.New()
	r.GET("/todos", h.GetAllTodos)

	req, _ := http.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test task")
}

func TestCreateTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockRepo{}
	service := services.NewTodoService(repo, nil, nil)
	h := NewTodoHandler(service)

	r := gin.New()
	r.POST("/todos", h.CreateTodo)

	body := `{"task": "New task"}`
	req, _ := http.NewRequest("POST", "/todos", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "New task")
}

func TestGetTodo_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockRepo{}
	service := services.NewTodoService(repo, nil, nil)
	h := NewTodoHandler(service)

	r := gin.New()
	r.GET("/todos/:id", h.GetTodo)

	req, _ := http.NewRequest("GET", "/todos/42", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Answer")
}

func TestGetTodo_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockRepo{}
	service := services.NewTodoService(repo, nil, nil)
	h := NewTodoHandler(service)

	r := gin.New()
	r.GET("/todos/:id", h.GetTodo)

	req, _ := http.NewRequest("GET", "/todos/99", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Todo not found")
}

func TestUpdateTodo_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockRepo{}
	service := services.NewTodoService(repo, nil, nil)
	h := NewTodoHandler(service)

	r := gin.New()
	r.PUT("/todos/:id", h.UpdateTodo)

	body := `{"task": "Updated"}`
	req, _ := http.NewRequest("PUT", "/todos/42", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Updated")
}

func TestUpdateTodo_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockRepo{}
	service := services.NewTodoService(repo, nil, nil)
	h := NewTodoHandler(service)

	r := gin.New()
	r.PUT("/todos/:id", h.UpdateTodo)

	body := `{"task": "Updated"}`
	req, _ := http.NewRequest("PUT", "/todos/99", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Todo not found")
}

func TestDeleteTodo_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockRepo{}
	service := services.NewTodoService(repo, nil, nil)
	h := NewTodoHandler(service)

	r := gin.New()
	r.DELETE("/todos/:id", h.DeleteTodo)

	req, _ := http.NewRequest("DELETE", "/todos/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteTodo_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockRepo{}
	service := services.NewTodoService(repo, nil, nil)
	h := NewTodoHandler(service)

	r := gin.New()
	r.DELETE("/todos/:id", h.DeleteTodo)

	req, _ := http.NewRequest("DELETE", "/todos/500", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to delete todo")
}

func TestCreateTodo_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockRepo{}
	service := NewTodoHandler(services.NewTodoService(repo, nil, nil))

	r := gin.New()
	r.POST("/todos", service.CreateTodo)

	body := `{"task": ""}` // пустой task
	req, _ := http.NewRequest("POST", "/todos", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Validation failed")
}

func TestUpdateTodo_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockRepo{}
	service := NewTodoHandler(services.NewTodoService(repo, nil, nil))

	r := gin.New()
	r.PUT("/todos/:id", service.UpdateTodo)

	body := `{"task": ""}` // пустой task
	req, _ := http.NewRequest("PUT", "/todos/42", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Validation failed")
}
