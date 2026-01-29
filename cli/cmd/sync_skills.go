package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/neev-kit/neev/core/foundation"
	"github.com/neev-kit/neev/core/tools"
	"github.com/spf13/cobra"
)

var syncSkillsCmd = &cobra.Command{
	Use:   "sync-skills",
	Short: "Regenerate skills for all detected AI tools",
	Long: `Regenerate and synchronize AI skills across all installed tools.

This command:
1. Detects installed AI tools (Claude, Cursor, Copilot, etc.)
2. Regenerates skills from your blueprints
3. Updates native skill directories (.claude/skills/, .cursor/skills/, etc.)
4. Generates natural language fallback for unsupported tools

Skills are generated in tool-native formats:
- Cursor: JSON format
- Claude: Markdown format
- GitHub Copilot: Markdown format
- Codeium: JSON format
- Others: Natural language markdown

Example:
  neev sync-skills                    # Regenerate all skills
  neev sync-skills --verbose          # Show detailed output`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return syncSkills(cmd)
	},
}

// syncSkills synchronizes skills across all detected tools
func syncSkills(cmd *cobra.Command) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Find project root
	projectRoot, err := foundation.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a Neev project. Run 'neev init' first: %w", err)
	}

	projectName := filepath.Base(projectRoot)

	// Detect tools
	detectedTools := tools.DetectInstalledTools()

	fmt.Println("\n🔍 Tool Detection")
	fmt.Println("─────────────────────────────────────────────────────────────")
	tools.PrintDetectionSummary(detectedTools)

	// Load blueprints as skills
	blueprints, err := loadBlueprintsAsSkills(projectRoot)
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not load all blueprints: %v\n", err)
		blueprints = []tools.SkillContent{} // Continue with empty blueprints
	}

	if len(blueprints) == 0 {
		fmt.Println("\n⚠️  No blueprints found. Skills will be empty.")
		fmt.Println("    Create blueprints in .neev/blueprints/ first.")
	}

	// Generate skills
	generator := tools.NewSkillsGenerator(projectName, projectRoot, detectedTools)
	if err := generator.GenerateSkills(blueprints); err != nil {
		return fmt.Errorf("failed to generate skills: %w", err)
	}

	// Print summary
	fmt.Println(generator.GenerateSummaryReport(blueprints))

	return nil
}

// loadBlueprintsAsSkills loads blueprint files and converts them to skills
func loadBlueprintsAsSkills(projectRoot string) ([]tools.SkillContent, error) {
	blueprintsDir := filepath.Join(projectRoot, ".neev", "blueprints")

	// Check if blueprints directory exists
	if _, err := os.Stat(blueprintsDir); os.IsNotExist(err) {
		return nil, nil // No blueprints, return empty slice
	}

	var skills []tools.SkillContent
	entries, err := os.ReadDir(blueprintsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read blueprints directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !isMarkdownFile(entry.Name()) {
			continue
		}

		filePath := filepath.Join(blueprintsDir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("⚠️  Could not read %s: %v\n", entry.Name(), err)
			continue
		}

		skill := tools.SkillContent{
			Name:        filenameToTitle(entry.Name()),
			Description: extractDescription(string(content)),
			Content:     string(content),
			Type:        "blueprint",
			Language:    "markdown",
			Version:     "1.0",
		}

		skills = append(skills, skill)
	}

	return skills, nil
}

// isMarkdownFile checks if a file is a markdown file
func isMarkdownFile(filename string) bool {
	return len(filename) > 3 && filename[len(filename)-3:] == ".md"
}

// filenameToTitle converts a filename to a title
func filenameToTitle(filename string) string {
	// Remove .md extension
	title := filename[:len(filename)-3]
	// Replace hyphens and underscores with spaces
	for i, c := range title {
		if c == '-' || c == '_' {
			title = title[:i] + " " + title[i+1:]
		}
	}
	return title
}

// extractDescription extracts the first line or heading from markdown content
func extractDescription(content string) string {
	lines := make([]rune, 0)
	inHeading := false

	for i, r := range content {
		if r == '\n' {
			if inHeading || i > 0 {
				break
			}
			continue
		}
		if i < len(content)-1 && content[i:i+1] == "#" {
			inHeading = true
		}
		if inHeading && r != '#' && r != ' ' {
			inHeading = false
			lines = append(lines, r)
		} else if !inHeading && r != '#' {
			lines = append(lines, r)
		}
	}

	description := string(lines)
	if len(description) > 80 {
		description = description[:80] + "..."
	}

	return description
}

// detectToolsCmd shows all detected AI tools
var detectToolsCmd = &cobra.Command{
	Use:   "detect-tools",
	Short: "Detect installed AI tools on this system",
	Long: `Detect and list all installed AI tools that Neev can integrate with.

Supported tools:
- Claude (VS Code extension or standalone app)
- Cursor IDE
- GitHub Copilot (VS Code)
- Codeium
- Supabase
- Perplexity AI

Example:
  neev detect-tools                    # List all detected tools
  neev detect-tools --verbose          # Show detailed information`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return detectTools(cmd)
	},
}

