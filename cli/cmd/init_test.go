package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestInitCmd_Properties(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"Use", initCmd.Use, "init"},
		{"Short", initCmd.Short, "Initialize the Neev foundation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("Expected %s to be %q, got %q", tt.name, tt.expected, tt.value)
			}
		})
	}
}

func TestInitCmd_HasLongDescription(t *testing.T) {
	if initCmd.Long == "" {
		t.Error("initCmd.Long should not be empty")
	}
}

func TestInitCmd_ExecuteSuccess(t *testing.T) {
	tmpDir := t.TempDir()

	// Change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	initCmd.Run(initCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Error("Expected output from initCmd.Run()")
	}

	// Verify structure was created
	neevPath := filepath.Join(tmpDir, ".neev")
	if _, err := os.Stat(neevPath); os.IsNotExist(err) {
		t.Error("Expected .neev directory to be created")
	}
}

func TestInitCmd_IsRegisteredWithRoot(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "init" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected init command to be registered with root command")
	}
}
