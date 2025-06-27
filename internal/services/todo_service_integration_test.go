package services

import (
	"context"
	"database/sql"
	"testing"

	"todo_app_go/internal/models"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestTodoService_SQLiteIntegration(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task TEXT NOT NULL,
		completed BOOLEAN NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)`)
	assert.NoError(t, err)

	repo := models.NewSQLiteTodoRepository(db)
	service := NewTodoService(repo, nil, nil)

	ctx := context.Background()
	// Create
	todo, err := service.CreateTodo(ctx, models.TodoCreateRequest{Task: "Integration task"})
	assert.NoError(t, err)
	assert.Equal(t, "Integration task", todo.Task)

	// Get by ID
	got, err := service.GetTodo(ctx, todo.ID)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, todo.ID, got.ID)

	// Update
	newTask := "Updated integration task"
	updateReq := models.TodoUpdateRequest{Task: &newTask}
	updated, err := service.UpdateTodo(ctx, todo.ID, updateReq)
	assert.NoError(t, err)
	assert.Equal(t, newTask, updated.Task)

	// Delete
	err = service.DeleteTodo(ctx, todo.ID)
	assert.NoError(t, err)

	// Get after delete
	deleted, err := service.GetTodo(ctx, todo.ID)
	assert.NoError(t, err)
	assert.Nil(t, deleted)
}
