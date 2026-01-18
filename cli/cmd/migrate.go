package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/neev-kit/neev/core/migration"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [flags]",
	Short: "Migrate existing project to Neev",
	Long:  "Convert existing projects (openspec, speckit) to Neev structure",
	Run: func(cmd *cobra.Command, args []string) {
		sourceType, _ := cmd.Flags().GetString("source")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		backup, _ := cmd.Flags().GetBool("backup")

		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error: could not determine current working directory: %v\n", err)
			return
		}

		// Map string to SourceType
		var source migration.SourceType
		switch sourceType {
		case "openspec":
			source = migration.SourceTypeOpenSpec
		case "speckit":
			source = migration.SourceTypeSpecKit
		case "auto":
			source = migration.SourceTypeAuto
		default:
			errorStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("1")).
				Margin(0, 0, 1, 0)

			fmt.Println(errorStyle.Render("❌ Invalid source type: " + sourceType))
			fmt.Println("   Valid options: openspec, speckit, auto")
			return
		}

		// Create migration config
		cfg := migration.MigrationConfig{
			RootDir:    cwd,
			SourceType: source,
			DryRun:     dryRun,
			BackupOld:  backup,
		}

		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("14")).
			Margin(0, 0, 1, 0)

		if dryRun {
			fmt.Println(headerStyle.Render("🔍 Dry-run migration from " + sourceType))
		} else {
			fmt.Println(headerStyle.Render("🚀 Migrating to Neev..."))
		}

		// Execute migration
		result, err := migration.Migrate(cfg)
		if err != nil {
			errorStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("1")).
				Margin(0, 0, 1, 0)

			fmt.Println(errorStyle.Render("❌ Migration failed: " + err.Error()))
			return
		}

		// Print messages
		infoStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("14"))

		for _, msg := range result.Messages {
			fmt.Println(infoStyle.Render("ℹ️  " + msg))
		}

		// Print errors if any
		if len(result.Errors) > 0 {
			warnStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("3"))

			for _, errMsg := range result.Errors {
				fmt.Println(warnStyle.Render("⚠️  " + errMsg))
			}
		}

		// Print summary
		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")).
			Margin(0, 0, 1, 0)

		summary := fmt.Sprintf("✅ Migration complete: %d files moved, %d directories created",
			result.FilesMovedCount, result.DirsCreatedCount)

		if dryRun {
			summary = "(dry-run) " + summary
		}

		fmt.Println(successStyle.Render(summary))

		if result.BackupDir != "" {
			backupMsg := fmt.Sprintf("💾 Backup created: %s", result.BackupDir)
			fmt.Println(infoStyle.Render(backupMsg))
		}
	},
}

func init() {
	migrateCmd.Flags().StringP("source", "s", "auto", "Source type: openspec, speckit, or auto")
	migrateCmd.Flags().Bool("dry-run", false, "Preview changes without applying them")
	migrateCmd.Flags().Bool("backup", false, "Create a backup of existing .neev directory")
	rootCmd.AddCommand(migrateCmd)
}
