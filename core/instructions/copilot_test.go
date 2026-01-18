package instructions

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopilotInstructions_NoFoundation(t *testing.T) {
	tmpDir := t.TempDir()

	instructions, err := CopilotInstructions(tmpDir)
	if err != nil {
		t.Fatalf("CopilotInstructions failed: %v", err)
	}

	if !strings.Contains(instructions, "GitHub Copilot Instructions") {
		t.Errorf("Expected header in instructions")
	}

	if !strings.Contains(instructions, "Development Guidelines") {
		t.Errorf("Expected development guidelines")
	}
}

func TestCopilotInstructions_WithFoundation(t *testing.T) {
	tmpDir := t.TempDir()
	foundationDir := filepath.Join(tmpDir, ".neev", "foundation")
	os.MkdirAll(foundationDir, 0755)

	// Create foundation modules
	os.WriteFile(filepath.Join(foundationDir, "auth.md"), []byte("# Auth Module"), 0644)
	os.WriteFile(filepath.Join(foundationDir, "api.md"), []byte("# API Module"), 0644)

	instructions, err := CopilotInstructions(tmpDir)
	if err != nil {
		t.Fatalf("CopilotInstructions failed: %v", err)
	}

	if !strings.Contains(instructions, "Foundation modules") {
		t.Errorf("Expected foundation modules section")
	}

	if !strings.Contains(instructions, "auth") {
		t.Errorf("Expected auth module in instructions")
	}

	if !strings.Contains(instructions, "api") {
		t.Errorf("Expected api module in instructions")
	}
}

func TestCopilotInstructions_WithBlueprints(t *testing.T) {
	tmpDir := t.TempDir()
	blueprintsDir := filepath.Join(tmpDir, ".neev", "blueprints")
	blueprintDir := filepath.Join(blueprintsDir, "user-auth")
	os.MkdirAll(blueprintDir, 0755)

	// Create blueprint with intent
	intentContent := `# User Authentication

Implement secure user authentication with JWT tokens.`
	os.WriteFile(filepath.Join(blueprintDir, "intent.md"), []byte(intentContent), 0644)

	instructions, err := CopilotInstructions(tmpDir)
	if err != nil {
		t.Fatalf("CopilotInstructions failed: %v", err)
	}

	if !strings.Contains(instructions, "Active Blueprints") {
		t.Errorf("Expected active blueprints section")
	}

	if !strings.Contains(instructions, "user-auth") {
		t.Errorf("Expected user-auth blueprint")
	}

	if !strings.Contains(instructions, "Intent") {
		t.Errorf("Expected intent section")
	}
}

func TestSaveCopilotInstructions(t *testing.T) {
	tmpDir := t.TempDir()
	foundationDir := filepath.Join(tmpDir, ".neev", "foundation")
	os.MkdirAll(foundationDir, 0755)

	os.WriteFile(filepath.Join(foundationDir, "test.md"), []byte("# Test"), 0644)

	err := SaveCopilotInstructions(tmpDir)
	if err != nil {
		t.Fatalf("SaveCopilotInstructions failed: %v", err)
	}

	// Verify file was created
	instructionsPath := filepath.Join(tmpDir, ".github", "copilot-instructions.md")
	if _, err := os.Stat(instructionsPath); os.IsNotExist(err) {
		t.Errorf("Instructions file was not created")
	}

	// Verify content
	content, err := os.ReadFile(instructionsPath)
	if err != nil {
		t.Fatalf("Failed to read instructions file: %v", err)
	}

	if !strings.Contains(string(content), "GitHub Copilot Instructions") {
		t.Errorf("Expected proper content in instructions file")
	}
}

func TestFormatForClaude(t *testing.T) {
	input := "# Project Foundation\n\nSome content here\n\n# Blueprints\n\nBlueprint content"

	formatted := FormatForClaude(input)

	if !strings.Contains(formatted, "CONTEXT FOR CLAUDE") {
		t.Errorf("Expected Claude header")
	}

	if !strings.Contains(formatted, "RULES AND CONSTRAINTS") {
		t.Errorf("Expected rules section")
	}

	if !strings.Contains(formatted, "CURRENT TASK") {
		t.Errorf("Expected current task section")
	}

	if !strings.Contains(formatted, "Project Foundation") {
		t.Errorf("Expected original content to be preserved")
	}
}

func TestClaudeContext(t *testing.T) {
	standardContext := "# Foundation\n\nContent"
	remoteContext := "# Remote\n\nRemote content"

	// Without remotes
	result := ClaudeContext(standardContext, false, remoteContext)
	if strings.Contains(result, "Remote content") {
		t.Errorf("Should not include remotes when includeRemotes is false")
	}

	// With remotes
	result = ClaudeContext(standardContext, true, remoteContext)
	if !strings.Contains(result, "Remote content") {
		t.Errorf("Should include remotes when includeRemotes is true")
	}

	// Verify Claude formatting applied
	if !strings.Contains(result, "CONTEXT FOR CLAUDE") {
		t.Errorf("Should have Claude formatting")
	}
}

func TestClaudeContext_EmptyRemote(t *testing.T) {
	standardContext := "# Foundation\n\nContent"

	result := ClaudeContext(standardContext, true, "")
	if !strings.Contains(result, "CONTEXT FOR CLAUDE") {
		t.Errorf("Should have Claude formatting even with empty remote")
	}
}
