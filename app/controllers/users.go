package controllers

import (
	// "html/template"
	"fmt"
	"net/http"
	// "forum/app/models"
)

// Create a home handler:
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Home Page!")
}


// Get users:
func UserHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "User Profile")
}