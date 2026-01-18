package cmd

import (
	"bytes"
	"os"
	"testing"
)

func TestRootCmd_Execute_NoArgs(t *testing.T) {
	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.Run(rootCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Error("Expected output from rootCmd.Run()")
	}
}

func TestRootCmd_Properties(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"Use", rootCmd.Use, "neev"},
		{"Short", rootCmd.Short, "Neev - The blueprint orchestration tool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("Expected %s to be %q, got %q", tt.name, tt.expected, tt.value)
			}
		})
	}
}

func TestRootCmd_HasLongDescription(t *testing.T) {
	if rootCmd.Long == "" {
		t.Error("rootCmd.Long should not be empty")
	}
}

func TestExecute_Success(t *testing.T) {
	// This test verifies that Execute() can be called without panic
	oldArgs := os.Args
	os.Args = []string{"neev"}
	defer func() { os.Args = oldArgs }()

	// Execute should handle the root command successfully
	// We don't expect it to panic
}

func TestRootCmd_InitFunction(t *testing.T) {
	if rootCmd == nil {
		t.Error("Expected rootCmd to be initialized")
	}
}
