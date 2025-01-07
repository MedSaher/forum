package models

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrUserNotFound = errors.New("user not found")

// Declare a model to represent the user and ease data exchange between backend and frontend:
type User struct {
	ID             int    `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	PasswordHash   string `json:"password"`
	ProfilePicture string `json:"profilePicture"`
}

// CRUD (Create, Read, Update, Delete) operations between Go and SQLite3:
// ----->> Create a new user:
func CreateUser(user *User) error {
	db, err := Connection() // Assuming Connection() returns *sql.DB
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO User (FirstName, LastName, Email, PasswordHash, ProfilePicture)
	          VALUES (?, ?, ?, ?, ?)`

	// Execute the query and get the result
	result, err := db.Exec(query, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.ProfilePicture)
	if err != nil {
		return err
	}

	// Retrieve the auto-generated ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id) // Update the user struct with the auto-generated ID
	fmt.Println(user)
	return nil
}

// Fetch all Users
func GetAllUsers() ([]*User, error) {
	db, err := Connection()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// Fetch Users from the database
	rows, err := db.Query("SELECT * FROM User")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.ProfilePicture); err != nil {
			return nil, err
		}
		Users = append(Users, user)
	}
	return Users, nil
}

func GetUserByID(id string) (*User, error) {
	db, err := Connection()
	if err != nil {
		return nil, err
	}
	var user User

	query := `
		SELECT id, first_name, last_name, email, password_hash, profile_picture 
		FROM users 
		WHERE id = ?`

	err = db.QueryRow(query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.ProfilePicture,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Helper function to check if a user exists by ID
func UserExists(id string) (bool, error) {
	db, err := Connection()
	if err != nil {
		return false, err
	}
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)"

	err = db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// For handling errors in your controllers, you can use this helper:
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrUserNotFound)
}
