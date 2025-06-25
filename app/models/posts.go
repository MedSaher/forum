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
	CategoryName    string `json:"categoryName"`
	AuthorFirstName string `json:"authorFirstName"`
	AuthorLastName  string `json:"authorLastName"`
}

// CRUD (Create, Read, Update, Delete) operations between Go and SQLite3:
// ----->> Create a new Post:
// CRUD (Create, Read, Update, Delete) operations between Go and SQLite3:
// ----->> Create a new Post:
func CreatePost(title, content string, authorId int) (*PostDTO, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Insert the new post into the Post table
	query := `INSERT INTO Post (Title, Content, AuthorID)
          VALUES (?, ?, ?)`
	result, err := db.Exec(query, title, content, authorId)
	if err != nil {
		return nil, err
	}

	// Get the last inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Fetch the newly inserted post details
	query = `SELECT Post.ID, Title, Content, AuthorID, Timestamp, LikeCount, DislikeCount, 
                    User.FirstName, User.LastName 
             FROM Post
             INNER JOIN User ON User.ID = Post.AuthorID
             WHERE Post.ID = ?`

	var post PostDTO
	err = db.QueryRow(query, lastID).Scan(
		&post.ID, &post.Title, &post.Content, &post.AuthorID,
		&post.Timestamp, &post.LikeCount, &post.DislikeCount,
		&post.AuthorFirstName, &post.AuthorLastName)
	if err != nil {
		return nil, err
	}

	// Return the full post details including author and timestamp
	return &post, nil
}

func GetAllPosts(page, limit int) ([]*PostDTO, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	offset := (page - 1) * limit
	query := `SELECT 
	    Post.ID, Title, Content, AuthorID, Timestamp, LikeCount, DislikeCount, 
	    GROUP_CONCAT(DISTINCT Category.Name) AS Categories, 
	    User.FirstName, User.LastName 
	FROM Post 
	INNER JOIN User ON User.ID = Post.AuthorID 
	INNER JOIN PostCategory ON Post.ID = PostCategory.PostID 
	INNER JOIN Category ON PostCategory.CategoryID = Category.ID 
	GROUP BY Post.ID, Title, Content, AuthorID, Timestamp, LikeCount, DislikeCount, User.FirstName, User.LastName 
	ORDER BY Post.ID DESC 
	LIMIT ? OFFSET ?;`

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*PostDTO
	for rows.Next() {
		post := &PostDTO{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.Timestamp, &post.LikeCount, &post.DislikeCount, &post.CategoryName, &post.AuthorFirstName, &post.AuthorLastName); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// // Get the liked posts from database:
func GetLikedPosts(userId int) (map[int]bool, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	liked := make(map[int]bool)
	query := `SELECT PostID FROM Vote WHERE UserID = ? AND Value = 1 AND PostID IS NOT NULL;`
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
