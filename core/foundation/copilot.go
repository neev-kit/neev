package foundation

import (
	"encoding/json"
	"fmt"
)

// SlashCommandManifest represents GitHub Copilot's slash command format
type SlashCommandManifest struct {
	Version      string                        `json:"version"`
	Commands     map[string]SlashCommandDef    `json:"commands"`
	Description  string                        `json:"description"`
	ProjectName  string                        `json:"project_name"`
}

// SlashCommandDef defines a single slash command for GitHub Copilot
type SlashCommandDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Prompt      string `json:"prompt"`
	Aliases     []string `json:"aliases,omitempty"`
	Context     string `json:"context,omitempty"`
}

// GenerateSlashCommandManifest creates GitHub Copilot-compatible slash command manifest
func GenerateSlashCommandManifest(projectName string) (string, error) {
	manifest := SlashCommandManifest{
		Version:     "1.0.0",
		ProjectName: projectName,
		Description: "Neev slash commands for spec-driven development",
		Commands: map[string]SlashCommandDef{
			"neev:bridge": {
				Name:        "bridge",
				Description: "Generate aggregated project context for AI",
				Prompt:      "Generate the project bridge context including architecture, blueprints, and specifications",
				Aliases:     []string{"bridge", "context"},
				Context:     "Use this to get full project context for complex implementation tasks",
			},
			"neev:draft": {
				Name:        "draft",
				Description: "Create a new blueprint",
				Prompt:      "Create a new blueprint for the following feature or component",
				Aliases:     []string{"draft", "blueprint"},
				Context:     "Use this to plan new features with intent, architecture, API spec, and security considerations",
			},
			"neev:inspect": {
				Name:        "inspect",
				Description: "Analyze project structure and find gaps",
				Prompt:      "Analyze the project structure for gaps and inconsistencies between specifications and implementation",
				Aliases:     []string{"inspect", "gaps", "drift"},
				Context:     "Use this to verify that your code matches the specifications",
			},
			"neev:cucumber": {
				Name:        "cucumber",
				Description: "Generate Cucumber/BDD test scaffolding",
				Prompt:      "Generate Cucumber/BDD tests and scenarios for the following API or feature",
				Aliases:     []string{"cucumber", "bdd", "tests"},
				Context:     "Use this to create behavior-driven test cases",
			},
			"neev:openapi": {
				Name:        "openapi",
				Description: "Generate OpenAPI specification",
				Prompt:      "Generate an OpenAPI specification for this blueprint or API design",
				Aliases:     []string{"openapi", "api-spec", "swagger"},
				Context:     "Use this to create formal API documentation",
			},
			"neev:handoff": {
				Name:        "handoff",
				Description: "Format context for AI handoff",
				Prompt:      "Prepare the project context and specifications for handing off to another AI agent or developer",
				Aliases:     []string{"handoff", "hand-off", "transfer"},
				Context:     "Use this to ensure continuity when switching between AI agents",
			},
		},
	}

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal slash command manifest: %w", err)
	}

	return string(data), nil
}

// GenerateCopilotChatInstructions creates detailed Copilot Chat instructions with slash command registration
func GenerateCopilotChatInstructions(projectName string) string {
	return `# GitHub Copilot Instructions for ` + projectName + `

This project uses Neev for spec-driven development.

## Neev Slash Commands

Use these commands in GitHub Copilot Chat to manage your spec-driven workflow:

### /neev:bridge
**Generate aggregated project context for AI**

Retrieve and summarize project structure, architecture, blueprints, and all relevant documentation. Use this to get full project context for understanding the system before making changes.

Example: "@Copilot /neev:bridge Show me the complete project context"

### /neev:draft
**Create a new blueprint**

Plan new features or components with intent, architecture, API specifications, and security considerations.

Example: "@Copilot /neev:draft Create a blueprint for user authentication"

### /neev:inspect
**Analyze project structure and find gaps**

Verify that your code matches the specifications. Identify inconsistencies and misalignments between specs and implementation.

Example: "@Copilot /neev:inspect Check if our implementation matches the specifications"

### /neev:cucumber
**Generate Cucumber/BDD test scaffolding**

Create behavior-driven test cases and scenarios for APIs and features.

Example: "@Copilot /neev:cucumber Generate BDD tests for the user API"

### /neev:openapi
**Generate OpenAPI specification**

Create formal API documentation from blueprint architecture and API design.

Example: "@Copilot /neev:openapi Generate OpenAPI spec for this blueprint"

### /neev:handoff
**Format context for AI handoff**

Prepare project context and specifications for transitioning between AI agents or developers.

Example: "@Copilot /neev:handoff Prepare context for handing off to another AI agent"

## How to Use

1. **In GitHub Copilot Chat**: Type any of the commands above (` + "`/neev:bridge`" + `, ` + "`/neev:draft`" + `, etc.)
2. **From Terminal**: Run commands directly (` + "`neev bridge`" + `, ` + "`neev draft`" + `, etc.)
3. **In Cursor IDE**: Use native slash commands with full IDE integration

## Development Workflow

1. **Establish Context**: Use ` + "`/neev:bridge`" + ` to understand the full project scope
2. **Plan Feature**: Use ` + "`/neev:draft`" + ` to create a blueprint with specifications
3. **Implement**: Follow the blueprint and use ` + "`/neev:cucumber`" + ` for tests
4. **Verify**: Use ` + "`/neev:inspect`" + ` to ensure spec compliance
5. **Document**: Use ` + "`/neev:openapi`" + ` for API documentation
6. **Handoff**: Use ` + "`/neev:handoff`" + ` when switching AI agents

## Terminal Commands

For direct CLI usage:

` + "```bash" + `
neev bridge       # Get full project context
neev draft        # Create new blueprint
neev inspect      # Check for specification drift
neev cucumber     # Generate BDD tests
neev openapi      # Generate API specification
neev handoff      # Prepare for handoff
` + "```" + `

## Learning More

- Read the foundation files in ` + "`.neev/foundation/`" + ` to understand project principles
- Check ` + "`.neev/commands/registry.yaml`" + ` for complete command registry
- Use ` + "`/neev:inspect`" + ` to verify spec compliance across your codebase
`
}
