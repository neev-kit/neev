package foundation

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/neev-kit/neev/core/commands"
	"github.com/neev-kit/neev/core/slash"
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

	// Create foundation files with templates
	foundationFiles := map[string]string{
		"stack.md":      "# Technology Stack\n\nDescribe the technologies used in your project (e.g., \"We use Go, PostgreSQL, Redis\")",
		"principles.md": "# Design Principles\n\nDocument your core design principles (e.g., \"Security first, simplicity second\")",
		"patterns.md":   "# Patterns & Practices\n\nOutline your architectural patterns and practices (e.g., \"Repository pattern, dependency injection\")",
	}

	for file, content := range foundationFiles {
		filePath := filepath.Join(foundationPath, file)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create foundation file %s: %w", filePath, err)
		}
	}

	// Generate command manifests for AI tools
	projectName := filepath.Base(cwd)
	registry := commands.NewRegistry(projectName, cwd)
	if err := registry.GenerateAllManifests(); err != nil {
		return fmt.Errorf("failed to generate command manifests: %w", err)
	}

	// Generate Copilot instructions and slash command manifest
	copilotInstructions := GenerateCopilotChatInstructions(projectName)
	copilotPath := filepath.Join(cwd, ".github", "copilot-instructions.md")
	if err := os.MkdirAll(filepath.Dir(copilotPath), 0755); err != nil {
		return fmt.Errorf("failed to create .github directory: %w", err)
	}
	if err := os.WriteFile(copilotPath, []byte(copilotInstructions), 0644); err != nil {
		return fmt.Errorf("failed to write copilot instructions: %w", err)
	}

	// Generate slash command manifest for Copilot
	slashCommandManifest, err := GenerateSlashCommandManifest(projectName)
	if err != nil {
		return fmt.Errorf("failed to generate slash command manifest: %w", err)
	}
	slashCommandPath := filepath.Join(cwd, ".github", "slash-commands.json")
	if err := os.WriteFile(slashCommandPath, []byte(slashCommandManifest), 0644); err != nil {
		return fmt.Errorf("failed to write slash command manifest: %w", err)
	}

	// Generate GitHub Copilot prompt files
	copilotPrompts := slash.GenerateGitHubCopilotPrompts(projectName)
	copilotPromptsBasePath := filepath.Join(cwd, ".github", "prompts", "neev")
	if err := os.MkdirAll(copilotPromptsBasePath, 0755); err != nil {
		return fmt.Errorf("failed to create .github/prompts/neev directory: %w", err)
	}
	for fileName, content := range copilotPrompts {
		filePath := filepath.Join(copilotPromptsBasePath, fileName)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write GitHub Copilot prompt file %s: %w", fileName, err)
		}
	}

	// Generate Claude Code slash command files
	claudeCommands := slash.GenerateClaudeSlashCommands(projectName)
	claudeBasePath := filepath.Join(cwd, ".claude", "commands", "neev")
	if err := os.MkdirAll(claudeBasePath, 0755); err != nil {
		return fmt.Errorf("failed to create .claude/commands/neev directory: %w", err)
	}
	for fileName, content := range claudeCommands {
		filePath := filepath.Join(claudeBasePath, fileName)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write Claude command file %s: %w", fileName, err)
		}
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
