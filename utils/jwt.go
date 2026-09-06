package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a JWT containing the user's ID and Admin status
func GenerateToken(userID string, isAdmin bool) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "sommerce-secret-key-default" // Fallback if env is missing
	}

	claims := jwt.MapClaims{
		"user_id":  userID,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 Hours Expiration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
