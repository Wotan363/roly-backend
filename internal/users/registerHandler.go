package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/roly-backend/internal/config"
	"github.com/roly-backend/internal/database"
)

// Handles the Register requests
func RegisterHandler(c *gin.Context) {
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
			slog.LogAttrs(context.Background(), slog.LevelError, "Error binding register request body",
				slog.String("error", err.Error()),
				slog.Any("client_message", raw),
			)
		} else {
			// If the JSON request is completely broken, log the raw body
			slog.LogAttrs(context.Background(), slog.LevelError, "Error binding register request body - invalid JSON",
				slog.String("error", err.Error()),
				slog.String("raw_body", string(bodyBytes)),
			)
		}
		// Sends error to client
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Checks minimum password length
	if len(input.Password) < config.MinPasswordLength {
		slog.LogAttrs(context.Background(), slog.LevelInfo, "Error with register request (Password too short)",
			slog.String("email", input.Email),
			slog.Int("length", len(input.Password)),
		)
		// Sends error to client
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least " +
				fmt.Sprintf("%d characters", config.MinPasswordLength),
		})
		return
	}

	// Save E-Mail in lower case so no duplicate signups with lower-upper-cased mixed happen
	input.Email = strings.ToLower(input.Email)

	// Hashes password
	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Error with register request (Could not hash password)",
			slog.String("error", err.Error()),
			slog.String("email", input.Email),
		)
		// Sends error to client
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Creates an object representing a user in the database
	newUser := database.User{
		ID:        uuid.New(),
		Email:     input.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	// Saves new user in Database
	if err := database.DB.Create(&newUser).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			slog.LogAttrs(context.Background(), slog.LevelInfo, "Attempted to register with existing email",
				slog.String("email", input.Email),
			)
			// Sends error to client
			c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
			return
		}
		slog.LogAttrs(context.Background(), slog.LevelError, "Error with register request (couldn't save new user in database)",
			slog.String("error", err.Error()),
			slog.String("email", input.Email),
		)
		// Sends error to client
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	slog.LogAttrs(context.Background(), slog.LevelInfo, "User registered successfully",
		slog.String("user_id", newUser.ID.String()),
		slog.String("email", newUser.Email),
	)

	// Sends user_id back to client
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": newUser.ID,
	})
}
