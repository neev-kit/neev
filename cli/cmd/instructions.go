package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/neev-kit/neev/core/instructions"
	"github.com/spf13/cobra"
)

var instructionsCmd = &cobra.Command{
	Use:   "instructions",
	Short: "Generate GitHub Copilot instructions",
	Long: `Generate GitHub Copilot instructions based on the current foundation and active blueprints.
This creates or updates .github/copilot-instructions.md with context about your project
that helps Copilot provide better suggestions.

The instructions include:
- Foundation module summary
- Active blueprint intents
- Development guidelines

Run this command after adding new blueprints or updating your foundation.`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			errorStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("1"))
			fmt.Println(errorStyle.Render("❌ Failed to determine current working directory: " + err.Error()))
			return
		}

		// Generate and save instructions
		err = instructions.SaveCopilotInstructions(cwd)
		if err != nil {
			errorStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("1"))
			fmt.Println(errorStyle.Render("❌ Failed to generate instructions: " + err.Error()))
			return
		}

		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2"))
		fmt.Println(successStyle.Render("✅ Copilot instructions generated successfully"))
		fmt.Println()
		fmt.Println("📄 File: .github/copilot-instructions.md")
		fmt.Println()
		fmt.Println("💡 Copilot will now use this context to provide better suggestions.")
		fmt.Println("   Run this command again after updating blueprints or foundation.")
	},
}

func init() {
	rootCmd.AddCommand(instructionsCmd)
}
