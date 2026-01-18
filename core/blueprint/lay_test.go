package blueprint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLay_Success(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup blueprint
	blueprintPath := filepath.Join(tmpDir, ".neev", "blueprints", "auth-system")
	if err := os.MkdirAll(blueprintPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprint: %v", err)
	}

	if err := os.WriteFile(filepath.Join(blueprintPath, "intent.md"), []byte("# Intent"), 0644); err != nil {
		t.Fatalf("Failed to write intent: %v", err)
	}

	if err := os.WriteFile(filepath.Join(blueprintPath, "architecture.md"), []byte("# Architecture"), 0644); err != nil {
		t.Fatalf("Failed to write architecture: %v", err)
	}

	// Change to temp directory
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Lay the blueprint
	err := Lay("auth-system")
	if err != nil {
		t.Errorf("Lay failed: %v", err)
	}

	// Verify files moved to archive
	archivePath := filepath.Join(tmpDir, ".neev", "foundation", "archive", "auth-system")
	if _, err := os.Stat(filepath.Join(archivePath, "intent.md")); os.IsNotExist(err) {
		t.Error("intent.md not found in archive")
	}

	if _, err := os.Stat(filepath.Join(archivePath, "architecture.md")); os.IsNotExist(err) {
		t.Error("architecture.md not found in archive")
	}

	// Verify blueprint folder deleted
	if _, err := os.Stat(blueprintPath); err == nil {
		t.Error("Blueprint folder should be deleted")
	}

	// Verify changelog updated
	changelogPath := filepath.Join(tmpDir, ".neev", "changelog.md")
	content, err := os.ReadFile(changelogPath)
	if err != nil {
		t.Errorf("Failed to read changelog: %v", err)
	}

	if !strings.Contains(string(content), "auth-system") {
		t.Error("Changelog should mention auth-system")
	}
}

func TestLay_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	// Change to temp directory
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Try to lay non-existent blueprint
	err := Lay("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent blueprint")
	}
}

func TestLay_CreatesArchiveDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup blueprint without pre-existing archive
	blueprintPath := filepath.Join(tmpDir, ".neev", "blueprints", "feature")
	if err := os.MkdirAll(blueprintPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprint: %v", err)
	}

	if err := os.WriteFile(filepath.Join(blueprintPath, "intent.md"), []byte("# Intent"), 0644); err != nil {
		t.Fatalf("Failed to write intent: %v", err)
	}

	if err := os.WriteFile(filepath.Join(blueprintPath, "architecture.md"), []byte("# Architecture"), 0644); err != nil {
		t.Fatalf("Failed to write architecture: %v", err)
	}

	// Change to temp directory
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Lay the blueprint
	err := Lay("feature")
	if err != nil {
		t.Errorf("Lay failed: %v", err)
	}

	// Verify archive directory was created
	archivePath := filepath.Join(tmpDir, ".neev", "foundation", "archive", "feature")
	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		t.Error("Archive directory not created")
	}
}

func TestLay_UpdatesChangelog(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup blueprint
	blueprintPath := filepath.Join(tmpDir, ".neev", "blueprints", "logging")
	if err := os.MkdirAll(blueprintPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprint: %v", err)
	}

	if err := os.WriteFile(filepath.Join(blueprintPath, "intent.md"), []byte("# Logging"), 0644); err != nil {
		t.Fatalf("Failed to write intent: %v", err)
	}

	if err := os.WriteFile(filepath.Join(blueprintPath, "architecture.md"), []byte("# Architecture"), 0644); err != nil {
		t.Fatalf("Failed to write architecture: %v", err)
	}

	// Change to temp directory
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Lay the blueprint
	err := Lay("logging")
	if err != nil {
		t.Errorf("Lay failed: %v", err)
	}

	// Verify changelog exists and has entry
	changelogPath := filepath.Join(tmpDir, ".neev", "changelog.md")
	if _, err := os.Stat(changelogPath); os.IsNotExist(err) {
		t.Error("Changelog file not created")
	}

	content, _ := os.ReadFile(changelogPath)
	contentStr := string(content)

	if !strings.Contains(contentStr, "logging") {
		t.Error("Changelog should contain blueprint name")
	}

	if !strings.Contains(contentStr, "laid blueprint") {
		t.Error("Changelog should contain 'laid blueprint' text")
	}
}

func TestLay_MultipleLays(t *testing.T) {
	tmpDir := t.TempDir()

	blueprints := []string{"auth", "api", "database"}

	for _, name := range blueprints {
		// Setup blueprint
		blueprintPath := filepath.Join(tmpDir, ".neev", "blueprints", name)
		if err := os.MkdirAll(blueprintPath, 0755); err != nil {
			t.Fatalf("Failed to create blueprint: %v", err)
		}

		if err := os.WriteFile(filepath.Join(blueprintPath, "intent.md"), []byte("# Intent"), 0644); err != nil {
			t.Fatalf("Failed to write intent: %v", err)
		}

		if err := os.WriteFile(filepath.Join(blueprintPath, "architecture.md"), []byte("# Architecture"), 0644); err != nil {
			t.Fatalf("Failed to write architecture: %v", err)
		}
	}

	// Change to temp directory
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Lay all blueprints
	for _, name := range blueprints {
		err := Lay(name)
		if err != nil {
			t.Errorf("Lay failed for %s: %v", name, err)
		}
	}

	// Verify all are archived
	for _, name := range blueprints {
		archivePath := filepath.Join(tmpDir, ".neev", "foundation", "archive", name)
		if _, err := os.Stat(archivePath); os.IsNotExist(err) {
			t.Errorf("Archive for %s not found", name)
		}
	}

	// Verify changelog has all entries
	changelogPath := filepath.Join(tmpDir, ".neev", "changelog.md")
	content, _ := os.ReadFile(changelogPath)
	contentStr := string(content)

	for _, name := range blueprints {
		if !strings.Contains(contentStr, name) {
			t.Errorf("Changelog missing entry for %s", name)
		}
	}
}

func TestLay_PartialFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup blueprint with only intent.md
	blueprintPath := filepath.Join(tmpDir, ".neev", "blueprints", "partial")
	if err := os.MkdirAll(blueprintPath, 0755); err != nil {
		t.Fatalf("Failed to create blueprint: %v", err)
	}

	if err := os.WriteFile(filepath.Join(blueprintPath, "intent.md"), []byte("# Intent"), 0644); err != nil {
		t.Fatalf("Failed to write intent: %v", err)
	}

	// Change to temp directory
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Lay the blueprint - should not fail
	err := Lay("partial")
	if err != nil {
		t.Errorf("Lay should handle partial files: %v", err)
	}

	// Verify intent.md was moved
	archivePath := filepath.Join(tmpDir, ".neev", "foundation", "archive", "partial")
	if _, err := os.Stat(filepath.Join(archivePath, "intent.md")); os.IsNotExist(err) {
		t.Error("intent.md not found in archive")
	}
}
