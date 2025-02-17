package models

import (
	"database/sql"
	"errors"
	"fmt"

	// "fmt"
	"net/http"
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

// Get session based on universal unique id:
func GetSessionByUUID(uuid string) (*Session, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	session := &Session{}
	// Ensure we pass the address of each field, so the values can be scanned correctly
	err = db.QueryRow("SELECT ID, UserID, UUID, ExpiresAt FROM Session WHERE UUID = ?", uuid).
		Scan(&session.ID, &session.UserID, &session.UUID, &session.ExpiresAt) // Use & to pass the address of each field
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("session not found")
		}
		return nil, err
	}
	// fmt.Println(session) // You can log the session to verify it's correct
	return session, nil
}

// Delete session whe log_out:
func DeleteSessionByUUID(uuid string) error {
	db, err := config.InitDB()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}
	defer db.Close() // Close the DB connection if it's a new one every time.

	// Execute DELETE query
	result, err := db.Exec("DELETE FROM Session WHERE UUID = ?", uuid)
	if err != nil {
		return fmt.Errorf("failed to delete session: %v", err)
	}

	// Check if a row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no session found with UUID: %s", uuid)
	}

	return nil
}

// Get user id from session:
func GetUserIdFromSession(rq *http.Request) (int, error) {
	cookie, err := rq.Cookie("session_token")
	if err != nil {
		return -1, err
	}
	session, err := GetSessionByUUID(cookie.Value)
	if err != nil {
		return -1, err
	}
	uuid := session.UUID
	user, err := GetUserByTocken(uuid)
	if err != nil {
		return -1, err
	}
	return user.ID, nil
}

// Get session by user id:
func GetSessionByUserID(userID int) (*Session, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var session Session
	err = db.QueryRow("SELECT ID, UserID, UUID, ExpiresAt FROM Session WHERE UserID = ? AND ExpiresAt > datetime('now')", userID).
		Scan(&session.ID, &session.UserID, &session.UUID, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func DeleteAllSessions() {
	db, err := config.InitDB()
	if err == nil {
		db.Exec("DELETE FROM Session")
	}
	defer db.Close()
	
}