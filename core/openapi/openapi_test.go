package openapi

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseArchitecture(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	archPath := filepath.Join(tmpDir, "architecture.md")
	
	content := `# Test API Architecture

## Endpoints

### GET /v1/users
List all users.

**Query Parameters:**
- ` + "`page`" + ` (default: 1)
- ` + "`limit`" + ` (default: 20)

### POST /v1/users
Create a new user.

**Request:**
` + "```json" + `
{
  "name": "John Doe",
  "email": "john@example.com"
}
` + "```" + `

### GET /v1/users/:id
Get user by ID.

### DELETE /v1/users/:id
Delete a user.
`
	
	if err := os.WriteFile(archPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	endpoints, err := ParseArchitecture(archPath)
	if err != nil {
		t.Fatalf("ParseArchitecture failed: %v", err)
	}
	
	if len(endpoints) != 4 {
		t.Errorf("Expected 4 endpoints, got %d", len(endpoints))
	}
	
	// Check first endpoint
	if endpoints[0].Method != "GET" {
		t.Errorf("Expected method GET, got %s", endpoints[0].Method)
	}
	if endpoints[0].Path != "/v1/users" {
		t.Errorf("Expected path /v1/users, got %s", endpoints[0].Path)
	}
	if endpoints[0].Description != "List all users." {
		t.Errorf("Expected description 'List all users.', got %s", endpoints[0].Description)
	}
	
	// Check parameters
	if len(endpoints[0].Parameters) < 2 {
		t.Errorf("Expected at least 2 parameters, got %d", len(endpoints[0].Parameters))
	}
	
	// Check POST endpoint
	if endpoints[1].Method != "POST" {
		t.Errorf("Expected method POST, got %s", endpoints[1].Method)
	}
	if endpoints[1].Request == "" {
		t.Error("Expected request body for POST, got empty")
	}
	
	// Check path parameters
	if len(endpoints[2].Parameters) < 1 {
		t.Errorf("Expected at least 1 parameter for path params, got %d", len(endpoints[2].Parameters))
	}
	
	// Check GET endpoint with path param
	if endpoints[2].Method != "GET" || endpoints[2].Path != "/v1/users/:id" {
		t.Errorf("Expected GET /v1/users/:id, got %s %s", endpoints[2].Method, endpoints[2].Path)
	}
	
	// Path parameter should be detected
	hasPathParam := false
	for _, param := range endpoints[2].Parameters {
		if param.Name == "id" && param.In == "path" {
			hasPathParam = true
			break
		}
	}
	if !hasPathParam {
		t.Error("Expected path parameter 'id' to be detected")
	}
}

func TestParseArchitecture_NoEndpoints(t *testing.T) {
	tmpDir := t.TempDir()
	archPath := filepath.Join(tmpDir, "architecture.md")
	
	content := `# Test API Architecture

This is just documentation with no endpoints.
`
	
	if err := os.WriteFile(archPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	endpoints, err := ParseArchitecture(archPath)
	if err != nil {
		t.Fatalf("ParseArchitecture failed: %v", err)
	}
	
	if len(endpoints) != 0 {
		t.Errorf("Expected 0 endpoints, got %d", len(endpoints))
	}
}

func TestParseArchitecture_FileNotFound(t *testing.T) {
	_, err := ParseArchitecture("/nonexistent/file.md")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
	if !strings.Contains(err.Error(), "failed to open architecture file") {
		t.Errorf("Expected 'failed to open architecture file' error, got: %v", err)
	}
}

func TestGenerateOpenAPISpec(t *testing.T) {
	endpoints := []Endpoint{
		{
			Method:      "GET",
			Path:        "/v1/users",
			Description: "List users",
			Parameters: []Parameter{
				{Name: "page", In: "query", Description: "Page number", Required: false},
			},
		},
		{
			Method:      "POST",
			Path:        "/v1/users",
			Description: "Create user",
			Request:     `{"name": "test"}`,
		},
		{
			Method:      "GET",
			Path:        "/v1/users/:id",
			Description: "Get user by ID",
			Parameters: []Parameter{
				{Name: "id", In: "path", Description: "User ID", Required: true},
			},
		},
	}
	
	spec, err := GenerateOpenAPISpec(endpoints, "test-api")
	if err != nil {
		t.Fatalf("GenerateOpenAPISpec failed: %v", err)
	}
	
	if spec.OpenAPI != "3.1.0" {
		t.Errorf("Expected OpenAPI version 3.1.0, got %s", spec.OpenAPI)
	}
	
	if spec.Info.Title != "Test Api" {
		t.Errorf("Expected title 'Test Api', got %s", spec.Info.Title)
	}
	
	if len(spec.Paths) != 2 {
		t.Errorf("Expected 2 paths, got %d", len(spec.Paths))
	}
	
	// Check /v1/users path
	usersPath, exists := spec.Paths["/v1/users"]
	if !exists {
		t.Error("Expected /v1/users path to exist")
	}
	if usersPath.Get == nil {
		t.Error("Expected GET operation for /v1/users")
	}
	if usersPath.Post == nil {
		t.Error("Expected POST operation for /v1/users")
	}
	
	// Check path parameters are converted
	userByIdPath, exists := spec.Paths["/v1/users/{id}"]
	if !exists {
		t.Error("Expected /v1/users/{id} path to exist (converted from :id)")
	}
	if userByIdPath.Get == nil {
		t.Error("Expected GET operation for /v1/users/{id}")
	}
	
	// Check parameter exists
	if len(userByIdPath.Get.Parameters) == 0 {
		t.Error("Expected path parameters for GET /v1/users/{id}")
	}
}

func TestGenerateYAML(t *testing.T) {
	spec := &OpenAPISpec{
		OpenAPI: "3.1.0",
		Info: Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: make(map[string]PathItem),
	}
	
	yamlData, err := GenerateYAML(spec)
	if err != nil {
		t.Fatalf("GenerateYAML failed: %v", err)
	}
	
	yamlStr := string(yamlData)
	if !strings.Contains(yamlStr, "openapi: 3.1.0") {
		t.Error("Expected YAML to contain 'openapi: 3.1.0'")
	}
	if !strings.Contains(yamlStr, "title: Test API") {
		t.Error("Expected YAML to contain 'title: Test API'")
	}
}

func TestGenerateOpenAPI_Integration(t *testing.T) {
	tmpDir := t.TempDir()
	archPath := filepath.Join(tmpDir, "architecture.md")
	
	content := `# API Architecture

## Endpoints

### GET /v1/health
Check API health status.

**Response (200):**
` + "```json" + `
{
  "status": "ok"
}
` + "```" + `
`
	
	if err := os.WriteFile(archPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	yamlData, err := GenerateOpenAPI(archPath, "health-api")
	if err != nil {
		t.Fatalf("GenerateOpenAPI failed: %v", err)
	}
	
	yamlStr := string(yamlData)
	if !strings.Contains(yamlStr, "openapi: 3.1.0") {
		t.Error("Expected YAML to contain 'openapi: 3.1.0'")
	}
	if !strings.Contains(yamlStr, "/v1/health") {
		t.Error("Expected YAML to contain '/v1/health' path")
	}
}

func TestGenerateOpenAPI_NoEndpoints(t *testing.T) {
	tmpDir := t.TempDir()
	archPath := filepath.Join(tmpDir, "architecture.md")
	
	content := `# API Architecture

Just some documentation.
`
	
	if err := os.WriteFile(archPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	_, err := GenerateOpenAPI(archPath, "test-api")
	if err == nil {
		t.Error("Expected error for no endpoints")
	}
	if !strings.Contains(err.Error(), "no API endpoints found") {
		t.Errorf("Expected 'no API endpoints found' error, got: %v", err)
	}
}

func TestConvertPathParams(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/v1/users/:id", "/v1/users/{id}"},
		{"/v1/posts/:postId/comments/:commentId", "/v1/posts/{postId}/comments/{commentId}"},
		{"/v1/users", "/v1/users"},
	}
	
	for _, tt := range tests {
		result := convertPathParams(tt.input)
		if result != tt.expected {
			t.Errorf("convertPathParams(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestFormatTitle(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test-api", "Test Api"},
		{"user-management", "User Management"},
		{"api", "Api"},
	}
	
	for _, tt := range tests {
		result := formatTitle(tt.input)
		if result != tt.expected {
			t.Errorf("formatTitle(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}
