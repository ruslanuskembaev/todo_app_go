package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Config holds the database configuration
type Config struct {
	Path string
}

// NewSQLiteDB creates a new SQLite database connection
func NewSQLiteDB(config Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", config.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %v", err)
	}

	return db, nil
}

// initSchema creates the necessary database tables
func initSchema(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			task TEXT NOT NULL,
			completed BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create todos table: %v", err)
	}

	return nil
}
