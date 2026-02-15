package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// NewLogger creates a new structured logger with zerolog
func NewLogger(level string) zerolog.Logger {
	// Parse log level
	logLevel := parseLogLevel(level)
	zerolog.SetGlobalLevel(logLevel)

	// Use pretty console output for development
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	logger := zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	return logger
}

// parseLogLevel converts string level to zerolog.Level
func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}