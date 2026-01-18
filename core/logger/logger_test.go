package logger

import "testing"

func TestInit(t *testing.T) {
	Init()
}

func TestInfo(t *testing.T) {
	Info("test message")
	Info("test with key", "key", "value")
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
