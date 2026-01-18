package blueprint

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Draft creates a new blueprint folder with templates.
func Draft(name string) error {
	// Sanitize the name
	sanitized := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	blueprintPath := filepath.Join(".neev", "blueprints", sanitized)

	// Check if the blueprint already exists
	if _, err := os.Stat(blueprintPath); !os.IsNotExist(err) {
		return fmt.Errorf("blueprint already exists: %s", blueprintPath)
	}

	// Create the blueprint directory
	if err := os.MkdirAll(blueprintPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create blueprint directory: %w", err)
	}

	// Create intent.md and architecture.md with default templates
	files := []string{"intent.md", "architecture.md"}
	for _, file := range files {
		filePath := filepath.Join(blueprintPath, file)
		if err := os.WriteFile(filePath, []byte("# Template for "+file), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", filePath, err)
		}
	}

	return nil
}
