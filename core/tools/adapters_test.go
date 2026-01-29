package tools

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestCursorAdapter tests the Cursor adapter
func TestCursorAdapter(t *testing.T) {
	tool := &Tool{
		Type:      ToolCursor,
		Name:      "Cursor",
		Installed: true,
		Config: ToolConfig{
			SkillsDir: "/home/user/.cursor/skills",
			Native:    true,
		},
	}

	adapter := NewCursorAdapter(tool)

	// Test Name
	if adapter.Name() != "Cursor" {
		t.Errorf("Expected 'Cursor', got '%s'", adapter.Name())
	}

	// Test GenerateSkill
	skill := SkillContent{
		Name:        "test-skill",
		Description: "Test skill",
		Content:     "Implementation",
		Type:        "command",
		Version:     "1.0",
	}

	result, err := adapter.GenerateSkill(skill)
	if err != nil {
		t.Fatalf("GenerateSkill failed: %v", err)
	}

	// Verify JSON format
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Errorf("Generated skill is not valid JSON: %v", err)
	}

	if data["name"] != "test-skill" {
		t.Errorf("Expected name 'test-skill', got '%v'", data["name"])
	}
}

// TestClaudeAdapter tests the Claude adapter
func TestClaudeAdapter(t *testing.T) {
	tool := &Tool{
		Type:      ToolClaude,
		Name:      "Claude",
		Installed: true,
		Config: ToolConfig{
			SkillsDir: "/home/user/.claude/skills",
			Native:    true,
		},
	}

	adapter := NewClaudeAdapter(tool)

	// Test Name
	if adapter.Name() != "Claude" {
		t.Errorf("Expected 'Claude', got '%s'", adapter.Name())
	}

	// Test GenerateSkill
	skill := SkillContent{
		Name:        "test-skill",
		Description: "Test skill",
		Content:     "Implementation",
		Type:        "command",
		Language:    "go",
		Version:     "1.0",
	}

	result, err := adapter.GenerateSkill(skill)
	if err != nil {
		t.Fatalf("GenerateSkill failed: %v", err)
	}

	// Verify markdown format
	if !strings.Contains(result, "# Skill: test-skill") {
		t.Error("Expected markdown heading")
	}
	if !strings.Contains(result, "**Description:**") {
		t.Error("Expected description field")
	}
}

// TestCopilotAdapter tests the GitHub Copilot adapter
func TestCopilotAdapter(t *testing.T) {
	tool := &Tool{
		Type:      ToolCopilot,
		Name:      "GitHub Copilot",
		Installed: true,
		Config: ToolConfig{
			SkillsDir: "/home/user/.copilot/skills",
			Native:    true,
		},
	}

	adapter := NewCopilotAdapter(tool)

	// Test Name
	if adapter.Name() != "GitHub Copilot" {
		t.Errorf("Expected 'GitHub Copilot', got '%s'", adapter.Name())
	}

	// Test GenerateSkill
	skill := SkillContent{
		Name:        "test-skill",
		Description: "Test skill",
		Content:     "console.log('test');",
		Type:        "command",
		Language:    "javascript",
		Version:     "1.0",
	}

	result, err := adapter.GenerateSkill(skill)
	if err != nil {
		t.Fatalf("GenerateSkill failed: %v", err)
	}

	// Verify markdown with code block
	if !strings.Contains(result, "# test-skill") {
		t.Error("Expected markdown heading")
	}
	if !strings.Contains(result, "```javascript") {
		t.Error("Expected code block with language")
	}
}

