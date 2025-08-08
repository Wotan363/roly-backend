package config

import (
	"fmt"
	"log/slog"
)

// Debug and Log Settings
var DebugMode bool = true                 // If set to true then debug messages are logged if false, then not -> so that way no unneccessary Mutext RLocks happen
var LogLevel slog.Level = slog.LevelDebug // Defines the Loglevel
var AllowedOrigins = map[string]bool{
	"https://roly.ai":                        true,
	fmt.Sprintf("http://localhost:%v", Port): true, // only for dev
}

var Port int = 8080
var MinPasswordLength int = 6
