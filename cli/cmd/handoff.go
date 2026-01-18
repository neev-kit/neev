package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/neev-kit/neev/core/bridge"
	"github.com/spf13/cobra"
)

var handoffCmd = &cobra.Command{
	Use:   "handoff <role>",
	Short: "Create a handoff prompt for an agent role",
	Long:  "Generate a structured handoff prompt with role-specific instructions from .neev/agents/<role>.md",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		role := args[0]

		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error: could not determine current working directory: %v\n", err)
			return
		}

		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("14")).
			Margin(0, 0, 1, 0)

		fmt.Println(headerStyle.Render(fmt.Sprintf("🤝 Generating handoff for role: %s", role)))

		// Build base context
		context, err := bridge.BuildContext("")
		if err != nil {
			errorStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("1")).
				Margin(0, 0, 1, 0)

			fmt.Println(errorStyle.Render("❌ Failed to build context: " + err.Error()))
			return
		}

		// Try to load role-specific instructions
		agentInstructions := ""
		agentPath := filepath.Join(cwd, ".neev", "agents", role+".md")

		if _, err := os.Stat(agentPath); err == nil {
			data, err := os.ReadFile(agentPath)
			if err == nil {
				agentInstructions = string(data)
				infoStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("14"))
				fmt.Println(infoStyle.Render(fmt.Sprintf("ℹ️  Loaded instructions from: %s", agentPath)))
			}
		} else {
			warnStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("3"))
			fmt.Println(warnStyle.Render(fmt.Sprintf("⚠️  No instructions found at: %s", agentPath)))
		}

		// Generate handoff prompt
		prompt := bridge.FormatHandoffPrompt(role, context, agentInstructions)

		// Optionally wrap in markdown fence
		markdown, _ := cmd.Flags().GetBool("markdown")
		if markdown {
			prompt = bridge.FormatHandoffMarkdown(prompt)
		}

		fmt.Println()
		fmt.Println(prompt)
	},
}

func init() {
	handoffCmd.Flags().Bool("markdown", true, "Wrap output in markdown code fence for copy-paste")
	rootCmd.AddCommand(handoffCmd)
}
