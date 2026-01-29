package tools

import (
	"os"
	"path/filepath"
	"testing"
)

// TestSkillsGenerator tests the skills generator
func TestSkillsGenerator(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()
	projectName := "test-project"

	// Create mock tools
	testTools := []Tool{
		{
			Type:      ToolCursor,
			Name:      "Cursor",
			Installed: true,
			Config: ToolConfig{
				SkillsDir: filepath.Join(tmpDir, "cursor-skills"),
				ConfigDir: filepath.Join(tmpDir, "cursor"),
				Native:    true,
			},
		},
	}

	// Create generator
	gen := NewSkillsGenerator(projectName, tmpDir, testTools)

	if gen == nil {
		t.Fatal("Failed to create generator")
	}

	if gen.projectName != projectName {
		t.Errorf("Generator project name mismatch")
	}

	if gen.projectRoot != tmpDir {
		t.Errorf("Generator project root mismatch")
	}
}

// TestGenerateSkills tests the full skill generation workflow
func TestGenerateSkills(t *testing.T) {
	tmpDir := t.TempDir()

	testTools := []Tool{
		{
			Type:      ToolCursor,
			Installed: true,
			Config: ToolConfig{
				SkillsDir: filepath.Join(tmpDir, "cursor-skills"),
			},
		},
	}

	testSkills := []SkillContent{
		{
			Name:        "test-skill-1",
			Description: "Test skill 1",
			Content:     "Implementation 1",
			Type:        "command",
			Version:     "1.0",
		},
		{
			Name:        "test-skill-2",
			Description: "Test skill 2",
			Content:     "Implementation 2",
			Type:        "snippet",
			Version:     "1.0",
		},
	}

	gen := NewSkillsGenerator("test", tmpDir, testTools)
	err := gen.GenerateSkills(testSkills)

	if err != nil {
		t.Fatalf("GenerateSkills failed: %v", err)
	}

	// Verify skill files were created
	skillsDir := testTools[0].Config.SkillsDir
	if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
		t.Errorf("Skills directory was not created: %s", skillsDir)
	}

	// Check if skill files exist
	for _, skill := range testSkills {
		skillPath := filepath.Join(skillsDir, skill.Name+".json")
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			t.Errorf("Skill file not created: %s", skillPath)
		}
	}

	// Verify README was created
	readmePath := filepath.Join(skillsDir, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		t.Error("README.md not created")
	}
}

// TestGenerateFallbackDocumentation tests fallback generation when no tools detected
func TestGenerateFallbackDocumentation(t *testing.T) {
	tmpDir := t.TempDir()

	testSkills := []SkillContent{
		{
			Name:        "test-skill",
			Description: "Test skill",
			Content:     "Implementation",
			Type:        "command",
			Version:     "1.0",
		},
	}

	gen := NewSkillsGenerator("test", tmpDir, []Tool{})
	err := gen.generateFallbackDocumentation(testSkills)

	if err != nil {
		t.Fatalf("generateFallbackDocumentation failed: %v", err)
	}

	// Verify fallback directory was created
	fallbackDir := filepath.Join(tmpDir, ".neev", "skills")
	if _, err := os.Stat(fallbackDir); os.IsNotExist(err) {
		t.Errorf("Fallback skills directory not created: %s", fallbackDir)
	}

	// Verify skill file exists
	skillPath := filepath.Join(fallbackDir, "test-skill.md")
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		t.Errorf("Fallback skill file not created: %s", skillPath)
	}
}

// TestGenerateSummaryReport tests report generation
func TestGenerateSummaryReport(t *testing.T) {
	tmpDir := t.TempDir()

	testTools := []Tool{
		{
			Type:      ToolCursor,
			Name:      "Cursor",
			Installed: true,
			Config: ToolConfig{
				SkillsDir: filepath.Join(tmpDir, "cursor-skills"),
				ConfigDir: filepath.Join(tmpDir, "cursor"),
			},
		},
		{
			Type:      ToolClaude,
			Name:      "Claude",
			Installed: true,
			Config: ToolConfig{
				SkillsDir: filepath.Join(tmpDir, "claude-skills"),
				ConfigDir: filepath.Join(tmpDir, "claude"),
			},
		},
	}

	testSkills := []SkillContent{
		{Name: "skill-1", Description: "Skill 1", Type: "command", Version: "1.0"},
		{Name: "skill-2", Description: "Skill 2", Type: "snippet", Version: "1.0"},
	}

	gen := NewSkillsGenerator("test", tmpDir, testTools)
	report := gen.GenerateSummaryReport(testSkills)

	// Verify report contains expected information
	if len(report) == 0 {
		t.Error("Report is empty")
	}

	requiredStrings := []string{
		"SKILLS GENERATION SUMMARY",
		"test",
		"TOOLS DETECTED",
		"Cursor",
		"Claude",
		"GENERATED SKILLS",
		"2",
		"skill-1",
		"skill-2",
	}

	for _, required := range requiredStrings {
		if !contains(report, required) {
			t.Errorf("Report missing expected string: %s", required)
		}
	}
}

