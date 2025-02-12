package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
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
}
