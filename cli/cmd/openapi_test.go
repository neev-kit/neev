package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenAPICmd_Properties(t *testing.T) {
	tests := []struct {
		name     string
		property string
		expected string
	}{
		{"Use", openapiCmd.Use, "openapi <blueprint>"},
		{"Short", openapiCmd.Short, "Generate OpenAPI specification from a blueprint"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.property != tt.expected {
				t.Errorf("Expected %s to be %q, got %q", tt.name, tt.expected, tt.property)
			}
		})
	}
}

func TestOpenAPICmd_HasLongDescription(t *testing.T) {
	if openapiCmd.Long == "" {
		t.Error("Expected Long description to be set")
	}
	if !strings.Contains(openapiCmd.Long, "OpenAPI") {
		t.Error("Expected Long description to mention OpenAPI")
	}
}

func TestOpenAPICmd_IsRegisteredWithRoot(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "openapi" {
			found = true
			break
		}
	}
	if !found {
		t.Error("openapi command not registered with root command")
	}
}

func TestOpenAPICmd_Execute(t *testing.T) {
	// Create a temporary .neev structure
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	
	os.Chdir(tmpDir)
	
	// Create blueprint structure
	blueprintPath := filepath.Join(".neev", "blueprints", "test-api")
	if err := os.MkdirAll(blueprintPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprint dir: %v", err)
	}
	
	// Create architecture.md
	archContent := `# Test API

## Endpoints

### GET /v1/test
Test endpoint.
`
	archPath := filepath.Join(blueprintPath, "architecture.md")
	if err := os.WriteFile(archPath, []byte(archContent), 0644); err != nil {
		t.Fatalf("Failed to write architecture.md: %v", err)
	}
	
	// Execute command
	openapiCmd.Run(openapiCmd, []string{"test-api"})
	
	// Check if openapi.yaml was created
	openapiPath := filepath.Join(blueprintPath, "openapi.yaml")
	if _, err := os.Stat(openapiPath); os.IsNotExist(err) {
		t.Error("openapi.yaml was not created")
	}
	
	// Read and verify content
	content, err := os.ReadFile(openapiPath)
	if err != nil {
		t.Fatalf("Failed to read openapi.yaml: %v", err)
	}
	
	contentStr := string(content)
	if !strings.Contains(contentStr, "openapi: 3.1.0") {
		t.Error("Expected openapi.yaml to contain 'openapi: 3.1.0'")
	}
	if !strings.Contains(contentStr, "/v1/test") {
		t.Error("Expected openapi.yaml to contain '/v1/test'")
	}
}
