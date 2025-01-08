package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"html/template"
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

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // Max upload size: 10 MB
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the file from form data
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file in the server's directory
	dir := "./uploads"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0o755)
	}

	dst, err := os.Create(filepath.Join(dir, header.Filename))
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the server
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s\n", header.Filename)
}
