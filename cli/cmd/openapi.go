package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/neev-kit/neev/core/openapi"
	"github.com/spf13/cobra"
)

var openapiCmd = &cobra.Command{
	Use:   "openapi <blueprint>",
	Short: "Generate OpenAPI specification from a blueprint",
	Long:  "Parse architecture.md from a blueprint and generate an OpenAPI 3.1 specification file (openapi.yaml)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		blueprintName := args[0]
		
		// Construct paths
		blueprintPath := filepath.Join(".neev", "blueprints", blueprintName)
		architecturePath := filepath.Join(blueprintPath, "architecture.md")
		outputPath := filepath.Join(blueprintPath, "openapi.yaml")
		
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
		
		// Generate OpenAPI spec
		yamlData, err := openapi.GenerateOpenAPI(architecturePath, blueprintName)
		if err != nil {
			fmt.Printf("❌ Failed to generate OpenAPI spec: %v\n", err)
			os.Exit(1)
		}
		
		// Write to file
		if err := os.WriteFile(outputPath, yamlData, 0644); err != nil {
			fmt.Printf("❌ Failed to write openapi.yaml: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("✅ Generated OpenAPI specification: %s\n", outputPath)
	},
}

func init() {
	rootCmd.AddCommand(openapiCmd)
}
