package main

import (
	"log/slog"
	"os"

	"github.com/roly-backend/internal/config"
	"github.com/roly-backend/internal/webSocket"
)

func main() {
	// Setting default Logger settings
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: config.LogLevel}))
	slog.SetDefault(logger)

	// Loading Environment variables
	config.LoadEnv()

	// Starts the Websocket server
	webSocket.StartServer()
}
