package foundation

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DefaultConfig is the default neev.yaml configuration
type DefaultConfig struct {
	Version     string `yaml:"version"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// Initialize creates the .neev directory structure and default config
func Initialize(cwd string) error {
	neevPath := filepath.Join(cwd, RootDir)

	// Check if .neev already exists
	if _, err := os.Stat(neevPath); err == nil {
		return fmt.Errorf(".neev directory already exists at %s", neevPath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error checking .neev directory: %w", err)
	}

	// Create .neev directory
	if err := os.MkdirAll(neevPath, 0755); err != nil {
		return fmt.Errorf("failed to create .neev directory: %w", err)
	}

	// Create .neev/blueprints directory
	blueprintsPath := filepath.Join(neevPath, BlueprintsDir)
	if err := os.MkdirAll(blueprintsPath, 0755); err != nil {
		return fmt.Errorf("failed to create blueprints directory: %w", err)
	}

	// Create .neev/foundation directory
	foundationPath := filepath.Join(neevPath, FoundationDir)
	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		return fmt.Errorf("failed to create foundation directory: %w", err)
	}

	// Create default neev.yaml
	defaultConfig := DefaultConfig{
		Version:     "1.0.0",
		Name:        "My Project",
		Description: "A Neev-managed project",
	}

	configPath := filepath.Join(neevPath, ConfigFile)
	configData, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %w", err)
	}

	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		return fmt.Errorf("failed to write neev.yaml: %w", err)
	}

	return nil
}
