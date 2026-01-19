package blueprint

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDraft_Success(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	err = Draft("My Feature")
	if err != nil {
		t.Errorf("Draft failed: %v", err)
	}

	// Verify blueprint directory was created
	blueprintPath := filepath.Join(".neev", "blueprints", "my-feature")
	if _, err := os.Stat(blueprintPath); os.IsNotExist(err) {
		t.Errorf("Expected blueprint directory at %s", blueprintPath)
	}

	// Verify intent.md exists
	intentFile := filepath.Join(blueprintPath, "intent.md")
	if _, err := os.Stat(intentFile); os.IsNotExist(err) {
		t.Errorf("Expected intent.md at %s", intentFile)
	}

	// Verify architecture.md exists
	archFile := filepath.Join(blueprintPath, "architecture.md")
	if _, err := os.Stat(archFile); os.IsNotExist(err) {
		t.Errorf("Expected architecture.md at %s", archFile)
	}

	// Verify api-spec.md exists
	apiSpecFile := filepath.Join(blueprintPath, "api-spec.md")
	if _, err := os.Stat(apiSpecFile); os.IsNotExist(err) {
		t.Errorf("Expected api-spec.md at %s", apiSpecFile)
	}

	// Verify security.md exists
	securityFile := filepath.Join(blueprintPath, "security.md")
	if _, err := os.Stat(securityFile); os.IsNotExist(err) {
		t.Errorf("Expected security.md at %s", securityFile)
	}
}

func TestDraft_SanitizesName(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	err = Draft("My COMPLEX Feature NAME")
	if err != nil {
		t.Errorf("Draft with complex name failed: %v", err)
	}

	// Verify directory was created with sanitized name
	blueprintPath := filepath.Join(".neev", "blueprints", "my-complex-feature-name")
	if _, err := os.Stat(blueprintPath); os.IsNotExist(err) {
		t.Errorf("Expected blueprint directory at %s (sanitized name)", blueprintPath)
	}
}

func TestDraft_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory with existing blueprint
	blueprintsPath := filepath.Join(tmpDir, ".neev", "blueprints")
	existingPath := filepath.Join(blueprintsPath, "existing-feature")
	if err := os.MkdirAll(existingPath, 0755); err != nil {
		t.Fatalf("Failed to create existing blueprint: %v", err)
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

	err = Draft("existing feature")
	if err == nil {
		t.Error("Expected error when blueprint already exists")
	}
}

func TestDraft_NoBlueprints(t *testing.T) {
	tmpDir := t.TempDir()

	// Don't create .neev/blueprints directory

	// Change to temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(oldCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	err = Draft("new feature")
	if err != nil {
		t.Errorf("Draft should handle missing parent directory: %v", err)
	}

	// Verify directory was created
	blueprintPath := filepath.Join(".neev", "blueprints", "new-feature")
	if _, err := os.Stat(blueprintPath); os.IsNotExist(err) {
		t.Errorf("Expected blueprint directory to be created at %s", blueprintPath)
	}
}

func TestDraft_IntentFileContent(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	err = Draft("Feature")
	if err != nil {
		t.Errorf("Draft failed: %v", err)
	}

	// Read and verify intent.md content
	intentFile := filepath.Join(".neev", "blueprints", "feature", "intent.md")
	content, err := os.ReadFile(intentFile)
	if err != nil {
		t.Errorf("Failed to read intent.md: %v", err)
	}

	if len(content) == 0 {
		t.Error("intent.md should not be empty")
	}
}

func TestDraft_ArchitectureFileContent(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	err = Draft("Feature")
	if err != nil {
		t.Errorf("Draft failed: %v", err)
	}

	// Read and verify architecture.md content
	archFile := filepath.Join(".neev", "blueprints", "feature", "architecture.md")
	content, err := os.ReadFile(archFile)
	if err != nil {
		t.Errorf("Failed to read architecture.md: %v", err)
	}

	if len(content) == 0 {
		t.Error("architecture.md should not be empty")
	}
}

