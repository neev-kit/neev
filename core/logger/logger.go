package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

var globalLogger *slog.Logger

// Init initializes the global logger based on environment variables
func Init() {
	logFormat := os.Getenv("NEEV_LOG")
	globalLogger = createLogger(logFormat)
}

// createLogger creates a logger based on the format
func createLogger(format string) *slog.Logger {
	var handler slog.Handler

	// Use JSON format if NEEV_LOG=json, otherwise use human-readable format
	if format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Human-readable format with colors and emojis
		handler = NewColoredHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	return slog.New(handler)
}

// ColoredHandler is a custom slog handler that outputs human-readable logs
type ColoredHandler struct {
	handler slog.Handler
}

// NewColoredHandler creates a new ColoredHandler
func NewColoredHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &ColoredHandler{
		handler: slog.NewTextHandler(w, opts),
	}
}

// Handle processes log records with color and emoji enhancements
func (h *ColoredHandler) Handle(ctx context.Context, r slog.Record) error {
	// Add emoji based on level
	emoji := ""
	switch r.Level {
	case slog.LevelDebug:
		emoji = "🔍 "
	case slog.LevelInfo:
		emoji = "ℹ️  "
	case slog.LevelWarn:
		emoji = "⚠️  "
	case slog.LevelError:
		emoji = "❌ "
	}

	// Prepend emoji to the message
	r.Message = emoji + r.Message

	return h.handler.Handle(ctx, r)
}

// Enabled reports whether the handler handles records at the given level
func (h *ColoredHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// WithAttrs returns a handler with attributes
func (h *ColoredHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ColoredHandler{
		handler: h.handler.WithAttrs(attrs),
	}
}

// WithGroup returns a handler with a group
func (h *ColoredHandler) WithGroup(name string) slog.Handler {
	return &ColoredHandler{
		handler: h.handler.WithGroup(name),
	}
}

// Info logs an info level message
func Info(msg string, args ...any) {
	if globalLogger == nil {
		Init()
	}
	globalLogger.Info(msg, args...)
}

// Debug logs a debug level message
func Debug(msg string, args ...any) {
	if globalLogger == nil {
		Init()
	}
	globalLogger.Debug(msg, args...)
}

// Warn logs a warning level message
func Warn(msg string, args ...any) {
	if globalLogger == nil {
		Init()
	}
	globalLogger.Warn(msg, args...)
}

// Error logs an error level message
func Error(msg string, args ...any) {
	if globalLogger == nil {
		Init()
	}
	globalLogger.Error(msg, args...)
}

// Printf is a convenience function for formatted output
func Printf(format string, args ...any) {
	fmt.Printf(format, args...)
}

// GetLogger returns the global logger instance
func GetLogger() *slog.Logger {
	if globalLogger == nil {
		Init()
	}
	return globalLogger
}
