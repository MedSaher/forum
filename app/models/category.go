package models

import (
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
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		category := &Category{}
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
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
