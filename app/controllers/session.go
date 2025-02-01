package controllers

import (
	"fmt"
	"net/http"
	"time"

	"forum/app/models"
	"forum/app/utils"
)

// Create a new session
func CreateSession(userId int, expiration time.Duration) (*models.Session, error) {
	uuid := utils.GenerateUUID()
	currentTime := time.Now()
	expiresAt := currentTime.Add(expiration)
	session := &models.Session{
		UserID:    userId,
		UUID:      uuid,
		ExpiresAt: expiresAt,
		CreatedAt: currentTime,
	}
	// Save the session to database:
	err := models.InsertSession(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// Logout and kill the session:
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
        http.Error(w, "No active session", http.StatusBadRequest)
        return
    }
	fmt.Println(cookie)
    // Expire the cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    "",
        MaxAge:   -1,
        HttpOnly: true,
        Path:     "/",
    })

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Logged out successfully"}`))
}

