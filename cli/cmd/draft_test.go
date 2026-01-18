package cmd

import (
	"bytes"
	"os"
	"testing"
)

func TestDraftCmd_Properties(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"Use", draftCmd.Use, "draft"},
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
	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	draftCmd.Run(draftCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Error("Expected output from draftCmd.Run()")
	}

	if !containsLower(output, "draft") || !containsLower(output, "blueprint") {
		t.Errorf("Expected output to mention 'draft' or 'blueprint', got: %s", output)
	}
}

func TestDraftCmd_IsRegisteredWithRoot(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "draft" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected draft command to be registered with root command")
	}
}
