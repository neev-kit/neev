package inspect

import (
	"os"
	"testing"
)

func TestGoDetector_Detect(t *testing.T) {
	detector := &GoDetector{}
	
	tests := []struct {
		path     string
		expected bool
	}{
		{"main.go", true},
		{"test.GO", true},
		{"handler.py", false},
		{"app.js", false},
	}
	
	for _, tt := range tests {
		result := detector.Detect(tt.path)
		if result != tt.expected {
			t.Errorf("Detect(%q) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}

func TestGoDetector_ExtractEndpoints(t *testing.T) {
	detector := &GoDetector{}
	
	code := `package main

import "github.com/gin-gonic/gin"

func setupRoutes(r *gin.Engine) {
	r.GET("/api/users", handlers.ListUsers)
	r.POST("/api/users", handlers.CreateUser)
	r.DELETE("/api/users/:id", handlers.DeleteUser)
}
`
	
	endpoints, err := detector.ExtractEndpoints("main.go", []byte(code))
	if err != nil {
		t.Fatalf("ExtractEndpoints failed: %v", err)
	}
	
	if len(endpoints) < 3 {
		t.Errorf("Expected at least 3 endpoints, got %d", len(endpoints))
	}
	
	// Check first endpoint
	foundGetUsers := false
	for _, ep := range endpoints {
		if ep.Method == "GET" && ep.Path == "/api/users" {
			foundGetUsers = true
			break
		}
	}
	
	if !foundGetUsers {
		t.Errorf("Expected to find GET /api/users endpoint")
	}
}

func TestGoDetector_ExtractFunctions(t *testing.T) {
	detector := &GoDetector{}
	
	code := `package main

func ListUsers(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func CreateUser(ctx context.Context, user *User) (*User, error) {
	return user, nil
}
`
	
	functions, err := detector.ExtractFunctions("handlers.go", []byte(code))
	if err != nil {
		t.Fatalf("ExtractFunctions failed: %v", err)
	}
	
	if len(functions) != 2 {
		t.Errorf("Expected 2 functions, got %d", len(functions))
	}
	
	// Check first function
	if functions[0].Name != "ListUsers" {
		t.Errorf("First function name: got %s, want ListUsers", functions[0].Name)
	}
	
	if functions[0].Visibility != "public" {
		t.Errorf("First function visibility: got %s, want public", functions[0].Visibility)
	}
}

func TestPythonDetector_Detect(t *testing.T) {
	detector := &PythonDetector{}
	
	tests := []struct {
		path     string
		expected bool
	}{
		{"app.py", true},
		{"test.PY", true},
		{"handler.go", false},
		{"app.js", false},
	}
	
	for _, tt := range tests {
		result := detector.Detect(tt.path)
		if result != tt.expected {
			t.Errorf("Detect(%q) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}

func TestPythonDetector_ExtractEndpoints(t *testing.T) {
	detector := &PythonDetector{}
	
	code := `from fastapi import FastAPI

app = FastAPI()

@app.get("/api/users")
async def list_users():
    return []

@app.post("/api/users")
async def create_user(user: User):
    return user
`
	
	endpoints, err := detector.ExtractEndpoints("main.py", []byte(code))
	if err != nil {
		t.Fatalf("ExtractEndpoints failed: %v", err)
	}
	
	if len(endpoints) < 2 {
		t.Errorf("Expected at least 2 endpoints, got %d", len(endpoints))
	}
	
	// Check for GET endpoint
	foundGetUsers := false
	for _, ep := range endpoints {
		if ep.Method == "GET" && ep.Path == "/api/users" {
			foundGetUsers = true
			break
		}
	}
	
	if !foundGetUsers {
		t.Errorf("Expected to find GET /api/users endpoint")
	}
}

func TestJavaScriptDetector_Detect(t *testing.T) {
	detector := &JavaScriptDetector{}
	
	tests := []struct {
		path     string
		expected bool
	}{
		{"app.js", true},
		{"server.ts", true},
		{"component.jsx", true},
		{"app.tsx", true},
		{"handler.go", false},
		{"app.py", false},
	}
	
	for _, tt := range tests {
		result := detector.Detect(tt.path)
		if result != tt.expected {
			t.Errorf("Detect(%q) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}

func TestJavaScriptDetector_ExtractEndpoints(t *testing.T) {
	detector := &JavaScriptDetector{}
	
	code := `const express = require('express');
const app = express();

app.get('/api/users', listUsers);
app.post('/api/users', createUser);
app.delete('/api/users/:id', deleteUser);
`
	
	endpoints, err := detector.ExtractEndpoints("app.js", []byte(code))
	if err != nil {
		t.Fatalf("ExtractEndpoints failed: %v", err)
	}
	
	if len(endpoints) != 3 {
		t.Errorf("Expected 3 endpoints, got %d", len(endpoints))
	}
	
	// Check first endpoint
	if endpoints[0].Method != "GET" || endpoints[0].Path != "/api/users" {
		t.Errorf("First endpoint: got %s %s, want GET /api/users", endpoints[0].Method, endpoints[0].Path)
	}
}

func TestPolyglotAnalyzer_DetectLanguages(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create test files
	files := map[string]string{
		"main.go":    "package main",
		"app.py":     "print('hello')",
		"server.js":  "console.log('hello')",
		"test.java":  "public class Test {}",
	}
	
	for name, content := range files {
		path := tmpDir + "/" + name
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}
	
	analyzer := NewPolyglotAnalyzer()
	analyzer.RegisterDetector(&GoDetector{})
	analyzer.RegisterDetector(&PythonDetector{})
	analyzer.RegisterDetector(&JavaScriptDetector{})
	analyzer.RegisterDetector(&JavaDetector{})
	
	languages, err := analyzer.DetectLanguages(tmpDir, map[string]bool{})
	if err != nil {
		t.Fatalf("DetectLanguages failed: %v", err)
	}
	
	expectedLangs := map[string]int{
		"go":         1,
		"python":     1,
		"javascript": 1,
		"java":       1,
	}
	
	for lang, expectedCount := range expectedLangs {
		if count, ok := languages[lang]; !ok || count != expectedCount {
			t.Errorf("Language %s: got %d files, want %d", lang, count, expectedCount)
		}
	}
}

