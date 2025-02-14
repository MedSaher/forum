package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"forum/app/models"
)
// Create a structure to represent the post vote:
type VotePost struct {
	PostId int `json:"postId"`
	Value  int `json:"value"`
}

// Create a structute to vote for a comment:
type VoteComment struct {
	Comment_id int `json:"comment_id"`
	Value  int `json:"value"`
}

func VoteForPost(wr http.ResponseWriter, rq *http.Request) {
	if rq.Method != http.MethodPost {
		http.Error(wr, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	vote := &VotePost{}
	// Parse the JSON request body
	if err := json.NewDecoder(rq.Body).Decode(vote); err != nil {
		http.Error(wr, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println(vote)
	// validate input:
	if vote.Value != 1 && vote.Value != -1 {
		http.Error(wr, "Invalid vote value", http.StatusBadRequest)
		return
	}

	cookie, err := rq.Cookie("session_token")
	if err != nil {
		http.Error(wr, "No active session", http.StatusBadRequest)
		return
	}

	session, err := models.GetSessionByUUID(cookie.Value)
	if err != nil {
		http.Error(wr, "No active session", http.StatusInternalServerError)
		return
	}
	fmt.Println(session)
	err = models.VoteForPost(session.UserID, vote.PostId, vote.Value)
	if err != nil {
		fmt.Println("Error voting for the post")
		http.Error(wr, "No active session", http.StatusInternalServerError)
		return
	}

	// Get the vote count:
	err = models.UpdateVoteCount(vote.PostId)
	if err != nil {
		fmt.Println("Error updating post count", err)
		http.Error(wr, "No active session", http.StatusInternalServerError)
		return
	}
	// Log success and send a response back to the client.
	log.Println("User voted successfully.")
	response := map[string]string{
		"message": "User voted successfully.",
	}
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(response) // Return success response in JSON format.
}


// Create a handler to return the vote count:
func VoteForComment(wr http.ResponseWriter, rq *http.Request) {
	if rq.Method != http.MethodPost {
		http.Error(wr, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	vote := &VoteComment{}
	// Parse the JSON request body
	if err := json.NewDecoder(rq.Body).Decode(vote); err != nil {
		http.Error(wr, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println(vote)
	// validate input:
	if vote.Value != 1 && vote.Value != -1 {
		http.Error(wr, "Invalid vote value", http.StatusBadRequest)
		return
	}

	cookie, err := rq.Cookie("session_token")
	if err != nil {
		http.Error(wr, "No active session", http.StatusBadRequest)
		return
	}

	session, err := models.GetSessionByUUID(cookie.Value)
	if err != nil {
		http.Error(wr, "No active session", http.StatusInternalServerError)
		return
	}
	fmt.Println(session)
	err = models.VoteForComment(session.UserID, vote.Comment_id, vote.Value)
	if err != nil {
		fmt.Println("Error voting for the post")
		http.Error(wr, "No active session", http.StatusInternalServerError)
		return
	}

	// Get the vote count:
	err = models.UpdateCommentVoteCount(vote.Comment_id)
	if err != nil {
		fmt.Println("Error updating post count", err)
		http.Error(wr, "No active session", http.StatusInternalServerError)
		return
	}
	// Log success and send a response back to the client.
	log.Println("User voted successfully.")

	response := map[string]string{
		"message": "User voted successfully.",
	}
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(response) // Return success response in JSON format.
}