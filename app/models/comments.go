package models

import "forum/app/config"

// comment model: [one to many with Comment] and [one to many with user]:
type Comment struct {
	ID           int    `json:"id"`
	Content      string `json:"content"`
	AuthorID     int    `json:"authorId"`
	PostID       int    `json:"postId"`
	CommentID    int    `json:"CommentId"`
	Timestamp    int    `json:"timeStamp"`
	LikeCount    int    `json:"likeCount"`
	DislikeCount int    `json:"dislikeCount"`
}

// CRUD (Create, Read, Update, Delete) operations between Go and SQLite3:
// ----->> Create a new Comment:
func CreateComment(title, content string, authorId int, postId int) error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()
	query := `INSERT INTO Comment (Title, Content, AuthorId, PostID)
          VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, title, content, authorId)
	if err != nil {
		return err
	}
	return nil
}

// Fetch all Comments
func GetAllComments() ([]*Comment, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// Fetch Comments from the database
	rows, err := db.Query("SELECT * FROM Comment")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Comments []*Comment
	for rows.Next() {
		comment := &Comment{}
		if err := rows.Scan(&comment.ID, comment.Content, comment.AuthorID, comment.PostID, comment.LikeCount, comment.DislikeCount); err != nil {
			return nil, err
		}
		Comments = append(Comments, comment)
	}
	return Comments, nil
}
