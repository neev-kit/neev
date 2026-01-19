package slash

import (
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
