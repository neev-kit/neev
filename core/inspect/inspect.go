package inspect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// InspectOptions configures the inspection behavior
type InspectOptions struct {
	RootDir        string
	FoundationPath string
	IgnoreDirs     map[string]bool
	UseDescriptors bool // If true, use .module.yaml files for detailed inspection
	Depth          int  // Analysis depth: 1=structure, 2=+API, 3=+signatures
	CheckAPI       bool // Enable OpenAPI validation (Level 2)
	CheckSignatures bool // Enable signature validation (Level 3)
}

// Inspect performs drift detection between foundation specs and code structure
func Inspect(opts InspectOptions) (*InspectResult, error) {
	result := &InspectResult{
		Success:  true,
		Warnings: []Warning{},
		Summary:  Summary{
			Languages: make(map[string]int),
		},
	}

	// Initialize polyglot analyzer
	analyzer := NewPolyglotAnalyzer()
	
	// Import detector types properly
	goDetector := &GoDetector{}
	pyDetector := &PythonDetector{}
	jsDetector := &JavaScriptDetector{}
	javaDetector := &JavaDetector{}
	csDetector := &CSharpDetector{}
	rbDetector := &RubyDetector{}
	
	analyzer.RegisterDetector(goDetector)
	analyzer.RegisterDetector(pyDetector)
	analyzer.RegisterDetector(jsDetector)
	analyzer.RegisterDetector(javaDetector)
	analyzer.RegisterDetector(csDetector)
	analyzer.RegisterDetector(rbDetector)

	// Detect languages in the codebase
	languages, err := analyzer.DetectLanguages(opts.RootDir, opts.IgnoreDirs)
	if err != nil {
		return nil, fmt.Errorf("failed to detect languages: %w", err)
	}
	result.Summary.Languages = languages

	// Get foundation modules
	foundationModules, descriptors, err := getFoundationModules(opts.FoundationPath, opts.UseDescriptors)
	if err != nil {
		return nil, fmt.Errorf("failed to read foundation modules: %w", err)
	}

	// Get code modules
	codeModules, err := getCodeModules(opts.RootDir, opts.IgnoreDirs)
	if err != nil {
		return nil, fmt.Errorf("failed to scan code modules: %w", err)
	}

	result.Summary.TotalModules = len(foundationModules)

	// Check for missing code directories
	for module := range foundationModules {
		if _, exists := codeModules[module]; !exists {
			warning := Warning{
				Type:        WarningMissingModule,
				Module:      module,
				Message:     fmt.Sprintf("Foundation spec '%s.md' exists but directory '%s/' not found in code", module, module),
				Severity:    "warning",
				Remediation: fmt.Sprintf("Create directory '%s/' or remove the foundation spec", module),
			}
			result.Warnings = append(result.Warnings, warning)
			result.Summary.MissingModules++
			result.Summary.WarningCount++
		} else {
			result.Summary.MatchingModules++

			// If descriptors are enabled, check file-level details
			if opts.UseDescriptors {
				if descriptor, hasDescriptor := descriptors[module]; hasDescriptor {
					fileWarnings := checkModuleFiles(opts.RootDir, module, descriptor, codeModules[module])
					result.Warnings = append(result.Warnings, fileWarnings...)
					for _, w := range fileWarnings {
						if w.Severity == "error" {
							result.Summary.ErrorCount++
						} else {
							result.Summary.WarningCount++
						}
					}
				}
			}
		}
	}

	// Check for orphaned code directories (code without specs)
	for module := range codeModules {
		if _, exists := foundationModules[module]; !exists {
			warning := Warning{
				Type:        WarningExtraCode,
				Module:      module,
				Message:     fmt.Sprintf("Code directory '%s/' exists but no foundation spec '%s.md' found", module, module),
				Severity:    "info",
				Remediation: fmt.Sprintf("Create foundation spec '%s.md' or remove the directory", module),
			}
			result.Warnings = append(result.Warnings, warning)
			result.Summary.ExtraCodeDirs++
			result.Summary.WarningCount++
		}
	}

	// Level 2: OpenAPI validation (if enabled)
	if opts.CheckAPI || opts.Depth >= 2 {
		apiWarnings, err := ValidateOpenAPIContracts(opts, analyzer)
		if err != nil {
			return nil, fmt.Errorf("failed to validate API contracts: %w", err)
		}
		
		for _, w := range apiWarnings {
			result.Warnings = append(result.Warnings, w)
			if w.Type == WarningMissingEndpoint {
				result.Summary.MissingEndpoints++
			} else if w.Type == WarningUndocumentedEndpoint {
				result.Summary.UndocumentedEnds++
			}
			
			if w.Severity == "error" {
				result.Summary.ErrorCount++
			} else {
				result.Summary.WarningCount++
			}
		}
	}

	// Level 3: Function signature validation (if enabled)
	if opts.CheckSignatures || opts.Depth >= 3 {
		sigWarnings, err := ValidateFunctionSignatures(opts, analyzer)
		if err != nil {
			return nil, fmt.Errorf("failed to validate function signatures: %w", err)
		}
		
		for _, w := range sigWarnings {
			result.Warnings = append(result.Warnings, w)
			if w.Type == WarningSignatureMismatch {
				result.Summary.SignatureMismatches++
			}
			
			if w.Severity == "error" {
				result.Summary.ErrorCount++
			} else {
				result.Summary.WarningCount++
			}
		}
	}

	result.Summary.TotalWarnings = len(result.Warnings)
	result.Success = result.Summary.ErrorCount == 0

	return result, nil
}

