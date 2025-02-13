package models

import (
	"fmt"

	"forum/app/config"
)

// Category represents a category with a name and description
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetAllCategories retrieves all categories from the database
func GetAllCategories() ([]*Category, error) {
	db, Err := config.InitDB()
	if Err != nil {
		return nil, Err
	}
	defer db.Close()
	query := "SELECT ID, Name, Description FROM Category"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		category := &Category{}
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}
	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}
	return categories, nil
}

// Extract the category id:
func GetCategoryId(category string) (int, error) {
	db, err := config.InitDB()
	if err != nil {
		return -1, err
	}
	var id int
	query := "SELECT ID FROM Category where Name = ?"

	err = db.QueryRow(query, category).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// Link a new inserted post to its category:
func LinkPostToCategory(postId, categoryId int) error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO PostCategory (PostID, CategoryID)
	          VALUES (?, ?)`

	// Execute the query and get the result
	_, er := db.Exec(query, postId, categoryId)
	if er != nil {
		return er
	}
	return nil
}

// Get all categories in relation to a post:
func GetPostCategories(Post_id int) ([]*Category, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	query := `
	SELECT Category.ID, Category.Name, Category.Description 
	FROM Post 
	INNER JOIN PostCategory ON Post.ID = PostCategory.PostID 
	INNER JOIN Category ON Category.ID = PostCategory.CategoryID 
	WHERE Post.ID = ?;
	`
	rows, err := db.Query(query, Post_id)
	if err != nil {
		fmt.Println("happened")
		return nil, err
	}
	defer rows.Close() // Close rows
	categories := []*Category{}
	for rows.Next() {
		category := &Category{}
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// Filter post based on category:
func GetPostCategoriesId(category string) (map[int]bool, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	query := `
		SELECT p.ID FROM Post p
		JOIN PostCategory pc ON p.ID = pc.PostID
		JOIN Category c ON pc.CategoryID = c.ID
		WHERE c.Name = ?`
	rows, err := db.Query(query, category)
	if err != nil {
		return nil, err
	}
	posts := make(map[int]bool)
	defer rows.Close()
	for rows.Next() {
		var postId int
		if err := rows.Scan(&postId); err != nil {
			return nil, err
		}
		posts[postId] = true
	}
	return posts, nil
}
