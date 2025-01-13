package models

import (
	"time"

	"forum/app/config"
)

// Create a structure to represent a session:
type Session struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	UUID      string    `json:"uuid"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// Create a new session:
func InsertSession(session *Session) error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	query := `
    INSERT INTO Session (UserID, UUID, ExpiresAt, CreatedAt)
    VALUES (?, ?, ?, ?);`
	_, err = db.Exec(query, session.UserID, session.UUID, session.ExpiresAt, session.CreatedAt)
	return err
}

// Updating the session:
func UpdateSession(id int, newExpiresAt time.Time) error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	query := `UPDATE Session SET ExpiresAt = ? WHERE ID = ?;`
	_, err = db.Exec(query, newExpiresAt, id)
	return err
}

// Delete a session:
func DeleteSession(id int) error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	query := `DELETE FROM Session WHERE ID = ?;`
	_, err = db.Exec(query, id)
	return err
}