// getFoundationModules returns module names from foundation specs and their descriptors if available
func getFoundationModules(foundationPath string, useDescriptors bool) (map[string]bool, map[string]ModuleDescriptor, error) {
	modules := make(map[string]bool)
	descriptors := make(map[string]ModuleDescriptor)

	// Check if foundation directory exists
	if _, err := os.Stat(foundationPath); os.IsNotExist(err) {
		return modules, descriptors, nil // No foundation yet, that's okay
	}

	entries, err := os.ReadDir(foundationPath)
	if err != nil {
		return nil, nil, err
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

			// Try to load descriptor if enabled
			if useDescriptors {
				descriptorPath := filepath.Join(foundationPath, moduleName+".module.yaml")
				if descriptor, err := loadModuleDescriptor(descriptorPath); err == nil {
					descriptors[moduleName] = descriptor
				}
			}
		}
	}

	return modules, descriptors, nil
}

// loadModuleDescriptor loads a module descriptor from a YAML file
func loadModuleDescriptor(path string) (ModuleDescriptor, error) {
	var descriptor ModuleDescriptor

	data, err := os.ReadFile(path)
	if err != nil {
		return descriptor, err
	}

	err = yaml.Unmarshal(data, &descriptor)
	if err != nil {
		return descriptor, fmt.Errorf("failed to parse module descriptor: %w", err)
	}

	return descriptor, nil
}

// getCodeModules returns a set of module directory names from the code structure
func getCodeModules(rootDir string, ignoreDirs map[string]bool) (map[string]string, error) {
	modules := make(map[string]string)

	// Try src/ first, fall back to root
	srcPath := filepath.Join(rootDir, "src")
	scanPath := srcPath
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		scanPath = rootDir
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
		if ignoreDirs[name] {
			continue
		}

		// Skip hidden directories
		if strings.HasPrefix(name, ".") {
			continue
		}

		modules[name] = filepath.Join(scanPath, name)
	}

	return modules, nil
}

// checkModuleFiles verifies that expected files from the descriptor exist
func checkModuleFiles(rootDir, moduleName string, descriptor ModuleDescriptor, modulePath string) []Warning {
	var warnings []Warning

	// Check expected files
	for _, expectedFile := range descriptor.ExpectedFiles {
		filePath := filepath.Join(modulePath, expectedFile)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			warning := Warning{
				Type:        WarningMissingFile,
				Module:      moduleName,
				Message:     fmt.Sprintf("Expected file '%s' not found in module '%s'", expectedFile, moduleName),
				Severity:    "warning",
				Remediation: fmt.Sprintf("Create file '%s' or update module descriptor", filePath),
			}
			warnings = append(warnings, warning)
		}
	}

	// Check expected directories
	for _, expectedDir := range descriptor.ExpectedDirs {
		dirPath := filepath.Join(modulePath, expectedDir)
		if stat, err := os.Stat(dirPath); os.IsNotExist(err) || !stat.IsDir() {
			warning := Warning{
				Type:        WarningMissingFile,
				Module:      moduleName,
				Message:     fmt.Sprintf("Expected directory '%s' not found in module '%s'", expectedDir, moduleName),
				Severity:    "warning",
				Remediation: fmt.Sprintf("Create directory '%s' or update module descriptor", dirPath),
			}
			warnings = append(warnings, warning)
		}
	}

	// Check patterns (glob matching)
	for _, pattern := range descriptor.Patterns {
		fullPattern := filepath.Join(modulePath, pattern)
		matches, err := filepath.Glob(fullPattern)
		if err != nil || len(matches) == 0 {
			warning := Warning{
				Type:        WarningMissingFile,
				Module:      moduleName,
				Message:     fmt.Sprintf("No files matching pattern '%s' in module '%s'", pattern, moduleName),
				Severity:    "info",
				Remediation: fmt.Sprintf("Add files matching pattern '%s' or update module descriptor", pattern),
			}
			warnings = append(warnings, warning)
		}
	}

	return warnings
}
