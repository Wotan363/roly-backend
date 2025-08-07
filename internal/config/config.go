package config

import "log/slog"

// Debug and Log Settings
var DebugMode bool = true                 // If set to true then debug messages are logged if false, then not -> so that way no unneccessary Mutext RLocks happen
var LogLevel slog.Level = slog.LevelDebug // Defines the Loglevel
var AllowedOrigins = map[string]bool{
	"https://roly.ai":       true,
	"http://localhost:3000": true, // only for dev
}
