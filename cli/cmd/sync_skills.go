package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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
3. Updates native skill directories
4. Generates natural language fallback for unsupported tools

Example:
  neev sync-skills`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return syncSkills(cmd)
	},
}

func syncSkills(cmd *cobra.Command) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Use current directory as project root
	projectRoot := cwd

	// Check if .neev directory exists
	neevDir := filepath.Join(projectRoot, ".neev")
	if _, err := os.Stat(neevDir); os.IsNotExist(err) {
		return fmt.Errorf("not in a Neev project. Run 'neev init' first")
	}

	projectName := filepath.Base(projectRoot)

	// Load blueprints as skills
	blueprints, err := loadBlueprintsAsSkills(projectRoot)
	if err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	if len(blueprints) == 0 {
		fmt.Println("No blueprints found in .neev/blueprints/")
		return nil
	}

	// Generate skills
	generator := tools.NewSkillsGenerator(projectName, projectRoot)
	if err := generator.GenerateSkills(blueprints); err != nil {
		return fmt.Errorf("failed to generate skills: %w", err)
	}

	fmt.Println("\n✅ Skills generated successfully!")
	fmt.Println(generator.GenerateSummaryReport(blueprints))

	return nil
}

func loadBlueprintsAsSkills(projectRoot string) ([]tools.SkillContent, error) {
	blueprintsDir := filepath.Join(projectRoot, ".neev", "blueprints")

	if _, err := os.Stat(blueprintsDir); os.IsNotExist(err) {
		return nil, nil
	}

	var skills []tools.SkillContent
	entries, err := os.ReadDir(blueprintsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read blueprints: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !isMarkdownFile(entry.Name()) {
			continue
		}

		filePath := filepath.Join(blueprintsDir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Warning: could not read %s\n", entry.Name())
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

func isMarkdownFile(filename string) bool {
	return len(filename) > 3 && filename[len(filename)-3:] == ".md"
}

func filenameToTitle(filename string) string {
	title := filename[:len(filename)-3]
	return title
}

func extractDescription(content string) string {
	lines := []rune(content)
	description := ""

	for i, r := range lines {
		if r == '\n' {
			break
		}
		if r != '#' && (i == 0 || lines[i-1] != '#') {
			description += string(r)
		}
	}

	if len(description) > 80 {
		description = description[:80] + "..."
	}

	return description
}

var detectToolsCmd = &cobra.Command{
	Use:   "detect-tools",
	Short: "Detect installed AI tools",
	Long:  `Detect and list all installed AI tools on this system.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return detectTools(cmd)
	},
}

func detectTools(cmd *cobra.Command) error {
	fmt.Println("\n🔍 Detecting AI Tools")
	fmt.Println("═════════════════════════════════════════════════════════════")

	detectedTools := tools.DetectInstalledTools()

	if len(detectedTools) == 0 {
		fmt.Println("\nNo AI tools detected on this system.")
		return nil
	}

	fmt.Println("\nDetected Tools:")
	for _, tool := range detectedTools {
		if tool.Installed {
			fmt.Printf("\n✓ %s\n", tool.Name)
			fmt.Printf("  Config: %s\n", tool.Config.ConfigDir)
			fmt.Printf("  Skills: %s\n", tool.Config.SkillsDir)
		}
	}

	fmt.Printf("\nFound %d tool(s).\n\n", len(detectedTools))

	return nil
}

var skillsStatusCmd = &cobra.Command{
	Use:   "skills-status",
	Short: "Show status of generated skills",
	Long:  `Display the status of skills for each detected tool.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return skillsStatus(cmd)
	},
}

func skillsStatus(cmd *cobra.Command) error {
	fmt.Println("\n📊 Skills Status")
	fmt.Println("═════════════════════════════════════════════════════════════")

	detectedTools := tools.DetectInstalledTools()

	if len(detectedTools) == 0 {
		fmt.Println("\nNo tools detected.\n")
		return nil
	}

	for _, tool := range detectedTools {
		if tool.Installed {
			count := countSkillFiles(tool.Config.SkillsDir)
			fmt.Printf("\n✓ %s: %d skills\n", tool.Name, count)
			fmt.Printf("  Directory: %s\n", tool.Config.SkillsDir)
		}
	}

	fmt.Println()
	return nil
}

func countSkillFiles(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			if isMarkdownFile(entry.Name()) || isJSONFile(entry.Name()) {
				count++
			}
		}
	}

	return count
}

func isJSONFile(filename string) bool {
	return len(filename) > 5 && filename[len(filename)-5:] == ".json"
}

func init() {
	rootCmd.AddCommand(syncSkillsCmd)
	rootCmd.AddCommand(detectToolsCmd)
	rootCmd.AddCommand(skillsStatusCmd)
}
