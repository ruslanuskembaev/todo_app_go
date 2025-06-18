package models

import (
	"database/sql"
	"time"
)

// Todo represents a todo item in the system
type Todo struct {
	ID        int64     `json:"id"`
	Task      string    `json:"task" validate:"required,min=1,max=500"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TodoCreateRequest represents a request to create a new todo
type TodoCreateRequest struct {
	Task string `json:"task" validate:"required,min=1,max=500"`
}

// TodoUpdateRequest represents a request to update a todo
type TodoUpdateRequest struct {
	Task      *string `json:"task,omitempty" validate:"omitempty,min=1,max=500"`
	Completed *bool   `json:"completed,omitempty"`
}

// TodoRepository defines the interface for todo storage operations
type TodoRepository interface {
	Create(task string) (*Todo, error)
	GetByID(id int64) (*Todo, error)
	GetAll() ([]Todo, error)
	Update(id int64, req TodoUpdateRequest) (*Todo, error)
	UpdateStatus(id int64, completed bool) error
	Delete(id int64) error
}

// SQLiteTodoRepository implements TodoRepository for SQLite
type SQLiteTodoRepository struct {
	db *sql.DB
}

// NewSQLiteTodoRepository creates a new SQLiteTodoRepository
func NewSQLiteTodoRepository(db *sql.DB) *SQLiteTodoRepository {
	return &SQLiteTodoRepository{db: db}
}

// Create adds a new todo to the database
func (r *SQLiteTodoRepository) Create(task string) (*Todo, error) {
	now := time.Now()
	result, err := r.db.Exec("INSERT INTO todos (task, completed, created_at, updated_at) VALUES (?, ?, ?, ?)",
		task, false, now, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Todo{
		ID:        id,
		Task:      task,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// GetByID retrieves a todo by ID
func (r *SQLiteTodoRepository) GetByID(id int64) (*Todo, error) {
	var todo Todo
	err := r.db.QueryRow("SELECT id, task, completed, created_at, updated_at FROM todos WHERE id = ?", id).
		Scan(&todo.ID, &todo.Task, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &todo, nil
}

// GetAll retrieves all todos from the database
func (r *SQLiteTodoRepository) GetAll() ([]Todo, error) {
	rows, err := r.db.Query("SELECT id, task, completed, created_at, updated_at FROM todos ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Task, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

// Update updates a todo
func (r *SQLiteTodoRepository) Update(id int64, req TodoUpdateRequest) (*Todo, error) {
	// Сначала получаем текущий todo
	todo, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, nil // Todo не найден
	}

	// Обновляем поля
	if req.Task != nil {
		todo.Task = *req.Task
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	todo.UpdatedAt = time.Now()

	// Обновляем в базе
	_, err = r.db.Exec("UPDATE todos SET task = ?, completed = ?, updated_at = ? WHERE id = ?",
		todo.Task, todo.Completed, todo.UpdatedAt, id)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// UpdateStatus updates the completion status of a todo
func (r *SQLiteTodoRepository) UpdateStatus(id int64, completed bool) error {
	_, err := r.db.Exec("UPDATE todos SET completed = ?, updated_at = ? WHERE id = ?",
		completed, time.Now(), id)
	return err
}

// Delete removes a todo from the database
func (r *SQLiteTodoRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}
