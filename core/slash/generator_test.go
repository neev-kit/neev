package slash

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGenerateAgentsMD(t *testing.T) {
	result := GenerateAgentsMD([]string{"claude-code", "cursor"}, "TestProject")

	if !strings.Contains(result, "# Neev AI Agent Instructions") {
		t.Error("Expected title in AGENTS.md")
	}

	if !strings.Contains(result, "TestProject") {
		t.Error("Expected project name in AGENTS.md")
	}

	if !strings.Contains(result, "/neev:bridge") {
		t.Error("Expected /neev:bridge command")
	}

	if !strings.Contains(result, "/neev:draft") {
		t.Error("Expected /neev:draft command")
	}

	if !strings.Contains(result, "/neev:inspect") {
		t.Error("Expected /neev:inspect command")
	}

	if !strings.Contains(result, "Claude Code") {
		t.Error("Expected Claude Code tool name")
	}
}

func TestGenerateAgentsMD_NoTools(t *testing.T) {
	result := GenerateAgentsMD([]string{}, "TestProject")

	if !strings.Contains(result, "No AI tools with native slash command support") {
		t.Error("Expected message about no tools configured")
	}
}

func TestGenerateSlashCommandManifest(t *testing.T) {
	result := GenerateSlashCommandManifest("claude-code")

	if !strings.Contains(result, "Claude Code") {
		t.Error("Expected tool name in manifest")
	}

	if !strings.Contains(result, "/neev:bridge") {
		t.Error("Expected /neev:bridge command in manifest")
	}

	if !strings.Contains(result, "Description") {
		t.Error("Expected description field")
	}
}

func TestGenerateInstructions(t *testing.T) {
	result := GenerateInstructions("MyProject")

	if !strings.Contains(result, "Neev Project Instructions") {
		t.Error("Expected instructions title")
	}

	if !strings.Contains(result, "MyProject") {
		t.Error("Expected project name in instructions")
	}

	if !strings.Contains(result, "Blueprints") {
		t.Error("Expected Blueprints concept")
	}

	if !strings.Contains(result, "/neev:bridge") {
		t.Error("Expected /neev:bridge in instructions")
	}
}

