package cmd

import (
	"fmt"

	"github.com/neev-kit/neev/core/bridge"
	"github.com/spf13/cobra"
)

var bridgeCmd = &cobra.Command{
	Use:   "bridge [flags]",
	Short: "Bridge to external systems",
	Long:  "Aggregate context for AI agents",
	Run: func(cmd *cobra.Command, args []string) {
		focus, _ := cmd.Flags().GetString("focus")
		context, err := bridge.BuildContext(focus)
		if err != nil {
			fmt.Println("❌", err)
			return
		}
		fmt.Println(context)
	},
}

func init() {
	bridgeCmd.Flags().StringP("focus", "f", "", "Focus on a specific context string")
	rootCmd.AddCommand(bridgeCmd)
}
