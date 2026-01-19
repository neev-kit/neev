package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/neev-kit/neev/core/cucumber"
	"github.com/spf13/cobra"
)

var (
	cucumberLang string
)

var cucumberCmd = &cobra.Command{
	Use:   "cucumber <blueprint>",
	Short: "Generate Cucumber/BDD test scaffolding from a blueprint",
	Long:  "Parse architecture.md from a blueprint and generate Cucumber feature files and step definition scaffolds for API testing",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		blueprintName := args[0]
		
		// Construct paths
		blueprintPath := filepath.Join(".neev", "blueprints", blueprintName)
		architecturePath := filepath.Join(blueprintPath, "architecture.md")
		testsPath := filepath.Join(blueprintPath, "tests")
		
		// Check if blueprint exists
		if _, err := os.Stat(blueprintPath); os.IsNotExist(err) {
			fmt.Printf("❌ Blueprint not found: %s\n", blueprintName)
			fmt.Printf("💡 Create it with: neev draft \"%s\"\n", blueprintName)
			os.Exit(1)
		}
		
		// Check if architecture.md exists
		if _, err := os.Stat(architecturePath); os.IsNotExist(err) {
			fmt.Printf("❌ architecture.md not found in blueprint: %s\n", blueprintName)
			os.Exit(1)
		}
		
		// Create tests directory
		if err := os.MkdirAll(testsPath, 0755); err != nil {
			fmt.Printf("❌ Failed to create tests directory: %v\n", err)
			os.Exit(1)
		}
		
		// Generate Cucumber tests
		if err := cucumber.GenerateCucumber(architecturePath, blueprintName, testsPath, cucumberLang); err != nil {
			fmt.Printf("❌ Failed to generate Cucumber tests: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("✅ Generated Cucumber tests in: %s\n", testsPath)
		fmt.Printf("   - Feature file: api.feature\n")
		if cucumberLang != "" {
			var ext string
			switch cucumberLang {
			case "go":
				ext = "go"
			case "javascript", "js":
				ext = "js"
			case "python":
				ext = "py"
			}
			fmt.Printf("   - Step definitions: steps.%s\n", ext)
		}
	},
}

func init() {
	cucumberCmd.Flags().StringVarP(&cucumberLang, "lang", "l", "", "Language for step definitions (go, javascript, python)")
	rootCmd.AddCommand(cucumberCmd)
}
