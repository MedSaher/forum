package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"forum/app/middleware"
	"forum/app/models"
	"forum/app/utils"
)

// Declare a golobale template variable:
var Tmpl *template.Template

// Create a loging credentials structure:
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

// Register a new user to my app:
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request is a GET request. If so, render the form for user registration.
	if r.Method == http.MethodGet {
		Tmpl.ExecuteTemplate(w, "user.html", nil)
		return
	}

	// Parse the multipart form with a maximum memory of 10MB for uploaded files.
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Failed to parse form: %v", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Extract form values for user details.
	firstName := strings.TrimSpace(r.FormValue("firstName")) // Remove leading/trailing spaces
	lastName := strings.TrimSpace(r.FormValue("lastName"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password") // No need to trim spaces from passwords.

	// Validate required fields to ensure none are empty.
	if firstName == "" || lastName == "" || email == "" || password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Check if the email already exists in the database.
	exists, err := models.CheckEmailExists(email)
	if err != nil {
		log.Printf("Error checking email: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Retrieve the uploaded profile picture file.
	file, handler, err := r.FormFile("profilePicture")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close() // Ensure the file is closed after processing.

	// Validate the file type to allow only JPG, JPEG, or PNG files.
	if !isAllowedFileType(handler.Filename) {
		http.Error(w, "Invalid file type. Only jpg, jpeg, png allowed", http.StatusBadRequest)
		return
	}

	// Define the directory where uploaded files will be saved.
	uploadDir := "./app/uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		// Create the uploads directory if it doesn't exist.
		if err := os.Mkdir(uploadDir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create upload directory: %v", err)
		}
	}

	// Define the full path where the file will be saved.
	filePath := filepath.Join(uploadDir, handler.Filename)
	destFile, err := os.Create(filePath) // Create the file on the server.
	if err != nil {
		log.Printf("Error creating file: %v", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer destFile.Close() // Ensure the file is closed after writing.

	// Save the file content to the newly created file.
	_, err = io.Copy(destFile, file)
	if err != nil {
		log.Printf("Error saving file: %v", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Hash the user's password for secure storage in the database.
	hashedPassword, err := utils.GenerateCryptoPassword(password)
	if err != nil {
		// Clean up the uploaded file if password hashing fails.
		os.Remove(filePath)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Create a new user object with all the required fields.
	user := &models.User{
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		PasswordHash:   string(hashedPassword), // Save the hashed password.
		ProfilePicture: handler.Filename,       // Store the file's relative path.
	}

	// Log the user data for debugging purposes.
	log.Printf("Attempting to save user: %+v", user)

	// Save the user to the database.
	if err := models.CreateUser(user); err != nil {
		log.Printf("Error saving user to database: %v", err)
		// Remove the file if the database operation fails.
		os.Remove(filePath)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Log success and send a response back to the client.
	log.Println("User successfully registered.")
	response := map[string]string{
		"message":   "User registered successfully",
		"image_url": user.ProfilePicture,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response) // Return success response in JSON format.
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

// Create a login handler:
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	credentials := &Credentials{}
	fmt.Println(credentials)
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println(credentials)
	// Fetch user by email
	user, err := models.GetUserByEmail(credentials.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	fmt.Println(user)
	// Validate password
	if !utils.ValidatePassword(user.PasswordHash, credentials.Password) {
		http.Error(w, "error password", http.StatusUnauthorized)
		return
	}

	// Instantiate anew session:
	session, err := middleware.CreateSession(user.ID, 4*time.Hour)
	if err != nil {
		http.Error(w, "Failure during session creation", http.StatusInternalServerError)
		return
	}

	// Set session coockies:
	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.UUID,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   true, // Use https
	})
	

	log.Println("User successfully loged in.")
	response := map[string]string{
		"message":   "User loged in successfully",
		"user_name": user.FirstName,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response) // Return success response in JSON format.
}
