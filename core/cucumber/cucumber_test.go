package cucumber

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/neev-kit/neev/core/openapi"
)

func TestGenerateFeatureFile(t *testing.T) {
	endpoints := []openapi.Endpoint{
		{
			Method:      "GET",
			Path:        "/v1/users",
			Description: "List all users",
			Parameters: []openapi.Parameter{
				{Name: "page", In: "query", Description: "Page number"},
			},
		},
		{
			Method:      "POST",
			Path:        "/v1/users",
			Description: "Create a new user",
			Request:     `{"name": "John"}`,
		},
	}

	feature, err := GenerateFeatureFile(endpoints, "user-api")
	if err != nil {
		t.Fatalf("GenerateFeatureFile failed: %v", err)
	}

	if !strings.Contains(feature, "Feature: User Api API") {
		t.Error("Expected feature to contain 'Feature: User Api API'")
	}

	if !strings.Contains(feature, "Scenario: GET /v1/users") {
		t.Error("Expected feature to contain 'Scenario: GET /v1/users'")
	}

	if !strings.Contains(feature, "Scenario: POST /v1/users") {
		t.Error("Expected feature to contain 'Scenario: POST /v1/users'")
	}

	if !strings.Contains(feature, "Given the API is available") {
		t.Error("Expected feature to contain 'Given the API is available'")
	}

	if !strings.Contains(feature, "I am authenticated") {
		t.Error("POST endpoint should require authentication")
	}

	if !strings.Contains(feature, `"name": "John"`) {
		t.Error("Expected feature to contain request body")
	}
}

func TestGenerateScenario_PathParameters(t *testing.T) {
	endpoint := openapi.Endpoint{
		Method:      "GET",
		Path:        "/v1/users/:id",
		Description: "Get user by ID",
		Parameters: []openapi.Parameter{
			{Name: "id", In: "path", Description: "User ID", Required: true},
		},
	}

	scenario := generateScenario(endpoint)

	if !strings.Contains(scenario, "I have a valid id") {
		t.Error("Expected scenario to set up path parameter")
	}

	if !strings.Contains(scenario, `I GET to "/v1/users/:id"`) {
		t.Error("Expected scenario to include path in When step")
	}
}

func TestGenerateStepDefinitions_Go(t *testing.T) {
	stepDefs, err := GenerateStepDefinitions("go")
	if err != nil {
		t.Fatalf("GenerateStepDefinitions failed: %v", err)
	}

	if !strings.Contains(stepDefs, "package steps") {
		t.Error("Expected Go step definitions to start with package declaration")
	}

	if !strings.Contains(stepDefs, "APIContext") {
		t.Error("Expected Go step definitions to define APIContext")
	}

	if !strings.Contains(stepDefs, "func InitializeScenario") {
		t.Error("Expected Go step definitions to define InitializeScenario")
	}
}

func TestGenerateStepDefinitions_JavaScript(t *testing.T) {
	stepDefs, err := GenerateStepDefinitions("javascript")
	if err != nil {
		t.Fatalf("GenerateStepDefinitions failed: %v", err)
	}

	if !strings.Contains(stepDefs, "require('@cucumber/cucumber')") {
		t.Error("Expected JavaScript step definitions to import cucumber")
	}

	if !strings.Contains(stepDefs, "Given('the API is available'") {
		t.Error("Expected JavaScript step definitions to define Given steps")
	}
}

func TestGenerateStepDefinitions_Python(t *testing.T) {
	stepDefs, err := GenerateStepDefinitions("python")
	if err != nil {
		t.Fatalf("GenerateStepDefinitions failed: %v", err)
	}

	if !strings.Contains(stepDefs, "from behave import") {
		t.Error("Expected Python step definitions to import behave")
	}

	if !strings.Contains(stepDefs, "@given('the API is available')") {
		t.Error("Expected Python step definitions to define given decorators")
	}
}

func TestGenerateStepDefinitions_UnsupportedLanguage(t *testing.T) {
	_, err := GenerateStepDefinitions("ruby")
	if err == nil {
		t.Error("Expected error for unsupported language")
	}

	if !strings.Contains(err.Error(), "unsupported language") {
		t.Errorf("Expected 'unsupported language' error, got: %v", err)
	}
}

