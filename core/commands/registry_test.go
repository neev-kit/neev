package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	reg := NewRegistry("test-project", "/tmp/test")

	if reg.projectName != "test-project" {
		t.Errorf("expected project name 'test-project', got '%s'", reg.projectName)
	}

	if len(reg.registry.Commands) == 0 {
		t.Error("expected commands to be initialized")
	}

	if len(reg.registry.ToolSupport) == 0 {
		t.Error("expected tool support to be initialized")
	}
}

func TestSaveRegistry(t *testing.T) {
	tmpDir := t.TempDir()
	reg := NewRegistry("test-project", tmpDir)

	if err := reg.SaveRegistry(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	registryPath := filepath.Join(tmpDir, ".neev", "commands", "registry.yaml")
	if _, err := os.Stat(registryPath); os.IsNotExist(err) {
		t.Fatalf("registry file not created at %s", registryPath)
	}
}

func TestGenerateCursorConfig(t *testing.T) {
	reg := NewRegistry("test-project", "/tmp/test")

	config, err := reg.GenerateCursorConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(config), &data); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	if data["version"] != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %v", data["version"])
	}
}

func TestGenerateVSCodeConfig(t *testing.T) {
	reg := NewRegistry("test-project", "/tmp/test")

	config, err := reg.GenerateVSCodeConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(config), &data); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	_, ok := data["contributes"].(map[string]interface{})
	if !ok {
		t.Fatal("contributes not found in config")
	}
}

func TestGenerateAllManifests(t *testing.T) {
	tmpDir := t.TempDir()
	reg := NewRegistry("test-project", tmpDir)

	if err := reg.GenerateAllManifests(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	registryPath := filepath.Join(tmpDir, ".neev", "commands", "registry.yaml")
	if _, err := os.Stat(registryPath); os.IsNotExist(err) {
		t.Fatalf("registry file not created")
	}
}

func TestGetCommand(t *testing.T) {
	reg := NewRegistry("test-project", "/tmp/test")

	cmd, err := reg.GetCommand("neev:bridge")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cmd.ID != "neev:bridge" {
		t.Errorf("expected ID 'neev:bridge', got '%s'", cmd.ID)
	}
}

func TestAddCommand(t *testing.T) {
	reg := NewRegistry("test-project", "/tmp/test")

	newCmd := Command{
		ID:          "custom:test",
		Name:        "Test",
		Description: "Test command",
		Prompt:      "Test prompt",
	}

	if err := reg.AddCommand(newCmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cmd, err := reg.GetCommand("custom:test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cmd.ID != "custom:test" {
		t.Errorf("expected ID 'custom:test', got '%s'", cmd.ID)
	}
}
