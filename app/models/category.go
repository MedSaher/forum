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
