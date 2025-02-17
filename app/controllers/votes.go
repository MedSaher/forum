// Backend Go Controllers (controllers/vote.go)
package controllers

import (
    "encoding/json"
    "net/http"
    "forum/app/models"
)

type VotePost struct {
    PostId int `json:"postId"`
    Value  int `json:"value"`
}

type VoteComment struct {
    CommentId int `json:"comment_id"`
    Value     int `json:"value"`
}

func VoteForPost(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var vote VotePost
    if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if vote.Value != 1 && vote.Value != -1 {
        http.Error(w, "Invalid vote value", http.StatusBadRequest)
        return
    }

    cookie, err := r.Cookie("session_token")
    if err != nil {
        http.Error(w, "No active session", http.StatusUnauthorized)
        return
    }

    session, err := models.GetSessionByUUID(cookie.Value)
    if err != nil {
        http.Error(w, "Invalid session", http.StatusUnauthorized)
        return
    }

    counts, err := models.VoteForPost(session.UserID, vote.PostId, vote.Value)
    if err != nil {
        http.Error(w, "Failed to process vote", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(counts)
}

func VoteForComment(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var vote VoteComment
    if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if vote.Value != 1 && vote.Value != -1 {
        http.Error(w, "Invalid vote value", http.StatusBadRequest)
        return
    }

    cookie, err := r.Cookie("session_token")
    if err != nil {
        http.Error(w, "No active session", http.StatusUnauthorized)
        return
    }

    session, err := models.GetSessionByUUID(cookie.Value)
    if err != nil {
        http.Error(w, "Invalid session", http.StatusUnauthorized)
        return
    }

    counts, err := models.VoteForComment(session.UserID, vote.CommentId, vote.Value)
    if err != nil {
        http.Error(w, "Failed to process vote", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(counts)
}