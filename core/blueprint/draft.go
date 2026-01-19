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

	// Create blueprint files with templates
	blueprintFiles := map[string]string{
		"intent.md":       "# Intent\n\nWhat and why - describe the purpose and motivation for this blueprint.",
		"architecture.md": "# Architecture\n\nHow it works - describe the architectural design and implementation details.",
		"api-spec.md":     "# API Specification\n\nAPI contracts - define the endpoints, request/response formats, and protocols.",
		"security.md":     "# Security Considerations\n\nSecurity considerations - document security best practices, threat models, and mitigations.",
	}

	for file, content := range blueprintFiles {
		filePath := filepath.Join(blueprintPath, file)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", filePath, err)
		}
	}

	// Create foundation directory and files if they don't exist
	foundationPath := filepath.Join(".neev", "foundation")
	if _, err := os.Stat(foundationPath); os.IsNotExist(err) {
		if err := os.MkdirAll(foundationPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create foundation directory: %w", err)
		}

		// Create foundation files with templates
		foundationFiles := map[string]string{
			"stack.md":      "# Technology Stack\n\nDescribe the technologies used in your project (e.g., \"We use Go, PostgreSQL, Redis\")",
			"principles.md": "# Design Principles\n\nDocument your core design principles (e.g., \"Security first, simplicity second\")",
			"patterns.md":   "# Patterns & Practices\n\nOutline your architectural patterns and practices (e.g., \"Repository pattern, dependency injection\")",
		}

		for file, content := range foundationFiles {
			filePath := filepath.Join(foundationPath, file)
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to create file %s: %w", filePath, err)
			}
		}
	}

	return nil
}
