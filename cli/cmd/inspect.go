package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/neev-kit/neev/core/config"
	"github.com/neev-kit/neev/core/foundation"
	"github.com/neev-kit/neev/core/inspect"
	"github.com/spf13/cobra"
)

var (
	jsonOutput      bool
	useDescriptors  bool
	strictMode      bool
	depth           int
	checkAPI        bool
	checkSignatures bool
	checkTests      bool
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
		if useDescriptors || jsonOutput || depth > 1 || checkAPI || checkSignatures {
			foundationPath := filepath.Join(cwd, ".neev", "foundation")

			opts := inspect.InspectOptions{
				RootDir:         cwd,
				FoundationPath:  foundationPath,
				IgnoreDirs:      cfg.GetIgnoreDirs(),
				UseDescriptors:  useDescriptors,
				Depth:           depth,
				CheckAPI:        checkAPI,
				CheckSignatures: checkSignatures,
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

				// Exit with error code if strict mode and drift found
				if strictMode && (!result.Success || len(result.Warnings) > 0) {
					os.Exit(1)
				}
				return
			}

			// Pretty print structured output
			printStructuredResult(result)

			// Exit with error code if strict mode and drift found
			if strictMode && (!result.Success || len(result.Warnings) > 0) {
				os.Exit(1)
			}
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

		// Exit with error code if strict mode and drift found
		if strictMode {
			os.Exit(1)
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
		
		// Print language breakdown if available
		if len(result.Summary.Languages) > 0 {
			fmt.Print("🔍 LANGUAGES DETECTED: ")
			var langs []string
			for lang, count := range result.Summary.Languages {
				langs = append(langs, fmt.Sprintf("%s (%d files)", lang, count))
			}
			fmt.Println(strings.Join(langs, " | "))
			fmt.Println()
		}
		
		fmt.Printf("📊 Summary: %d modules checked, all in sync\n", result.Summary.TotalModules)
		return
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("3")).
		Margin(0, 0, 1, 0)

	fmt.Println(titleStyle.Render("⚠️  Foundation drift detected:"))
	fmt.Println()
	
	// Print language breakdown if available
	if len(result.Summary.Languages) > 0 {
		fmt.Println("═════ LANGUAGES DETECTED =====")
		var langs []string
		for lang, count := range result.Summary.Languages {
			langs = append(langs, fmt.Sprintf("%s (%d files)", lang, count))
		}
		fmt.Println(strings.Join(langs, " | "))
		fmt.Println()
	}

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
	
	// Print API contract summary if applicable
	if result.Summary.MissingEndpoints > 0 || result.Summary.UndocumentedEnds > 0 {
		fmt.Printf("  Missing endpoints: %d\n", result.Summary.MissingEndpoints)
		fmt.Printf("  Undocumented endpoints: %d\n", result.Summary.UndocumentedEnds)
	}
	
	// Print signature mismatch summary if applicable
	if result.Summary.SignatureMismatches > 0 {
		fmt.Printf("  Signature mismatches: %d\n", result.Summary.SignatureMismatches)
	}
	
	fmt.Printf("  Total warnings: %d (errors: %d, warnings: %d)\n",
		result.Summary.TotalWarnings, result.Summary.ErrorCount, result.Summary.WarningCount)
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
	inspectCmd.Flags().BoolVar(&useDescriptors, "use-descriptors", false, "Use .module.yaml files for detailed inspection")
	inspectCmd.Flags().BoolVar(&strictMode, "strict", false, "Exit with code 1 if any drift is detected (for CI pipelines)")
	inspectCmd.Flags().IntVar(&depth, "depth", 1, "Depth of analysis (1=structure, 2=+API, 3=+signatures)")
	inspectCmd.Flags().BoolVar(&checkAPI, "check-api", false, "Validate OpenAPI specs (enables Level 2)")
	inspectCmd.Flags().BoolVar(&checkSignatures, "check-signatures", false, "Validate function signatures (enables Level 3)")
	inspectCmd.Flags().BoolVar(&checkTests, "check-tests", false, "Validate BDD test coverage (not yet implemented)")
}