// TestGenerateIndexFile tests index file generation
func TestGenerateIndexFile(t *testing.T) {
	tmpDir := t.TempDir()
	skillsDir := filepath.Join(tmpDir, "skills")
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tool := &Tool{
		Type:      ToolCursor,
		Installed: true,
		Config: ToolConfig{SkillsDir: skillsDir},
	}

	adapter := NewCursorAdapter(tool)

	testSkills := []SkillContent{
		{Name: "skill-1", Description: "Skill 1", Type: "command", Version: "1.0"},
		{Name: "skill-2", Description: "Skill 2", Type: "snippet", Version: "1.0"},
	}

	gen := NewSkillsGenerator("test", tmpDir, []Tool{tool})
	err := gen.generateIndexFile(adapter, skillsDir, testSkills)

	if err != nil {
		t.Fatalf("generateIndexFile failed: %v", err)
	}

	// Verify index file was created
	indexPath := filepath.Join(skillsDir, "INDEX.md")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read index file: %v", err)
	}

	indexContent := string(content)
	if !contains(indexContent, "Skills Index") {
		t.Error("Index missing header")
	}
	if !contains(indexContent, "skill-1") {
		t.Error("Index missing skill-1")
	}
	if !contains(indexContent, "skill-2") {
		t.Error("Index missing skill-2")
	}
}

// TestWriteSkillToFile tests writing skills to disk
func TestWriteSkillToFile(t *testing.T) {
	tmpDir := t.TempDir()

	tool := &Tool{
		Type:      ToolCursor,
		Installed: true,
		Config: ToolConfig{SkillsDir: filepath.Join(tmpDir, "skills")},
	}

	adapter := NewCursorAdapter(tool)

	skill := SkillContent{
		Name:        "test-skill",
		Description: "Test skill",
		Content:     "Implementation",
		Type:        "command",
		Version:     "1.0",
	}

	err := WriteSkillToFile(adapter, skill, tool.Config.SkillsDir)

	if err != nil {
		t.Fatalf("WriteSkillToFile failed: %v", err)
	}

	// Verify file was created
	skillPath := filepath.Join(tool.Config.SkillsDir, skill.Name+".json")
	content, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("Failed to read written skill: %v", err)
	}

	if len(content) == 0 {
		t.Error("Written skill file is empty")
	}
}

// TestMultipleToolGeneration tests generating skills for multiple tools
func TestMultipleToolGeneration(t *testing.T) {
	tmpDir := t.TempDir()

	testTools := []Tool{
		{
			Type:      ToolCursor,
			Name:      "Cursor",
			Installed: true,
			Config: ToolConfig{
				SkillsDir: filepath.Join(tmpDir, "cursor-skills"),
			},
		},
		{
			Type:      ToolClaude,
			Name:      "Claude",
			Installed: true,
			Config: ToolConfig{
				SkillsDir: filepath.Join(tmpDir, "claude-skills"),
			},
		},
	}

	testSkills := []SkillContent{
		{
			Name:        "test-skill",
			Description: "Test skill",
			Content:     "Implementation",
			Type:        "command",
			Version:     "1.0",
		},
	}

	gen := NewSkillsGenerator("test", tmpDir, testTools)
	err := gen.GenerateSkills(testSkills)

	if err != nil {
		t.Fatalf("GenerateSkills failed: %v", err)
	}

	// Verify skills generated for both tools in different formats
	cursorPath := filepath.Join(tmpDir, "cursor-skills", "test-skill.json")
	claudePath := filepath.Join(tmpDir, "claude-skills", "test-skill.md")

	if _, err := os.Stat(cursorPath); os.IsNotExist(err) {
		t.Error("Cursor skill not generated")
	}

	if _, err := os.Stat(claudePath); os.IsNotExist(err) {
		t.Error("Claude skill not generated")
	}
}

// TestEmptySkillsList tests handling of empty skills list
func TestEmptySkillsList(t *testing.T) {
	tmpDir := t.TempDir()

	testTools := []Tool{
		{
			Type:      ToolCursor,
			Installed: true,
			Config: ToolConfig{
				SkillsDir: filepath.Join(tmpDir, "skills"),
			},
		},
	}

	gen := NewSkillsGenerator("test", tmpDir, testTools)
	err := gen.GenerateSkills([]SkillContent{})

	// Should not fail with empty skills
	if err != nil {
		t.Fatalf("GenerateSkills failed with empty skills: %v", err)
	}
}

// TestSkillsGeneratorEdgeCases tests edge cases
func TestSkillsGeneratorEdgeCases(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name   string
		skills []SkillContent
	}{
		{
			name: "Special characters in skill name",
			skills: []SkillContent{
				{
					Name:        "skill-with-special_chars",
					Description: "Test",
					Content:     "test",
					Type:        "command",
					Version:     "1.0",
				},
			},
		},
		{
			name: "Long skill description",
			skills: []SkillContent{
				{
					Name:        "test",
					Description: "This is a very long description " + strings.Repeat("that repeats ", 10),
					Content:     "test",
					Type:        "command",
					Version:     "1.0",
				},
			},
		},
	}

	for _, test := range tests {
		testTools := []Tool{
			{
				Type:      ToolCursor,
				Installed: true,
				Config: ToolConfig{
					SkillsDir: filepath.Join(tmpDir, test.name),
				},
			},
		}

		gen := NewSkillsGenerator("test", tmpDir, testTools)
		if err := gen.GenerateSkills(test.skills); err != nil {
			t.Errorf("Test '%s' failed: %v", test.name, err)
		}
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s[len(s)-len(substr):] == substr
}

// BenchmarkGenerateSkills benchmarks skill generation
func BenchmarkGenerateSkills(b *testing.B) {
	tmpDir := b.TempDir()

	testTools := []Tool{
		{
			Type:      ToolCursor,
			Installed: true,
			Config: ToolConfig{
				SkillsDir: filepath.Join(tmpDir, "skills"),
			},
		},
	}

	testSkills := []SkillContent{
		{Name: "skill1", Description: "Test 1", Content: "test", Type: "command", Version: "1.0"},
		{Name: "skill2", Description: "Test 2", Content: "test", Type: "snippet", Version: "1.0"},
	}

	gen := NewSkillsGenerator("test", tmpDir, testTools)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gen.GenerateSkills(testSkills)
	}
}

// Import strings package for test
import "strings"
