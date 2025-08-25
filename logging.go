package machineid

import (
	"log/slog"
	"os"
	"strings"
)

// LogLevel represents the logging level for the machineid package
type LogLevel int

const (
	// LogLevelError shows only error messages
	LogLevelError LogLevel = iota
	// LogLevelInfo shows info and error messages
	LogLevelInfo
	// LogLevelDebug shows all messages including debug
	LogLevelDebug
)

// SetLogLevel configures the global log level for the machineid package.
// By default, logging is set to Error level to avoid noise in production.
func SetLogLevel(level LogLevel) {
	var slogLevel slog.Level

	switch level {
	case LogLevelError:
		slogLevel = slog.LevelError
	case LogLevelInfo:
		slogLevel = slog.LevelInfo
	case LogLevelDebug:
		slogLevel = slog.LevelDebug
	default:
		slogLevel = slog.LevelError
	}

	// Create a new handler with the specified level
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slogLevel,
	})

	// Set the default logger
	slog.SetDefault(slog.New(handler))
}

// SetLogLevelFromEnv configures the log level from the MACHINEID_LOG_LEVEL environment variable.
// Valid values are: "error", "info", "debug" (case insensitive).
// If the environment variable is not set or invalid, defaults to Error level.
func SetLogLevelFromEnv() {
	envLevel := strings.ToLower(os.Getenv("MACHINEID_LOG_LEVEL"))

	switch envLevel {
	case "debug":
		SetLogLevel(LogLevelDebug)
	case "info":
		SetLogLevel(LogLevelInfo)
	case "error":
		SetLogLevel(LogLevelError)
	default:
		// Default to error level for production use
		SetLogLevel(LogLevelError)
	}
}

func init() {
	// Initialize with environment-based configuration
	SetLogLevelFromEnv()
}
