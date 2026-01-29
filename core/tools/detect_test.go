package tools

import (
	"testing"
)

func TestDetectInstalledTools(t *testing.T) {
	tools := DetectInstalledTools()
	if tools == nil {
		t.Fatal("Expected tools list, got nil")
	}
	t.Logf("Detected %d tools", len(tools))
}

func TestHasAnyTool(t *testing.T) {
	tools := DetectInstalledTools()
	if len(tools) == 0 {
		t.Log("No tools detected on this system")
		return
	}
	// Verify we have valid tool entries
	for _, tool := range tools {
		if len(tool.Name) == 0 {
			t.Error("Tool name is empty")
		}
	}
}

func TestGetInstalledToolsNames(t *testing.T) {
	tools := DetectInstalledTools()
	for _, tool := range tools {
		if tool.Installed {
			if tool.Name == "" {
				t.Error("Installed tool has empty name")
			}
		}
	}
}

func TestDetectToolsTypes(t *testing.T) {
	tools := DetectInstalledTools()
	if len(tools) > 0 {
		// Verify tool types are set (ToolType is int, so just check it exists)
		for _, tool := range tools {
			if tool.Type < 0 {
				t.Error("Tool type appears invalid")
			}
		}
	}
}

func TestDetectPlatformTools(t *testing.T) {
	tools := DetectInstalledTools()
	if len(tools) == 0 {
		t.Skip("No tools detected, skipping detailed test")
	}
	t.Logf("Detected tools: %d", len(tools))
}
