package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"forum/app/models"
)

type PostCategoryId struct {
	Post_id int `json:"post_id"`
}

// A controller the get all the categories:
func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := models.GetAllCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with json:
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// Get categories of a specific post:
func GetPostCategories(w http.ResponseWriter, r *http.Request) {
	var post PostCategoryId

	// Read and log the request body for debugging
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Always close the request body

	fmt.Println("Received request body:", string(body))

	// Decode JSON into struct
	if err := json.Unmarshal(body, &post); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	fmt.Println("The post ID I want categories for:", post.Post_id)

	// Fetch categories from the database
	categories, err := models.GetPostCategories(post.Post_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Fetched categories:", categories)

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// Fetch posts based on categories:
func GetPostsByCategory(w http.ResponseWriter, r *http.Request) {
	// Retrieve category name from the query parameters
	category := r.URL.Query().Get("category")
	if category == "" {
		http.Error(w, "Category is required", http.StatusBadRequest)
		return
	}
	fmt.Println("the category is category", category)
	posts, err := models.GetPostCategoriesId(category)
	if err != nil {
		http.Error(w, "Category is required", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}
