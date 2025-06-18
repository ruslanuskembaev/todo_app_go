package database

import (
	"database/sql"
)

func EnsureTodosTableAndColumn(db *sql.DB) error {
	// Создаём таблицу, если её нет
	createTable := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL
	);
	`
	if _, err := db.Exec(createTable); err != nil {
		return err
	}

	// Проверяем, есть ли столбец updated_at
	found := false
	rows, err := db.Query("PRAGMA table_info(todos);")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dfltValue interface{}
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			return err
		}
		if name == "updated_at" {
			found = true
			break
		}
	}
	if !found {
		_, err := db.Exec(`ALTER TABLE todos ADD COLUMN updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP;`)
		if err != nil {
			return err
		}
	}
	return nil
}
