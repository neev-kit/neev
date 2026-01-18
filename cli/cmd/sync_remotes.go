package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/neev-kit/neev/core/config"
	"github.com/neev-kit/neev/core/remotes"
	"github.com/spf13/cobra"
)

var (
	syncJSON bool
)

var syncRemotesCmd = &cobra.Command{
	Use:   "sync-remotes",
	Short: "Sync remote foundations to local .neev/remotes",
	Long: `Synchronize remote foundation sources defined in neev.yaml to the local
.neev/remotes directory. This allows you to reference external foundations
from other repositories.

Example neev.yaml configuration:

  remotes:
    - name: api
      path: "../my-api-repo/.neev/foundation"
      public_only: true
    - name: shared
      path: "../shared-lib/.neev/foundation"
      public_only: false
`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			errorStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("1"))
			fmt.Println(errorStyle.Render("❌ Failed to determine current working directory: " + err.Error()))
			return
		}

		// Load config
		cfg, err := config.LoadConfig(cwd)
		if err != nil {
			errorStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("1"))
			fmt.Println(errorStyle.Render("❌ Failed to load config: " + err.Error()))
			return
		}

		// Check if remotes are configured
		if len(cfg.Remotes) == 0 {
			infoStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("6"))
			fmt.Println(infoStyle.Render("ℹ️  No remotes configured in neev.yaml"))
			fmt.Println()
			fmt.Println("To add remotes, edit neev.yaml and add:")
			fmt.Print(`
remotes:
  - name: api
    path: "../my-api-repo/.neev/foundation"
    public_only: true
`)
			return
		}

		// Convert config remotes to remotes package type
		remotesToSync := make([]remotes.Remote, len(cfg.Remotes))
		for i, r := range cfg.Remotes {
			remotesToSync[i] = remotes.Remote{
				Name:       r.Name,
				Path:       r.Path,
				PublicOnly: r.PublicOnly,
			}
		}

		// Sync remotes
		result, err := remotes.Sync(cwd, remotesToSync)
		if err != nil {
			errorStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("1"))
			fmt.Println(errorStyle.Render("❌ Sync failed: " + err.Error()))
			return
		}

		// Output in JSON if requested
		if syncJSON {
			jsonData, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				fmt.Printf("Error: Failed to generate JSON: %v\n", err)
				return
			}
			fmt.Println(string(jsonData))
			return
		}

		// Pretty print results
		if result.Success {
			successStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("2"))
			fmt.Println(successStyle.Render("✅ Remotes synced successfully"))
		} else {
			warningStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("3"))
			fmt.Println(warningStyle.Render("⚠️  Sync completed with errors"))
		}

		fmt.Println()
		fmt.Printf("📁 Synced remotes: %d\n", len(result.SyncedRemotes))
		for _, name := range result.SyncedRemotes {
			fmt.Printf("  ✓ %s\n", name)
		}

		if len(result.Errors) > 0 {
			fmt.Println()
			errorStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("1"))
			fmt.Println(errorStyle.Render("❌ Errors:"))
			for name, errMsg := range result.Errors {
				fmt.Printf("  ✗ %s: %s\n", name, errMsg)
			}
		}

		fmt.Printf("\n📄 Total files copied: %d\n", result.FilesCopied)
		fmt.Printf("📂 Remotes directory: .neev/remotes/\n")
	},
}

func init() {
	rootCmd.AddCommand(syncRemotesCmd)
	syncRemotesCmd.Flags().BoolVar(&syncJSON, "json", false, "Output results in JSON format")
}
