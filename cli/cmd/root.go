package cmd

import (
	"fmt"
	"os"

	"github.com/neev-kit/neev/core/errors"
	"github.com/neev-kit/neev/core/logger"
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
	// Initialize logger
	logger.Init()

	if err := rootCmd.Execute(); err != nil {
		// Check if it's a NeevError for user-friendly output
		if neevErr, ok := err.(*errors.NeevError); ok {
			fmt.Printf("Error: %v\n", neevErr)
			fmt.Printf("💡 %s\n", neevErr.GetSolutionHint())
		} else {
			fmt.Printf("Error: %v\n", err)
		}
		os.Exit(1)
	}
}

func init() {
	// Root command initialization
}