func TestGenerateCucumber_Integration(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create architecture.md
	archPath := filepath.Join(tmpDir, "architecture.md")
	archContent := `# Test API

## Endpoints

### GET /v1/health
Check health status.

### POST /v1/users
Create a user.
`
	if err := os.WriteFile(archPath, []byte(archContent), 0644); err != nil {
		t.Fatalf("Failed to write architecture.md: %v", err)
	}

	// Generate cucumber tests
	outputPath := filepath.Join(tmpDir, "tests")
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	err := GenerateCucumber(archPath, "test-api", outputPath, "go")
	if err != nil {
		t.Fatalf("GenerateCucumber failed: %v", err)
	}

	// Check feature file exists
	featurePath := filepath.Join(outputPath, "api.feature")
	if _, err := os.Stat(featurePath); os.IsNotExist(err) {
		t.Error("Feature file was not created")
	}

	// Check step definitions exist
	stepsPath := filepath.Join(outputPath, "steps.go")
	if _, err := os.Stat(stepsPath); os.IsNotExist(err) {
		t.Error("Step definitions file was not created")
	}

	// Verify feature file content
	featureContent, err := os.ReadFile(featurePath)
	if err != nil {
		t.Fatalf("Failed to read feature file: %v", err)
	}

	if !strings.Contains(string(featureContent), "Feature: Test Api API") {
		t.Error("Feature file doesn't contain expected feature name")
	}

	if !strings.Contains(string(featureContent), "Scenario: GET /v1/health") {
		t.Error("Feature file doesn't contain expected scenario")
	}
}

func TestGenerateCucumber_NoLanguage(t *testing.T) {
	tmpDir := t.TempDir()
	
	archPath := filepath.Join(tmpDir, "architecture.md")
	archContent := `# Test API

## Endpoints

### GET /v1/test
Test endpoint.
`
	if err := os.WriteFile(archPath, []byte(archContent), 0644); err != nil {
		t.Fatalf("Failed to write architecture.md: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "tests")
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Generate without language (no step definitions)
	err := GenerateCucumber(archPath, "test-api", outputPath, "")
	if err != nil {
		t.Fatalf("GenerateCucumber failed: %v", err)
	}

	// Check only feature file exists
	featurePath := filepath.Join(outputPath, "api.feature")
	if _, err := os.Stat(featurePath); os.IsNotExist(err) {
		t.Error("Feature file was not created")
	}

	// Step definitions should not exist
	stepsPath := filepath.Join(outputPath, "steps.go")
	if _, err := os.Stat(stepsPath); err == nil {
		t.Error("Step definitions file should not be created without language")
	}
}

func TestHasPathParams(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/v1/users/:id", true},
		{"/v1/users/{id}", true},
		{"/v1/users", false},
		{"/v1/posts/:postId/comments/:commentId", true},
	}

	for _, tt := range tests {
		result := hasPathParams(tt.path)
		if result != tt.expected {
			t.Errorf("hasPathParams(%s) = %v, expected %v", tt.path, result, tt.expected)
		}
	}
}

func TestExtractPathParams(t *testing.T) {
	tests := []struct {
		path     string
		expected []string
	}{
		{"/v1/users/:id", []string{"id"}},
		{"/v1/users/{id}", []string{"id"}},
		{"/v1/posts/:postId/comments/:commentId", []string{"postId", "commentId"}},
		{"/v1/users", []string{}},
	}

	for _, tt := range tests {
		result := extractPathParams(tt.path)
		if len(result) != len(tt.expected) {
			t.Errorf("extractPathParams(%s) returned %d params, expected %d", 
				tt.path, len(result), len(tt.expected))
			continue
		}
		for i, param := range result {
			if param != tt.expected[i] {
				t.Errorf("extractPathParams(%s)[%d] = %s, expected %s", 
					tt.path, i, param, tt.expected[i])
			}
		}
	}
}

func TestGetExpectedStatus(t *testing.T) {
	tests := []struct {
		method   string
		expected string
	}{
		{"GET", "200"},
		{"POST", "201"},
		{"PUT", "200"},
		{"DELETE", "204"},
		{"PATCH", "200"},
	}

	for _, tt := range tests {
		result := getExpectedStatus(tt.method)
		if result != tt.expected {
			t.Errorf("getExpectedStatus(%s) = %s, expected %s", 
				tt.method, result, tt.expected)
		}
	}
}
