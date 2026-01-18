package blueprint

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Lay archives a blueprint by moving it to the foundation archive and updating changelog
func Lay(blueprintName string) error {
	// Locate blueprint
	blueprintPath := filepath.Join(".neev", "blueprints", blueprintName)

	// Check if blueprint exists
	if _, err := os.Stat(blueprintPath); os.IsNotExist(err) {
		return fmt.Errorf("blueprint '%s' not found at %s", blueprintName, blueprintPath)
	}

	// Create archive directory if it doesn't exist
	archivePath := filepath.Join(".neev", "foundation", "archive")
	if err := os.MkdirAll(archivePath, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	// Create archive subdirectory for this blueprint
	blueprintArchivePath := filepath.Join(archivePath, blueprintName)
	if err := os.MkdirAll(blueprintArchivePath, 0755); err != nil {
		return fmt.Errorf("failed to create blueprint archive directory: %w", err)
	}

	// Move intent.md to archive
	intentSrc := filepath.Join(blueprintPath, "intent.md")
	intentDest := filepath.Join(blueprintArchivePath, "intent.md")
	if err := moveFile(intentSrc, intentDest); err != nil {
		return fmt.Errorf("failed to move intent.md: %w", err)
	}

	// Move architecture.md to archive
	archSrc := filepath.Join(blueprintPath, "architecture.md")
	archDest := filepath.Join(blueprintArchivePath, "architecture.md")
	if err := moveFile(archSrc, archDest); err != nil {
		return fmt.Errorf("failed to move architecture.md: %w", err)
	}

	// Update changelog
	if err := appendToChangelog(blueprintName); err != nil {
		return fmt.Errorf("failed to update changelog: %w", err)
	}

	// Delete the original blueprint folder
	if err := os.RemoveAll(blueprintPath); err != nil {
		return fmt.Errorf("failed to delete blueprint folder: %w", err)
	}

	return nil
}

// moveFile moves a file from src to dst, handling non-existent source gracefully
func moveFile(src, dst string) error {
	// Check if source exists
	if _, err := os.Stat(src); os.IsNotExist(err) {
		// File doesn't exist, skip it (not all blueprints may have all files)
		return nil
	}

	// Read the source file
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Write to destination
	if err := os.WriteFile(dst, content, 0644); err != nil {
		return fmt.Errorf("failed to write destination file: %w", err)
	}

	// Delete source
	if err := os.Remove(src); err != nil {
		return fmt.Errorf("failed to remove source file: %w", err)
	}

	return nil
}

// appendToChangelog adds an entry to the changelog
func appendToChangelog(blueprintName string) error {
	changelogPath := filepath.Join(".neev", "changelog.md")

	// Create changelog entry
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	entry := fmt.Sprintf("- laid blueprint '%s' on %s\n", blueprintName, timestamp)

	// Read existing content
	var content []byte
	if _, err := os.Stat(changelogPath); err == nil {
		// File exists, read it
		var readErr error
		content, readErr = os.ReadFile(changelogPath)
		if readErr != nil {
			return fmt.Errorf("failed to read changelog: %w", readErr)
		}
	}

	// Prepend header if file doesn't exist or is empty
	if len(content) == 0 {
		content = []byte("# Neev Blueprint Changelog\n\n")
	}

	// Append new entry
	content = append([]byte(entry), content...)

	// Write back
	if err := os.WriteFile(changelogPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write changelog: %w", err)
	}

	return nil
}
