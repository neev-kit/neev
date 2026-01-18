package migration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrateBasic(t *testing.T) {
	// Test basic migration functionality
	cfg := MigrationConfig{
		RootDir:    "/tmp/test",
		SourceType: SourceTypeAuto,
		DryRun:     true,
	}

	result, err := Migrate(cfg)
	if err == nil && result != nil {
		// Expected behavior - at least it shouldn't panic
		return
	}
}

func TestDetectSourceType(t *testing.T) {
	tmpDir := t.TempDir()

	// Test auto-detection with openspec structure
	if err := os.Mkdir(filepath.Join(tmpDir, "specs"), 0755); err != nil {
		t.Fatalf("Failed to create specs dir: %v", err)
	}
	sourceType, err := detectSourceType(tmpDir)
	if err != nil || sourceType != SourceTypeOpenSpec {
		t.Errorf("Expected SourceTypeOpenSpec, got %s with error %v", sourceType, err)
	}

	// Test with changes directory
	tmpDir2 := t.TempDir()
	if err := os.Mkdir(filepath.Join(tmpDir2, "changes"), 0755); err != nil {
		t.Fatalf("Failed to create changes dir: %v", err)
	}
	sourceType, err = detectSourceType(tmpDir2)
	if err != nil || sourceType != SourceTypeOpenSpec {
		t.Errorf("Expected SourceTypeOpenSpec for changes dir, got %s with error %v", sourceType, err)
	}

	// Test with speckit structure
	tmpDir3 := t.TempDir()
	specifyDir := filepath.Join(tmpDir3, ".specify")
	if err := os.MkdirAll(specifyDir, 0755); err != nil {
		t.Fatalf("Failed to create .specify dir: %v", err)
	}
	specFile := filepath.Join(specifyDir, "spec.md")
	if err := os.WriteFile(specFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create spec.md: %v", err)
	}
	sourceType, err = detectSourceType(tmpDir3)
	if err != nil || sourceType != SourceTypeSpecKit {
		t.Errorf("Expected SourceTypeSpecKit, got %s with error %v", sourceType, err)
	}

	// Test with no recognized structure
	tmpDir4 := t.TempDir()
	sourceType, err = detectSourceType(tmpDir4)
	if err == nil {
		t.Errorf("Expected error when no project structure found, got sourceType: %s", sourceType)
	}
}

func TestMigrateOpenSpec(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test structure
	specsDir := filepath.Join(tmpDir, "specs")
	changesDir := filepath.Join(tmpDir, "changes")
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("Failed to create specs dir: %v", err)
	}
	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatalf("Failed to create changes dir: %v", err)
	}

	// Create test spec file
	specFile := filepath.Join(specsDir, "foundation.md")
	if err := os.WriteFile(specFile, []byte("# Test Spec"), 0644); err != nil {
		t.Fatalf("Failed to create spec file: %v", err)
	}

	// Create test blueprint
	blueprintDir := filepath.Join(changesDir, "feature1")
	if err := os.MkdirAll(blueprintDir, 0755); err != nil {
		t.Fatalf("Failed to create blueprint dir: %v", err)
	}
	blueprintFile := filepath.Join(blueprintDir, "blueprint.md")
	if err := os.WriteFile(blueprintFile, []byte("# Test Blueprint"), 0644); err != nil {
		t.Fatalf("Failed to create blueprint file: %v", err)
	}

	cfg := MigrationConfig{
		RootDir:    tmpDir,
		SourceType: SourceTypeOpenSpec,
		DryRun:     false,
		BackupOld:  false,
	}

	result, err := Migrate(cfg)
	if err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected successful migration, got Success=%v, Errors=%v", result.Success, result.Errors)
	}

	// Verify files were moved
	neevFoundationDir := filepath.Join(tmpDir, ".neev", "foundation")
	if _, err := os.Stat(filepath.Join(neevFoundationDir, "foundation.md")); err != nil {
		t.Errorf("Expected spec file to be migrated: %v", err)
	}

	neevBlueprintDir := filepath.Join(tmpDir, ".neev", "blueprints", "feature1")
	if _, err := os.Stat(filepath.Join(neevBlueprintDir, "blueprint.md")); err != nil {
		t.Errorf("Expected blueprint file to be migrated: %v", err)
	}
}

func TestMigrateSpecKit(t *testing.T) {
	tmpDir := t.TempDir()

	// Create speckit structure
	specifyDir := filepath.Join(tmpDir, ".specify")
	if err := os.MkdirAll(specifyDir, 0755); err != nil {
		t.Fatalf("Failed to create .specify dir: %v", err)
	}

	specFile := filepath.Join(specifyDir, "spec.md")
	if err := os.WriteFile(specFile, []byte("# SpecKit Test"), 0644); err != nil {
		t.Fatalf("Failed to create spec.md: %v", err)
	}

	cfg := MigrationConfig{
		RootDir:    tmpDir,
		SourceType: SourceTypeSpecKit,
		DryRun:     false,
		BackupOld:  false,
	}

	result, err := Migrate(cfg)
	if err != nil {
		t.Fatalf("SpecKit migration failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected successful migration, got Success=%v, Errors=%v", result.Success, result.Errors)
	}

	// Verify file was migrated
	neevFoundationDir := filepath.Join(tmpDir, ".neev", "foundation")
	coreFile := filepath.Join(neevFoundationDir, "core.md")
	if _, err := os.Stat(coreFile); err != nil {
		t.Errorf("Expected spec.md to be migrated to core.md: %v", err)
	}

	// Verify file contents
	data, err := os.ReadFile(coreFile)
	if err != nil || string(data) != "# SpecKit Test" {
		t.Errorf("Migrated file has incorrect contents: %v", err)
	}
}

