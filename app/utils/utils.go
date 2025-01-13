package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

// creating a secure hash of the user's password before storing it in a database.
func GenerateCryptoPassword(pass string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return hashedPassword, err
}

// Decifer the password stored in the database and compare it with the input:
func ValidatePassword(hashedPassword, plainPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)) == nil
}

// Generate a universal unique session identifier:
func GenerateUUID() string {
	uuid := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		panic(err)
	}

	// Set the version (4) and variant (10xx)
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10xx

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// Generate a CSRF:
// CSRF is an attack where a malicious website tricks a user's browser into performing 
// unwanted actions on another website where the user is authenticated.
func GenerateCSRFToken() (string, error) {
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// Encode the bytes to a base64 string
	return base64.URLEncoding.EncodeToString(bytes), nil
}