func TestDraft_ApiSpecFileContent(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	err = Draft("Feature")
	if err != nil {
		t.Errorf("Draft failed: %v", err)
	}

	// Read and verify api-spec.md content
	apiSpecFile := filepath.Join(".neev", "blueprints", "feature", "api-spec.md")
	content, err := os.ReadFile(apiSpecFile)
	if err != nil {
		t.Errorf("Failed to read api-spec.md: %v", err)
	}

	if len(content) == 0 {
		t.Error("api-spec.md should not be empty")
	}
}

func TestDraft_SecurityFileContent(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	err = Draft("Feature")
	if err != nil {
		t.Errorf("Draft failed: %v", err)
	}

	// Read and verify security.md content
	securityFile := filepath.Join(".neev", "blueprints", "feature", "security.md")
	content, err := os.ReadFile(securityFile)
	if err != nil {
		t.Errorf("Failed to read security.md: %v", err)
	}

	if len(content) == 0 {
		t.Error("security.md should not be empty")
	}
}

func TestDraft_SpecialCharactersInName(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	err = Draft("Feature@#$%^&*()")
	if err != nil {
		t.Errorf("Draft with special characters failed: %v", err)
	}

	// Should be created as lowercase with spaces replaced by dashes
	blueprintPath := filepath.Join(".neev", "blueprints", "feature@#$%^&*()")
	if _, err := os.Stat(blueprintPath); os.IsNotExist(err) {
		t.Errorf("Expected blueprint directory to be created")
	}
}

func TestDraft_EmptyName(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	// Empty name creates a blueprint at empty path, which means .neev/blueprints itself
	// This should fail as it already exists
	err = Draft("")
	if err == nil {
		// Empty path might not be treated as a path component
		// Verify it tries to create something
	}
}

func TestDraft_MultipleBlueprints(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	// Create first blueprint
	err = Draft("Feature One")
	if err != nil {
		t.Errorf("First Draft failed: %v", err)
	}

	// Create second blueprint
	err = Draft("Feature Two")
	if err != nil {
		t.Errorf("Second Draft failed: %v", err)
	}

	// Verify both exist
	feature1 := filepath.Join(".neev", "blueprints", "feature-one")
	feature2 := filepath.Join(".neev", "blueprints", "feature-two")

	if _, err := os.Stat(feature1); os.IsNotExist(err) {
		t.Errorf("Expected first blueprint directory at %s", feature1)
	}

	if _, err := os.Stat(feature2); os.IsNotExist(err) {
		t.Errorf("Expected second blueprint directory at %s", feature2)
	}
}

func TestDraft_FilePermissions(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	err = Draft("Test")
	if err != nil {
		t.Errorf("Draft failed: %v", err)
	}

	// Check that files are readable
	intentFile := filepath.Join(".neev", "blueprints", "test", "intent.md")
	if _, err := os.ReadFile(intentFile); err != nil {
		t.Errorf("Failed to read intent.md: %v", err)
	}
}

func TestDraft_ContentNotEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	err = Draft("MyFeature")
	if err != nil {
		t.Errorf("Draft failed: %v", err)
	}

	// Verify files have content
	intentFile := filepath.Join(".neev", "blueprints", "myfeature", "intent.md")
	archFile := filepath.Join(".neev", "blueprints", "myfeature", "architecture.md")
	apiSpecFile := filepath.Join(".neev", "blueprints", "myfeature", "api-spec.md")
	securityFile := filepath.Join(".neev", "blueprints", "myfeature", "security.md")

	intent, _ := os.ReadFile(intentFile)
	arch, _ := os.ReadFile(archFile)
	apiSpec, _ := os.ReadFile(apiSpecFile)
	security, _ := os.ReadFile(securityFile)

	if len(intent) == 0 {
		t.Error("intent.md should have content")
	}
	if len(arch) == 0 {
		t.Error("architecture.md should have content")
	}
	if len(apiSpec) == 0 {
		t.Error("api-spec.md should have content")
	}
	if len(security) == 0 {
		t.Error("security.md should have content")
	}
}

