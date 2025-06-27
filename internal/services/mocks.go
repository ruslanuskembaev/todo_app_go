package services

import (
	"context"
	"todo_app_go/internal/models"
)

type MockTodoService struct {
	GetAllTodosFunc func(ctx context.Context) ([]models.Todo, error)
	CreateTodoFunc  func(ctx context.Context, req models.TodoCreateRequest) (*models.Todo, error)
	GetTodoFunc     func(ctx context.Context, id int64) (*models.Todo, error)
	UpdateTodoFunc  func(ctx context.Context, id int64, req models.TodoUpdateRequest) (*models.Todo, error)
	DeleteTodoFunc  func(ctx context.Context, id int64) error
}

func (m *MockTodoService) GetAllTodos(ctx context.Context) ([]models.Todo, error) {
	return m.GetAllTodosFunc(ctx)
}
func (m *MockTodoService) CreateTodo(ctx context.Context, req models.TodoCreateRequest) (*models.Todo, error) {
	return m.CreateTodoFunc(ctx, req)
}
func (m *MockTodoService) GetTodo(ctx context.Context, id int64) (*models.Todo, error) {
	return m.GetTodoFunc(ctx, id)
}
func (m *MockTodoService) UpdateTodo(ctx context.Context, id int64, req models.TodoUpdateRequest) (*models.Todo, error) {
	return m.UpdateTodoFunc(ctx, id, req)
}
func (m *MockTodoService) DeleteTodo(ctx context.Context, id int64) error {
	return m.DeleteTodoFunc(ctx, id)
}
