package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/neev-kit/neev/core/remotes"
	"gopkg.in/yaml.v3"
)

// Config represents the neev.yaml configuration structure
type Config struct {
	ProjectName    string           `yaml:"project_name"`
	IgnoreDirs     []string         `yaml:"ignore_dirs"`
	FoundationPath string           `yaml:"foundation_path"`
	Remotes        []remotes.Remote `yaml:"remotes,omitempty"`
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		ProjectName:    "My App",
		FoundationPath: ".neev",
		IgnoreDirs: []string{
			"node_modules",
			"dist",
			"build",
			"vendor",
			".git",
			".env",
			"bin",
			"obj",
			".idea",
			".vscode",
			"target",
		},
	}
}

// LoadConfig loads the neev.yaml configuration from the given directory
func LoadConfig(cwd string) (*Config, error) {
	configPath := filepath.Join(cwd, "neev.yaml")

	// If config doesn't exist, return defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read neev.yaml: %w", err)
	}

	// Parse YAML
	cfg := DefaultConfig() // Start with defaults
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse neev.yaml: %w", err)
	}

	// Validate
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// SaveConfig writes the configuration to neev.yaml
func SaveConfig(cwd string, cfg *Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	configPath := filepath.Join(cwd, "neev.yaml")
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write neev.yaml: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.ProjectName == "" {
		return fmt.Errorf("project_name cannot be empty")
	}

	if c.FoundationPath == "" {
		return fmt.Errorf("foundation_path cannot be empty")
	}

	// Ensure foundation_path is a relative path
	if filepath.IsAbs(c.FoundationPath) {
		return fmt.Errorf("foundation_path must be a relative path, got: %s", c.FoundationPath)
	}

	// Validate remotes
	remoteNames := make(map[string]bool)
	for _, remote := range c.Remotes {
		if remote.Name == "" {
			return fmt.Errorf("remote name cannot be empty")
		}
		// Validate remote name doesn't contain path separators or traversal sequences
		if remote.Name != filepath.Base(remote.Name) ||
			strings.Contains(remote.Name, "..") ||
			strings.ContainsAny(remote.Name, `/\`) {
			return fmt.Errorf("invalid remote name '%s': must be a simple name without path separators", remote.Name)
		}
		if remote.Path == "" {
			return fmt.Errorf("remote path cannot be empty for remote '%s'", remote.Name)
		}
		if remoteNames[remote.Name] {
			return fmt.Errorf("duplicate remote name: %s", remote.Name)
		}
		remoteNames[remote.Name] = true
	}

	return nil
}

// GetIgnoreDirs returns the list of directories to ignore, including defaults
func (c *Config) GetIgnoreDirs() map[string]bool {
	ignoredMap := make(map[string]bool)
	for _, dir := range c.IgnoreDirs {
		ignoredMap[dir] = true
	}
	return ignoredMap
}
