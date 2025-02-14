package models

import (
	"forum/app/config"
)

// Declare a model to represent the Post and ease data exchange between backend and frontend:
type Post struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	AuthorID     int    `json:"authorId"`
	Timestamp    string `json:"time"`
	LikeCount    int    `json:"likeCount"`
	DislikeCount int    `json:"dislikeCount"`
}

type PostDTO struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	AuthorID        int    `json:"authorId"`
	Timestamp       string `json:"time"`
	LikeCount       int    `json:"likeCount"`
	DislikeCount    int    `json:"dislikeCount"`
	AuthorFirstName string `json:"authorFirstName"`
	AuthorLastName  string `json:"authorLastName"`
}

// CRUD (Create, Read, Update, Delete) operations between Go and SQLite3:
// ----->> Create a new Post:
func CreatePost(title, content string, authorId int) (int, error) {
	db, err := config.InitDB()
	if err != nil {
		return -1, err
	}
	defer db.Close()
	query := `INSERT INTO Post (Title, Content, AuthorId)
          VALUES (?, ?, ?)`
	result, Err := db.Exec(query, title, content, authorId)
	if Err != nil {
		return -1, err
	}

	// Get the last inserted ID
	lastID, er := result.LastInsertId()
	if er != nil {
		return -1, err
	}
	return int(lastID), nil
}

// Fetch all Posts:
func GetAllPosts() ([]*PostDTO, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// The query:
	query := `SELECT 
    Post.ID, 
    Title, 
    Content, 
    AuthorID, 
    Timestamp, 
    LikeCount, 
    DislikeCount, 
    User.FirstName,
    User.LastName 
FROM User 
INNER JOIN Post ON User.ID = Post.AuthorID 
`

	// Fetch Posts from the database
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Posts []*PostDTO
	for rows.Next() {
		post := &PostDTO{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.Timestamp, &post.LikeCount, &post.DislikeCount, &post.AuthorFirstName, &post.AuthorLastName); err != nil {
			return nil, err
		}
		Posts = append(Posts, post)
	}
	return Posts, nil
}

// // Get the liked posts from database:
func GetLikedPosts(userId int) (map[int]bool, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	liked := make(map[int]bool)
	query := `SELECT PostID FROM Vote WHERE UserID = ? AND Value = 1`
	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var postId int
		if err := rows.Scan(&postId); err != nil {
			return nil, err
		}
		liked[postId] = true
	}
	return liked, nil
}

// Get the owned posts from database:
func GetOwnedPosts(userId int) (map[int]bool, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	liked := make(map[int]bool)
	query := `SELECT ID FROM Post WHERE AuthorID = ?`
	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var postId int
		if err := rows.Scan(&postId); err != nil {
			return nil, err
		}
		liked[postId] = true
	}
	return liked, nil
}
