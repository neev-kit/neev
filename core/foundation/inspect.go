package foundation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/neev-kit/neev/core/config"
	neevErr "github.com/neev-kit/neev/core/errors"
)

// ignoredDirs contains default directories to skip during walk
var ignoredDirs = map[string]bool{
	".git":         true,
	".neev":        true,
	"node_modules": true,
	"dist":         true,
	"vendor":       true,
	"build":        true,
	".env":         true,
	"bin":          true,
	"obj":          true,
	".idea":        true,
	".vscode":      true,
	"target":       true,
}

// InspectWithConfig checks for drift using a configuration object
func InspectWithConfig(cwd string, cfg *config.Config) ([]string, error) {
	// Convert config ignore dirs to map
	ignoredDirsMap := make(map[string]bool)
	for _, dir := range cfg.IgnoreDirs {
		ignoredDirsMap[dir] = true
	}

	return inspectInternal(cwd, ignoredDirsMap)
}

// Inspect checks for drift between foundation specs and actual code structure
func Inspect(cwd string) ([]string, error) {
	return inspectInternal(cwd, ignoredDirs)
}

// inspectInternal performs the actual inspection with custom ignored directories
func inspectInternal(cwd string, ignored map[string]bool) ([]string, error) {
	var warnings []string

	// Get foundation modules
	foundationPath := filepath.Join(cwd, RootDir, FoundationDir)
	foundationModules, err := getFoundationModules(foundationPath)
	if err != nil {
		return nil, neevErr.NewNeevError(
			neevErr.ErrTypeFoundation,
			"failed to read foundation modules",
			err,
		)
	}

	// Get code modules (directories in src/ or root)
	codeModules, err := getCodeModulesWithIgnored(cwd, ignored)
	if err != nil {
		return nil, neevErr.NewNeevError(
			neevErr.ErrTypeIO,
			"failed to scan code modules",
			err,
		)
	}

	// Check for missing code directories
	for module := range foundationModules {
		if _, exists := codeModules[module]; !exists {
			warnings = append(warnings, fmt.Sprintf("⚠️  Foundation spec '%s.md' exists but directory '%s/' not found in code", module, module))
		}
	}

	// Check for orphaned code directories (code without specs)
	for module := range codeModules {
		if _, exists := foundationModules[module]; !exists {
			warnings = append(warnings, fmt.Sprintf("⚠️  Code directory '%s/' exists but no foundation spec '%s.md' found", module, module))
		}
	}

	return warnings, nil
}

// getFoundationModules returns a set of module names from foundation specs
func getFoundationModules(foundationPath string) (map[string]bool, error) {
	modules := make(map[string]bool)

	// Check if foundation directory exists
	if _, err := os.Stat(foundationPath); os.IsNotExist(err) {
		return modules, nil // No foundation yet, that's okay
	}

	entries, err := os.ReadDir(foundationPath)
	if err != nil {
		return nil, err
	}

	// Skip the archive directory
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == "archive" {
			continue
		}

		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			// Remove .md extension to get module name
			moduleName := strings.TrimSuffix(entry.Name(), ".md")
			modules[moduleName] = true
		}
	}

	return modules, nil
}

// getCodeModules returns a set of module directory names from the code structure
func getCodeModules(cwd string) (map[string]bool, error) {
	return getCodeModulesWithIgnored(cwd, ignoredDirs)
}

// getCodeModulesWithIgnored returns a set of module directory names, respecting ignored dirs
func getCodeModulesWithIgnored(cwd string, ignored map[string]bool) (map[string]bool, error) {
	modules := make(map[string]bool)

	// Try src/ first, fall back to root
	srcPath := filepath.Join(cwd, "src")
	scanPath := srcPath
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		scanPath = cwd
	}

	entries, err := os.ReadDir(scanPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Skip ignored directories
		if ignored[name] {
			continue
		}

		// Skip hidden directories
		if strings.HasPrefix(name, ".") {
			continue
		}

		modules[name] = true
	}

	return modules, nil
}
