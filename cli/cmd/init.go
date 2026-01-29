package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/neev-kit/neev/core/foundation"
	"github.com/neev-kit/neev/core/slash"
	"github.com/neev-kit/neev/core/tools"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the Neev foundation",
	Long:  "Initialize a new Neev project by creating the foundation structure",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()

		// Stylized output using Lipgloss
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("14")). // Cyan
			Margin(0, 0, 1, 0)

		fmt.Println(headerStyle.Render("🏗️  Laying foundation in " + cwd))

		if err := foundation.Initialize(cwd); err != nil {
			errorStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("1")). // Red
				Margin(0, 0, 1, 0)

			fmt.Println(errorStyle.Render("❌ Error: " + err.Error()))
			os.Exit(1)
		}

		// Generate AGENTS.md for AI tool integration
		projectName := filepath.Base(cwd)
		agentsMD := slash.GenerateAgentsMD(slash.SupportedAITools, projectName)
		agentsMDPath := filepath.Join(cwd, "AGENTS.md")

		if err := os.WriteFile(agentsMDPath, []byte(agentsMD), 0644); err != nil {
			fmt.Printf("Warning: Failed to create AGENTS.md: %v\n", err)
		} else {
			fmt.Println("📋 Created AGENTS.md with slash command definitions")
		}

		// Verify Copilot prompt files were created
		copilotPromptsBasePath := filepath.Join(cwd, ".github", "prompts", "neev")
		if _, err := os.Stat(copilotPromptsBasePath); err == nil {
			fmt.Println("📝 Generated prompt files for GitHub Copilot")
		}

		// Verify Claude commands were created
		claudeBasePath := filepath.Join(cwd, ".claude", "commands", "neev")
		if _, err := os.Stat(claudeBasePath); err == nil {
			fmt.Println("🤖 Generated slash commands for Claude Code")
		}

		// Detect and initialize skills for installed tools
		fmt.Println("\n🔍 Detecting installed AI tools...")
		detectedTools := tools.DetectInstalledTools()

		if len(detectedTools) > 0 {
			fmt.Printf("✓ Found %d AI tool(s)\n", len(detectedTools))
			fmt.Println("\n💡 Tip: Run 'neev sync-skills' to generate skills for your tools.")
		} else {
			fmt.Println("ℹ️  No AI tools detected yet.")
			fmt.Println("    When you install Claude, Cursor, or Copilot,")
			fmt.Println("    run 'neev sync-skills' to generate skills.")
		}

		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")). // Green
			Margin(0, 0, 1, 0)

		fmt.Println(successStyle.Render("✅ Foundation laid successfully!"))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