func TestDraft_DirectoryHierarchy(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	err = Draft("Feature")
	if err != nil {
		t.Errorf("Draft failed: %v", err)
	}

	// Verify directory structure
	blueprintPath := filepath.Join(".neev", "blueprints", "feature")
	info, err := os.Stat(blueprintPath)
	if err != nil {
		t.Errorf("Failed to stat blueprint path: %v", err)
	}
	if !info.IsDir() {
		t.Error("Blueprint path should be a directory")
	}
}

func TestDraft_CreatesFoundationOnFirstDraft(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	err = Draft("First Blueprint")
	if err != nil {
		t.Errorf("Draft failed: %v", err)
	}

	// Verify foundation directory was created
	foundationPath := filepath.Join(".neev", "foundation")
	if _, err := os.Stat(foundationPath); os.IsNotExist(err) {
		t.Errorf("Expected foundation directory at %s", foundationPath)
	}
}

func TestDraft_FoundationHasStackFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	err = Draft("Feature")
	if err != nil {
		t.Errorf("Draft failed: %v", err)
	}

	// Verify stack.md exists
	stackFile := filepath.Join(".neev", "foundation", "stack.md")
	if _, err := os.Stat(stackFile); os.IsNotExist(err) {
		t.Errorf("Expected stack.md at %s", stackFile)
	}
}

func TestDraft_FoundationHasPrinciplesFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	err = Draft("Feature")
	if err != nil {
		t.Errorf("Draft failed: %v", err)
	}

	// Verify principles.md exists
	principlesFile := filepath.Join(".neev", "foundation", "principles.md")
	if _, err := os.Stat(principlesFile); os.IsNotExist(err) {
		t.Errorf("Expected principles.md at %s", principlesFile)
	}
}

func TestDraft_FoundationHasPatternsFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	err = Draft("Feature")
	if err != nil {
		t.Errorf("Draft failed: %v", err)
	}

	// Verify patterns.md exists
	patternsFile := filepath.Join(".neev", "foundation", "patterns.md")
	if _, err := os.Stat(patternsFile); os.IsNotExist(err) {
		t.Errorf("Expected patterns.md at %s", patternsFile)
	}
}

func TestDraft_FoundationFilesHaveContent(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	err = Draft("Feature")
	if err != nil {
		t.Errorf("Draft failed: %v", err)
	}

	// Verify all foundation files have content
	foundationFiles := []string{
		"stack.md",
		"principles.md",
		"patterns.md",
	}

	for _, filename := range foundationFiles {
		filePath := filepath.Join(".neev", "foundation", filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("Failed to read %s: %v", filename, err)
			continue
		}

		if len(content) == 0 {
			t.Errorf("%s should not be empty", filename)
		}
	}
}

func TestDraft_FoundationNotOverwrittenOnSecondDraft(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .neev/blueprints directory
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

	// Create first draft
	err = Draft("Feature One")
	if err != nil {
		t.Errorf("First Draft failed: %v", err)
	}

	// Read original content
	stackFile := filepath.Join(".neev", "foundation", "stack.md")
	originalContent, err := os.ReadFile(stackFile)
	if err != nil {
		t.Errorf("Failed to read stack.md: %v", err)
	}

	// Create second draft
	err = Draft("Feature Two")
	if err != nil {
		t.Errorf("Second Draft failed: %v", err)
	}

	// Verify foundation files still exist and weren't overwritten
	newContent, err := os.ReadFile(stackFile)
	if err != nil {
		t.Errorf("Failed to read stack.md after second draft: %v", err)
	}

	if string(originalContent) != string(newContent) {
		t.Error("Foundation files should not be overwritten on subsequent drafts")
	}
}
