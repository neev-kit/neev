package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Migrate detects the source type (if auto) and executes the migration.
func Migrate(cfg MigrationConfig) (*MigrationResult, error) {
	result := &MigrationResult{
		Success:    false,
		SourceType: cfg.SourceType,
	}

	// Auto-detect source type if needed
	if cfg.SourceType == SourceTypeAuto {
		detected, err := detectSourceType(cfg.RootDir)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Auto-detection failed: %v", err))
			return result, err
		}
		result.SourceType = detected
		result.Messages = append(result.Messages, fmt.Sprintf("Auto-detected source type: %s", detected))
	}

	// Create backup if requested and not in dry-run mode
	if cfg.BackupOld && !cfg.DryRun {
		backupDir, err := createBackup(cfg.RootDir)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Backup failed: %v", err))
			return result, err
		}
		result.BackupDir = backupDir
		result.Messages = append(result.Messages, fmt.Sprintf("Backup created at: %s", backupDir))
	}

	// Execute migration based on source type
	var err error
	switch result.SourceType {
	case SourceTypeOpenSpec:
		err = migrateOpenSpec(cfg, result)
	case SourceTypeSpecKit:
		err = migrateSpecKit(cfg, result)
	default:
		return result, fmt.Errorf("unsupported source type: %s", result.SourceType)
	}

	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Migration failed: %v", err))
		return result, err
	}

	result.Success = true
	return result, nil
}

// detectSourceType attempts to auto-detect the project type.
func detectSourceType(rootDir string) (SourceType, error) {
	// Check for openspec structure (specs/ and changes/ dirs)
	specsPath := filepath.Join(rootDir, "specs")
	changesPath := filepath.Join(rootDir, "changes")

	specsExists := false
	changesExists := false

	if _, err := os.Stat(specsPath); err == nil {
		specsExists = true
	}
	if _, err := os.Stat(changesPath); err == nil {
		changesExists = true
	}

	if specsExists || changesExists {
		return SourceTypeOpenSpec, nil
	}

	// Check for speckit structure (.specify/spec.md)
	specifyPath := filepath.Join(rootDir, ".specify", "spec.md")
	if _, err := os.Stat(specifyPath); err == nil {
		return SourceTypeSpecKit, nil
	}

	return "", fmt.Errorf("could not auto-detect project type - neither openspec nor speckit structure found")
}

// migrateOpenSpec handles migration from openspec projects.
// Moves specs/*.md -> .neev/foundation/, changes/* -> .neev/blueprints/
func migrateOpenSpec(cfg MigrationConfig, result *MigrationResult) error {
	neevDir := filepath.Join(cfg.RootDir, ".neev")
	foundationDir := filepath.Join(neevDir, "foundation")
	blueprintsDir := filepath.Join(neevDir, "blueprints")

	// Create .neev structure if it doesn't exist
	for _, dir := range []string{neevDir, foundationDir, blueprintsDir} {
		if !cfg.DryRun {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
			result.DirsCreatedCount++
		}
	}

	// Migrate specs/ -> .neev/foundation/
	specsPath := filepath.Join(cfg.RootDir, "specs")
	if _, err := os.Stat(specsPath); err == nil {
		entries, err := os.ReadDir(specsPath)
		if err != nil {
			return fmt.Errorf("failed to read specs directory: %w", err)
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				src := filepath.Join(specsPath, entry.Name())
				dst := filepath.Join(foundationDir, entry.Name())

				if !cfg.DryRun {
					data, err := os.ReadFile(src)
					if err != nil {
						result.Errors = append(result.Errors, fmt.Sprintf("Failed to read %s: %v", src, err))
						continue
					}

					if err := os.WriteFile(dst, data, 0644); err != nil {
						result.Errors = append(result.Errors, fmt.Sprintf("Failed to write %s: %v", dst, err))
						continue
					}
				}

				result.FilesMovedCount++
				result.Messages = append(result.Messages, fmt.Sprintf("Migrated: %s -> %s", src, dst))
			}
		}
	}

	// Migrate changes/ -> .neev/blueprints/
	changesPath := filepath.Join(cfg.RootDir, "changes")
	if _, err := os.Stat(changesPath); err == nil {
		entries, err := os.ReadDir(changesPath)
		if err != nil {
			return fmt.Errorf("failed to read changes directory: %w", err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				blueprintName := entry.Name()
				srcDir := filepath.Join(changesPath, blueprintName)
				dstDir := filepath.Join(blueprintsDir, blueprintName)

				if !cfg.DryRun {
					if err := copyDirectory(srcDir, dstDir); err != nil {
						result.Errors = append(result.Errors, fmt.Sprintf("Failed to copy blueprint %s: %v", blueprintName, err))
						continue
					}
				}

				result.FilesMovedCount++
				result.Messages = append(result.Messages, fmt.Sprintf("Migrated: %s -> %s", srcDir, dstDir))
			}
		}
	}

	return nil
}

// migrateSpecKit handles migration from speckit projects.
// Moves .specify/spec.md -> .neev/foundation/core.md
func migrateSpecKit(cfg MigrationConfig, result *MigrationResult) error {
	neevDir := filepath.Join(cfg.RootDir, ".neev")
	foundationDir := filepath.Join(neevDir, "foundation")

	// Create .neev/foundation if it doesn't exist
	for _, dir := range []string{neevDir, foundationDir} {
		if !cfg.DryRun {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
			result.DirsCreatedCount++
		}
	}

	// Migrate .specify/spec.md -> .neev/foundation/core.md
	srcPath := filepath.Join(cfg.RootDir, ".specify", "spec.md")
	if _, err := os.Stat(srcPath); err != nil {
		return fmt.Errorf("speckit spec.md not found at %s: %w", srcPath, err)
	}

	dstPath := filepath.Join(foundationDir, "core.md")

	if !cfg.DryRun {
		data, err := os.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("failed to read spec.md: %w", err)
		}

		if err := os.WriteFile(dstPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write core.md: %w", err)
		}
	}

	result.FilesMovedCount++
	result.Messages = append(result.Messages, fmt.Sprintf("Migrated: %s -> %s", srcPath, dstPath))

	return nil
}

// createBackup creates a timestamped backup of the .neev directory if it exists.
func createBackup(rootDir string) (string, error) {
	neevDir := filepath.Join(rootDir, ".neev")

	// Check if .neev already exists
	if _, err := os.Stat(neevDir); os.IsNotExist(err) {
		return "", nil // No backup needed
	}

	timestamp := time.Now().Format("20060102_150405")
	backupDir := filepath.Join(rootDir, fmt.Sprintf(".neev.backup_%s", timestamp))

	if err := copyDirectory(neevDir, backupDir); err != nil {
		return "", fmt.Errorf("failed to create backup: %w", err)
	}

	return backupDir, nil
}

// copyDirectory recursively copies a directory.
func copyDirectory(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate the relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		// Copy file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, 0644)
	})
}
