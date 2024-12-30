package models

import (
	
)

// Declare a model to represent the user and ease data exchange between backend and frontend:
type User struct {
	ID             int    `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	UserName       string `json:"userName"`
	Email          string `json:"email"`
	PasswordHash       string `json:"password"`
	ProfilePicture string `json:"profilePicture"`
	Role           string `json:"role"`
}

// CRUD (Create, Read, Update, Delete) operations between Go and SQLite3:
// ----->> Create a new user:
func CreateUser(firstName, lastName, email, password, profilePicture, role string) error {
	db, err := Connection()
	if err != nil {
		return err
	}
	query := `INSERT INTO User (FirstName, LastName, Email, PasswordHash, ProfilePicture, Role)
          VALUES (?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(query, firstName, lastName, email, password, profilePicture, role)
	if err != nil {
		return err
	}
	return nil
}

// Read users from database:
func GetUsers() ([]*User, error){
	db, err := Connection()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user = &User{}
		err = rows.Scan(user.ID, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.ProfilePicture, user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}


