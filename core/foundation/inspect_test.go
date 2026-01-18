package foundation

import (
	"os"
	"path/filepath"
	"testing"
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


