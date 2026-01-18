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
	// Split on lines starting with "# " to find markdown headers
	lines := strings.Split(context, "\n")
	inSection := false
	sawHeader := false
	var preHeaderLines []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Check if this is a markdown header (starts with # at beginning of line)
		if strings.HasPrefix(trimmed, "# ") {
			// If this is the first header and we have pre-header content,
			// emit it as a general context section before the first header.
			if !sawHeader && len(preHeaderLines) > 0 {
				builder.WriteString("## 📚 GENERAL CONTEXT\n")
				for _, preLine := range preHeaderLines {
					builder.WriteString(preLine)
					builder.WriteString("\n")
				}
				builder.WriteString("\n")
				preHeaderLines = nil
			} else if inSection {
				builder.WriteString("\n")
			}
			// Add clear section separator with emoji
			builder.WriteString("## 📚 ")
			builder.WriteString(strings.TrimPrefix(trimmed, "# "))
			builder.WriteString("\n")
			inSection = true
			sawHeader = true
		} else {
			if !sawHeader {
				// Buffer lines before the first header so they are not lost
				preHeaderLines = append(preHeaderLines, line)
			} else if inSection {
				builder.WriteString(line)
				builder.WriteString("\n")
			}
		}
	}
	
	// If no headers were found at all, emit any collected lines as general context.
	if !sawHeader && len(preHeaderLines) > 0 {
		builder.WriteString("## 📚 GENERAL CONTEXT\n")
		for _, preLine := range preHeaderLines {
			builder.WriteString(preLine)
			builder.WriteString("\n")
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
