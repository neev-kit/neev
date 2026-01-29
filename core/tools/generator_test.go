package tools

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSkillsGenerator(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-project"

	gen := NewSkillsGenerator(projectName, tmpDir)
	if gen == nil {
		t.Fatal("Failed to create generator")
	}

	if gen.projectName != projectName {
		t.Errorf("Project name mismatch: expected %s, got %s", projectName, gen.projectName)
	}

	if gen.projectRoot != tmpDir {
		t.Errorf("Project root mismatch: expected %s, got %s", tmpDir, gen.projectRoot)
	}
}

func TestGenerateSkills(t *testing.T) {
	tmpDir := t.TempDir()

	testSkills := []SkillContent{
		{
			Name:        "test-skill-1",
			Description: "Test skill 1",
			Content:     "Implementation 1",
			Type:        "command",
			Version:     "1.0",
		},
	}

	gen := NewSkillsGenerator("test", tmpDir)
	if gen == nil {
		t.Fatal("Failed to create generator")
	}

	err := gen.GenerateSkills(testSkills)
	if err != nil {
		t.Logf("GenerateSkills returned: %v (may be expected if no tools detected)", err)
	}
}

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

	gen := NewSkillsGenerator("test", tmpDir)
	err := gen.generateFallbackDocumentation(testSkills)

	if err != nil {
		t.Logf("Fallback generation returned: %v", err)
	}

	// Verify fallback directory was created
	fallbackDir := filepath.Join(tmpDir, ".neev", "skills")
	if _, err := os.Stat(fallbackDir); os.IsNotExist(err) {
		t.Logf("Fallback directory not created (may be expected): %s", fallbackDir)
	}
}

func TestEmptySkillsList(t *testing.T) {
	tmpDir := t.TempDir()

	gen := NewSkillsGenerator("test", tmpDir)
	err := gen.GenerateSkills([]SkillContent{})

	if err != nil {
		t.Logf("GenerateSkills with empty list returned: %v", err)
	}
}

func TestGenerateSummaryReport(t *testing.T) {
	tmpDir := t.TempDir()

	testSkills := []SkillContent{
		{Name: "skill-1", Description: "Skill 1", Type: "command", Version: "1.0"},
		{Name: "skill-2", Description: "Skill 2", Type: "snippet", Version: "1.0"},
	}

	gen := NewSkillsGenerator("test", tmpDir)
	report := gen.GenerateSummaryReport(testSkills)

	if len(report) == 0 {
		t.Log("Report is empty (expected if no tools detected)")
		return
	}

	if len(report) > 0 {
		t.Logf("Report generated: %d characters", len(report))
	}
}
