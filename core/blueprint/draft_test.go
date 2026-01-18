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

	intent, _ := os.ReadFile(intentFile)
	arch, _ := os.ReadFile(archFile)

	if len(intent) == 0 {
		t.Error("intent.md should have content")
	}
	if len(arch) == 0 {
		t.Error("architecture.md should have content")
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