func TestFormatToolName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"claude-code", "Claude Code"},
		{"cursor", "Cursor"},
		{"codebuddy", "CodeBuddy"},
		{"opencode", "OpenCode"},
	}

	for _, test := range tests {
		result := formatToolName(test.input)
		if result != test.expected {
			t.Errorf("formatToolName(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestAllDefaultCommandsIncluded(t *testing.T) {
	expectedCommands := []string{"bridge", "draft", "inspect", "cucumber", "openapi", "handoff"}

	for _, cmd := range expectedCommands {
		if _, exists := DefaultSlashCommands[cmd]; !exists {
			t.Errorf("Expected command %s in DefaultSlashCommands", cmd)
		}
	}
}

func TestSlashCommandStructure(t *testing.T) {
	for name, cmd := range DefaultSlashCommands {
		if cmd.Name == "" {
			t.Errorf("Command %s has empty Name field", name)
		}

		if cmd.Description == "" {
			t.Errorf("Command %s has empty Description field", name)
		}

		if cmd.Prompt == "" {
			t.Errorf("Command %s has empty Prompt field", name)
		}
	}
}

func TestGenerateGitHubCopilotManifest(t *testing.T) {
	result, err := GenerateGitHubCopilotManifest("TestProject")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !strings.Contains(result, "TestProject") {
		t.Error("Expected project name in manifest")
	}

	if !strings.Contains(result, "version") {
		t.Error("Expected version field in JSON")
	}

	// Verify it's valid JSON
	var manifest GitHubCopilotManifest
	if err := json.Unmarshal([]byte(result), &manifest); err != nil {
		t.Errorf("Expected valid JSON, got error: %v", err)
	}

	if manifest.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", manifest.Version)
	}

	if manifest.ProjectName != "TestProject" {
		t.Errorf("Expected project name TestProject, got %s", manifest.ProjectName)
	}

	// Check all commands are present
	expectedCommands := []string{"neev:bridge", "neev:draft", "neev:inspect", "neev:cucumber", "neev:openapi", "neev:handoff"}
	for _, cmd := range expectedCommands {
		if _, exists := manifest.Commands[cmd]; !exists {
			t.Errorf("Expected command %s in manifest", cmd)
		}
	}
}

func TestGenerateGitHubCopilotManifest_CommandMetadata(t *testing.T) {
	result, _ := GenerateGitHubCopilotManifest("TestProject")

	var manifest GitHubCopilotManifest
	json.Unmarshal([]byte(result), &manifest)

	// Verify bridge command has all required metadata
	bridgeCmd, exists := manifest.Commands["neev:bridge"]
	if !exists {
		t.Fatal("Expected neev:bridge command")
	}

	if bridgeCmd.Name != "bridge" {
		t.Errorf("Expected name 'bridge', got '%s'", bridgeCmd.Name)
	}

	if bridgeCmd.Description == "" {
		t.Error("Expected non-empty description")
	}

	if bridgeCmd.Prompt == "" {
		t.Error("Expected non-empty prompt")
	}

	if bridgeCmd.Icon == "" {
		t.Error("Expected icon for command")
	}

	if len(bridgeCmd.Aliases) == 0 {
		t.Error("Expected at least one alias")
	}
}

func TestGenerateClaudeSlashCommandFile(t *testing.T) {
	cmd := SlashCommand{
		Name:        "Test",
		Description: "Test command",
		Prompt:      "Run test",
	}

	result := GenerateClaudeSlashCommandFile("test", cmd)

	if !strings.Contains(result, "---") {
		t.Error("Expected YAML frontmatter")
	}

	if !strings.Contains(result, "name: Neev: Test") {
		t.Error("Expected correct name in frontmatter")
	}

	if !strings.Contains(result, "description: Test command") {
		t.Error("Expected description in frontmatter")
	}

	if !strings.Contains(result, "category: Neev") {
		t.Error("Expected category in frontmatter")
	}

	if !strings.Contains(result, "<!-- NEEV:START -->") {
		t.Error("Expected start marker")
	}

	if !strings.Contains(result, "<!-- NEEV:END -->") {
		t.Error("Expected end marker")
	}

	if !strings.Contains(result, "# Neev Test Command") {
		t.Error("Expected command title")
	}
}

func TestGenerateGitHubCopilotPromptFile(t *testing.T) {
	cmd := SlashCommand{
		Name:        "Test",
		Description: "Test command",
		Prompt:      "Run test",
	}

	result := GenerateGitHubCopilotPromptFile("test", cmd)

	if !strings.Contains(result, "---") {
		t.Error("Expected YAML frontmatter")
	}

	if !strings.Contains(result, "description: Test command") {
		t.Error("Expected description in frontmatter")
	}

	if !strings.Contains(result, "<!-- NEEV:START -->") {
		t.Error("Expected start marker")
	}

	if !strings.Contains(result, "<!-- NEEV:END -->") {
		t.Error("Expected end marker")
	}

	if !strings.Contains(result, "# Neev Test Command") {
		t.Error("Expected command title")
	}

	if !strings.Contains(result, "$ARGUMENTS") {
		t.Error("Expected $ARGUMENTS placeholder")
	}
}

func TestGenerateGitHubCopilotPrompts(t *testing.T) {
	result := GenerateGitHubCopilotPrompts("TestProject")

	if len(result) != 6 {
		t.Errorf("Expected 6 command files, got %d", len(result))
	}

	expectedFiles := []string{"bridge.prompt.md", "draft.prompt.md", "inspect.prompt.md", "cucumber.prompt.md", "openapi.prompt.md", "handoff.prompt.md"}
	for _, file := range expectedFiles {
		if _, exists := result[file]; !exists {
			t.Errorf("Expected file %s not found", file)
		}
	}

	// Check bridge.prompt.md content
	bridgeContent := result["bridge.prompt.md"]
	if !strings.Contains(bridgeContent, "Analyze the project structure") {
		t.Error("Expected bridge-specific instructions")
	}
}
