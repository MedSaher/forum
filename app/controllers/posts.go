package controllers

import (
	"encoding/json"
	"fmt"
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
	posts, err := models.GetAllPosts()
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

// Add a newpost:
func AddPost(wr http.ResponseWriter, rq *http.Request) {
	if rq.Method == http.MethodGet {
		if err := Tmpl.ExecuteTemplate(wr, "post_form.html", nil); err != nil {
			http.Error(wr, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	cookie, err := rq.Cookie("session_token")
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}
	title := rq.FormValue("post_title")
	content := rq.FormValue("post_content")
	categories := rq.Form["chosen_categories[]"] // Retrieves multiple values

	fmt.Println(categories)
	session, er := models.GetSessionByUUID(cookie.Value)
	if er != nil {
		fmt.Println("Error:", err)
		return
	}
	uuid := session.UUID
	user, Err := models.GetUserByTocken(uuid)
	if Err != nil {
		fmt.Println("Error:", Err)
		return
	}
	fmt.Println(user)
	// Create the post in the database:
	post_id, err := models.CreatePost(title, content, user.ID)
	if err != nil {
		fmt.Println("Error:", er)
		return
	}

	for _, category := range categories {
		category_id, er := models.GetCategoryId(category)
		if er != nil {
			fmt.Println("Error:", er)
			return
		}
		// link the new inserted post to its category in database:
		if err := models.LinkPostToCategory(post_id, category_id); err != nil {
			fmt.Println("Error:", er)
			return
		}
	}
	// Respond with a success message
	response := map[string]string{
		"message": "User logged in successfully",
	}
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(response)
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
	fmt.Printf("user id : %v\n",  userId)
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

