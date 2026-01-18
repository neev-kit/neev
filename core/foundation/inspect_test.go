package foundation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/neev-kit/neev/core/config"
)

func TestInspect_NoFoundation(t *testing.T) {
	tmpDir := t.TempDir()

	warnings, err := Inspect(tmpDir)
	if err != nil {
		t.Errorf("Inspect should not fail with no foundation: %v", err)
	}

	if len(warnings) > 0 {
		t.Errorf("Expected no warnings with empty directory, got %d", len(warnings))
	}
}

func TestInspect_FoundationDriftMissing(t *testing.T) {
	tmpDir := t.TempDir()

	// Create foundation spec
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("Failed to create foundation: %v", err)
	}

	// Create auth.md spec
	if err := os.WriteFile(filepath.Join(foundationPath, "auth.md"), []byte("# Auth"), 0644); err != nil {
		t.Fatalf("Failed to write spec: %v", err)
	}

	// Run inspect - should warn about missing auth/ directory
	warnings, err := Inspect(tmpDir)
	if err != nil {
		t.Errorf("Inspect failed: %v", err)
	}

	if len(warnings) == 0 {
		t.Error("Expected warnings about missing auth directory")
	}

	found := false
	for _, w := range warnings {
		if contains(w, "auth") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected warning about auth directory")
	}
}

func TestInspect_CodeDriftMissing(t *testing.T) {
	tmpDir := t.TempDir()

	// Create code directory without spec
	codePath := filepath.Join(tmpDir, "services")
	if err := os.MkdirAll(codePath, 0755); err != nil {
		t.Fatalf("Failed to create code dir: %v", err)
	}

	// Run inspect - should warn about missing spec
	warnings, err := Inspect(tmpDir)
	if err != nil {
		t.Errorf("Inspect failed: %v", err)
	}

	if len(warnings) == 0 {
		t.Error("Expected warnings about missing spec")
	}

	found := false
	for _, w := range warnings {
		if contains(w, "services") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected warning about services directory")
	}
}

func TestInspect_Balanced(t *testing.T) {
	tmpDir := t.TempDir()

	// Create matching spec and code
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("Failed to create foundation: %v", err)
	}

	if err := os.WriteFile(filepath.Join(foundationPath, "auth.md"), []byte("# Auth"), 0644); err != nil {
		t.Fatalf("Failed to write spec: %v", err)
	}

	authDir := filepath.Join(tmpDir, "auth")
	if err := os.MkdirAll(authDir, 0755); err != nil {
		t.Fatalf("Failed to create auth dir: %v", err)
	}

	// Run inspect - should have no warnings
	warnings, err := Inspect(tmpDir)
	if err != nil {
		t.Errorf("Inspect failed: %v", err)
	}

	if len(warnings) > 0 {
		t.Errorf("Expected no warnings with balanced structure, got %d: %v", len(warnings), warnings)
	}
}

func TestInspect_IgnoresCommonDirs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create common ignored directories
	ignoredPaths := []string{"node_modules", "dist", "vendor", ".git", "build"}
	for _, dir := range ignoredPaths {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatalf("Failed to create %s: %v", dir, err)
		}
	}

	// Run inspect - should not complain about ignored directories
	warnings, err := Inspect(tmpDir)
	if err != nil {
		t.Errorf("Inspect failed: %v", err)
	}

	for _, w := range warnings {
		for _, ignored := range ignoredPaths {
			if contains(w, ignored) {
				t.Errorf("Inspect should ignore '%s' directory, but got warning: %s", ignored, w)
			}
		}
	}
}

func TestInspect_WithSrcDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create src/auth directory
	authPath := filepath.Join(tmpDir, "src", "auth")
	if err := os.MkdirAll(authPath, 0755); err != nil {
		t.Fatalf("Failed to create src/auth: %v", err)
	}

	// Create matching spec
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("Failed to create foundation: %v", err)
	}

	if err := os.WriteFile(filepath.Join(foundationPath, "auth.md"), []byte("# Auth"), 0644); err != nil {
		t.Fatalf("Failed to write spec: %v", err)
	}

	// Run inspect
	warnings, err := Inspect(tmpDir)
	if err != nil {
		t.Errorf("Inspect failed: %v", err)
	}

	if len(warnings) > 0 {
		t.Errorf("Expected no warnings with matching src/auth, got: %v", warnings)
	}
}

func TestInspect_MultipleModules(t *testing.T) {
	tmpDir := t.TempDir()

	// Create foundation with multiple specs
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("Failed to create foundation: %v", err)
	}

	specs := []string{"auth", "api", "database"}
	for _, spec := range specs {
		if err := os.WriteFile(filepath.Join(foundationPath, spec+".md"), []byte("# "+spec), 0644); err != nil {
			t.Fatalf("Failed to write spec: %v", err)
		}
	}

	// Create matching directories
	for _, spec := range specs {
		if err := os.MkdirAll(filepath.Join(tmpDir, spec), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
	}

	// Run inspect
	warnings, err := Inspect(tmpDir)
	if err != nil {
		t.Errorf("Inspect failed: %v", err)
	}

	if len(warnings) > 0 {
		t.Errorf("Expected no warnings with all matching specs, got: %v", warnings)
	}
}

