package users

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/roly-backend/internal/config"
	"github.com/roly-backend/internal/database"
	"golang.org/x/crypto/bcrypt"
)

// Hashes the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Checks if password matches with the given password hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generates a new JWT token and returns it with the expiration time
func getNewJWTToken(user database.User) (string, time.Time, error) {
	// Creates JWT Claims
	expirationTime := time.Now().Add(24 * time.Hour) // Valid for one day
	claims := &Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "roly-backend",
		},
	}

	// Signs the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Env.JWTSecret))
	if err != nil {
		return "", expirationTime, err
	}

	return tokenString, expirationTime, nil
}
