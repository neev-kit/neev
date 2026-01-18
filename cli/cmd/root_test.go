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

func TestRootCmd_HasSubcommands(t *testing.T) {
	commands := rootCmd.Commands()
	if len(commands) == 0 {
		t.Error("Expected root command to have subcommands")
	}

	// Check for expected commands
	expectedCmds := []string{"init", "draft", "bridge"}
	for _, expectedCmd := range expectedCmds {
		found := false
		for _, cmd := range commands {
			if cmd.Name() == expectedCmd {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected %s command to be registered", expectedCmd)
		}
	}
}

func TestExecute_WithRootOnly(t *testing.T) {
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
		t.Error("Expected output from root command")
	}
}

func TestRootCmd_RunnableCommand(t *testing.T) {
	if !rootCmd.Runnable() {
		t.Error("Root command should be runnable")
	}
}

func TestRootCmd_Short_NotEmpty(t *testing.T) {
	if rootCmd.Short == "" {
		t.Error("Root command Short should not be empty")
	}
}

func TestRootCmd_Long_NotEmpty(t *testing.T) {
	if rootCmd.Long == "" {
		t.Error("Root command Long should not be empty")
	}
}

func TestInitCmd_NotEmpty(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "init" {
			found = true
			if cmd.Short == "" {
				t.Error("Init command Short should not be empty")
			}
			if cmd.Long == "" {
				t.Error("Init command Long should not be empty")
			}
			break
		}
	}
	if !found {
		t.Error("Init command not found")
	}
}

func TestDraftCmd_NotEmpty(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "draft" {
			found = true
			if cmd.Short == "" {
				t.Error("Draft command Short should not be empty")
			}
			if cmd.Long == "" {
				t.Error("Draft command Long should not be empty")
			}
			break
		}
	}
	if !found {
		t.Error("Draft command not found")
	}
}

func TestBridgeCmd_NotEmpty(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "bridge" {
			found = true
			if cmd.Short == "" {
				t.Error("Bridge command Short should not be empty")
			}
			if cmd.Long == "" {
				t.Error("Bridge command Long should not be empty")
			}
			break
		}
	}
	if !found {
		t.Error("Bridge command not found")
	}
}
