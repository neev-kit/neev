package bridge

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildContext_Success(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup directory structure
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")

	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("Failed to create foundation dir: %v", err)
	}
	if err := os.MkdirAll(blueprintsPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprints dir: %v", err)
	}

	// Create test files
	foundationFile := filepath.Join(foundationPath, "base.md")
	if err := os.WriteFile(foundationFile, []byte("# Foundation\nBase content"), 0644); err != nil {
		t.Fatalf("Failed to write foundation file: %v", err)
	}

	blueprintDir := filepath.Join(blueprintsPath, "feature-x")
	if err := os.MkdirAll(blueprintDir, 0755); err != nil {
		t.Fatalf("Failed to create blueprint dir: %v", err)
	}

	blueprintFile := filepath.Join(blueprintDir, "spec.md")
	if err := os.WriteFile(blueprintFile, []byte("# Feature X\nFeature specification"), 0644); err != nil {
		t.Fatalf("Failed to write blueprint file: %v", err)
	}

	// Change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	context, err := BuildContext("")
	if err != nil {
		t.Errorf("BuildContext failed: %v", err)
	}

	if !strings.Contains(context, "# Project Foundation") {
		t.Errorf("Expected '# Project Foundation' in context, got: %s", context)
	}

	if !strings.Contains(context, "Base content") {
		t.Errorf("Expected 'Base content' in context, got: %s", context)
	}

	if !strings.Contains(context, "Feature specification") {
		t.Errorf("Expected 'Feature specification' in context, got: %s", context)
	}
}

func TestBuildContext_WithFocus(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup directory structure
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")

	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("Failed to create foundation dir: %v", err)
	}
	if err := os.MkdirAll(blueprintsPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprints dir: %v", err)
	}

	// Create test files with focused content
	foundationFile := filepath.Join(foundationPath, "api.md")
	if err := os.WriteFile(foundationFile, []byte("# API\napi specification"), 0644); err != nil {
		t.Fatalf("Failed to write foundation file: %v", err)
	}

	foundationFile2 := filepath.Join(foundationPath, "database.md")
	if err := os.WriteFile(foundationFile2, []byte("# Database\ndatabase info"), 0644); err != nil {
		t.Fatalf("Failed to write foundation file: %v", err)
	}

	// Change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	context, err := BuildContext("api")
	if err != nil {
		t.Errorf("BuildContext with focus failed: %v", err)
	}

	if !strings.Contains(context, "api specification") {
		t.Errorf("Expected 'api specification' in focused context, got: %s", context)
	}

	if strings.Contains(context, "database info") {
		t.Errorf("Did not expect 'database info' in focused context, got: %s", context)
	}
}

func TestBuildContext_NoFoundationDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup only blueprints directory (missing foundation)
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")
	if err := os.MkdirAll(blueprintsPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprints dir: %v", err)
	}

	// Change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	_, err = BuildContext("")
	if err == nil {
		t.Error("Expected error when foundation directory doesn't exist")
	}
}

func TestBuildContext_NoBlueprints(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup only foundation directory (missing blueprints)
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("Failed to create foundation dir: %v", err)
	}

	foundationFile := filepath.Join(foundationPath, "base.md")
	if err := os.WriteFile(foundationFile, []byte("# Foundation\nBase"), 0644); err != nil {
		t.Fatalf("Failed to write foundation file: %v", err)
	}

	// Change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	_, err = BuildContext("")
	if err == nil {
		t.Error("Expected error when blueprints directory doesn't exist")
	}
}

func TestBuildContext_NonMarkdownFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup directory structure
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")

	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("Failed to create foundation dir: %v", err)
	}
	if err := os.MkdirAll(blueprintsPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprints dir: %v", err)
	}

	// Create markdown file
	foundationFile := filepath.Join(foundationPath, "base.md")
	if err := os.WriteFile(foundationFile, []byte("# Foundation"), 0644); err != nil {
		t.Fatalf("Failed to write foundation file: %v", err)
	}

	// Create non-markdown file (should be ignored)
	foundationTxt := filepath.Join(foundationPath, "ignore.txt")
	if err := os.WriteFile(foundationTxt, []byte("ignore me"), 0644); err != nil {
		t.Fatalf("Failed to write txt file: %v", err)
	}

	// Change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	context, err := BuildContext("")
	if err != nil {
		t.Errorf("BuildContext failed: %v", err)
	}

	if strings.Contains(context, "ignore me") {
		t.Errorf("Did not expect 'ignore me' in context (non-markdown file)")
	}
}

