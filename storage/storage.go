package storage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "todos.db")
	if err != nil {
		log.Fatal(err)
	}
	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}
	createTable()
}

func createTable() {
	query := `
		CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT
		task TEXT NOT NULL
		)
	`

	if _, err := DB.Exec(query); err != nil {
		log.Fatal(err)
	}
}
