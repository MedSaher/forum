package middleware

import "time"

// Create a structure to represent a session:
type Session struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	UUID      string    `json:"uuid"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// Create a new session
func CreateSession() {
	
}