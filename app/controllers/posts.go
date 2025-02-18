package controllers

import (
	"encoding/json"
	"html"
	"net/http"
	"strconv"
	"strings"

	"forum/app/models"
)

func GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read pagination parameters from query
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Convert them to integers (default: page = 1, limit = 10)
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Fetch paginated posts from the database
	posts, err := models.GetAllPosts(page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set headers and encode response
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// Test:
func PostsHandler(w http.ResponseWriter, r *http.Request) {
	if err := Tmpl.ExecuteTemplate(w, "posts.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Post represents the structure of a post in the system.
// AddPost function to create a new post and send the newly created post as a response.
func AddPost(wr http.ResponseWriter, rq *http.Request) {
	if rq.Method == http.MethodGet {
		if err := Tmpl.ExecuteTemplate(wr, "post_form.html", nil); err != nil {
			http.Error(wr, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// Get session token from cookie
	cookie, err := rq.Cookie("session_token")
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve form values
	title := strings.TrimSpace(html.EscapeString(rq.FormValue("post_title")))
	content := strings.TrimSpace(html.EscapeString(rq.FormValue("post_content")))
	categories := rq.Form["chosen_categories[]"] // Retrieves multiple values
	if title == "" || content == "" {
		http.Error(wr, "Pleasr fill in all the fields", http.StatusBadRequest)
		return
	}

	for _, cat := range categories {
		cat = strings.TrimSpace(cat)
		if cat == "" {
			http.Error(wr, "Pleasr fill in all the fields", http.StatusBadRequest)
			return
		}
	}

	// Get session details
	session, err := models.GetSessionByUUID(cookie.Value)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get user details
	user, err := models.GetUserByTocken(session.UUID)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the post in the database and fetch post details
	newPost, err := models.CreatePost(title, content, user.ID)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	// Link post to categories
	for _, category := range categories {
		categoryID, err := models.GetCategoryId(category)
		if err != nil {
			http.Error(wr, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := models.LinkPostToCategory(newPost.ID, categoryID); err != nil {
			http.Error(wr, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// Respond with new post details
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(newPost)
}

// A handler to get liked posts:
func GetLikedPosts(wr http.ResponseWriter, rq *http.Request) {
	userId, err := models.GetUserIdFromSession(rq)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}
	liked, err := models.GetLikedPosts(userId)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(liked)
}

// A handler to get owned posts:
func GetOwnedPosts(wr http.ResponseWriter, rq *http.Request) {
	userId, err := models.GetUserIdFromSession(rq)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}
	liked, err := models.GetOwnedPosts(userId)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(liked)
}
