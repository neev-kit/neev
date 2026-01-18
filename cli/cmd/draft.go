package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var draftCmd = &cobra.Command{
	Use:   "draft",
	Short: "Draft a new blueprint",
	Long:  "Create a draft blueprint for your project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📝 Drafting blueprint...")
	},
}

func init() {
	rootCmd.AddCommand(draftCmd)
}
