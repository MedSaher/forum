package models

import (
	"database/sql"
	"errors"
	"fmt"
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

// Get session based on univeral unique id:
// Get session based on universal unique id:
func GetSessionByUUID(uuid string) (*Session, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	var session = &Session{}
	// Ensure we pass the address of each field, so the values can be scanned correctly
	err = db.QueryRow("SELECT ID, UserID, UUID, ExpiresAt FROM Session WHERE UUID = ?", uuid).
		Scan(&session.ID, &session.UserID, &session.UUID, &session.ExpiresAt)  // Use & to pass the address of each field
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("session not found")
		}
		return nil, err
	}
	fmt.Println(session)  // You can log the session to verify it's correct
	return session, nil
}

// Delete session whe log_out:
func DeleteSession(uuid string) error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM Session WHERE UUID = ?", uuid)
	return err
}
