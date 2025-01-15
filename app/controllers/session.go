package controllers

import (
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
