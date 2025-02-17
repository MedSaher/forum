package controllers

import (
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

func Logout(w http.ResponseWriter, r *http.Request) {
    // Get session token from cookie
    cookie, err := r.Cookie("session_token")
    if err != nil {
        http.Error(w, "No active session", http.StatusBadRequest)
        return
    }

    // Delete session from the database
    err = models.DeleteSessionByUUID(cookie.Value)
    if err != nil {
        http.Error(w, "Failed to logout", http.StatusInternalServerError)
        return
    }

    // Expire the session cookie
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
