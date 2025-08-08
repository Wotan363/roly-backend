package users

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/roly-backend/internal/database"
)

// JWT Claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Handles the Login requests
func LoginHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Buffers body in case we need to log it when a error happens
	bodyBytes, err1 := io.ReadAll(c.Request.Body)
	if err1 != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Error buffering body on login request",
			slog.String("error", err1.Error()),
		)
		// Sends error to client
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Reset Body so ShouldBindJSON still works
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// project JSON request on input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		// Save Body in a map for logging
		var raw map[string]any
		if err2 := json.Unmarshal(bodyBytes, &raw); err2 == nil {
			// Obfuscate password so we don't log any sensitive data
			if _, ok := raw["password"]; ok {
				raw["password"] = "***PASSWORD REDACTED***"
			}
			slog.LogAttrs(context.Background(), slog.LevelError, "Error binding login request body",
				slog.String("error", err.Error()),
				slog.Any("client_message", raw),
			)
		} else {
			// If the JSON request is completely broken, log the raw body
			slog.LogAttrs(context.Background(), slog.LevelError, "Error binding login request body - invalid JSON",
				slog.String("error", err.Error()),
				slog.String("raw_body", string(bodyBytes)),
			)
		}
		// Sends error to client
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Conver E-Mail to lower case since all mails are saved in lower case in the database
	input.Email = strings.ToLower(input.Email)

	// Searches for the user in the database
	var user database.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		slog.LogAttrs(context.Background(), slog.LevelInfo, "Login failed - user not found",
			slog.String("email", input.Email),
		)
		// Sends error to client
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Checks if the password is correct
	if !CheckPasswordHash(input.Password, user.Password) {
		slog.LogAttrs(context.Background(), slog.LevelInfo, "Login failed - wrong password",
			slog.String("email", input.Email),
		)
		// Sends error to client
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	tokenString, expirationTime, err2 := getNewJWTToken(user)
	if err2 != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Error generating JWT token",
			slog.String("error", err2.Error()),
			slog.String("email", input.Email),
		)
		// Sends error to client
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	slog.LogAttrs(context.Background(), slog.LevelInfo, "User logged in successfully",
		slog.String("user_id", user.ID.String()),
		slog.String("email", user.Email),
	)

	// Sends token back to the client
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
		"expires": expirationTime,
	})
}
