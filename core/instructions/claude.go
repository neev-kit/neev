package instructions

import (
	"strings"
)

// FormatForClaude formats context text with Claude-optimized structure
func FormatForClaude(context string) string {
	var builder strings.Builder

	builder.WriteString("# CONTEXT FOR CLAUDE\n\n")
	builder.WriteString("This document contains the complete project context in a format optimized for Claude.\n\n")
	builder.WriteString("---\n\n")

	// Add explicit section headers that Claude understands well
	builder.WriteString("## 📋 RULES AND CONSTRAINTS\n\n")
	builder.WriteString("- Follow all specifications exactly as defined\n")
	builder.WriteString("- Maintain consistency with existing code patterns\n")
	builder.WriteString("- Ask for clarification if requirements conflict\n")
	builder.WriteString("- Preserve existing functionality unless explicitly changing it\n\n")
	builder.WriteString("---\n\n")

	// Process the context to add clear section markers
	sections := strings.Split(context, "# ")
	for i, section := range sections {
		if i == 0 && strings.TrimSpace(section) == "" {
			continue
		}

		if i > 0 {
			// Add clear section separators
			builder.WriteString("## 📚 ")
			builder.WriteString(section)
			builder.WriteString("\n")
		} else {
			builder.WriteString(section)
		}
	}

	builder.WriteString("\n---\n\n")
	builder.WriteString("## 🎯 CURRENT TASK\n\n")
	builder.WriteString("Review the specifications above and implement features according to:\n")
	builder.WriteString("1. The ARCHITECTURE defined in foundation specs\n")
	builder.WriteString("2. The INTENT described in active blueprints\n")
	builder.WriteString("3. The RULES AND CONSTRAINTS at the top of this document\n\n")

	return builder.String()
}

// ClaudeContext wraps the standard context with Claude-optimized formatting
func ClaudeContext(standardContext string, includeRemotes bool, remoteContext string) string {
	var fullContext strings.Builder

	fullContext.WriteString(standardContext)
	
	if includeRemotes && remoteContext != "" {
		fullContext.WriteString("\n\n")
		fullContext.WriteString(remoteContext)
	}

	return FormatForClaude(fullContext.String())
}
