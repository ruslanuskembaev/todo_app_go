package services

import (
	"context"
	"testing"
	"time"

	"todo_app_go/internal/logger"
	"todo_app_go/internal/models"

	"github.com/stretchr/testify/assert"
)

func init() {
	_ = logger.Init("debug", "console")
}

type mockRepo struct{}

func (m *mockRepo) Create(task string) (*models.Todo, error) {
	return &models.Todo{
		ID:        1,
		Task:      task,
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *mockRepo) GetByID(id int64) (*models.Todo, error) {
	if id == 42 {
		return &models.Todo{ID: 42, Task: "Answer"}, nil
	}
	if id == 99 {
		return nil, nil // not found
	}
	if id == 500 {
		return nil, assert.AnError
	}
	return nil, nil
}

func (m *mockRepo) GetAll() ([]models.Todo, error) { return nil, nil }
func (m *mockRepo) Update(id int64, req models.TodoUpdateRequest) (*models.Todo, error) {
	if id == 42 {
		return &models.Todo{ID: 42, Task: "Updated"}, nil
	}
	if id == 99 {
		return nil, nil // not found
	}
	if id == 500 {
		return nil, assert.AnError
	}
	return nil, nil
}
func (m *mockRepo) UpdateStatus(id int64, completed bool) error { return nil }
func (m *mockRepo) Delete(id int64) error {
	if id == 500 {
		return assert.AnError
	}
	return nil
}

func TestCreateTodo(t *testing.T) {
	repo := &mockRepo{}
	service := &TodoService{repo: repo}
	ctx := context.Background()
	request := models.TodoCreateRequest{Task: "Test task"}
	todo, err := service.CreateTodo(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, "Test task", todo.Task)
	assert.Equal(t, int64(1), todo.ID)
}

func TestGetTodo_Success(t *testing.T) {
	repo := &mockRepo{}
	service := &TodoService{repo: repo}
	todo, err := service.GetTodo(context.Background(), 42)
	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.Equal(t, int64(42), todo.ID)
}

func TestGetTodo_NotFound(t *testing.T) {
	repo := &mockRepo{}
	service := &TodoService{repo: repo}
	todo, err := service.GetTodo(context.Background(), 99)
	assert.NoError(t, err)
	assert.Nil(t, todo)
}

func TestGetTodo_Error(t *testing.T) {
	repo := &mockRepo{}
	service := &TodoService{repo: repo}
	todo, err := service.GetTodo(context.Background(), 500)
	assert.Error(t, err)
	assert.Nil(t, todo)
}

func TestUpdateTodo_Success(t *testing.T) {
	repo := &mockRepo{}
	service := &TodoService{repo: repo}
	req := models.TodoUpdateRequest{Task: ptrString("Updated")}
	todo, err := service.UpdateTodo(context.Background(), 42, req)
	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.Equal(t, "Updated", todo.Task)
}

func TestUpdateTodo_NotFound(t *testing.T) {
	repo := &mockRepo{}
	service := &TodoService{repo: repo}
	req := models.TodoUpdateRequest{Task: ptrString("Updated")}
	todo, err := service.UpdateTodo(context.Background(), 99, req)
	assert.NoError(t, err)
	assert.Nil(t, todo)
}

func TestUpdateTodo_Error(t *testing.T) {
	repo := &mockRepo{}
	service := &TodoService{repo: repo}
	req := models.TodoUpdateRequest{Task: ptrString("Updated")}
	todo, err := service.UpdateTodo(context.Background(), 500, req)
	assert.Error(t, err)
	assert.Nil(t, todo)
}

func TestDeleteTodo_Success(t *testing.T) {
	repo := &mockRepo{}
	service := &TodoService{repo: repo}
	err := service.DeleteTodo(context.Background(), 1)
	assert.NoError(t, err)
}

func TestDeleteTodo_Error(t *testing.T) {
	repo := &mockRepo{}
	service := &TodoService{repo: repo}
	err := service.DeleteTodo(context.Background(), 500)
	assert.Error(t, err)
}

func ptrString(s string) *string { return &s }
