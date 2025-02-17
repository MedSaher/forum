// Backend Go Models (models/vote.go)
package models

import (
	"database/sql"
	"fmt"
	"time"

	"forum/app/config"
)

type Vote struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    *int      `json:"post_id,omitempty"`
	CommentID *int      `json:"comment_id,omitempty"`
	Value     int       `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type VoteCounts struct {
	LikeCount    int `json:"likeCount"`
	DislikeCount int `json:"dislikeCount"`
}

func VoteForPost(userID int, postID int, value int) (*VoteCounts, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Check if the user has already voted
	var existingValue int
	query := `SELECT Value FROM Vote WHERE UserID = ? AND PostID = ?`
	err = db.QueryRow(query, userID, postID).Scan(&existingValue)

	if err == nil {
		// If the same vote exists, return current counts
		if existingValue == value {
			return GetPostVoteCounts(postID)
		}

		// Update existing vote
		updateQuery := `UPDATE Vote SET Value = ?, Timestamp = ? WHERE UserID = ? AND PostID = ?`
		_, err = db.Exec(updateQuery, value, time.Now(), userID, postID)
		if err != nil {
			return nil, fmt.Errorf("failed to update vote: %w", err)
		}
	} else if err == sql.ErrNoRows {
		// Insert new vote
		insertQuery := `INSERT INTO Vote (UserID, PostID, Value, Timestamp) VALUES (?, ?, ?, ?)`
		_, err := db.Exec(insertQuery, userID, postID, value, time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to insert vote: %w", err)
		}
	} else {
		return nil, fmt.Errorf("failed to check existing vote: %w", err)
	}

	// Update and return the new vote counts
	return UpdateAndGetPostVoteCounts(postID)
}

func VoteForComment(userID int, commentID int, value int) (*VoteCounts, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Check if the user has already voted
	var existingValue int
	query := `SELECT Value FROM Vote WHERE UserID = ? AND CommentID = ?`
	err = db.QueryRow(query, userID, commentID).Scan(&existingValue)

	if err == nil {
		// If the same vote exists, return current counts
		if existingValue == value {
			return GetCommentVoteCounts(commentID)
		}

		// Update existing vote
		updateQuery := `UPDATE Vote SET Value = ?, Timestamp = ? WHERE UserID = ? AND CommentID = ?`
		_, err = db.Exec(updateQuery, value, time.Now(), userID, commentID)
		if err != nil {
			return nil, fmt.Errorf("failed to update vote: %w", err)
		}
	} else if err == sql.ErrNoRows {
		// Insert new vote
		insertQuery := `INSERT INTO Vote (UserID, CommentID, Value, Timestamp) VALUES (?, ?, ?, ?)`
		_, err := db.Exec(insertQuery, userID, commentID, value, time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to insert vote: %w", err)
		}
	} else {
		return nil, fmt.Errorf("failed to check existing vote: %w", err)
	}

	// Update and return the new vote counts
	return UpdateAndGetCommentVoteCounts(commentID)
}

func UpdateAndGetPostVoteCounts(postID int) (*VoteCounts, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Update the counts
	updateQuery := `
        UPDATE Post
        SET LikeCount = (SELECT COUNT(*) FROM Vote WHERE PostID = ? AND Value = 1),
            DislikeCount = (SELECT COUNT(*) FROM Vote WHERE PostID = ? AND Value = -1)
        WHERE ID = ?
    `
	_, err = db.Exec(updateQuery, postID, postID, postID)
	if err != nil {
		return nil, err
	}

	// Get the updated counts
	return GetPostVoteCounts(postID)
}

func UpdateAndGetCommentVoteCounts(commentID int) (*VoteCounts, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Update the counts
	updateQuery := `
        UPDATE Comment
        SET LikeCount = (SELECT COUNT(*) FROM Vote WHERE CommentID = ? AND Value = 1),
            DislikeCount = (SELECT COUNT(*) FROM Vote WHERE CommentID = ? AND Value = -1)
        WHERE ID = ?
    `
	_, err = db.Exec(updateQuery, commentID, commentID, commentID)
	if err != nil {
		return nil, err
	}

	// Get the updated counts
	return GetCommentVoteCounts(commentID)
}

func GetPostVoteCounts(postID int) (*VoteCounts, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	counts := &VoteCounts{}
	query := `SELECT LikeCount, DislikeCount FROM Post WHERE ID = ?`
	err = db.QueryRow(query, postID).Scan(&counts.LikeCount, &counts.DislikeCount)
	if err != nil {
		return nil, err
	}

	return counts, nil
}

func GetCommentVoteCounts(commentID int) (*VoteCounts, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	counts := &VoteCounts{}
	query := `SELECT LikeCount, DislikeCount FROM Comment WHERE ID = ?`
	err = db.QueryRow(query, commentID).Scan(&counts.LikeCount, &counts.DislikeCount)
	if err != nil {
		return nil, err
	}

	return counts, nil
}
