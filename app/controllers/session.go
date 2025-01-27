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

// Extract the loged in user:
func LogedInUser(wr http.ResponseWriter, rq *http.Request) {
    cookie, err := rq.Cookie("session_token")
    if err != nil {
        // Handle error (e.g., cookie not found)
        fmt.Println("Error:", err)
        return
    }
	session, er := models.GetSessionByUUID(cookie.Value)
	if er != nil {
		fmt.Println("Error:", err)
        return
	}
    fmt.Println(session.UUID)
	uuid := session.UUID
	
}
