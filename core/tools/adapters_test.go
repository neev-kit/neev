package tools

import (
	"testing"
)

func TestGetAdapter(t *testing.T) {
	tool := &Tool{
		Type: ToolCursor,
		Name: "Cursor",
	}
	adapter := GetAdapter(tool)
	if adapter == nil {
		t.Fatal("Expected adapter, got nil")
	}
}

func TestCursorAdapter(t *testing.T) {
	tool := &Tool{
		Type: ToolCursor,
		Name: "Cursor",
	}
	adapter := NewCursorAdapter(tool)
	if adapter == nil {
		t.Fatal("Failed to create Cursor adapter")
	}
	if adapter.Name() == "" {
		t.Error("Adapter name is empty")
	}
}

func TestClaudeAdapter(t *testing.T) {
	tool := &Tool{
		Type: ToolClaude,
		Name: "Claude",
	}
	adapter := NewClaudeAdapter(tool)
	if adapter == nil {
		t.Fatal("Failed to create Claude adapter")
	}
	if adapter.Name() == "" {
		t.Error("Adapter name is empty")
	}
}

func TestGetAdapters(t *testing.T) {
	tools := []Tool{
		{Type: ToolCursor, Name: "Cursor"},
		{Type: ToolClaude, Name: "Claude"},
	}
	for _, tool := range tools {
		adapter := GetAdapter(&tool)
		if adapter == nil {
			t.Errorf("Failed to get adapter for %s", tool.Name)
		}
	}
}
