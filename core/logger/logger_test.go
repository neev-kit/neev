package logger

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	Init()
	logger := GetLogger()
	if logger == nil {
		t.Fatal("GetLogger returned nil after Init()")
	}
}

func TestInfo(t *testing.T) {
	Info("test message")
	Info("test with key", "key", "value")
	Info("test with multiple", "key1", "value1", "key2", "value2")
}

func TestDebug(t *testing.T) {
	Debug("debug message")
	Debug("debug with key", "key", "value")
}

func TestWarn(t *testing.T) {
	Warn("warning message")
	Warn("warning with key", "key", "value")
}

func TestError(t *testing.T) {
	Error("error message")
	Error("error with key", "key", "value")
}

func TestPrintf(t *testing.T) {
	// Printf should not panic
	Printf("formatted: %d\n", 42)
	Printf("test %s", "string")
}

func TestGetLogger(t *testing.T) {
	logger := GetLogger()
	if logger == nil {
		t.Fatal("GetLogger returned nil")
	}

	// Second call should return the same logger
	logger2 := GetLogger()
	if logger != logger2 {
		t.Error("GetLogger should return same logger instance")
	}
}

func TestColoredHandler_Handle(t *testing.T) {
	var buf bytes.Buffer
	handler := NewColoredHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	ctx := context.Background()

	// Test Debug level
	record := slog.NewRecord(
		time.Now(), slog.LevelDebug, "debug msg", 0,
	)
	err := handler.Handle(ctx, record)
	if err != nil {
		t.Errorf("Handle failed: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "debug msg") {
		t.Errorf("Expected 'debug msg' in output, got: %s", output)
	}

	// Test Info level
	buf.Reset()
	record = slog.NewRecord(
		time.Now(), slog.LevelInfo, "info msg", 0,
	)
	_ = handler.Handle(ctx, record)
	output = buf.String()
	if !strings.Contains(output, "info msg") {
		t.Errorf("Expected 'info msg' in output, got: %s", output)
	}

	// Test Warn level
	buf.Reset()
	record = slog.NewRecord(
		time.Now(), slog.LevelWarn, "warn msg", 0,
	)
	_ = handler.Handle(ctx, record)
	output = buf.String()
	if !strings.Contains(output, "warn msg") {
		t.Errorf("Expected 'warn msg' in output, got: %s", output)
	}

	// Test Error level
	buf.Reset()
	record = slog.NewRecord(
		time.Now(), slog.LevelError, "error msg", 0,
	)
	_ = handler.Handle(ctx, record)
	output = buf.String()
	if !strings.Contains(output, "error msg") {
		t.Errorf("Expected 'error msg' in output, got: %s", output)
	}
}

func TestColoredHandler_Enabled(t *testing.T) {
	var buf bytes.Buffer
	handler := NewColoredHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	ctx := context.Background()

	// Info level should be enabled
	if !handler.Enabled(ctx, slog.LevelInfo) {
		t.Error("Expected Info level to be enabled")
	}

	// Debug level should be disabled (below Info)
	if handler.Enabled(ctx, slog.LevelDebug) {
		t.Error("Expected Debug level to be disabled")
	}

	// Error level should be enabled (above Info)
	if !handler.Enabled(ctx, slog.LevelError) {
		t.Error("Expected Error level to be enabled")
	}
}

func TestColoredHandler_WithAttrs(t *testing.T) {
	var buf bytes.Buffer
	handler := NewColoredHandler(&buf, nil)

	attrs := []slog.Attr{
		slog.String("key", "value"),
	}

	newHandler := handler.WithAttrs(attrs)
	if newHandler == nil {
		t.Fatal("WithAttrs returned nil")
	}

	// Verify it's a ColoredHandler
	_, ok := newHandler.(*ColoredHandler)
	if !ok {
		t.Error("WithAttrs should return a ColoredHandler")
	}
}

func TestColoredHandler_WithGroup(t *testing.T) {
	var buf bytes.Buffer
	handler := NewColoredHandler(&buf, nil)

	newHandler := handler.WithGroup("testgroup")
	if newHandler == nil {
		t.Fatal("WithGroup returned nil")
	}

	// Verify it's a ColoredHandler
	_, ok := newHandler.(*ColoredHandler)
	if !ok {
		t.Error("WithGroup should return a ColoredHandler")
	}
}

func TestNewColoredHandler(t *testing.T) {
	var buf bytes.Buffer

	// Test with nil options
	handler1 := NewColoredHandler(&buf, nil)
	if handler1 == nil {
		t.Fatal("NewColoredHandler with nil options returned nil")
	}

	// Test with custom options
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler2 := NewColoredHandler(&buf, opts)
	if handler2 == nil {
		t.Fatal("NewColoredHandler with options returned nil")
	}
}

func TestCreateLoggerJSON(t *testing.T) {
	// Save original env
	origFormat := os.Getenv("NEEV_LOG")
	defer os.Setenv("NEEV_LOG", origFormat)

	os.Setenv("NEEV_LOG", "json")
	logger := createLogger("json")
	if logger == nil {
		t.Fatal("createLogger(json) returned nil")
	}
}

func TestCreateLoggerText(t *testing.T) {
	logger := createLogger("text")
	if logger == nil {
		t.Fatal("createLogger(text) returned nil")
	}
}

func TestCreateLoggerDefault(t *testing.T) {
	logger := createLogger("")
	if logger == nil {
		t.Fatal("createLogger(empty) returned nil")
	}
}

func TestLoggerFunctions_WithoutInit(t *testing.T) {
	// Reset global logger to test lazy initialization
	globalLogger = nil

	// These should not panic even with uninitialized global logger
	Info("test")
	Debug("test")
	Warn("test")
	Error("test")

	// GetLogger should initialize if needed
	logger := GetLogger()
	if logger == nil {
		t.Fatal("GetLogger returned nil")
	}
}

func TestEmojiLevels(t *testing.T) {
	var buf bytes.Buffer
	handler := NewColoredHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	ctx := context.Background()

	testCases := []struct {
		level   slog.Level
		message string
	}{
		{slog.LevelDebug, "debug message"},
		{slog.LevelInfo, "info message"},
		{slog.LevelWarn, "warning message"},
		{slog.LevelError, "error message"},
	}

	for _, tc := range testCases {
		buf.Reset()
		record := slog.NewRecord(time.Now(), tc.level, tc.message, 0)
		err := handler.Handle(ctx, record)
		if err != nil {
			t.Errorf("Handle failed for level %s: %v", tc.level, err)
		}

		output := buf.String()
		if !strings.Contains(output, tc.message) {
			t.Errorf("Expected '%s' in output for level %s, got: %s", tc.message, tc.level, output)
		}
	}
}
