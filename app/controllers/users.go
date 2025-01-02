package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"forum/app/models"
)

var Tmpl *template.Template

// Create a home handler:
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")
	err := Tmpl.ExecuteTemplate(w, "home.html", "Home page!!!")
	if err != nil {
		log.Fatal(err)
	}
}

// Get all users:
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// Fetch books from the database
	users, err := models.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
