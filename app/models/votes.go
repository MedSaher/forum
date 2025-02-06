package models

import (
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
func CreateVote(userID int, postID *int, commentID *int, value int) (*Vote, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	
}
