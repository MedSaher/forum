package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type voteData struct {
	PostId int `json:"postId"`
	Value  int `json:"value"`
}

func CreateVote(wr http.ResponseWriter, rq *http.Request) {
	if rq.Method != http.MethodPost {
		http.Error(wr, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	vote := &voteData{}
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
	
}
