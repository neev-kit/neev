package foundation

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/neev-kit/neev/core/commands"
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

	// Generate Copilot instructions file
	copilotInstructions := generateCopilotInstructions(projectName)
	copilotPath := filepath.Join(cwd, ".github", "copilot-instructions.md")
	if err := os.MkdirAll(filepath.Dir(copilotPath), 0755); err != nil {
		return fmt.Errorf("failed to create .github directory: %w", err)
	}
	if err := os.WriteFile(copilotPath, []byte(copilotInstructions), 0644); err != nil {
		return fmt.Errorf("failed to write copilot instructions: %w", err)
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

// generateCopilotInstructions creates GitHub Copilot Chat instructions
func generateCopilotInstructions(projectName string) string {
	return `# GitHub Copilot Instructions for ` + projectName + `

This project uses Neev for spec-driven development.

## Development Guidelines

- Follow the architecture defined in foundation specifications
- Implement features according to blueprint intent and architecture
- Use ` + "`neev bridge`" + ` to get full context for complex tasks
- Run ` + "`neev inspect`" + ` to check for drift between specs and code

## Neev Slash Commands for Copilot Chat

You can use these slash commands in GitHub Copilot Chat:

### ` + "`/neev:bridge`" + `
Generate aggregated project context for AI. Retrieves and summarizes project structure, architecture, blueprints, and all relevant documentation.

**Usage:** ` + "`/neev:bridge`" + ` - Get full project context

### ` + "`/neev:draft`" + `
Create a new blueprint for planning features or components.

**Usage:** ` + "`/neev:draft`" + ` - Create a new blueprint

### ` + "`/neev:inspect`" + `
Analyze project structure and find gaps or inconsistencies between specifications and implementation.

**Usage:** ` + "`/neev:inspect`" + ` - Check for drift between specs and code

### ` + "`/neev:cucumber`" + `
Generate Cucumber/BDD test scaffolding and scenarios.

**Usage:** ` + "`/neev:cucumber`" + ` - Generate BDD tests for this feature

### ` + "`/neev:openapi`" + `
Generate OpenAPI specification from blueprint architecture and API design.

**Usage:** ` + "`/neev:openapi`" + ` - Generate API spec for this blueprint

### ` + "`/neev:handoff`" + `
Format context and specifications for AI agent handoff or team handover.

**Usage:** ` + "`/neev:handoff`" + ` - Prepare context for handoff

## How to Run Commands

In **VS Code with GitHub Copilot Chat**:
1. Open Copilot Chat (Cmd+Shift+I on Mac, Ctrl+Shift+I on Windows/Linux)
2. Type any of the slash commands above (` + "`/neev:bridge`" + `, ` + "`/neev:draft`" + `, etc.)
3. Copilot will execute the Neev CLI command and provide results

From **Terminal**:
` + "```bash" + `
neev bridge       # Generate project context
neev draft        # Create new blueprint
neev inspect      # Analyze project structure
neev cucumber     # Generate BDD tests
neev openapi      # Generate API spec
neev handoff      # Prepare for handoff
` + "```" + `
`
}

