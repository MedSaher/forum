package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func Connection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return db, nil
}
