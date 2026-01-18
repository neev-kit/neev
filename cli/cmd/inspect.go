package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/neev-kit/neev/core/config"
	"github.com/neev-kit/neev/core/foundation"
	"github.com/neev-kit/neev/core/inspect"
	"github.com/spf13/cobra"
)

var (
	jsonOutput     bool
	useDescriptors bool
)

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect the foundation for drift",
	Long:  "Check if the project structure matches the foundation specifications",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error: could not determine current working directory: %v\n", err)
			return
		}

		// Load config to get foundation path and ignore dirs
		cfg, err := config.LoadConfig(cwd)
		if err != nil {
			fmt.Printf("Warning: Could not load config, using defaults: %v\n", err)
			cfg = config.DefaultConfig()
		}

		// Use new structured inspect if descriptors are enabled or JSON output requested
		if useDescriptors || jsonOutput {
			foundationPath := filepath.Join(cwd, ".neev", "foundation")
			
			opts := inspect.InspectOptions{
				RootDir:        cwd,
				FoundationPath: foundationPath,
				IgnoreDirs:     cfg.GetIgnoreDirs(),
				UseDescriptors: useDescriptors,
			}

			result, err := inspect.Inspect(opts)
			if err != nil {
				errorStyle := lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("1")).
					Margin(0, 0, 1, 0)

				fmt.Println(errorStyle.Render("❌ Inspection failed: " + err.Error()))
				return
			}

			// JSON output mode
			if jsonOutput {
				jsonData, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					fmt.Printf("Error: Failed to generate JSON: %v\n", err)
					return
				}
				fmt.Println(string(jsonData))
				return
			}

			// Pretty print structured output
			printStructuredResult(result)
			return
		}

		// Fall back to legacy inspect for backwards compatibility
		warnings, err := foundation.Inspect(cwd)
		if err != nil {
			errorStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("1")).
				Margin(0, 0, 1, 0)

			fmt.Println(errorStyle.Render("❌ Inspection failed: " + err.Error()))
			return
		}

		if len(warnings) == 0 {
			successStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("2")).
				Margin(0, 0, 1, 0)

			fmt.Println(successStyle.Render("✅ Foundation is solid."))
			return
		}

		warningStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("3")).
			Margin(0, 0, 1, 0)

		fmt.Println(warningStyle.Render("⚠️  Foundation drift detected:"))
		fmt.Println()

		for _, warning := range warnings {
			fmt.Println("  " + warning)
		}
	},
}

func printStructuredResult(result *inspect.InspectResult) {
	if result.Success && len(result.Warnings) == 0 {
		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")).
			Margin(0, 0, 1, 0)

		fmt.Println(successStyle.Render("✅ Foundation is solid."))
		fmt.Println()
		fmt.Printf("📊 Summary: %d modules checked, all in sync\n", result.Summary.TotalModules)
		return
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("3")).
		Margin(0, 0, 1, 0)

	fmt.Println(titleStyle.Render("⚠️  Foundation drift detected:"))
	fmt.Println()

	// Group warnings by severity
	errors := []inspect.Warning{}
	warnings := []inspect.Warning{}
	infos := []inspect.Warning{}

	for _, w := range result.Warnings {
		switch w.Severity {
		case "error":
			errors = append(errors, w)
		case "warning":
			warnings = append(warnings, w)
		case "info":
			infos = append(infos, w)
		}
	}

	// Print errors first
	if len(errors) > 0 {
		errorStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("1"))
		fmt.Println(errorStyle.Render("🔴 Errors:"))
		for _, w := range errors {
			fmt.Printf("  [%s] %s: %s\n", w.Type, w.Module, w.Message)
			if w.Remediation != "" {
				fmt.Printf("    💡 %s\n", w.Remediation)
			}
		}
		fmt.Println()
	}

	// Print warnings
	if len(warnings) > 0 {
		warningStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("3"))
		fmt.Println(warningStyle.Render("🟡 Warnings:"))
		for _, w := range warnings {
			fmt.Printf("  [%s] %s: %s\n", w.Type, w.Module, w.Message)
			if w.Remediation != "" {
				fmt.Printf("    💡 %s\n", w.Remediation)
			}
		}
		fmt.Println()
	}

	// Print infos
	if len(infos) > 0 {
		infoStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("6"))
		fmt.Println(infoStyle.Render("ℹ️  Info:"))
		for _, w := range infos {
			fmt.Printf("  [%s] %s: %s\n", w.Type, w.Module, w.Message)
		}
		fmt.Println()
	}

	// Print summary
	fmt.Println("📊 Summary:")
	fmt.Printf("  Total modules: %d\n", result.Summary.TotalModules)
	fmt.Printf("  Matching: %d\n", result.Summary.MatchingModules)
	fmt.Printf("  Missing: %d\n", result.Summary.MissingModules)
	fmt.Printf("  Extra code dirs: %d\n", result.Summary.ExtraCodeDirs)
	fmt.Printf("  Total warnings: %d (errors: %d, warnings: %d)\n",
		result.Summary.TotalWarnings, result.Summary.ErrorCount, result.Summary.WarningCount)
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
	inspectCmd.Flags().BoolVar(&useDescriptors, "use-descriptors", false, "Use .module.yaml files for detailed inspection")
}
