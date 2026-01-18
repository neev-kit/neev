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
		withRemotes, _ := cmd.Flags().GetBool("with-remotes")
		
		context, err := bridge.BuildContext(focus)
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		// If with-remotes flag is set, append remote contexts
		if withRemotes {
			remoteContext, err := bridge.BuildRemoteContext()
			if err != nil {
				fmt.Printf("Warning: Failed to include remotes: %v\n", err)
			} else if remoteContext != "" {
				context += "\n\n" + remoteContext
			}
		}

		fmt.Println(context)
	},
}

func init() {
	bridgeCmd.Flags().StringP("focus", "f", "", "Focus on a specific context string")
	bridgeCmd.Flags().Bool("with-remotes", false, "Include synced remote foundations in context")
	rootCmd.AddCommand(bridgeCmd)
}
