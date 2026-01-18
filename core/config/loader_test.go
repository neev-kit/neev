package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/neev-kit/neev/core/remotes"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}

	if cfg.ProjectName != "My App" {
		t.Errorf("Expected ProjectName 'My App', got '%s'", cfg.ProjectName)
	}

	if cfg.FoundationPath != ".neev" {
		t.Errorf("Expected FoundationPath '.neev', got '%s'", cfg.FoundationPath)
	}

	if len(cfg.IgnoreDirs) == 0 {
		t.Error("Expected default IgnoreDirs to be populated")
	}
}

func TestValidate(t *testing.T) {
	cfg := &Config{
		ProjectName:    "test",
		FoundationPath: ".neev",
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
}

func TestValidateEmptyProjectName(t *testing.T) {
	cfg := &Config{
		ProjectName:    "",
		FoundationPath: ".neev",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("Expected error for empty ProjectName")
	}
}

func TestValidateEmptyFoundationPath(t *testing.T) {
	cfg := &Config{
		ProjectName:    "test",
		FoundationPath: "",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("Expected error for empty FoundationPath")
	}
}

func TestValidateAbsoluteFoundationPath(t *testing.T) {
	cfg := &Config{
		ProjectName:    "test",
		FoundationPath: "/absolute/path",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("Expected error for absolute FoundationPath")
	}
}

func TestValidateRemotes(t *testing.T) {
	cfg := &Config{
		ProjectName:    "test",
		FoundationPath: ".neev",
		Remotes: []remotes.Remote{
			{Name: "origin", Path: "/path/to/origin"},
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
}

func TestValidateRemoteEmptyName(t *testing.T) {
	cfg := &Config{
		ProjectName:    "test",
		FoundationPath: ".neev",
		Remotes: []remotes.Remote{
			{Name: "", Path: "/path/to/origin"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Error("Expected error for empty remote name")
	}
}

func TestValidateRemoteEmptyPath(t *testing.T) {
	cfg := &Config{
		ProjectName:    "test",
		FoundationPath: ".neev",
		Remotes: []remotes.Remote{
			{Name: "origin", Path: ""},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Error("Expected error for empty remote path")
	}
}

func TestValidateRemoteDuplicateName(t *testing.T) {
	cfg := &Config{
		ProjectName:    "test",
		FoundationPath: ".neev",
		Remotes: []remotes.Remote{
			{Name: "origin", Path: "/path1"},
			{Name: "origin", Path: "/path2"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Error("Expected error for duplicate remote name")
	}
}

func TestValidateRemoteWithPathSeparator(t *testing.T) {
	cfg := &Config{
		ProjectName:    "test",
		FoundationPath: ".neev",
		Remotes: []remotes.Remote{
			{Name: "origin/test", Path: "/path/to/origin"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Error("Expected error for remote name with path separator")
	}
}

func TestValidateRemoteWithTraversal(t *testing.T) {
	cfg := &Config{
		ProjectName:    "test",
		FoundationPath: ".neev",
		Remotes: []remotes.Remote{
			{Name: "../origin", Path: "/path/to/origin"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Error("Expected error for remote name with traversal sequence")
	}
}

func TestLoadConfigNotExist(t *testing.T) {
	tmpDir := t.TempDir()

	cfg, err := LoadConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadConfig should not error when neev.yaml doesn't exist: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadConfig should return default config when file doesn't exist")
	}

	if cfg.ProjectName != "My App" {
		t.Errorf("Expected default ProjectName, got '%s'", cfg.ProjectName)
	}
}

func TestLoadConfigValid(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a valid neev.yaml
	configPath := filepath.Join(tmpDir, "neev.yaml")
	configContent := `project_name: TestProject
foundation_path: .neev
ignore_dirs:
  - node_modules
  - dist
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	cfg, err := LoadConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.ProjectName != "TestProject" {
		t.Errorf("Expected ProjectName 'TestProject', got '%s'", cfg.ProjectName)
	}
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()

	// Create an invalid neev.yaml
	configPath := filepath.Join(tmpDir, "neev.yaml")
	configContent := `project_name: TestProject
invalid yaml content [[[`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	_, err := LoadConfig(tmpDir)
	if err == nil {
		t.Error("LoadConfig should error on invalid YAML")
	}
}

func TestLoadConfigInvalidConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a config with invalid values
	configPath := filepath.Join(tmpDir, "neev.yaml")
	configContent := `project_name: ""
foundation_path: .neev
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	_, err := LoadConfig(tmpDir)
	if err == nil {
		t.Error("LoadConfig should error on invalid config")
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &Config{
		ProjectName:    "SaveTest",
		FoundationPath: ".neev",
		IgnoreDirs:     []string{"node_modules", "dist"},
	}

	err := SaveConfig(tmpDir, cfg)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify file exists
	configPath := filepath.Join(tmpDir, "neev.yaml")
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("Config file was not created: %v", err)
	}

	// Reload and verify
	loadedCfg, err := LoadConfig(tmpDir)
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	if loadedCfg.ProjectName != "SaveTest" {
		t.Errorf("Reloaded config has wrong ProjectName: %s", loadedCfg.ProjectName)
	}
}

func TestSaveConfigInvalid(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &Config{
		ProjectName:    "",
		FoundationPath: ".neev",
	}

	err := SaveConfig(tmpDir, cfg)
	if err == nil {
		t.Error("SaveConfig should error for invalid config")
	}
}

func TestGetIgnoreDirs(t *testing.T) {
	cfg := &Config{
		ProjectName:    "test",
		FoundationPath: ".neev",
		IgnoreDirs:     []string{"node_modules", "dist", "build"},
	}

	ignoredMap := cfg.GetIgnoreDirs()

	if !ignoredMap["node_modules"] {
		t.Error("Expected node_modules to be in ignored map")
	}

	if !ignoredMap["dist"] {
		t.Error("Expected dist to be in ignored map")
	}

	if !ignoredMap["build"] {
		t.Error("Expected build to be in ignored map")
	}

	if ignoredMap["src"] {
		t.Error("Unexpected src to be in ignored map")
	}
}

func TestLoadConfigWithRemotes(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a neev.yaml with remotes
	configPath := filepath.Join(tmpDir, "neev.yaml")
	configContent := `project_name: TestProject
foundation_path: .neev
remotes:
  - name: origin
    path: /path/to/origin
  - name: backup
    path: /path/to/backup
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	cfg, err := LoadConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if len(cfg.Remotes) != 2 {
		t.Errorf("Expected 2 remotes, got %d", len(cfg.Remotes))
	}

	if cfg.Remotes[0].Name != "origin" {
		t.Errorf("Expected first remote name 'origin', got '%s'", cfg.Remotes[0].Name)
	}
}

func TestLoadConfigReadError(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a directory instead of a file (will cause read error)
	configPath := filepath.Join(tmpDir, "neev.yaml")
	if err := os.Mkdir(configPath, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	_, err := LoadConfig(tmpDir)
	if err == nil {
		t.Error("LoadConfig should error when neev.yaml is a directory")
	}
}
