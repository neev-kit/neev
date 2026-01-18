package cmd

import (
	"bytes"
	"os"
	"testing"
)

func TestBridgeCmd_Properties(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"Use", bridgeCmd.Use, "bridge"},
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

	if !containsLower(output, "bridge") || !containsLower(output, "building") {
		t.Errorf("Expected output to mention 'bridge', got: %s", output)
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
		if cmd.Use == "bridge" {
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
