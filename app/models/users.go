package models

// Declare a model to represent the user and ease data exchange between backend and frontend:
type User struct {
	ID             int    `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	PasswordHash   string `json:"password"`
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
	defer db.Close()
	query := `INSERT INTO User (FirstName, LastName, Email, PasswordHash, ProfilePicture, Role)
          VALUES (?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(query, firstName, lastName, email, password, profilePicture, role)
	if err != nil {
		return err
	}
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
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.ProfilePicture, &user.Role); err != nil {
			return nil, err
		}
		Users = append(Users, user)
	}
	return Users, nil
}
