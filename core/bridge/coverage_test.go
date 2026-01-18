package bridge

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Additional coverage tests for edge cases

func TestBuildContext_EdgeCases(t *testing.T) {
	tmpDir := t.TempDir()

	// Create minimal structure
	os.MkdirAll(filepath.Join(tmpDir, ".neev", "foundation"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, ".neev", "blueprints"), 0755)

	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Test with empty focus
	context, err := BuildContext("")
	if err != nil {
		t.Errorf("BuildContext failed: %v", err)
	}

	if context == "" {
		t.Error("BuildContext should return non-empty context")
	}

	// Test with specific focus
	os.WriteFile(filepath.Join(tmpDir, ".neev", "foundation", "test.md"), []byte("# Test\nfocusword"), 0644)
	context, err = BuildContext("focusword")
	if err != nil {
		t.Errorf("BuildContext with focus failed: %v", err)
	}
}

func TestBuildContext_MissingDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Test with no .neev directory
	_, err := BuildContext("")
	// Should fail or handle gracefully
	_ = err
}

func TestBuildRemoteContext_Empty(t *testing.T) {
	tmpDir := t.TempDir()

	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Test with no remotes directory
	context, err := BuildRemoteContext()
	if err != nil {
		t.Errorf("BuildRemoteContext failed: %v", err)
	}

	// Should return empty string for no remotes
	if context != "" && !strings.Contains(context, "Remote") {
		t.Error("Expected empty or Remote-containing context")
	}
}

func TestBuildRemoteContext_WithRemotes(t *testing.T) {
	tmpDir := t.TempDir()

	// Create remotes structure
	remotePath := filepath.Join(tmpDir, ".neev", "remotes", "origin")
	os.MkdirAll(remotePath, 0755)
	os.WriteFile(filepath.Join(remotePath, "foundation.md"), []byte("# Remote Foundation"), 0644)

	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	context, err := BuildRemoteContext()
	if err != nil {
		t.Errorf("BuildRemoteContext failed: %v", err)
	}

	if !strings.Contains(context, "Remote") {
		t.Error("Expected 'Remote' in context")
	}
}

func TestReadFilesInDir_NonMD(t *testing.T) {
	tmpDir := t.TempDir()

	// Create non-.md files
	os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file.json"), []byte("{}"), 0644)

	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// readFilesInDir should only read .md files
	var builder strings.Builder
	_ = readFilesInDir(tmpDir, &builder, "")

	// Shouldn't include non-.md content
	output := builder.String()
	if strings.Contains(output, "readme.txt") {
		t.Error("Should not include .txt files")
	}
}

func TestFormatSlashCommand_LongContext(t *testing.T) {
	longContext := strings.Repeat("This is a long context line\n", 100)
	result := FormatSlashCommand(longContext)

	if !strings.Contains(result, "/context") {
		t.Error("Should contain /context")
	}

	// Should still be valid markdown
	if !strings.Contains(result, "```markdown") {
		t.Error("Should have markdown fence")
	}
}

func TestFormatHandoffPrompt_AllCombinations(t *testing.T) {
	testCases := []struct {
		role         string
		context      string
		instructions string
	}{
		{"Developer", "ctx", "instr"},
		{"", "", ""},
		{"QA", "longer context\nwith multiple lines", "numbered\n1. first\n2. second"},
	}

	for _, tc := range testCases {
		result := FormatHandoffPrompt(tc.role, tc.context, tc.instructions)
		if !strings.Contains(result, "/neev-handoff-") {
			t.Error("Should contain handoff command")
		}
	}
}
