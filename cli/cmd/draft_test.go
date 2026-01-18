package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestDraftCmd_Properties(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"Use", draftCmd.Use, "draft <title>"},
		{"Short", draftCmd.Short, "Draft a new blueprint"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("Expected %s to be %q, got %q", tt.name, tt.expected, tt.value)
			}
		})
	}
}

func TestDraftCmd_HasLongDescription(t *testing.T) {
	if draftCmd.Long == "" {
		t.Error("draftCmd.Long should not be empty")
	}
}

func TestDraftCmd_Execute(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")
	if err := os.MkdirAll(blueprintsPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprints dir: %v", err)
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

	draftCmd.Run(draftCmd, []string{"test-feature"})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Error("Expected output from draftCmd.Run()")
	}
}

func TestDraftCmd_IsRegisteredWithRoot(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "draft" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected draft command to be registered with root command")
	}
}

func TestDraftCmd_ExecuteWithError(t *testing.T) {
	// Don't create .neev structure - should show error

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	draftCmd.Run(draftCmd, []string{"test-feature"})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Error("Expected error output from draftCmd.Run()")
	}
}

func TestDraftCmd_HasCorrectArgs(t *testing.T) {
	if draftCmd.Args == nil {
		t.Error("Expected draft command to validate arguments")
	}
}

func TestDraftCmd_HasCorrectUse(t *testing.T) {
	if draftCmd.Use != "draft <title>" {
		t.Errorf("Expected Use to be 'draft <title>', got %s", draftCmd.Use)
	}
}

func TestDraftCmd_HasCorrectShort(t *testing.T) {
	if draftCmd.Short != "Draft a new blueprint" {
		t.Errorf("Expected Short to be 'Draft a new blueprint', got %s", draftCmd.Short)
	}
}

func TestDraftCmd_HasCorrectLong(t *testing.T) {
	if draftCmd.Long == "" {
		t.Error("draftCmd.Long should not be empty")
	}
}

func TestDraftCmd_ExecuteCreatesFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")
	if err := os.MkdirAll(blueprintsPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprints dir: %v", err)
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

	draftCmd.Run(draftCmd, []string{"new-blueprint"})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check output contains success or creation message
	if output == "" {
		t.Error("Expected output from draftCmd.Run()")
	}

	// Verify blueprint was created
	blueprintPath := filepath.Join(".neev", "blueprints", "new-blueprint")
	if info, err := os.Stat(blueprintPath); err != nil || !info.IsDir() {
		t.Errorf("Expected blueprint directory to be created at %s", blueprintPath)
	}
}
