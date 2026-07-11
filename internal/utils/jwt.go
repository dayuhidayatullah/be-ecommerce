package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT generates a JWT token for the given user ID and role
func GenerateJWT(userID, role string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default_secret_key" // Fallback, shouldn't happen in prod with proper .env
	}

	expHoursStr := os.Getenv("JWT_EXPIRATION_HOURS")
	expHours, err := strconv.Atoi(expHoursStr)
	if err != nil || expHours <= 0 {
		expHours = 24 // default 24 hours
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * time.Duration(expHours)).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
