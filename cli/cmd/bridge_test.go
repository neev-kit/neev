package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBridgeCmd_Properties(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"Use", bridgeCmd.Use, "bridge [flags]"},
		{"Short", bridgeCmd.Short, "Bridge to external systems"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("Expected %s to be %q, got %q", tt.name, tt.expected, tt.value)
			}
		})
	}
}

func TestBridgeCmd_HasLongDescription(t *testing.T) {
	if bridgeCmd.Long == "" {
		t.Error("bridgeCmd.Long should not be empty")
	}
}

func TestBridgeCmd_Execute(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev structure
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")

	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("Failed to create foundation dir: %v", err)
	}
	if err := os.MkdirAll(blueprintsPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprints dir: %v", err)
	}

	// Create a test file
	if err := os.WriteFile(filepath.Join(foundationPath, "test.md"), []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	bridgeCmd.Run(bridgeCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Error("Expected output from bridgeCmd.Run()")
	}
}

func containsLower(s, substr string) bool {
	lowerS := toLower(s)
	lowerSubstr := toLower(substr)
	for i := 0; i <= len(lowerS)-len(lowerSubstr); i++ {
		if lowerS[i:i+len(lowerSubstr)] == lowerSubstr {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			b[i] = c + ('a' - 'A')
		} else {
			b[i] = c
		}
	}
	return string(b)
}

func TestBridgeCmd_IsRegisteredWithRoot(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "bridge" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected bridge command to be registered with root command")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestBridgeCmd_ExecuteWithFocus(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev structure with specific focus content
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")

	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("Failed to create foundation dir: %v", err)
	}
	if err := os.MkdirAll(blueprintsPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprints dir: %v", err)
	}

	// Create a test file
	if err := os.WriteFile(filepath.Join(foundationPath, "test.md"), []byte("api endpoint"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a new command with the focus flag set
	bridgeCmd.Flags().Set("focus", "api")
	bridgeCmd.Run(bridgeCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Error("Expected output from bridgeCmd.Run() with focus")
	}
}

func TestBridgeCmd_HasFlags(t *testing.T) {
	flags := bridgeCmd.Flags()
	if flags == nil {
		t.Error("Expected bridgeCmd to have flags")
	}

	focusFlag := flags.Lookup("focus")
	if focusFlag == nil {
		t.Error("Expected 'focus' flag to be defined")
	}
}

func TestBridgeCmd_LongDescription_NotEmpty(t *testing.T) {
	if bridgeCmd.Long == "" {
		t.Error("bridgeCmd.Long should not be empty")
	}
}

func TestBridgeCmd_Execute_CheckOutput(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev structure
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")

	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("Failed to create foundation dir: %v", err)
	}
	if err := os.MkdirAll(blueprintsPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprints dir: %v", err)
	}

	// Create multiple test files
	if err := os.WriteFile(filepath.Join(foundationPath, "api.md"), []byte("# API\nendpoints"), 0644); err != nil {
		t.Fatalf("Failed to write api file: %v", err)
	}

	if err := os.WriteFile(filepath.Join(foundationPath, "db.md"), []byte("# Database\nschema"), 0644); err != nil {
		t.Fatalf("Failed to write db file: %v", err)
	}

	// Change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	bridgeCmd.Run(bridgeCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Error("Expected output from bridgeCmd.Run()")
	}

	if !strings.Contains(output, "# Project Foundation") {
		t.Errorf("Expected '# Project Foundation' in output")
	}
}
