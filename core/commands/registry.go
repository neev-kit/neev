package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Registry manages command registration and generation
type Registry struct {
	projectName string
	projectPath string
	registry    *CommandRegistry
}

// NewRegistry creates a new command registry
func NewRegistry(projectName, projectPath string) *Registry {
	registry := &CommandRegistry{
		Version:     "1.0.0",
		ProjectName: projectName,
		Commands:    DefaultCommands,
		ToolSupport: make(map[string]bool),
	}

	for toolID := range DefaultTools {
		registry.ToolSupport[toolID] = true
	}

	return &Registry{
		projectName: projectName,
		projectPath: projectPath,
		registry:    registry,
	}
}

// SaveRegistry persists the command registry to disk
func (r *Registry) SaveRegistry() error {
	commandsDir := filepath.Join(r.projectPath, ".neev", "commands")
	if err := os.MkdirAll(commandsDir, 0755); err != nil {
		return fmt.Errorf("failed to create commands directory: %w", err)
	}

	registryPath := filepath.Join(commandsDir, "registry.yaml")
	data, err := yaml.Marshal(r.registry)
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	if err := os.WriteFile(registryPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write registry file: %w", err)
	}

	return nil
}

// GenerateCursorConfig generates Cursor-specific command configuration
func (r *Registry) GenerateCursorConfig() (string, error) {
	commands := make([]map[string]interface{}, 0)

	for _, cmd := range r.registry.Commands {
		commands = append(commands, map[string]interface{}{
			"id":          cmd.ID,
			"name":        cmd.Name,
			"description": cmd.Description,
			"prompt":      cmd.Prompt,
			"icon":        cmd.Icon,
			"category":    cmd.Category,
		})
	}

	config := map[string]interface{}{
		"version":  "1.0.0",
		"project":  r.projectName,
		"commands": commands,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal Cursor config: %w", err)
	}

	return string(data), nil
}

// GenerateVSCodeConfig generates VS Code command palette integration
func (r *Registry) GenerateVSCodeConfig() (string, error) {
	commands := make([]map[string]interface{}, 0)

	for _, cmd := range r.registry.Commands {
		commandID := strings.ReplaceAll(cmd.ID, ":", ".")

		commands = append(commands, map[string]interface{}{
			"command":     "neev." + commandID,
			"title":       fmt.Sprintf("Neev: %s", cmd.Name),
			"description": cmd.Description,
			"category":    "Neev",
		})
	}

	config := map[string]interface{}{
		"contributes": map[string]interface{}{
			"commands": commands,
		},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal VS Code config: %w", err)
	}

	return string(data), nil
}

// GenerateAllManifests creates all tool-specific configuration files
func (r *Registry) GenerateAllManifests() error {
	if err := r.SaveRegistry(); err != nil {
		return err
	}

	// Generate Cursor config
	if tool, ok := DefaultTools["cursor"]; ok {
		config, err := r.GenerateCursorConfig()
		if err != nil {
			return err
		}

		configPath := filepath.Join(r.projectPath, tool.ConfigPath)
		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			return err
		}

		if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
			return fmt.Errorf("failed to write Cursor config: %w", err)
		}
	}

	// Generate VS Code config
	if tool, ok := DefaultTools["vscode"]; ok {
		config, err := r.GenerateVSCodeConfig()
		if err != nil {
			return err
		}

		configPath := filepath.Join(r.projectPath, tool.ConfigPath)
		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			return err
		}

		if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
			return fmt.Errorf("failed to write VS Code config: %w", err)
		}
	}

	return nil
}

// ListCommands returns all registered commands
func (r *Registry) ListCommands() map[string]Command {
	return r.registry.Commands
}

// GetCommand retrieves a specific command
func (r *Registry) GetCommand(id string) (Command, error) {
	if cmd, ok := r.registry.Commands[id]; ok {
		return cmd, nil
	}
	return Command{}, fmt.Errorf("command not found: %s", id)
}

// AddCommand adds a new command to the registry
func (r *Registry) AddCommand(cmd Command) error {
	if _, exists := r.registry.Commands[cmd.ID]; exists {
		return fmt.Errorf("command already exists: %s", cmd.ID)
	}
	r.registry.Commands[cmd.ID] = cmd
	return nil
}
