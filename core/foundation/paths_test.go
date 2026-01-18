package foundation

import "testing"

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		actual   string
		expected string
	}{
		{"RootDir", RootDir, ".neev"},
		{"BlueprintsDir", BlueprintsDir, "blueprints"},
		{"FoundationDir", FoundationDir, "foundation"},
		{"ConfigFile", ConfigFile, "neev.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("Expected %s to be %q, got %q", tt.name, tt.expected, tt.actual)
			}
		})
	}
}