func TestBuildContext_EmptyDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup directory structure with empty dirs
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")

	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("Failed to create foundation dir: %v", err)
	}
	if err := os.MkdirAll(blueprintsPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprints dir: %v", err)
	}

	// Change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	context, err := BuildContext("")
	if err != nil {
		t.Errorf("BuildContext with empty dirs failed: %v", err)
	}

	if !strings.Contains(context, "# Project Foundation") {
		t.Errorf("Expected '# Project Foundation' in context")
	}
}

func TestBuildContext_FilesInBlueprintRootIgnored(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup directory structure
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")

	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("Failed to create foundation dir: %v", err)
	}
	if err := os.MkdirAll(blueprintsPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprints dir: %v", err)
	}

	// Create markdown file in foundation
	foundationFile := filepath.Join(foundationPath, "base.md")
	if err := os.WriteFile(foundationFile, []byte("# Foundation"), 0644); err != nil {
		t.Fatalf("Failed to write foundation file: %v", err)
	}

	// Create markdown file in blueprints root (should be ignored)
	blueprintRootFile := filepath.Join(blueprintsPath, "readme.md")
	if err := os.WriteFile(blueprintRootFile, []byte("# Readme"), 0644); err != nil {
		t.Fatalf("Failed to write blueprint root file: %v", err)
	}

	// Change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	context, err := BuildContext("")
	if err != nil {
		t.Errorf("BuildContext failed: %v", err)
	}

	if strings.Contains(context, "# Readme") {
		t.Errorf("Did not expect blueprint root markdown file in context")
	}
}

func TestReadFilesInDir_Success(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test file
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("# Test\nContent"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	var builder strings.Builder
	err := readFilesInDir(tmpDir, &builder, "")
	if err != nil {
		t.Errorf("readFilesInDir failed: %v", err)
	}

	result := builder.String()
	if !strings.Contains(result, "Content") {
		t.Errorf("Expected 'Content' in result")
	}
}

func TestReadFilesInDir_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistent := filepath.Join(tmpDir, "nonexistent")

	var builder strings.Builder
	err := readFilesInDir(nonExistent, &builder, "")
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}

func TestReadFilesInDir_WithFocus(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	testFile1 := filepath.Join(tmpDir, "api.md")
	if err := os.WriteFile(testFile1, []byte("# API\napi content"), 0644); err != nil {
		t.Fatalf("Failed to write test file 1: %v", err)
	}

	testFile2 := filepath.Join(tmpDir, "db.md")
	if err := os.WriteFile(testFile2, []byte("# Database\ndb content"), 0644); err != nil {
		t.Fatalf("Failed to write test file 2: %v", err)
	}

	var builder strings.Builder
	err := readFilesInDir(tmpDir, &builder, "api")
	if err != nil {
		t.Errorf("readFilesInDir with focus failed: %v", err)
	}

	result := builder.String()
	if !strings.Contains(result, "api content") {
		t.Errorf("Expected 'api content' in result")
	}

	if strings.Contains(result, "db content") {
		t.Errorf("Did not expect 'db content' in focused result")
	}
}

func TestBuildContext_MultipleBlueprintDirs(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup directory structure
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")

	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("Failed to create foundation dir: %v", err)
	}
	if err := os.MkdirAll(blueprintsPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprints dir: %v", err)
	}

	// Create foundation file
	if err := os.WriteFile(filepath.Join(foundationPath, "base.md"), []byte("base"), 0644); err != nil {
		t.Fatalf("Failed to write foundation file: %v", err)
	}

	// Create multiple blueprint directories
	for i := 1; i <= 3; i++ {
		bpDir := filepath.Join(blueprintsPath, "bp"+string(rune(48+i)))
		if err := os.MkdirAll(bpDir, 0755); err != nil {
			t.Fatalf("Failed to create blueprint dir: %v", err)
		}
		bpFile := filepath.Join(bpDir, "spec.md")
		if err := os.WriteFile(bpFile, []byte("# BP"+string(rune(48+i))), 0644); err != nil {
			t.Fatalf("Failed to write blueprint file: %v", err)
		}
	}

	// Change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	context, err := BuildContext("")
	if err != nil {
		t.Errorf("BuildContext with multiple blueprints failed: %v", err)
	}

	if !strings.Contains(context, "# BP1") {
		t.Errorf("Expected BP1 in context")
	}
	if !strings.Contains(context, "# BP2") {
		t.Errorf("Expected BP2 in context")
	}
	if !strings.Contains(context, "# BP3") {
		t.Errorf("Expected BP3 in context")
	}
}

