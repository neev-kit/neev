package slash

// SlashCommand represents a slash command definition
type SlashCommand struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Prompt      string `yaml:"prompt"`
}

// AIToolConfig defines slash commands for a specific AI tool
type AIToolConfig struct {
	Name          string          `yaml:"name"`
	Enabled       bool            `yaml:"enabled"`
	SlashCommands []SlashCommand  `yaml:"slash_commands"`
}

// SlashCommandConfig contains all supported slash commands
type SlashCommandConfig struct {
	Version  string
	Commands map[string]SlashCommand
}

// SupportedAITools lists AI tools with native slash command support
var SupportedAITools = []string{
	"claude-code",
	"cursor",
	"codebuddy",
	"opencode",
	"qoder",
	"codex",
	"roocode",
}

// DefaultSlashCommands defines the default slash commands for Neev
var DefaultSlashCommands = map[string]SlashCommand{
	"bridge": {
		Name:        "bridge",
		Description: "Generate aggregated project context for AI",
		Prompt:      "Generate the project context for me to review",
	},
	"draft": {
		Name:        "draft",
		Description: "Create a new blueprint",
		Prompt:      "Create a new blueprint for this feature",
	},
	"inspect": {
		Name:        "inspect",
		Description: "Analyze project structure and find gaps",
		Prompt:      "Analyze the project structure and show me any gaps",
	},
	"cucumber": {
		Name:        "cucumber",
		Description: "Generate Cucumber/BDD test scaffolding",
		Prompt:      "Generate Cucumber tests for this API",
	},
	"openapi": {
		Name:        "openapi",
		Description: "Generate OpenAPI specification",
		Prompt:      "Generate an OpenAPI specification for this blueprint",
	},
	"handoff": {
		Name:        "handoff",
		Description: "Format context for AI handoff",
		Prompt:      "Prepare the context for handoff to an AI agent",
	},
}
