package models

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

// CRUD (Create, Read, Update, Delete) operations between Go and SQLite3:
// ----->> Create a new Post:
func CreatePost(title, content string, authorId string) error {
	db, err := Connection()
	if err != nil {
		return err
	}
	defer db.Close()
	query := `INSERT INTO Post (Title, Content, AuthorId)
          VALUES (?, ?, ?)`
	_, err = db.Exec(query, title, content, authorId)
	if err != nil {
		return err
	}
	return nil
}

// Fetch all Posts
func GetAllPosts() ([]*Post, error) {
	db, err := Connection()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// Fetch Posts from the database
	rows, err := db.Query("SELECT * FROM Post")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Posts []*Post
	for rows.Next() {
		post := &Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.LikeCount, &post.DislikeCount); err != nil {
			return nil, err
		}
		Posts = append(Posts, post)
	}
	return Posts, nil
}
