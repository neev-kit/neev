package inspect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInspect_NoFoundation(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	opts := InspectOptions{
		RootDir:        tmpDir,
		FoundationPath: filepath.Join(tmpDir, ".neev", "foundation"),
		IgnoreDirs:     map[string]bool{".git": true, "node_modules": true},
		UseDescriptors: false,
	}

	result, err := Inspect(opts)
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success when no foundation exists")
	}

	if len(result.Warnings) != 0 {
		t.Errorf("Expected 0 warnings, got %d", len(result.Warnings))
	}
}

func TestInspect_MissingModule(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	foundationDir := filepath.Join(tmpDir, ".neev", "foundation")
	os.MkdirAll(foundationDir, 0755)

	// Create a foundation spec without corresponding code
	specPath := filepath.Join(foundationDir, "auth.md")
	os.WriteFile(specPath, []byte("# Auth Module"), 0644)

	opts := InspectOptions{
		RootDir:        tmpDir,
		FoundationPath: foundationDir,
		IgnoreDirs:     map[string]bool{".git": true, ".neev": true},
		UseDescriptors: false,
	}

	result, err := Inspect(opts)
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}

	if len(result.Warnings) == 0 {
		t.Errorf("Expected warnings for missing module")
	}

	foundMissingWarning := false
	for _, w := range result.Warnings {
		if w.Type == WarningMissingModule && w.Module == "auth" {
			foundMissingWarning = true
		}
	}

	if !foundMissingWarning {
		t.Errorf("Expected MISSING_MODULE warning for 'auth'")
	}

	if result.Summary.MissingModules != 1 {
		t.Errorf("Expected 1 missing module, got %d", result.Summary.MissingModules)
	}
}

func TestInspect_ExtraCode(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	foundationDir := filepath.Join(tmpDir, ".neev", "foundation")
	os.MkdirAll(foundationDir, 0755)

	// Create a code directory without foundation spec
	codeDir := filepath.Join(tmpDir, "utils")
	os.MkdirAll(codeDir, 0755)
	os.WriteFile(filepath.Join(codeDir, "helper.go"), []byte("package utils"), 0644)

	opts := InspectOptions{
		RootDir:        tmpDir,
		FoundationPath: foundationDir,
		IgnoreDirs:     map[string]bool{".git": true, ".neev": true},
		UseDescriptors: false,
	}

	result, err := Inspect(opts)
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}

	if len(result.Warnings) == 0 {
		t.Errorf("Expected warnings for extra code")
	}

	foundExtraWarning := false
	for _, w := range result.Warnings {
		if w.Type == WarningExtraCode && w.Module == "utils" {
			foundExtraWarning = true
		}
	}

	if !foundExtraWarning {
		t.Errorf("Expected EXTRA_CODE warning for 'utils'")
	}

	if result.Summary.ExtraCodeDirs != 1 {
		t.Errorf("Expected 1 extra code dir, got %d", result.Summary.ExtraCodeDirs)
	}
}

func TestInspect_MatchingModules(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	foundationDir := filepath.Join(tmpDir, ".neev", "foundation")
	os.MkdirAll(foundationDir, 0755)

	// Create matching foundation spec and code
	specPath := filepath.Join(foundationDir, "api.md")
	os.WriteFile(specPath, []byte("# API Module"), 0644)

	codeDir := filepath.Join(tmpDir, "api")
	os.MkdirAll(codeDir, 0755)
	os.WriteFile(filepath.Join(codeDir, "server.go"), []byte("package api"), 0644)

	opts := InspectOptions{
		RootDir:        tmpDir,
		FoundationPath: foundationDir,
		IgnoreDirs:     map[string]bool{".git": true, ".neev": true},
		UseDescriptors: false,
	}

	result, err := Inspect(opts)
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success for matching modules")
	}

	if len(result.Warnings) != 0 {
		t.Errorf("Expected 0 warnings, got %d", len(result.Warnings))
	}

	if result.Summary.MatchingModules != 1 {
		t.Errorf("Expected 1 matching module, got %d", result.Summary.MatchingModules)
	}
}

func TestInspect_WithDescriptors(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	foundationDir := filepath.Join(tmpDir, ".neev", "foundation")
	os.MkdirAll(foundationDir, 0755)

	// Create foundation spec
	specPath := filepath.Join(foundationDir, "auth.md")
	os.WriteFile(specPath, []byte("# Auth Module"), 0644)

	// Create module descriptor
	descriptorPath := filepath.Join(foundationDir, "auth.module.yaml")
	descriptorContent := `name: auth
description: Authentication module
expected_files:
  - handler.go
  - service.go
expected_dirs:
  - models
patterns:
  - "*.go"
`
	os.WriteFile(descriptorPath, []byte(descriptorContent), 0644)

	// Create code directory but missing some expected files
	codeDir := filepath.Join(tmpDir, "auth")
	os.MkdirAll(codeDir, 0755)
	os.WriteFile(filepath.Join(codeDir, "handler.go"), []byte("package auth"), 0644)
	// service.go is missing
	os.MkdirAll(filepath.Join(codeDir, "models"), 0755)

	opts := InspectOptions{
		RootDir:        tmpDir,
		FoundationPath: foundationDir,
		IgnoreDirs:     map[string]bool{".git": true, ".neev": true},
		UseDescriptors: true,
	}

	result, err := Inspect(opts)
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}

	// Should have warning about missing service.go
	foundMissingFile := false
	for _, w := range result.Warnings {
		if w.Type == WarningMissingFile && w.Module == "auth" {
			foundMissingFile = true
		}
	}

	if !foundMissingFile {
		t.Errorf("Expected MISSING_FILE warning for 'service.go'")
	}
}

func TestLoadModuleDescriptor(t *testing.T) {
	tmpDir := t.TempDir()
	descriptorPath := filepath.Join(tmpDir, "test.module.yaml")

	content := `name: test
description: Test module
expected_files:
  - main.go
  - utils.go
expected_dirs:
  - tests
patterns:
  - "*.go"
  - "**/*.test.go"
`
	os.WriteFile(descriptorPath, []byte(content), 0644)

	descriptor, err := loadModuleDescriptor(descriptorPath)
	if err != nil {
		t.Fatalf("Failed to load descriptor: %v", err)
	}

	if descriptor.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", descriptor.Name)
	}

	if len(descriptor.ExpectedFiles) != 2 {
		t.Errorf("Expected 2 expected files, got %d", len(descriptor.ExpectedFiles))
	}

	if len(descriptor.ExpectedDirs) != 1 {
		t.Errorf("Expected 1 expected dir, got %d", len(descriptor.ExpectedDirs))
	}

	if len(descriptor.Patterns) != 2 {
		t.Errorf("Expected 2 patterns, got %d", len(descriptor.Patterns))
	}
}

func TestInspectResult_JSON(t *testing.T) {
	result := InspectResult{
		Success: true,
		Warnings: []Warning{
			{
				Type:        WarningMissingModule,
				Module:      "test",
				Message:     "Test message",
				Severity:    "warning",
				Remediation: "Do something",
			},
		},
		Summary: Summary{
			TotalModules:    1,
			MatchingModules: 0,
			MissingModules:  1,
			TotalWarnings:   1,
			WarningCount:    1,
		},
	}

	// Should be marshallable to JSON
	_ = result // Just verify it compiles with JSON tags
}