// detectTools detects and displays available tools
func detectTools(cmd *cobra.Command) error {
	fmt.Println("\n🔍 Detecting AI Tools")
	fmt.Println("═════════════════════════════════════════════════════════════")

	detectedTools := tools.DetectInstalledTools()

	if len(detectedTools) == 0 {
		fmt.Println("\n⚠️  No AI tools detected on this system.\n")
		fmt.Println("Supported tools:")
		fmt.Println("  - Claude (VS Code extension or standalone app)")
		fmt.Println("  - Cursor IDE")
		fmt.Println("  - GitHub Copilot (VS Code)")
		fmt.Println("  - Codeium")
		fmt.Println("  - Supabase")
		fmt.Println("  - Perplexity AI")
		fmt.Println("\nInstall one of these tools to enable skill generation.")
		return nil
	}

	fmt.Println("\nDetected Tools:")
	fmt.Println("───────────────────────────────────────────────────────────────")

	for _, tool := range detectedTools {
		if tool.Installed {
			fmt.Printf("\n✓ %s\n", tool.Name)
			fmt.Printf("  Config Dir:  %s\n", tool.Config.ConfigDir)
			fmt.Printf("  Skills Dir:  %s\n", tool.Config.SkillsDir)
			fmt.Printf("  Native:      %v\n", tool.Config.Native)
			if tool.Path != "" {
				fmt.Printf("  Location:    %s\n", tool.Path)
			}
		}
	}

	fmt.Println("\n═════════════════════════════════════════════════════════════")
	fmt.Printf("\nFound %d AI tool(s). Ready for skill generation!\n", len(detectedTools))
	fmt.Println("Run 'neev sync-skills' to generate skills for these tools.\n")

	return nil
}

// skillsStatusCmd shows the status of generated skills
var skillsStatusCmd = &cobra.Command{
	Use:   "skills-status",
	Short: "Show status of generated skills",
	Long: `Display the status of skills generated for each detected AI tool.

Shows:
- Which tools have detected skills
- Location of skill directories
- Number of available skills
- Last generation time

Example:
  neev skills-status                   # Show skills status`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return skillsStatus(cmd)
	},
}

// skillsStatus shows the status of generated skills
func skillsStatus(cmd *cobra.Command) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Find project root
	projectRoot, err := foundation.FindProjectRoot(cwd)
	if err != nil {
		return fmt.Errorf("not in a Neev project: %w", err)
	}

	fmt.Println("\n📊 Skills Status Report")
	fmt.Println("═════════════════════════════════════════════════════════════")

	detectedTools := tools.DetectInstalledTools()

	if len(detectedTools) == 0 {
		fmt.Println("\nNo tools detected. Checking for fallback skills...")

		fallbackDir := filepath.Join(projectRoot, ".neev", "skills")
		if _, err := os.Stat(fallbackDir); os.IsNotExist(err) {
			fmt.Println("✗ No skills generated.")
			fmt.Println("\nRun 'neev sync-skills' to generate skills.")
			return nil
		}

		countSkills := countSkillFiles(fallbackDir)
		fmt.Printf("\n✓ Fallback skills: %d files in %s\n", countSkills, fallbackDir)
		return nil
	}

	allHaveSkills := true
	for _, tool := range detectedTools {
		if tool.Installed {
			skillCount := countSkillFiles(tool.Config.SkillsDir)
			if skillCount > 0 {
				fmt.Printf("\n✓ %s\n", tool.Name)
				fmt.Printf("  Skills Dir: %s\n", tool.Config.SkillsDir)
				fmt.Printf("  Skill Count: %d\n", skillCount)
			} else {
				fmt.Printf("\n✗ %s\n", tool.Name)
				fmt.Printf("  Skills Dir: %s\n", tool.Config.SkillsDir)
				fmt.Printf("  Skill Count: 0 (not generated)\n")
				allHaveSkills = false
			}
		}
	}

	fmt.Println("\n═════════════════════════════════════════════════════════════")

	if allHaveSkills {
		fmt.Println("\n✓ All tools have generated skills.")
	} else {
		fmt.Println("\n⚠️  Some tools are missing skills.")
		fmt.Println("Run 'neev sync-skills' to generate skills for all tools.\n")
	}

	return nil
}

// countSkillFiles counts the number of skill files in a directory
func countSkillFiles(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() && (isMarkdownFile(entry.Name()) || isJSONFile(entry.Name())) {
			count++
		}
	}

	return count
}

// isJSONFile checks if a file is a JSON file
func isJSONFile(filename string) bool {
	return len(filename) > 5 && filename[len(filename)-5:] == ".json"
}

func init() {
	rootCmd.AddCommand(syncSkillsCmd)
	rootCmd.AddCommand(detectToolsCmd)
	rootCmd.AddCommand(skillsStatusCmd)
}
