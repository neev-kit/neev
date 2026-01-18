package cmd

import (
	"fmt"
	"strings"

	"github.com/neev-kit/neev/core/blueprint"
	"github.com/spf13/cobra"
)

var draftCmd = &cobra.Command{
	Use:   "draft <title>",
	Short: "Draft a new blueprint",
	Long:  "Create a draft blueprint for your project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		sanitized := strings.ToLower(strings.ReplaceAll(name, " ", "-"))

		if err := blueprint.Draft(name); err != nil {
			fmt.Println("❌", err)
		} else {
			fmt.Printf("✅ Created blueprint at .neev/blueprints/%s\n", sanitized)
		}
	},
}

func init() {
	rootCmd.AddCommand(draftCmd)
}
