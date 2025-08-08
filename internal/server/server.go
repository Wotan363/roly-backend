package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/roly-backend/internal/config"
)

// Starts the server
func Start() {
	ginEngine := SetupRouter()

	slog.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Server listening on port %v", config.Port))

	// Starts websocket and user auth server
	err := ginEngine.Run(fmt.Sprintf(":%v", config.Port))
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Error with server",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}
