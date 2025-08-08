package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/roly-backend/internal/config"
	"github.com/roly-backend/internal/database"
	"github.com/roly-backend/internal/server"
)

func main() {

	// Setting default Logger settings
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: config.LogLevel}))
	slog.SetDefault(logger)

	// Loading Environment variables
	config.LoadEnv()

	slog.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Application started in %v mode", config.Env.AppEnv))

	// Connects to the Database
	database.Connect()

	// Starts the websocket and user auth server
	server.Start()
}
