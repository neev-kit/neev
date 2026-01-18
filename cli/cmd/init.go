package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/neev-kit/neev/core/foundation"
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
