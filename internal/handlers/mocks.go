package handlers

import (
	"time"
	"todo_app_go/internal/models"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct{}

func (m *mockRepo) Create(task string) (*models.Todo, error) {
	return &models.Todo{
		ID:        123,
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
	if id == 1 {
		return &models.Todo{ID: 1, Task: "Test task"}, nil
	}
	if id == 500 {
		return nil, assert.AnError
	}
	return nil, nil // not found
}

func (m *mockRepo) Update(id int64, req models.TodoUpdateRequest) (*models.Todo, error) {
	if id == 42 {
		return &models.Todo{ID: 42, Task: "Updated"}, nil
	}
	if id == 500 {
		return nil, assert.AnError
	}
	return nil, nil // not found
}

func (m *mockRepo) Delete(id int64) error {
	if id == 500 {
		return assert.AnError
	}
	return nil
}

func (m *mockRepo) GetAll() ([]models.Todo, error) {
	return []models.Todo{{ID: 1, Task: "Test task"}}, nil
}
func (m *mockRepo) UpdateStatus(id int64, completed bool) error { return nil }
