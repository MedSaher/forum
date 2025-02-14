package models

import (
	"database/sql"
	"fmt"
	"time"

	"forum/app/config"
)

// Create a structure to represent the vote(like/dislike):
type Vote struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    *int      `json:"post_id,omitempty"`
	CommentID *int      `json:"comment_id,omitempty"`
	Value     int       `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}


// Vote for post or comments:
func VoteForPostAndComment(userID int, postID *int, commentID *int, value int) error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	var existingValue int
	var query string
	var args []interface{}

	if postID != nil {
		query = `SELECT Value FROM Vote WHERE UserID = ? AND PostID = ?`
		args = append(args, userID, *postID)
	} else if commentID != nil {
		query = `SELECT Value FROM Vote WHERE UserID = ? AND CommentID = ?`
		args = append(args, userID, *commentID)
	} else {
		return fmt.Errorf("either postID or commentID must be provided")
	}

	err = db.QueryRow(query, args...).Scan(&existingValue)

	if err == nil {
		// If the same vote exists, return an error
		if existingValue == value {
			return fmt.Errorf("duplicate vote: user already voted the same way")
		}

		// Update existing vote
		var updateQuery string
		if postID != nil {
			updateQuery = `UPDATE Vote SET Value = ?, Timestamp = ? WHERE UserID = ? AND PostID = ?`
			args = []interface{}{value, time.Now().Format("2006-01-02 15:04:05"), userID, *postID}
		} else {
			updateQuery = `UPDATE Vote SET Value = ?, Timestamp = ? WHERE UserID = ? AND CommentID = ?`
			args = []interface{}{value, time.Now().Format("2006-01-02 15:04:05"), userID, *commentID}
		}
		_, err = db.Exec(updateQuery, args...)
		return err
	} else if err == sql.ErrNoRows {
		// Insert new vote
		insertQuery := `INSERT INTO Vote (UserID, PostID, CommentID, Value) VALUES (?, ?, ?, ?)`
		_, err := db.Exec(insertQuery, userID, ptrIntValue(postID), ptrIntValue(commentID), value)
		return err
	}

	return err // Return other unexpected errors
}


// Helper function to handle nil pointers
func ptrIntValue(ptr *int) interface{} {
	if ptr == nil {
		return nil
	}
	return *ptr
}

// Vote for post:
func VoteForPost(userID int, postID int, value int) error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	var existingValue int
	query := `SELECT Value FROM Vote WHERE UserID = ? AND PostID = ?`
	err = db.QueryRow(query, userID, postID).Scan(&existingValue)

	if err == nil {
		if existingValue == value {
			return fmt.Errorf("duplicate vote: user already voted the same way")
		}

		// Update the vote
		updateQuery := `UPDATE Vote SET Value = ?, Timestamp = ? WHERE UserID = ? AND PostID = ?`
		_, err = db.Exec(updateQuery, value, time.Now().Format("2006-01-02 15:04:05"), userID, postID)
		return err
	} else if err == sql.ErrNoRows {
		// Insert new vote
		insertQuery := `INSERT INTO Vote (UserID, PostID, Value) VALUES (?, ?, ?)`
		_, err := db.Exec(insertQuery, userID, postID, value)
		return err
	}

	return err
}


// Vot for comment:
func VoteForComment(userID int, commentID int, value int) error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	var existingValue int
	query := `SELECT Value FROM Vote WHERE UserID = ? AND CommentID = ?`
	err = db.QueryRow(query, userID, commentID).Scan(&existingValue)

	if err == nil {
		if existingValue == value {
			return fmt.Errorf("duplicate vote: user already voted the same way")
		}

		// Update the vote
		updateQuery := `UPDATE Vote SET Value = ?, Timestamp = ? WHERE UserID = ? AND CommentID = ?`
		_, err = db.Exec(updateQuery, value, time.Now().Format("2006-01-02 15:04:05"), userID, commentID)
		return err
	} else if err == sql.ErrNoRows {
		// Insert new vote
		insertQuery := `INSERT INTO Vote (UserID, CommentID, Value) VALUES (?, ?, ?)`
		_, err := db.Exec(insertQuery, userID, commentID, value)
		return err
	}

	return err
}


// UpdateVoteCount updates the like and dislike counts for a given post in the POSTS table
func UpdateVoteCount(postID int) error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	// Update query to set like and dislike counts
	updateQuery := `
		UPDATE Post
		SET LikeCount = (SELECT COUNT(*) FROM VOTE WHERE PostID = ? AND Value = 1),
		    DislikeCount = (SELECT COUNT(*) FROM VOTE WHERE PostID = ? AND Value = -1)
		WHERE ID = ?;
	`
	// Execute the updated query
	_, err = db.Exec(updateQuery, postID, postID, postID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateCommentVoteCount updates the like and dislike counts for a given comment
func UpdateCommentVoteCount(commentID int) error {
	
fmt.Println("------------>", commentID)
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	query := `
		UPDATE Comment
		SET LikeCount = (SELECT COUNT(*) FROM Vote WHERE CommentID = ? AND Value = 1),
		    DislikeCount = (SELECT COUNT(*) FROM Vote WHERE CommentID = ? AND Value = -1)
		WHERE ID = ?;
	`

	_, err = db.Exec(query, commentID, commentID, commentID)
	if err != nil {
		return fmt.Errorf("failed to update comment vote count: %w", err)
	}
	return nil
}
