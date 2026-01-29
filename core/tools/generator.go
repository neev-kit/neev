package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SkillsGenerator orchestrates skill generation for all detected tools
type SkillsGenerator struct {
	projectName string
	projectRoot string
	tools       []Tool
	adapters    []Adapter
}

// NewSkillsGenerator creates a new skills generator
func NewSkillsGenerator(projectName string, projectRoot string) *SkillsGenerator {
	return &SkillsGenerator{
		projectName: projectName,
		projectRoot: projectRoot,
	}
}

// GenerateSkills generates skills for all detected tools
func (sg *SkillsGenerator) GenerateSkills(blueprints []SkillContent) error {
	// Detect installed tools
	sg.tools = DetectInstalledTools()
	sg.adapters = GetAdapters(sg.tools)

	// If no tools detected, generate fallback
	if !HasAnyTool(sg.tools) {
		if err := sg.generateFallbackDocumentation(blueprints); err != nil {
			return fmt.Errorf("failed to generate fallback documentation: %w", err)
		}
		return nil
	}

	// Generate skills for each tool
	for _, adapter := range sg.adapters {
		if err := sg.generateForTool(adapter, blueprints); err != nil {
			return fmt.Errorf("failed to generate for %s: %w", adapter.Name(), err)
		}
	}

	// Generate index file
	if err := sg.generateIndexFile(blueprints); err != nil {
		return fmt.Errorf("failed to generate index: %w", err)
	}

	return nil
}

// generateForTool generates skills for a specific tool
func (sg *SkillsGenerator) generateForTool(adapter Adapter, blueprints []SkillContent) error {
	tool := FindTool(sg.getToolTypeFromAdapter(adapter), sg.tools)
	if tool == nil || !tool.Installed {
		return nil
	}

	skillsDir := tool.Config.SkillsDir

	// Create skills directory
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		return fmt.Errorf("failed to create skills directory: %w", err)
	}

	// Write each skill
	for _, blueprint := range blueprints {
		if err := WriteSkillToFile(adapter, blueprint, skillsDir); err != nil {
			return err
		}
	}

	// Generate config file
	configContent, err := adapter.GenerateConfigFile(sg.projectName, blueprints)
	if err != nil {
		return err
	}

	configPath := filepath.Join(skillsDir, "README.md")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return err
	}

	return nil
}

// generateFallbackDocumentation generates fallback documentation
func (sg *SkillsGenerator) generateFallbackDocumentation(blueprints []SkillContent) error {
	skillsDir := filepath.Join(sg.projectRoot, ".neev", "skills")

	// Create directory
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		return fmt.Errorf("failed to create skills directory: %w", err)
	}

	// Create fallback tool for documentation
	fallbackTool := Tool{
		Name: "Generic Tool",
		Config: ToolConfig{
			SkillsDir: skillsDir,
			ConfigDir: filepath.Join(sg.projectRoot, ".neev"),
		},
	}

	adapter := NewFallbackAdapter(&fallbackTool)

	// Write each skill
	for _, blueprint := range blueprints {
		if err := WriteSkillToFile(adapter, blueprint, skillsDir); err != nil {
			return err
		}
	}

	// Generate config file
	configContent, err := adapter.GenerateConfigFile(sg.projectName, blueprints)
	if err != nil {
		return err
	}

	configPath := filepath.Join(skillsDir, "README.md")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return err
	}

	return nil
}

// generateIndexFile generates an index file for all skills
func (sg *SkillsGenerator) generateIndexFile(blueprints []SkillContent) error {
	indexContent := fmt.Sprintf("# Skills Index for %s\n\n", sg.projectName)
	indexContent += "## Generated Skills\n\n"

	for i, skill := range blueprints {
		indexContent += fmt.Sprintf("%d. **%s** (%s)\n   - Description: %s\n   - Version: %s\n\n",
			i+1, skill.Name, skill.Type, skill.Description, skill.Version)
	}

	indexContent += "## Tools\n\n"
	for _, tool := range sg.tools {
		if tool.Installed {
			indexContent += fmt.Sprintf("- %s: %s/skills/\n", tool.Name, tool.Config.SkillsDir)
		}
	}

	neevDir := filepath.Join(sg.projectRoot, ".neev")
	indexPath := filepath.Join(neevDir, "SKILLS_INDEX.md")

	if err := os.WriteFile(indexPath, []byte(indexContent), 0644); err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}

	return nil
}

// GenerateSummaryReport generates a summary report of generated skills
func (sg *SkillsGenerator) GenerateSummaryReport(blueprints []SkillContent) string {
	report := fmt.Sprintf("Skills Generation Report for %s\n", sg.projectName)
	report += strings.Repeat("=", len(report)) + "\n\n"

	report += fmt.Sprintf("Detected Tools: %d\n", len(sg.tools))
	for _, tool := range sg.tools {
		status := "Not Installed"
		if tool.Installed {
			status = "Installed"
		}
		report += fmt.Sprintf("  - %s: %s\n", tool.Name, status)
	}

	report += fmt.Sprintf("\nBlueprints Converted: %d\n", len(blueprints))
	for _, bp := range blueprints {
		report += fmt.Sprintf("  - %s\n", bp.Name)
	}

	report += fmt.Sprintf("\nAdapters Used: %d\n", len(sg.adapters))
	for _, adapter := range sg.adapters {
		metadata := adapter.GetMetadata()
		reportFormat, _ := metadata["formatType"].(string)
		report += fmt.Sprintf("  - %s (%s)\n", adapter.Name(), reportFormat)
	}

	return report
}

// getToolTypeFromAdapter gets the tool type from an adapter
func (sg *SkillsGenerator) getToolTypeFromAdapter(adapter Adapter) ToolType {
	switch adapter.Name() {
	case "Cursor":
		return ToolCursor
	case "Claude":
		return ToolClaude
	case "GitHub Copilot":
		return ToolCopilot
	case "Codeium":
		return ToolCodeium
	case "Supabase":
		return ToolSupabase
	case "Perplexity":
		return ToolPerplexity
	default:
		return ToolUnknown
	}
}
