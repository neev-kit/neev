package foundation

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestInitialize_Success(t *testing.T) {
	tmpDir := t.TempDir()

	err := Initialize(tmpDir)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// Verify .neev directory exists
	neevPath := filepath.Join(tmpDir, RootDir)
	if _, err := os.Stat(neevPath); os.IsNotExist(err) {
		t.Errorf("Expected .neev directory to exist at %s", neevPath)
	}

	// Verify blueprints subdirectory
	blueprintsPath := filepath.Join(neevPath, BlueprintsDir)
	if _, err := os.Stat(blueprintsPath); os.IsNotExist(err) {
		t.Errorf("Expected blueprints directory to exist at %s", blueprintsPath)
	}

	// Verify foundation subdirectory
	foundationPath := filepath.Join(neevPath, FoundationDir)
	if _, err := os.Stat(foundationPath); os.IsNotExist(err) {
		t.Errorf("Expected foundation directory to exist at %s", foundationPath)
	}

	// Verify config file exists
	configPath := filepath.Join(neevPath, ConfigFile)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Expected config file to exist at %s", configPath)
	}
}

func TestInitialize_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Initialize once
	if err := Initialize(tmpDir); err != nil {
		t.Fatalf("First Initialize() failed: %v", err)
	}

	// Try to initialize again
	err := Initialize(tmpDir)
	if err == nil {
		t.Error("Expected error when .neev directory already exists")
	}

	if err.Error() != fmt.Sprintf(".neev directory already exists at %s", filepath.Join(tmpDir, RootDir)) {
		t.Errorf("Got unexpected error message: %v", err)
	}
}

func TestInitialize_PermissionDenied(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	tmpDir := t.TempDir()
	readOnlyDir := filepath.Join(tmpDir, "readonly")

	if err := os.Mkdir(readOnlyDir, 0555); err != nil {
		t.Fatalf("Failed to create read-only directory: %v", err)
	}
	defer os.Chmod(readOnlyDir, 0755)

	err := Initialize(readOnlyDir)
	if err == nil {
		t.Error("Expected error due to permission denied")
	}
}

func TestInitialize_FailsToCreateBlueprints(t *testing.T) {
	tmpDir := t.TempDir()
	neevPath := filepath.Join(tmpDir, RootDir)

	if err := os.MkdirAll(neevPath, 0755); err != nil {
		t.Fatalf("Failed to create .neev directory: %v", err)
	}

	// Create blueprints as a file (not a directory) to cause the second mkdir to fail
	blueprintsFile := filepath.Join(neevPath, BlueprintsDir)
	if err := os.WriteFile(blueprintsFile, []byte("file"), 0644); err != nil {
		t.Fatalf("Failed to create blueprints file: %v", err)
	}

	// Now .neev exists, so Initialize should fail on MkdirAll for blueprints
	err := Initialize(tmpDir)
	if err == nil {
		t.Error("Expected error when blueprints directory already exists as a file")
	}
}

func TestInitialize_BadDirectoryState(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, RootDir), []byte("file"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	err := Initialize(tmpDir)
	if err == nil {
		t.Error("Expected error when .neev exists as a file")
	}
}

func TestInitialize_FailsToCreateFoundation(t *testing.T) {
	tmpDir := t.TempDir()
	neevPath := filepath.Join(tmpDir, RootDir)

	if err := os.MkdirAll(neevPath, 0755); err != nil {
		t.Fatalf("Failed to create .neev directory: %v", err)
	}

	// Create foundation as a file (not a directory) to cause the third mkdir to fail
	foundationFile := filepath.Join(neevPath, FoundationDir)
	if err := os.WriteFile(foundationFile, []byte("file"), 0644); err != nil {
		t.Fatalf("Failed to create foundation file: %v", err)
	}

	err := Initialize(tmpDir)
	if err == nil {
		t.Error("Expected error when foundation directory already exists as a file")
	}
}

func TestInitialize_FailsToWriteConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	neevPath := filepath.Join(tmpDir, RootDir)

	if err := os.MkdirAll(neevPath, 0755); err != nil {
		t.Fatalf("Failed to create .neev directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(neevPath, BlueprintsDir), 0755); err != nil {
		t.Fatalf("Failed to create blueprints directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(neevPath, FoundationDir), 0755); err != nil {
		t.Fatalf("Failed to create foundation directory: %v", err)
	}

	// Make the directory read-only to prevent file writing
	if err := os.Chmod(neevPath, 0555); err != nil {
		t.Fatalf("Failed to make directory read-only: %v", err)
	}
	defer os.Chmod(neevPath, 0755)

	err := Initialize(tmpDir)
	if err == nil {
		t.Error("Expected error when unable to write config file")
	}
}

func TestDefaultConfig_Fields(t *testing.T) {
	config := DefaultConfig{
		Version:     "1.0.0",
		Name:        "Test Project",
		Description: "A test project",
	}

	if config.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", config.Version)
	}
	if config.Name != "Test Project" {
		t.Errorf("Expected name 'Test Project', got %s", config.Name)
	}
	if config.Description != "A test project" {
		t.Errorf("Expected description 'A test project', got %s", config.Description)
	}
}

func TestInitialize_ConfigContent(t *testing.T) {
	tmpDir := t.TempDir()

	err := Initialize(tmpDir)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	configPath := filepath.Join(tmpDir, RootDir, ConfigFile)
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	contentStr := string(content)
	if contentStr == "" {
		t.Error("Config file is empty")
	}

	if !contains(contentStr, "version") {
		t.Error("Config file should contain 'version'")
	}
	if !contains(contentStr, "name") {
		t.Error("Config file should contain 'name'")
	}
	if !contains(contentStr, "description") {
		t.Error("Config file should contain 'description'")
	}
}

func TestProject_Struct(t *testing.T) {
	project := Project{
		RepoRoot: "/path/to/repo",
	}

	if project.RepoRoot != "/path/to/repo" {
		t.Errorf("Expected RepoRoot to be '/path/to/repo', got %s", project.RepoRoot)
	}
}

func TestInitialize_DirectoriesHaveCorrectPermissions(t *testing.T) {
	tmpDir := t.TempDir()

	err := Initialize(tmpDir)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	neevPath := filepath.Join(tmpDir, RootDir)
	info, err := os.Stat(neevPath)
	if err != nil {
		t.Fatalf("Failed to stat .neev directory: %v", err)
	}

	if !info.IsDir() {
		t.Error("Expected .neev to be a directory")
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
