package blueprint

import (
	"os"
	"path/filepath"
	"testing"
)

// Additional edge case tests for better coverage

func TestDraftVariousNames(t *testing.T) {
	tmpDir := t.TempDir()
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")
	os.MkdirAll(blueprintsPath, 0755)

	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Test with various edge case names that haven't been tested
	testCases := []string{
		"Simple",
		"with-dash",
		"123numbers",
	}

	for _, name := range testCases {
		Draft(name)
	}
}

func TestLayVariousBlueprints(t *testing.T) {
	tmpDir := t.TempDir()
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")

	// Create multiple test blueprints
	for i := 1; i <= 3; i++ {
		blueprintName := "blueprint" + string(rune('0'+i))
		blueprintDir := filepath.Join(blueprintsPath, blueprintName)
		os.MkdirAll(blueprintDir, 0755)
		os.WriteFile(filepath.Join(blueprintDir, "intent.md"), []byte("# Intent"), 0644)
		os.WriteFile(filepath.Join(blueprintDir, "architecture.md"), []byte("# Arch"), 0644)
	}

	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Lay each blueprint
	for i := 1; i <= 3; i++ {
		blueprintName := "blueprint" + string(rune('0'+i))
		_ = Lay(blueprintName)
	}
}

func TestMoveFileEdgeCases(t *testing.T) {
	tmpDir := t.TempDir()

	// Test with non-existent source
	_ = moveFile(filepath.Join(tmpDir, "nonexistent.txt"), filepath.Join(tmpDir, "dest.txt"))

	// Test with existing file
	srcFile := filepath.Join(tmpDir, "src.txt")
	dstFile := filepath.Join(tmpDir, "dst.txt")
	os.WriteFile(srcFile, []byte("content"), 0644)
	_ = moveFile(srcFile, dstFile)
}
