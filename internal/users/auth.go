package users

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/roly-backend/internal/config"
)

// ContextKey for User-Data in Gin/Websocket
type ctxKey string

const UserContextKey ctxKey = "user"

// Middleware for Gin HTTP-Routs to authenticate the client with JWT
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extracts token from header
		tokenString, err := extractTokenFromHeader(c.Request)
		if err != nil {
			slog.LogAttrs(context.Background(), slog.LevelInfo, "Unauthorized HTTP request (no/invalid token)",
				slog.String("path", c.FullPath()),
				slog.String("error", err.Error()),
			)
			// Sends error to client
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Validates the JWT token
		claims, err := validateJWT(tokenString)
		if err != nil {
			slog.LogAttrs(context.Background(), slog.LevelInfo, "Unauthorized HTTP request - invalid JWT",
				slog.String("path", c.FullPath()),
				slog.String("error", err.Error()),
			)
			// Sends error to client
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Puts claims in the Context (connection is approved)
		c.Set(string(UserContextKey), claims)
		c.Next()
	}
}

// Extracts the "Bearer <token>" from the header
func extractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no Authorization header provided")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	return parts[1], nil
}

// validates JWT
func validateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	// Decodes Header and Payload and then checks if signature is correct
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Env.JWTSecret), nil // This just gives the ParseWithClaims function the secret so it can validate the singature
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// Checks if the JWT is correct for the WebSocket-Handshake and then returns the claims (user data)
func ValidateWebSocketJWT(r *http.Request) (*Claims, error) {
	tokenString, err := extractTokenFromHeader(r)
	if err != nil {
		return nil, err
	}
	return validateJWT(tokenString)
}