// TestCodeiumAdapter tests the Codeium adapter
func TestCodeiumAdapter(t *testing.T) {
	tool := &Tool{
		Type:      ToolCodeium,
		Name:      "Codeium",
		Installed: true,
		Config: ToolConfig{
			SkillsDir: "/home/user/.codeium/skills",
			Native:    true,
		},
	}

	adapter := NewCodeiumAdapter(tool)

	// Test Name
	if adapter.Name() != "Codeium" {
		t.Errorf("Expected 'Codeium', got '%s'", adapter.Name())
	}

	// Test GenerateSkill
	skill := SkillContent{
		Name:        "test-skill",
		Description: "Test skill",
		Content:     "Implementation",
		Type:        "command",
		Language:    "python",
		Version:     "1.0",
	}

	result, err := adapter.GenerateSkill(skill)
	if err != nil {
		t.Fatalf("GenerateSkill failed: %v", err)
	}

	// Verify JSON format
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Errorf("Generated skill is not valid JSON: %v", err)
	}

	if data["id"] != "test-skill" {
		t.Errorf("Expected id 'test-skill', got '%v'", data["id"])
	}
}

// TestFallbackAdapter tests the fallback adapter
func TestFallbackAdapter(t *testing.T) {
	tool := &Tool{
		Type:      ToolType("unknown"),
		Name:      "Unknown Tool",
		Installed: true,
		Config: ToolConfig{
			SkillsDir: "/home/user/.unknown/skills",
			Native:    false,
		},
	}

	adapter := NewFallbackAdapter(tool)

	// Test Name
	if adapter.Name() != "Natural Language Fallback" {
		t.Errorf("Expected 'Natural Language Fallback', got '%s'", adapter.Name())
	}

	// Test GenerateSkill
	skill := SkillContent{
		Name:        "test-skill",
		Description: "Test skill",
		Content:     "Implementation",
		Type:        "command",
		Version:     "1.0",
	}

	result, err := adapter.GenerateSkill(skill)
	if err != nil {
		t.Fatalf("GenerateSkill failed: %v", err)
	}

	// Verify markdown with documentation
	if !strings.Contains(result, "# Skill: test-skill") {
		t.Error("Expected markdown heading")
	}
	if !strings.Contains(result, "How to Use") {
		t.Error("Expected usage documentation")
	}
	if !strings.Contains(result, "Unknown Tool") {
		t.Error("Expected tool name in documentation")
	}
}

// TestGetAdapter tests adapter selection
func TestGetAdapter(t *testing.T) {
	tests := []struct {
		toolType     ToolType
		expectedName string
		expectNative bool
	}{
		{ToolCursor, "Cursor", true},
		{ToolClaude, "Claude", true},
		{ToolCopilot, "GitHub Copilot", true},
		{ToolCodeium, "Codeium", true},
		{ToolSupabase, "Natural Language Fallback", false},
		{ToolPerplexity, "Natural Language Fallback", false},
	}

	for _, test := range tests {
		tool := &Tool{
			Type:      test.toolType,
			Installed: true,
			Config: ToolConfig{
				SkillsDir: "/tmp/skills",
				Native:    test.expectNative,
			},
		}

		adapter := GetAdapter(tool)
		if adapter.Name() != test.expectedName {
			t.Errorf("For tool type %s: expected '%s', got '%s'",
				test.toolType, test.expectedName, adapter.Name())
		}
	}
}

// TestGetAdapters tests getting adapters for multiple tools
func TestGetAdapters(t *testing.T) {
	tools := []Tool{
		{
			Type:      ToolCursor,
			Installed: true,
			Config:    ToolConfig{SkillsDir: "/tmp"},
		},
		{
			Type:      ToolClaude,
			Installed: true,
			Config:    ToolConfig{SkillsDir: "/tmp"},
		},
		{
			Type:      ToolCopilot,
			Installed: false, // Should be skipped
			Config:    ToolConfig{SkillsDir: "/tmp"},
		},
	}

	adapters := GetAdapters(tools)

	if len(adapters) != 2 {
		t.Errorf("Expected 2 adapters, got %d", len(adapters))
	}

	// Verify we got the right adapters
	names := map[string]bool{}
	for _, adapter := range adapters {
		names[adapter.Name()] = true
	}

	if !names["Cursor"] {
		t.Error("Missing Cursor adapter")
	}
	if !names["Claude"] {
		t.Error("Missing Claude adapter")
	}
}

