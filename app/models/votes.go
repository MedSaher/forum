package models

import (
	"database/sql"
	"fmt"
	"time"

	"forum/app/config"
)

// Create a structure to represent the vote(like/dislike):
type Vote struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    *int      `json:"post_id,omitempty"`
	CommentID *int      `json:"comment_id,omitempty"`
	Value     int       `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

// VoteForPost toggles a vote if it exists; otherwise, inserts a new vote.
func VoteForPost(userID int, postID *int, commentID *int, value int) error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	var existingValue int

	// Query to check if a vote exists and toggle it
	query := `SELECT Value FROM Vote
	WHERE UserID = ? AND (PostID = COALESCE(?, PostID) OR CommentID = COALESCE(?, CommentID))`
	err = db.QueryRow(query, userID, ptrIntValue(postID), ptrIntValue(commentID)).Scan(&existingValue)

	if err == sql.ErrNoRows {
		// Insert a new vote
		insertQuery := `INSERT INTO Vote (UserID, PostID, CommentID, Value) VALUES (?, ?, ?, ?)`
		_, err := db.Exec(insertQuery, userID, ptrIntValue(postID), ptrIntValue(commentID), value)
		if err != nil {
			return err
		}
	} else if err == nil {
		// Vote exists, toggle the value
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		updateQuery := `UPDATE Vote
		SET Value = ?, Timestamp = ? 
		WHERE UserID = ? AND (PostID = COALESCE(?, PostID) OR CommentID = COALESCE(?, CommentID))`
		_, err = db.Exec(updateQuery, value, timestamp, userID, ptrIntValue(postID), ptrIntValue(commentID))
		if err != nil {
			return err
		}
	}

	fmt.Println("Existing Value:", existingValue)
	return nil
}

// Helper function to handle nil pointers
func ptrIntValue(ptr *int) interface{} {
	if ptr == nil {
		return nil
	}
	return *ptr
}
