package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/neev-kit/neev/core/slash"
	"github.com/spf13/cobra"
)

var slashCommandsCmd = &cobra.Command{
	Use:   "slash-commands [flags]",
	Short: "Manage slash commands for AI tools",
	Long:  "Configure and manage slash commands for AI coding assistants like Claude Code, Cursor, CodeBuddy, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		list, _ := cmd.Flags().GetBool("list")
		update, _ := cmd.Flags().GetBool("update")
		tool, _ := cmd.Flags().GetString("tool")

		cwd, _ := os.Getwd()

		if list {
			listSlashCommands()
		} else if update {
			updateSlashCommands(cwd)
		} else if tool != "" {
			showToolCommands(tool)
		} else {
			cmd.Help()
		}
	},
}

func listSlashCommands() {
	fmt.Println("📋 Available Neev Slash Commands:")
	fmt.Println()

	for cmd, details := range slash.DefaultSlashCommands {
		fmt.Printf("  /neev:%s\n", cmd)
		fmt.Printf("    %s\n\n", details.Description)
	}

	fmt.Println("🔧 Supported AI Tools:")
	for _, tool := range slash.SupportedAITools {
		fmt.Printf("  • %s\n", tool)
	}
}

func updateSlashCommands(cwd string) {
	projectName := filepath.Base(cwd)
	agentsMD := slash.GenerateAgentsMD(slash.SupportedAITools, projectName)
	agentsMDPath := filepath.Join(cwd, "AGENTS.md")

	if err := os.WriteFile(agentsMDPath, []byte(agentsMD), 0644); err != nil {
		fmt.Printf("❌ Failed to update AGENTS.md: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Updated AGENTS.md with latest slash commands")
}

func showToolCommands(toolName string) {
	manifest := slash.GenerateSlashCommandManifest(toolName)
	fmt.Println(manifest)
}

func init() {
	slashCommandsCmd.Flags().Bool("list", false, "List all available slash commands")
	slashCommandsCmd.Flags().Bool("update", false, "Update AGENTS.md with latest commands")
	slashCommandsCmd.Flags().String("tool", "", "Show commands for a specific AI tool")
	rootCmd.AddCommand(slashCommandsCmd)
}
