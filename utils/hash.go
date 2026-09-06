package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword generates a bcrypt hash for a plain text password
func HashPassword(password string) (string, error) {
	// Cost 12 is a good balance between security and performance
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// CheckPasswordHash compares a plain text password with a bcrypt hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
