package controllers

import (
	"encoding/json"
	"net/http"

	"forum/app/models"
)

// Get all users:
func GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// Fetch books from the database
	users, err := models.GetAllPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Test:
func PostsHandler(w http.ResponseWriter, r *http.Request) {
	Tmpl.ExecuteTemplate(w, "posts.html", nil)
}
