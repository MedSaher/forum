package config

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Initialize a connection with database:
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		return nil, err
	}

	// Ensure the connection is valid
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
