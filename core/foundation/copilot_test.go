package foundation

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestGenerateSlashCommandManifest tests slash command manifest generation
func TestGenerateSlashCommandManifest(t *testing.T) {
	projectName := "test-project"
	manifest, err := GenerateSlashCommandManifest(projectName)

	if err != nil {
		t.Fatalf("GenerateSlashCommandManifest failed: %v", err)
	}

	if manifest == "" {
		t.Fatal("manifest is empty")
	}

	// Verify JSON is valid
	var data SlashCommandManifest
	if err := json.Unmarshal([]byte(manifest), &data); err != nil {
		t.Fatalf("failed to parse manifest as JSON: %v", err)
	}

	// Verify manifest structure
	if data.Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", data.Version)
	}

	if data.ProjectName != projectName {
		t.Errorf("expected project name %s, got %s", projectName, data.ProjectName)
	}

	// Verify all commands are present
	expectedCommands := []string{"neev:bridge", "neev:draft", "neev:inspect", "neev:cucumber", "neev:openapi", "neev:handoff"}
	for _, cmd := range expectedCommands {
		if _, exists := data.Commands[cmd]; !exists {
			t.Errorf("command %s not found in manifest", cmd)
		}
	}

	// Verify command structure
	for cmd, def := range data.Commands {
		if def.Name == "" {
			t.Errorf("command %s has empty name", cmd)
		}
		if def.Description == "" {
			t.Errorf("command %s has empty description", cmd)
		}
		if def.Prompt == "" {
			t.Errorf("command %s has empty prompt", cmd)
		}
	}
}

// TestGenerateSlashCommandManifestComplete verifies all command details
func TestGenerateSlashCommandManifestComplete(t *testing.T) {
	manifest, _ := GenerateSlashCommandManifest("test")
	var data SlashCommandManifest
	json.Unmarshal([]byte(manifest), &data)

	// Verify bridge command
	bridge := data.Commands["neev:bridge"]
	if bridge.Name != "bridge" {
		t.Errorf("bridge command name incorrect")
	}
	if !strings.Contains(bridge.Description, "aggregated") {
		t.Errorf("bridge command description doesn't mention context")
	}

	// Verify draft command
	draft := data.Commands["neev:draft"]
	if draft.Name != "draft" {
		t.Errorf("draft command name incorrect")
	}
	if !strings.Contains(draft.Description, "blueprint") {
		t.Errorf("draft command description doesn't mention blueprint")
	}

	// Verify all commands have aliases
	for cmd, def := range data.Commands {
		if len(def.Aliases) == 0 {
			t.Errorf("command %s has no aliases", cmd)
		}
	}
}

// TestGenerateCopilotChatInstructions tests Copilot instructions generation
func TestGenerateCopilotChatInstructions(t *testing.T) {
	projectName := "test-project"
	instructions := GenerateCopilotChatInstructions(projectName)

	if instructions == "" {
		t.Fatal("instructions is empty")
	}

	// Verify content includes project name
	if !strings.Contains(instructions, projectName) {
		t.Errorf("instructions don't include project name %s", projectName)
	}

	// Verify all slash commands are documented
	expectedCommands := []string{"/neev:bridge", "/neev:draft", "/neev:inspect", "/neev:cucumber", "/neev:openapi", "/neev:handoff"}
	for _, cmd := range expectedCommands {
		if !strings.Contains(instructions, cmd) {
			t.Errorf("instructions missing command %s", cmd)
		}
	}

	// Verify section headers
	if !strings.Contains(instructions, "## Neev Slash Commands") {
		t.Error("instructions missing Neev Slash Commands section")
	}

	if !strings.Contains(instructions, "## How to Use") {
		t.Error("instructions missing How to Use section")
	}

	if !strings.Contains(instructions, "## Development Workflow") {
		t.Error("instructions missing Development Workflow section")
	}
}

// TestCopilotInstructionsContainExamples verifies usage examples
func TestCopilotInstructionsContainExamples(t *testing.T) {
	instructions := GenerateCopilotChatInstructions("test")

	// Verify examples are included
	if !strings.Contains(instructions, "@Copilot /neev:bridge") {
		t.Error("instructions missing bridge example")
	}

	if !strings.Contains(instructions, "@Copilot /neev:draft") {
		t.Error("instructions missing draft example")
	}

	if !strings.Contains(instructions, "@Copilot /neev:inspect") {
		t.Error("instructions missing inspect example")
	}

	if !strings.Contains(instructions, "@Copilot /neev:cucumber") {
		t.Error("instructions missing cucumber example")
	}

	if !strings.Contains(instructions, "@Copilot /neev:openapi") {
		t.Error("instructions missing openapi example")
	}

	if !strings.Contains(instructions, "@Copilot /neev:handoff") {
		t.Error("instructions missing handoff example")
	}
}

// TestCopilotInstructionsIncludesTerminalCommands verifies terminal equivalents
func TestCopilotInstructionsIncludesTerminalCommands(t *testing.T) {
	instructions := GenerateCopilotChatInstructions("test")

	terminalCommands := []string{
		"neev bridge",
		"neev draft",
		"neev inspect",
		"neev cucumber",
		"neev openapi",
		"neev handoff",
	}

	for _, cmd := range terminalCommands {
		if !strings.Contains(instructions, cmd) {
			t.Errorf("instructions missing terminal command %s", cmd)
		}
	}
}
