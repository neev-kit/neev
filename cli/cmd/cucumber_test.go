package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCucumberCmd_Properties(t *testing.T) {
	tests := []struct {
		name     string
		property string
		expected string
	}{
		{"Use", cucumberCmd.Use, "cucumber <blueprint>"},
		{"Short", cucumberCmd.Short, "Generate Cucumber/BDD test scaffolding from a blueprint"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.property != tt.expected {
				t.Errorf("Expected %s to be %q, got %q", tt.name, tt.expected, tt.property)
			}
		})
	}
}

func TestCucumberCmd_HasLongDescription(t *testing.T) {
	if cucumberCmd.Long == "" {
		t.Error("Expected Long description to be set")
	}
	if !strings.Contains(cucumberCmd.Long, "Cucumber") {
		t.Error("Expected Long description to mention Cucumber")
	}
}

func TestCucumberCmd_IsRegisteredWithRoot(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "cucumber" {
			found = true
			break
		}
	}
	if !found {
		t.Error("cucumber command not registered with root command")
	}
}

func TestCucumberCmd_Execute(t *testing.T) {
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

### POST /v1/test
Create test resource.
`
	archPath := filepath.Join(blueprintPath, "architecture.md")
	if err := os.WriteFile(archPath, []byte(archContent), 0644); err != nil {
		t.Fatalf("Failed to write architecture.md: %v", err)
	}
	
	// Execute command with language flag
	cucumberLang = "go"
	cucumberCmd.Run(cucumberCmd, []string{"test-api"})
	
	// Check if tests directory was created
	testsPath := filepath.Join(blueprintPath, "tests")
	if _, err := os.Stat(testsPath); os.IsNotExist(err) {
		t.Error("tests directory was not created")
	}
	
	// Check if api.feature was created
	featurePath := filepath.Join(testsPath, "api.feature")
	if _, err := os.Stat(featurePath); os.IsNotExist(err) {
		t.Error("api.feature was not created")
	}
	
	// Check if steps.go was created
	stepsPath := filepath.Join(testsPath, "steps.go")
	if _, err := os.Stat(stepsPath); os.IsNotExist(err) {
		t.Error("steps.go was not created")
	}
	
	// Read and verify feature content
	featureContent, err := os.ReadFile(featurePath)
	if err != nil {
		t.Fatalf("Failed to read api.feature: %v", err)
	}
	
	contentStr := string(featureContent)
	if !strings.Contains(contentStr, "Feature: Test Api API") {
		t.Error("Expected api.feature to contain 'Feature: Test Api API'")
	}
	if !strings.Contains(contentStr, "Scenario: GET /v1/test") {
		t.Error("Expected api.feature to contain 'Scenario: GET /v1/test'")
	}
	if !strings.Contains(contentStr, "Scenario: POST /v1/test") {
		t.Error("Expected api.feature to contain 'Scenario: POST /v1/test'")
	}
	
	// Verify steps content
	stepsContent, err := os.ReadFile(stepsPath)
	if err != nil {
		t.Fatalf("Failed to read steps.go: %v", err)
	}
	
	stepsStr := string(stepsContent)
	if !strings.Contains(stepsStr, "package steps") {
		t.Error("Expected steps.go to contain 'package steps'")
	}
	if !strings.Contains(stepsStr, "APIContext") {
		t.Error("Expected steps.go to contain 'APIContext'")
	}
}

func TestCucumberCmd_WithoutLanguageFlag(t *testing.T) {
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
	
	// Execute command without language flag
	cucumberLang = ""
	cucumberCmd.Run(cucumberCmd, []string{"test-api"})
	
	// Check if api.feature was created
	testsPath := filepath.Join(blueprintPath, "tests")
	featurePath := filepath.Join(testsPath, "api.feature")
	if _, err := os.Stat(featurePath); os.IsNotExist(err) {
		t.Error("api.feature was not created")
	}
	
	// Step definitions should not be created
	stepsPath := filepath.Join(testsPath, "steps.go")
	if _, err := os.Stat(stepsPath); err == nil {
		t.Error("steps.go should not be created without language flag")
	}
}
