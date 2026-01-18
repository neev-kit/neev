package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var bridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Bridge to external systems",
	Long:  "Connect Neev with external systems and services",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🌉 Building bridge...")
	},
}

func init() {
	rootCmd.AddCommand(bridgeCmd)
}