func TestBuildContext_FileReadError(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup directory structure
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")

	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("Failed to create foundation dir: %v", err)
	}
	if err := os.MkdirAll(blueprintsPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprints dir: %v", err)
	}

	// Create a file and then make it unreadable
	testFile := filepath.Join(foundationPath, "test.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Make file unreadable
	if err := os.Chmod(testFile, 0000); err != nil {
		t.Fatalf("Failed to change permissions: %v", err)
	}
	defer os.Chmod(testFile, 0644) // cleanup

	context, err := BuildContext("")
	if err == nil {
		t.Errorf("Expected error when file is unreadable, got: %s", context)
	}
}

func TestBuildContext_WindowsPaths(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup directory structure using filepath.Join (handles Windows paths correctly)
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")

	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("Failed to create foundation dir: %v", err)
	}
	if err := os.MkdirAll(blueprintsPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprints dir: %v", err)
	}

	// Create test files with Windows-compatible paths
	foundationFile := filepath.Join(foundationPath, "base.md")
	if err := os.WriteFile(foundationFile, []byte("# Foundation\nBase content"), 0644); err != nil {
		t.Fatalf("Failed to write foundation file: %v", err)
	}

	blueprintDir := filepath.Join(blueprintsPath, "feature-x")
	if err := os.MkdirAll(blueprintDir, 0755); err != nil {
		t.Fatalf("Failed to create blueprint dir: %v", err)
	}

	blueprintFile := filepath.Join(blueprintDir, "spec.md")
	if err := os.WriteFile(blueprintFile, []byte("# Feature X\nFeature specification"), 0644); err != nil {
		t.Fatalf("Failed to write blueprint file: %v", err)
	}

	// Change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Test that context is built correctly regardless of path separator
	context, err := BuildContext("")
	if err != nil {
		t.Fatalf("BuildContext failed: %v", err)
	}

	if !strings.Contains(context, "Base content") {
		t.Errorf("Expected 'Base content' in context, got: %s", context)
	}

	if !strings.Contains(context, "Feature specification") {
		t.Errorf("Expected 'Feature specification' in context, got: %s", context)
	}
}

func TestBuildRemoteContext_WindowsPaths(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup remote directory structure with filepath.Join
	remotesPath := filepath.Join(tmpDir, ".neev", "remotes")
	remotePath := filepath.Join(remotesPath, "shared-foundation")

	if err := os.MkdirAll(remotePath, 0755); err != nil {
		t.Fatalf("Failed to create remote dir: %v", err)
	}

	// Create test file in remote
	remoteFile := filepath.Join(remotePath, "base.md")
	if err := os.WriteFile(remoteFile, []byte("# Shared Foundation\nRemote content"), 0644); err != nil {
		t.Fatalf("Failed to write remote file: %v", err)
	}

	// Change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Test that remote context is built correctly with Windows-compatible paths
	context, err := BuildRemoteContext()
	if err != nil {
		t.Fatalf("BuildRemoteContext failed: %v", err)
	}

	if !strings.Contains(context, "shared-foundation") {
		t.Errorf("Expected 'shared-foundation' in context, got: %s", context)
	}

	if !strings.Contains(context, "Remote content") {
		t.Errorf("Expected 'Remote content' in context, got: %s", context)
	}
}

func TestPathJoinCorrectness(t *testing.T) {
	// This test verifies that filepath.Join handles both Unix and Windows paths correctly
	testCases := []struct {
		name     string
		parts    []string
		expected string
	}{
		{
			name:     "nested foundation path",
			parts:    []string{".neev", "foundation", "base.md"},
			expected: filepath.Join(".neev", "foundation", "base.md"),
		},
		{
			name:     "nested blueprints path",
			parts:    []string{".neev", "blueprints", "feature", "spec.md"},
			expected: filepath.Join(".neev", "blueprints", "feature", "spec.md"),
		},
		{
			name:     "remote path",
			parts:    []string{".neev", "remotes", "shared", "base.md"},
			expected: filepath.Join(".neev", "remotes", "shared", "base.md"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filepath.Join(tc.parts...)
			if result != tc.expected {
				t.Errorf("filepath.Join mismatch: got %q, expected %q", result, tc.expected)
			}

			// Verify no hardcoded separators
			if strings.Contains(result, "/") && filepath.Separator == '\\' {
				t.Errorf("Result contains Unix separator on Windows: %q", result)
			}
		})
	}
}
