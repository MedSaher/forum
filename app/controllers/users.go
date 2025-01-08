package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"forum/app/models"

	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := models.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form with 10MB limit
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get and validate form fields
	firstName := strings.TrimSpace(r.FormValue("firstName"))
	lastName := strings.TrimSpace(r.FormValue("lastName"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	// Validate required fields
	if firstName == "" || lastName == "" || email == "" || password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Check if email already exists
	exists, err := models.CheckEmailExists(email)
	if err != nil {
		http.Error(w, "Error happens here"+ err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Handle file upload
	file, header, err := r.FormFile("profilePicture")
	if err != nil {
		http.Error(w, "Profile picture is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	if !isAllowedFileType(header.Filename) {
		http.Error(w, "Invalid file type. Only jpg, jpeg, png allowed", http.StatusBadRequest)
		return
	}

	// Create unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
	uploadPath := filepath.Join("app", "assets", "images")

	// Ensure upload directory exists
	if err := os.MkdirAll(uploadPath, 0o755); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Create destination file
	filePath := filepath.Join(uploadPath, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filePath) // Clean up on error
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Create user model
	user := &models.User{
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		PasswordHash:   string(hashedPassword),
		ProfilePicture: filename,
	}

	// Save to database
	if err := models.CreateUser(user); err != nil {
		os.Remove(filePath) // Clean up on error
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
    fmt.Println(user)
	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message: fmt.Sprintf("User %s %s registered successfully!", firstName, lastName),
	})
}

func isAllowedFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	return validTypes[ext]
}
