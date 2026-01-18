package migration

import (
	"testing"
)

func TestMigrateBasic(t *testing.T) {
	// Test basic migration functionality
	cfg := MigrationConfig{
		RootDir:    "/tmp/test",
		SourceType: SourceTypeAuto,
		DryRun:     true,
	}

	result, err := Migrate(cfg)
	if err == nil && result != nil {
		// Expected behavior - at least it shouldn't panic
		return
	}
}