func TestMigrateAutoDetect(t *testing.T) {
	// Test auto-detection
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "specs")
	if err := os.Mkdir(specsDir, 0755); err != nil {
		t.Fatalf("Failed to create specs dir: %v", err)
	}

	cfg := MigrationConfig{
		RootDir:    tmpDir,
		SourceType: SourceTypeAuto,
		DryRun:     true,
		BackupOld:  false,
	}

	result, err := Migrate(cfg)
	if err != nil {
		t.Fatalf("Auto-detect migration failed: %v", err)
	}

	if result.SourceType != SourceTypeOpenSpec {
		t.Errorf("Expected auto-detection to find SourceTypeOpenSpec, got %s", result.SourceType)
	}
}

func TestMigrateDryRun(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test structure
	specsDir := filepath.Join(tmpDir, "specs")
	if err := os.Mkdir(specsDir, 0755); err != nil {
		t.Fatalf("Failed to create specs dir: %v", err)
	}

	specFile := filepath.Join(specsDir, "test.md")
	if err := os.WriteFile(specFile, []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create spec file: %v", err)
	}

	cfg := MigrationConfig{
		RootDir:    tmpDir,
		SourceType: SourceTypeOpenSpec,
		DryRun:     true,
		BackupOld:  false,
	}

	result, err := Migrate(cfg)
	if err != nil {
		t.Fatalf("Dry-run migration failed: %v", err)
	}

	// In dry-run mode, files should NOT be actually created
	neevDir := filepath.Join(tmpDir, ".neev")
	if _, err := os.Stat(neevDir); err == nil {
		t.Errorf("Expected .neev dir not to be created in dry-run mode")
	}

	// But result should still be successful
	if !result.Success {
		t.Errorf("Expected dry-run to succeed, got Success=%v, Errors=%v", result.Success, result.Errors)
	}
}

func TestMigrateWithBackup(t *testing.T) {
	tmpDir := t.TempDir()

	// Create existing .neev directory
	neevDir := filepath.Join(tmpDir, ".neev")
	if err := os.MkdirAll(neevDir, 0755); err != nil {
		t.Fatalf("Failed to create .neev dir: %v", err)
	}

	neevFile := filepath.Join(neevDir, "existing.md")
	if err := os.WriteFile(neevFile, []byte("existing"), 0644); err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	// Create specs for migration
	specsDir := filepath.Join(tmpDir, "specs")
	if err := os.Mkdir(specsDir, 0755); err != nil {
		t.Fatalf("Failed to create specs dir: %v", err)
	}

	cfg := MigrationConfig{
		RootDir:    tmpDir,
		SourceType: SourceTypeOpenSpec,
		DryRun:     false,
		BackupOld:  true,
	}

	result, err := Migrate(cfg)
	if err != nil {
		t.Fatalf("Migration with backup failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected successful migration, got Success=%v, Errors=%v", result.Success, result.Errors)
	}

	// Check that backup was created
	if result.BackupDir == "" {
		t.Errorf("Expected backup directory to be created")
	}

	// Verify backup exists
	if _, err := os.Stat(result.BackupDir); err != nil {
		t.Errorf("Backup directory does not exist: %v", err)
	}

	// Verify backup contains the old file
	backupFile := filepath.Join(result.BackupDir, "existing.md")
	if _, err := os.Stat(backupFile); err != nil {
		t.Errorf("Backup file not found: %v", err)
	}
}

func TestMigrateInvalidSourceType(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := MigrationConfig{
		RootDir:    tmpDir,
		SourceType: "invalid",
		DryRun:     true,
	}

	result, err := Migrate(cfg)
	if err == nil {
		t.Errorf("Expected error for invalid source type, got nil error")
	}

	if result.Success {
		t.Errorf("Expected migration to fail for invalid source type")
	}
}

func TestMigrateSpecKitMissingSpec(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := MigrationConfig{
		RootDir:    tmpDir,
		SourceType: SourceTypeSpecKit,
		DryRun:     false,
	}

	result, err := Migrate(cfg)
	if err == nil {
		t.Errorf("Expected error when spec.md is missing")
	}

	if result.Success {
		t.Errorf("Expected migration to fail")
	}
}

func TestCopyDirectory(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create source structure
	subDir := filepath.Join(srcDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	file1 := filepath.Join(srcDir, "file1.txt")
	file2 := filepath.Join(subDir, "file2.txt")
	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	dstLocation := filepath.Join(dstDir, "copy")
	if err := copyDirectory(srcDir, dstLocation); err != nil {
		t.Fatalf("copyDirectory failed: %v", err)
	}

	// Verify files were copied
	if _, err := os.Stat(filepath.Join(dstLocation, "file1.txt")); err != nil {
		t.Errorf("file1.txt not copied: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dstLocation, "subdir", "file2.txt")); err != nil {
		t.Errorf("file2.txt not copied: %v", err)
	}

	// Verify content
	data, _ := os.ReadFile(filepath.Join(dstLocation, "file1.txt"))
	if string(data) != "content1" {
		t.Errorf("Copied file has incorrect content")
	}
}

