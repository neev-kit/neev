package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/neev-kit/neev/core/foundation"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect the foundation for drift",
	Long:  "Check if the project structure matches the foundation specifications",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()

		warnings, err := foundation.Inspect(cwd)
		if err != nil {
			errorStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("1")).
				Margin(0, 0, 1, 0)

			fmt.Println(errorStyle.Render("❌ Inspection failed: " + err.Error()))
			return
		}

		if len(warnings) == 0 {
			successStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("2")).
				Margin(0, 0, 1, 0)

			fmt.Println(successStyle.Render("✅ Foundation is solid."))
			return
		}

		warningStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("3")).
			Margin(0, 0, 1, 0)

		fmt.Println(warningStyle.Render("⚠️  Foundation drift detected:"))
		fmt.Println()

		for _, warning := range warnings {
			fmt.Println("  " + warning)
		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
