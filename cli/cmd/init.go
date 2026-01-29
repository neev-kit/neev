package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/neev-kit/neev/core/foundation"
	"github.com/neev-kit/neev/core/slash"
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

		// Generate GitHub Copilot slash-commands.json
		slashCommandsJSON, err := slash.GenerateGitHubCopilotManifest(projectName)
		if err != nil {
			fmt.Printf("Warning: Failed to generate slash-commands manifest: %v\n", err)
		} else {
			slashCommandsPath := filepath.Join(cwd, ".github", "slash-commands.json")
			if err := os.MkdirAll(filepath.Dir(slashCommandsPath), 0755); err != nil {
				fmt.Printf("Warning: Failed to create .github directory: %v\n", err)
			} else if err := os.WriteFile(slashCommandsPath, []byte(slashCommandsJSON), 0644); err != nil {
				fmt.Printf("Warning: Failed to write slash-commands.json: %v\n", err)
			} else {
				fmt.Println("🔗 Registered slash commands with GitHub Copilot")
			}
		}

		// Generate GitHub Copilot prompt files
		copilotPrompts := slash.GenerateGitHubCopilotPrompts(projectName)
		copilotPromptsBasePath := filepath.Join(cwd, ".github", "prompts", "neev")
		if err := os.MkdirAll(copilotPromptsBasePath, 0755); err != nil {
			fmt.Printf("Warning: Failed to create .github/prompts/neev directory: %v\n", err)
		} else {
			for fileName, content := range copilotPrompts {
				filePath := filepath.Join(copilotPromptsBasePath, fileName)
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					fmt.Printf("Warning: Failed to write GitHub Copilot prompt file %s: %v\n", fileName, err)
				}
			}
			fmt.Println("📝 Generated prompt files for GitHub Copilot")
		}

		// Generate Claude Code slash command files
		claudeCommands := slash.GenerateClaudeSlashCommands(projectName)
		claudeBasePath := filepath.Join(cwd, ".claude", "commands", "neev")
		if err := os.MkdirAll(claudeBasePath, 0755); err != nil {
			fmt.Printf("Warning: Failed to create .claude/commands/neev directory: %v\n", err)
		} else {
			for fileName, content := range claudeCommands {
				filePath := filepath.Join(claudeBasePath, fileName)
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					fmt.Printf("Warning: Failed to write Claude command file %s: %v\n", fileName, err)
				}
			}
			fmt.Println("🤖 Generated slash commands for Claude Code")
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