// TestGenerateConfigFile tests config file generation
func TestGenerateConfigFile(t *testing.T) {
	tool := &Tool{
		Type:      ToolCursor,
		Installed: true,
		Config:    ToolConfig{SkillsDir: "/tmp"},
	}

	adapter := NewCursorAdapter(tool)
	skills := []SkillContent{
		{Name: "skill1", Description: "Skill 1"},
		{Name: "skill2", Description: "Skill 2"},
	}

	result, err := adapter.GenerateConfigFile("test-project", skills)
	if err != nil {
		t.Fatalf("GenerateConfigFile failed: %v", err)
	}

	// Verify config contains expected fields
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(result), &config); err != nil {
		t.Errorf("Config is not valid JSON: %v", err)
	}

	if config["project"] != "test-project" {
		t.Error("Config missing project name")
	}
	if config["skills"] != 2 {
		t.Error("Config missing skill count")
	}
}

// TestGetMetadata tests metadata generation
func TestGetMetadata(t *testing.T) {
	tool := &Tool{
		Type:   ToolCursor,
		Config: ToolConfig{SkillsDir: "/path/to/skills", Native: true},
	}

	adapter := NewCursorAdapter(tool)
	metadata := adapter.GetMetadata()

	if metadata["adapter"] != "Cursor" {
		t.Error("Missing adapter in metadata")
	}
	if metadata["native"] != true {
		t.Error("Missing native flag in metadata")
	}
	if metadata["formatType"] != "json" {
		t.Error("Missing formatType in metadata")
	}
}

// TestAdapterEdgeCases tests edge cases
func TestAdapterEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		skill      SkillContent
		shouldFail bool
	}{
		{
			name: "Empty skill name",
			skill: SkillContent{
				Name:        "",
				Description: "Test",
				Content:     "test",
				Type:        "command",
				Version:     "1.0",
			},
			shouldFail: false, // Should still generate
		},
		{
			name: "Special characters",
			skill: SkillContent{
				Name:        "test/skill",
				Description: "Test with \"quotes\"",
				Content:     "Line 1\nLine 2\n\"quoted\"",
				Type:        "command",
				Version:     "1.0",
			},
			shouldFail: false,
		},
	}

	tool := &Tool{
		Type:      ToolCursor,
		Installed: true,
		Config:    ToolConfig{SkillsDir: "/tmp"},
	}
	adapter := NewCursorAdapter(tool)

	for _, test := range tests {
		result, err := adapter.GenerateSkill(test.skill)
		if test.shouldFail && err == nil {
			t.Errorf("Test '%s' should have failed", test.name)
		}
		if !test.shouldFail && err != nil {
			t.Errorf("Test '%s' failed: %v", test.name, err)
		}
		if result == "" && !test.shouldFail {
			t.Errorf("Test '%s' returned empty result", test.name)
		}
	}
}

// BenchmarkCursorAdapter benchmarks Cursor adapter
func BenchmarkCursorAdapter(b *testing.B) {
	tool := &Tool{Type: ToolCursor, Installed: true}
	adapter := NewCursorAdapter(tool)
	skill := SkillContent{
		Name:        "test",
		Description: "Test skill",
		Content:     "Implementation",
		Type:        "command",
		Version:     "1.0",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		adapter.GenerateSkill(skill)
	}
}

// BenchmarkClaudeAdapter benchmarks Claude adapter
func BenchmarkClaudeAdapter(b *testing.B) {
	tool := &Tool{Type: ToolClaude, Installed: true}
	adapter := NewClaudeAdapter(tool)
	skill := SkillContent{
		Name:        "test",
		Description: "Test skill",
		Content:     "Implementation",
		Type:        "command",
		Version:     "1.0",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		adapter.GenerateSkill(skill)
	}
}
