package instructions

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// CopilotInstructions generates GitHub Copilot instructions based on foundation and active blueprints
func CopilotInstructions(rootDir string) (string, error) {
	var builder strings.Builder

	builder.WriteString("# GitHub Copilot Instructions\n\n")
	builder.WriteString("This project uses Neev for spec-driven development.\n\n")

	// Read foundation summary
	foundationPath := filepath.Join(rootDir, ".neev", "foundation")
	if _, err := os.Stat(foundationPath); err == nil {
		builder.WriteString("## Project Foundation\n\n")
		
		entries, err := os.ReadDir(foundationPath)
		if err != nil {
			return "", fmt.Errorf("failed to read foundation: %w", err)
		}

		// List foundation modules
		builder.WriteString("Foundation modules:\n")
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				moduleName := strings.TrimSuffix(entry.Name(), ".md")
				builder.WriteString(fmt.Sprintf("- %s\n", moduleName))
			}
		}
		builder.WriteString("\n")
	}

	// Read active blueprints
	blueprintsPath := filepath.Join(rootDir, ".neev", "blueprints")
	if _, err := os.Stat(blueprintsPath); err == nil {
		entries, err := os.ReadDir(blueprintsPath)
		if err != nil {
			return "", fmt.Errorf("failed to read blueprints: %w", err)
		}

		if len(entries) > 0 {
			builder.WriteString("## Active Blueprints\n\n")
			
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}

				blueprintName := entry.Name()
				blueprintPath := filepath.Join(blueprintsPath, blueprintName)

				// Read intent if available
				intentPath := filepath.Join(blueprintPath, "intent.md")
				if intentData, err := ioutil.ReadFile(intentPath); err == nil {
					builder.WriteString(fmt.Sprintf("### Blueprint: %s\n\n", blueprintName))
					
					// Extract first paragraph or first 200 chars as summary
					content := string(intentData)
					lines := strings.Split(content, "\n")
					summary := ""
					for _, line := range lines {
						trimmed := strings.TrimSpace(line)
						if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
							summary = trimmed
							break
						}
					}
					if summary == "" && len(content) > 0 {
						summary = content
						if len(summary) > 200 {
							summary = summary[:200] + "..."
						}
					}
					
					builder.WriteString(fmt.Sprintf("**Intent**: %s\n\n", summary))
				}
			}
		}
	}

	// Add general guidelines
	builder.WriteString("## Development Guidelines\n\n")
	builder.WriteString("- Follow the architecture defined in foundation specifications\n")
	builder.WriteString("- Implement features according to blueprint intent and architecture\n")
	builder.WriteString("- Use `neev bridge` to get full context for complex tasks\n")
	builder.WriteString("- Run `neev inspect` to check for drift between specs and code\n")

	return builder.String(), nil
}

// SaveCopilotInstructions saves Copilot instructions to .github/copilot-instructions.md
func SaveCopilotInstructions(rootDir string) error {
	instructions, err := CopilotInstructions(rootDir)
	if err != nil {
		return err
	}

	githubDir := filepath.Join(rootDir, ".github")
	if err := os.MkdirAll(githubDir, 0755); err != nil {
		return fmt.Errorf("failed to create .github directory: %w", err)
	}

	instructionsPath := filepath.Join(githubDir, "copilot-instructions.md")
	if err := os.WriteFile(instructionsPath, []byte(instructions), 0644); err != nil {
		return fmt.Errorf("failed to write instructions file: %w", err)
	}

	return nil
}
