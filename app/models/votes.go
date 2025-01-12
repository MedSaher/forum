package models

import (
	"fmt"
	"forum/app/config"
)

// Vote represents the user vote (like or dislike)
type Vote struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId"`
	PostID    int    `json:"postId"`
	CommentID int    `json:"commentId"`
	Value     int    `json:"value"`
	Timestamp string `json:"timestamp"`
}

// CreateVote inserts a new vote into the database
func CreateVote(userID, postID, commentID, value int) error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()
	query := `INSERT INTO Vote (UserID, PostID, CommentID, Value) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, userID, postID, commentID, value)
	if err != nil {
		return fmt.Errorf("failed to insert vote: %w", err)
	}
	return nil
}

// GetAllVotes retrieves all votes from the database
func GetAllVotes() ([]*Vote, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	query := "SELECT ID, UserID, PostID, CommentID, Value, Timestamp FROM Vote"
	rows, Err := db.Query(query)
	if Err != nil {
		return nil, fmt.Errorf("failed to query votes: %w", err)
	}
	defer rows.Close()

	var votes []*Vote
	for rows.Next() {
		vote := &Vote{}
		if err := rows.Scan(&vote.ID, &vote.UserID, &vote.PostID, &vote.CommentID, &vote.Value, &vote.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan vote: %w", err)
		}
		votes = append(votes, vote)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return votes, nil
}
