package config

import (
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type ENV struct {
	AppEnv       string
	Port         string
	RolyAPIURL   string
	DBURL        string
	OpenAIAPIKey string
}

var Env ENV

// Loads the Environment Variables
func LoadEnv() {
	env := os.Getenv("APP_ENV")
	var envFile string

	switch env {
	case "production":
		envFile = ".env.production"
	default:
		envFile = ".env.development"
	}

	// Loads the db env values
	if err := godotenv.Load(".env.db"); err != nil {
		slog.LogAttrs(context.Background(), slog.LevelWarn, "Warning loading .env.db",
			slog.String("error", err.Error()),
		)
	}

	// Loads the secrets env values (API Keys)
	if err := godotenv.Load(".env.secrets"); err != nil {
		slog.LogAttrs(context.Background(), slog.LevelWarn, "Warning loading .env.db",
			slog.String("error", err.Error()),
		)
	}

	// loads the other env values
	if err := godotenv.Load(envFile); err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Error while loading .env file",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	Env = ENV{
		AppEnv:       os.Getenv("APP_ENV"),
		Port:         os.Getenv("PORT"),
		RolyAPIURL:   os.Getenv("API_URL"),
		DBURL:        os.Getenv("DATABASE_URL"),
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
	}
}
