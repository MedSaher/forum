package utils

import "golang.org/x/crypto/bcrypt"

// creating a secure hash of the user's password before storing it in a database.
func GenerateCryptoPassword(pass string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return hashedPassword, err
}

// Decifer the password stored in the database and compare it with the input:
func ValidatePassword(hashedPassword, plainPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)) == nil
}