// TestInspect_FindsMissingCodeDirectories tests that Inspect detects when
// foundation specs exist but corresponding code directories don't
func TestInspect_FindsMissingCodeDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Set up foundation directory with specs
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("failed to create foundation dir: %v", err)
	}

	// Create foundation specs (modules)
	specs := []string{"auth", "api", "database"}
	for _, spec := range specs {
		specFile := filepath.Join(foundationPath, spec+".md")
		if err := os.WriteFile(specFile, []byte("# "+spec+" Module\n"), 0644); err != nil {
			t.Fatalf("failed to create spec %s: %v", spec, err)
		}
	}

	// Create only some code directories
	srcPath := filepath.Join(tmpDir, "src")
	if err := os.MkdirAll(srcPath, 0755); err != nil {
		t.Fatalf("failed to create src dir: %v", err)
	}

	// Create only "auth" and "api" directories, missing "database"
	for _, dir := range []string{"auth", "api"} {
		dirPath := filepath.Join(srcPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("failed to create code dir %s: %v", dir, err)
		}
	}

	// Run inspection
	warnings, err := Inspect(tmpDir)
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}

	// Verify warnings contain the missing directory
	if len(warnings) == 0 {
		t.Error("expected at least one warning, got none")
	}

	foundWarning := false
	for _, warning := range warnings {
		if strings.Contains(warning, "database") && strings.Contains(warning, "not found") {
			foundWarning = true
			break
		}
	}

	if !foundWarning {
		t.Errorf("expected warning about missing 'database' directory, got: %v", warnings)
	}
}

// TestInspect_WithCustomConfig tests that Inspect respects custom config ignore dirs
func TestInspect_WithCustomConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Set up foundation directory
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("failed to create foundation dir: %v", err)
	}

	// Create specs
	for _, spec := range []string{"app", "configs"} {
		specFile := filepath.Join(foundationPath, spec+".md")
		if err := os.WriteFile(specFile, []byte("# "+spec+" Module\n"), 0644); err != nil {
			t.Fatalf("failed to create spec: %v", err)
		}
	}

	// Create code directories
	srcPath := filepath.Join(tmpDir, "src")
	if err := os.MkdirAll(srcPath, 0755); err != nil {
		t.Fatalf("failed to create src dir: %v", err)
	}

	for _, dir := range []string{"app", "configs", "custom_ignore_me"} {
		dirPath := filepath.Join(srcPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("failed to create code dir: %v", err)
		}
	}

	// Create custom config that ignores "custom_ignore_me"
	cfg := &config.Config{
		ProjectName:    "Test Project",
		FoundationPath: ".neev",
		IgnoreDirs: []string{
			"custom_ignore_me",
			"node_modules",
			"dist",
		},
	}

	// Run inspection with custom config
	warnings, err := InspectWithConfig(tmpDir, cfg)
	if err != nil {
		t.Fatalf("InspectWithConfig failed: %v", err)
	}

	// Verify that custom_ignore_me doesn't appear in warnings
	for _, warning := range warnings {
		if strings.Contains(warning, "custom_ignore_me") {
			t.Errorf("custom ignore dir appeared in warnings: %s", warning)
		}
	}
}

// TestInspect_HiddenDirsIgnored tests that hidden directories are ignored
func TestInspect_HiddenDirsIgnored(t *testing.T) {
	tmpDir := t.TempDir()

	// Set up foundation directory
	foundationPath := filepath.Join(tmpDir, ".neev", "foundation")
	if err := os.MkdirAll(foundationPath, 0755); err != nil {
		t.Fatalf("failed to create foundation dir: %v", err)
	}

	// Create one spec
	specFile := filepath.Join(foundationPath, "app.md")
	if err := os.WriteFile(specFile, []byte("# App\n"), 0644); err != nil {
		t.Fatalf("failed to create spec: %v", err)
	}

	// Create code directories including hidden ones
	srcPath := filepath.Join(tmpDir, "src")
	if err := os.MkdirAll(srcPath, 0755); err != nil {
		t.Fatalf("failed to create src dir: %v", err)
	}

	for _, dir := range []string{"app", ".vscode", ".idea"} {
		dirPath := filepath.Join(srcPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("failed to create code dir: %v", err)
		}
	}

	// Run inspection
	warnings, err := Inspect(tmpDir)
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}

	// Verify hidden directories don't appear in warnings
	for _, warning := range warnings {
		if strings.Contains(warning, ".vscode") || strings.Contains(warning, ".idea") {
			t.Errorf("hidden directory appeared in warnings: %s", warning)
		}
	}
}
