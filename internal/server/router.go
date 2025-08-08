package server

import (
	"github.com/gin-gonic/gin"
	"github.com/roly-backend/internal/users"
	"github.com/roly-backend/internal/webSocket"
)

// Registers all API routes (user auth and websocket) and returns the ginEngine
func SetupRouter() *gin.Engine {
	ginEngine := gin.Default()

	// Defines the available REST-API-Routes
	api := ginEngine.Group("/api")
	{
		api.POST("/register", users.RegisterHandler)
		api.POST("/login", users.LoginHandler)
	}

	// If we ever need a JWT protected HTTP Route, here is an example how to do the authentication since it is already implemented
	// authGroup := ginEngine.Group("/api")
	// authGroup.Use(users.JWTAuthMiddleware())
	// authGroup.GET("/profile", func(c *gin.Context) {
	// 	claims, _ := c.Get(string(users.UserContextKey))
	// 	c.JSON(200, gin.H{"user": claims})
	// })

	// WebSocket-Route Gin provides the http.ResponseWriter und *http.Request
	ginEngine.GET("/ws", func(c *gin.Context) {
		webSocket.HandleWebSocket(c.Writer, c.Request)
	})

	return ginEngine
}
