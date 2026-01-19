package commands

// Command represents a slash command definition
type Command struct {
	ID          string                 `yaml:"id"`
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Prompt      string                 `yaml:"prompt"`
	Icon        string                 `yaml:"icon,omitempty"`
	Category    string                 `yaml:"category,omitempty"`
}

// CommandRegistry stores all registered commands
type CommandRegistry struct {
	Version     string            `yaml:"version"`
	ProjectName string            `yaml:"project_name"`
	Commands    map[string]Command `yaml:"commands"`
	ToolSupport map[string]bool    `yaml:"tool_support"`
}

// SupportedTool represents an IDE/AI tool with command support
type SupportedTool struct {
	ID             string
	Name           string
	CommandFormat  string
	ConfigPath     string
	SupportsNative bool
}

// DefaultTools lists all tools with potential command support
var DefaultTools = map[string]SupportedTool{
	"cursor": {
		ID:             "cursor",
		Name:           "Cursor",
		CommandFormat:  "cursor",
		ConfigPath:     ".cursor/commands.json",
		SupportsNative: true,
	},
	"vscode": {
		ID:             "vscode",
		Name:           "VS Code",
		CommandFormat:  "vscode",
		ConfigPath:     ".vscode/commands.json",
		SupportsNative: true,
	},
	"github-copilot": {
		ID:             "github-copilot",
		Name:           "GitHub Copilot",
		CommandFormat:  "vscode",
		ConfigPath:     ".vscode/commands.json",
		SupportsNative: true,
	},
}

// DefaultCommands defines the base set of Neev commands
var DefaultCommands = map[string]Command{
	"neev:bridge": {
		ID:          "neev:bridge",
		Name:        "Bridge",
		Description: "Generate aggregated project context for AI",
		Prompt:      "Generate the project context for me to review",
		Icon:        "🌉",
		Category:    "Context",
	},
	"neev:draft": {
		ID:          "neev:draft",
		Name:        "Draft",
		Description: "Create a new blueprint",
		Prompt:      "Create a new blueprint for this feature",
		Icon:        "📋",
		Category:    "Planning",
	},
	"neev:inspect": {
		ID:          "neev:inspect",
		Name:        "Inspect",
		Description: "Analyze project structure and find gaps",
		Prompt:      "Analyze the project structure and show me any gaps",
		Icon:        "🔍",
		Category:    "Analysis",
	},
	"neev:cucumber": {
		ID:          "neev:cucumber",
		Name:        "Cucumber",
		Description: "Generate Cucumber/BDD test scaffolding",
		Prompt:      "Generate Cucumber tests for this API",
		Icon:        "🥒",
		Category:    "Testing",
	},
	"neev:openapi": {
		ID:          "neev:openapi",
		Name:        "OpenAPI",
		Description: "Generate OpenAPI specification",
		Prompt:      "Generate an OpenAPI specification for this blueprint",
		Icon:        "📖",
		Category:    "Documentation",
	},
	"neev:handoff": {
		ID:          "neev:handoff",
		Name:        "Handoff",
		Description: "Format context for AI handoff",
		Prompt:      "Prepare the context for handoff to an AI agent",
		Icon:        "🤝",
		Category:    "Collaboration",
	},
}
