package models

import (
	"errors"
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

// CreateVote function
func VoteForPost(userID int, postID *int, commentID *int, value int) error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	query := "INSERT INTO Vote (UserID, PostID, CommentID, Value) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(query, userID, postID, commentID, value)
	if err != nil {
		return err
	}
	return nil
}

// Check if the user has alraedy voted for the:
func ChekUserVote(id int, vote_id int, table string) (bool, error) {
	db, err := config.InitDB()
	if err != nil {
		return false, err
	}
	query := ""
	var count int
	switch table {
	case "Comment":
		query = "SELECT COUNT(*) FROM Vote WHERE UserId = ? AND PostID = ?"
	case "Post":
		query = "SELECT COUNT(*) FROM Vote WHERE UserId = ? AND PostID = ?"
	default:
		return false, errors.New("wrong table name")
	}

	err = db.QueryRow(query, id, vote_id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update a vote if it is already updated:
func UpdatePost() error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}

	return nil
}
