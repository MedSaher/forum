package models

import (
	"time"

	"forum/app/config"
)

type Comment struct {
	ID           int       `json:"id"`
	Content      string    `json:"content"`
	AuthorID     int       `json:"authorId"`
	PostID       int       `json:"postId"`
	CommentID    int       `json:"commentId"`
	Timestamp    time.Time `json:"timestamp"`
	LikeCount    int       `json:"likeCount"`
	DislikeCount int       `json:"dislikeCount"`
}

// Create a comment response structure
type CommentDTO struct {
	ID           int    `json:"id"`
	Content      string `json:"content"`
	Timestamp    string `json:"timestamp"`
	LikeCount    int    `json:"likeCount"`
	DislikeCount int    `json:"dislikeCount"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}
// Create a new comment:
func CreateComment(content string, authorId, postId int) (*CommentDTO, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `INSERT INTO Comment (Content, AuthorID, PostID, Timestamp, LikeCount, DislikeCount) VALUES (?, ?, ?, ?, 0, 0)`
	result, err := db.Exec(query, content, authorId, postId, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var comment CommentDTO
	queryLastInsert := `
		SELECT c.ID, c.Content, c.Timestamp, c.LikeCount, c.DislikeCount, u.FirstName, u.LastName 
		FROM Comment c 
		JOIN User u ON c.AuthorID = u.ID 
		WHERE c.ID = ?
	`

	// Use QueryRow instead of Exec, and Scan properly
	err = db.QueryRow(queryLastInsert, lastId).Scan(
		&comment.ID, &comment.Content, &comment.Timestamp, 
		&comment.LikeCount, &comment.DislikeCount, 
		&comment.FirstName, &comment.LastName,
	)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}


func GetAllComments(postID int) ([]*CommentDTO, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT c.ID, c.Content, c.Timestamp, c.LikeCount, c.DislikeCount, u.FirstName, u.LastName 
		FROM Comment c 
		JOIN User u ON c.AuthorID = u.ID 
		WHERE c.PostID = ?
		ORDER BY c.Timestamp DESC
	`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*CommentDTO
	for rows.Next() {
		comment := &CommentDTO{}
		var timestamp time.Time
		if err := rows.Scan(&comment.ID, &comment.Content, &timestamp, &comment.LikeCount, &comment.DislikeCount, &comment.FirstName, &comment.LastName); err != nil {
			return nil, err
		}
		comment.Timestamp = timestamp.Format(time.RFC3339)
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
