package cmd

import (
	"fmt"

	"github.com/neev-kit/neev/core/bridge"
	"github.com/neev-kit/neev/core/instructions"
	"github.com/spf13/cobra"
)

var bridgeCmd = &cobra.Command{
	Use:   "bridge [flags]",
	Short: "Bridge to external systems",
	Long:  "Aggregate context for AI agents",
	Run: func(cmd *cobra.Command, args []string) {
		focus, _ := cmd.Flags().GetString("focus")
		withRemotes, _ := cmd.Flags().GetBool("with-remotes")
		claudeMode, _ := cmd.Flags().GetBool("claude")
		slashMode, _ := cmd.Flags().GetBool("slash")

		context, err := bridge.BuildContext(focus)
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		remoteContext := ""
		// If with-remotes flag is set, append remote contexts
		if withRemotes {
			remoteContext, err = bridge.BuildRemoteContext()
			if err != nil {
				fmt.Printf("Warning: Failed to include remotes: %v\n", err)
			}
		}

		// Format for Claude if requested
		if claudeMode {
			context = instructions.ClaudeContext(context, withRemotes, remoteContext)
		} else if withRemotes && remoteContext != "" {
			context += "\n\n" + remoteContext
		}

		// Format for slash command if requested
		if slashMode {
			context = bridge.FormatSlashCommand(context)
		}

		fmt.Println(context)
	},
}

func init() {
	bridgeCmd.Flags().StringP("focus", "f", "", "Focus on a specific context string")
	bridgeCmd.Flags().Bool("with-remotes", false, "Include synced remote foundations in context")
	bridgeCmd.Flags().Bool("claude", false, "Format output optimized for Claude AI")
	bridgeCmd.Flags().Bool("slash", false, "Format output for IDE slash commands")
	rootCmd.AddCommand(bridgeCmd)
}
