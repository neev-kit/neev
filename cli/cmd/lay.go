package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/neev-kit/neev/core/blueprint"
	"github.com/spf13/cobra"
)

var layCmd = &cobra.Command{
	Use:   "lay <blueprint_name>",
	Short: "Archive a blueprint into the foundation",
	Long:  "Move a completed blueprint to the foundation archive and update changelog",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		blueprintName := args[0]

		if err := blueprint.Lay(blueprintName); err != nil {
			errorStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("1")).
				Margin(0, 0, 1, 0)

			fmt.Println(errorStyle.Render("❌ Failed to lay blueprint: " + err.Error()))
			os.Exit(1)
		}

		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")).
			Margin(0, 0, 1, 0)

		fmt.Println(successStyle.Render(fmt.Sprintf("🧱 Laid blueprint '%s' into the foundation.", blueprintName)))
	},
}

func init() {
	rootCmd.AddCommand(layCmd)
}
