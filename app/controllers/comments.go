package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum/app/models"
)

// Create a structure to get comments data:
type Comment struct {
	PostID  int    `json:"postId"`
	Content string `json:"content"`
}

// Create a handler to create a comment:
func CreateComment(wr http.ResponseWriter, rq *http.Request) {
	userId, err := models.GetUserIdFromSession(rq)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(userId)
	var comment Comment
	decoder := json.NewDecoder(rq.Body)
	if err := decoder.Decode(&comment); err != nil {
		http.Error(wr, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Trim the content to prevent empty spaces from being considered as valid input
	comment.Content = strings.TrimSpace(comment.Content)
	if comment.Content == "" || comment.PostID <= 0 {
		http.Error(wr, "Comment or Post ID cannot be empty", http.StatusBadRequest)
		return
	}
	fmt.Println(comment)
	// Create a comment:
	err = models.CreateComment(comment.Content, userId, comment.PostID)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusBadRequest)
		return
	}
}

// Create a handler function to get all comments in relation to a post:
func GetComments(wr http.ResponseWriter, rq *http.Request) {
	// Retrieve post_id from the query parameters
	post_id_str := rq.URL.Query().Get("post_id")
	if post_id_str == "" {
		http.Error(wr, "post_id is required", http.StatusBadRequest)
		return
	}
	// Convert post_id to an integer
	post_id, err := strconv.Atoi(post_id_str)
	if err != nil || post_id <= 0 {
		http.Error(wr, "valid post_id number is required", http.StatusBadRequest)
		return
	}

	// Print post_id for debugging (you may remove this in production)
	fmt.Println("Post ID:", post_id)
	comments, err := models.GetAllComments(post_id)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(comments)
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(comments)
}
