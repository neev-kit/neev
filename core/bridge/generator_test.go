package bridge

import (
	"strings"
	"testing"
)

func TestFormatSlashCommand(t *testing.T) {
	context := "Some context\nWith multiple lines"
	result := FormatSlashCommand(context)

	if !strings.Contains(result, "```markdown") {
		t.Error("Expected markdown code fence in result")
	}

	if !strings.Contains(result, "/context") {
		t.Error("Expected /context slash command in result")
	}

	if !strings.Contains(result, context) {
		t.Error("Expected context to be included in result")
	}

	// Verify it starts and ends with proper formatting
	if !strings.HasPrefix(strings.TrimSpace(result), "```markdown") {
		t.Error("Expected result to start with markdown fence")
	}

	if !strings.HasSuffix(strings.TrimSpace(result), "```") {
		t.Error("Expected result to end with markdown fence")
	}
}

func TestFormatSlashCommandEmpty(t *testing.T) {
	result := FormatSlashCommand("")

	if !strings.Contains(result, "```markdown") {
		t.Error("Expected markdown code fence even with empty context")
	}

	if !strings.Contains(result, "/context") {
		t.Error("Expected /context in result")
	}
}

func TestFormatHandoffPrompt(t *testing.T) {
	role := "Developer"
	context := "Project context here"
	instructions := "Follow these instructions"

	result := FormatHandoffPrompt(role, context, instructions)

	if !strings.Contains(result, "/neev-handoff-developer") {
		t.Error("Expected handoff slash command with role")
	}

	if !strings.Contains(result, "Developer") {
		t.Error("Expected role to be in result")
	}

	if !strings.Contains(result, "## Context") {
		t.Error("Expected Context section")
	}

	if !strings.Contains(result, context) {
		t.Error("Expected context in result")
	}

	if !strings.Contains(result, "## Instructions") {
		t.Error("Expected Instructions section")
	}

	if !strings.Contains(result, instructions) {
		t.Error("Expected instructions in result")
	}

	if !strings.Contains(result, "spec-driven development") {
		t.Error("Expected Neev footer text")
	}
}

func TestFormatHandoffPromptNoInstructions(t *testing.T) {
	role := "Architect"
	context := "Architecture context"

	result := FormatHandoffPrompt(role, context, "")

	if !strings.Contains(result, "/neev-handoff-architect") {
		t.Error("Expected handoff slash command with role")
	}

	if !strings.Contains(result, "## Context") {
		t.Error("Expected Context section")
	}

	if strings.Contains(result, "## Instructions") {
		t.Error("Did not expect Instructions section when none provided")
	}

	if !strings.Contains(result, context) {
		t.Error("Expected context in result")
	}
}

func TestFormatHandoffPromptEmptyRole(t *testing.T) {
	result := FormatHandoffPrompt("", "context", "instructions")

	// Should still work with empty role, it gets converted to lowercase
	if !strings.Contains(result, "/neev-handoff-") {
		t.Error("Expected handoff slash command even with empty role")
	}
}

func TestFormatHandoffPromptSpecialCharactersInRole(t *testing.T) {
	result := FormatHandoffPrompt("QA Tester", "context", "")

	// Role should be converted to lowercase with spaces
	if !strings.Contains(result, "/neev-handoff-qa tester") {
		t.Error("Expected role to be lowercase in slash command")
	}
}

func TestFormatHandoffMarkdown(t *testing.T) {
	prompt := "Some prompt content"
	result := FormatHandoffMarkdown(prompt)

	if !strings.HasPrefix(strings.TrimSpace(result), "```markdown") {
		t.Error("Expected to start with markdown fence")
	}

	if !strings.HasSuffix(strings.TrimSpace(result), "```") {
		t.Error("Expected to end with markdown fence")
	}

	if !strings.Contains(result, prompt) {
		t.Error("Expected prompt content in result")
	}
}

func TestFormatHandoffMarkdownEmpty(t *testing.T) {
	result := FormatHandoffMarkdown("")

	if !strings.Contains(result, "```markdown") {
		t.Error("Expected markdown fence")
	}
}

func TestFormatHandoffMarkdownMultiline(t *testing.T) {
	prompt := "Line 1\nLine 2\nLine 3"
	result := FormatHandoffMarkdown(prompt)

	if !strings.Contains(result, "Line 1") {
		t.Error("Expected Line 1 in result")
	}

	if !strings.Contains(result, "Line 2") {
		t.Error("Expected Line 2 in result")
	}

	if !strings.Contains(result, "Line 3") {
		t.Error("Expected Line 3 in result")
	}
}

func TestFormatHandoffPromptComplexContext(t *testing.T) {
	role := "Code Reviewer"
	context := "# Foundation\n\n## Section 1\nContent here\n\n## Section 2\nMore content"
	instructions := "1. Check for code quality\n2. Verify tests\n3. Approve or request changes"

	result := FormatHandoffPrompt(role, context, instructions)

	if !strings.Contains(result, "Code Reviewer") {
		t.Error("Expected role in result")
	}

	if !strings.Contains(result, "# Foundation") {
		t.Error("Expected complex context in result")
	}

	if !strings.Contains(result, "1. Check for code quality") {
		t.Error("Expected numbered instructions in result")
	}
}
