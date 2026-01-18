package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "neev",
	Short: "Neev - The blueprint orchestration tool",
	Long:  "Neev is a tool for managing blueprints and automating project initialization.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Neev - Blueprint Orchestration Tool")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Root command initialization
}
