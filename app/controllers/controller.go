package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"forum/app/models"

	"golang.org/x/crypto/bcrypt"
)

var Tmpl *template.Template

// Create a home handler:
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")
	err := Tmpl.ExecuteTemplate(w, "user.html", "Home page!!!")
	if err != nil {
		log.Fatal(err)
	}
}

type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit upload size to 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		respondWithError(w, "Could not parse form data", http.StatusBadRequest)
		return
	}

	// Validate form fields
	firstName := strings.TrimSpace(r.FormValue("firstName"))
	lastName := strings.TrimSpace(r.FormValue("lastName"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	if firstName == "" || lastName == "" || email == "" || password == "" {
		respondWithError(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Check if email already exists
	exists, err := models.EmailExists(email)
	if err != nil {
		respondWithError(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		respondWithError(w, "Email already registered", http.StatusConflict)
		return
	}

	// Handle profile picture
	file, header, err := r.FormFile("profilePicture")
	if err != nil {
		respondWithError(w, "Profile picture is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	if !isValidImageType(header.Filename) {
		respondWithError(w, "Invalid file type. Only jpg, jpeg, png allowed", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	filename := generateUniqueFilename(header.Filename)

	// Save file
	imagePath, err := saveImage(file, filename)
	if err != nil {
		respondWithError(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Create user
	user := &models.User{
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		PasswordHash:   string(hashedPassword),
		ProfilePicture: imagePath,
	}

	// Save user to database
	err = models.CreateUser(user)
	if err != nil {
		// Clean up uploaded file if database operation fails
		os.Remove(imagePath)
		respondWithError(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, Response{
		Message: "User registered successfully",
	}, http.StatusCreated)
}

func isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	return validExtensions[ext]
}

func generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
}

func saveImage(file io.Reader, filename string) (string, error) {
	uploadDir := "./app/assets/images"
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", err
	}

	filepath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return filepath, nil
}

func GetUserProfilePicture(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, user.ProfilePicture)
}

func respondWithError(w http.ResponseWriter, message string, code int) {
	respondWithJSON(w, Response{Error: message}, code)
}

func respondWithJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
