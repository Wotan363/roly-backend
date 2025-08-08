package database

import (
	"context"
	"log/slog"
	"os"

	"github.com/roly-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// This connects to the Database (has to be called once on startup and not per user session)
func Connect() {

	// Checks if Database URL is set in the ENV
	if config.Env.DBURL == "" {
		slog.LogAttrs(context.Background(), slog.LevelError, "Error connecting to database",
			slog.String("error", "DATABASE_URL is not set in your ENV"),
		)
		os.Exit(1)
	}

	// Opens the Database
	var err error
	DB, err = gorm.Open(postgres.Open(config.Env.DBURL), &gorm.Config{})
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Error connecting to database",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// Creates all neccessary tables on startup. When they already exist it does nothing
	err = DB.AutoMigrate(
		&User{},
		&Role{},
		&Chat{},
		&Message{},
		&RoleSnapshot{},
	)
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Error creating tables with AutoMigrate for database after connecting",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	slog.LogAttrs(context.Background(), slog.LevelInfo, "Connected to database successfully")
}
