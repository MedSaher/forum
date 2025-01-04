package controllers

import (
	"encoding/json"
	"net/http"

	"forum/app/models"
)

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
