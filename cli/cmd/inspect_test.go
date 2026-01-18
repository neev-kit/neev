package cmd

import (
	"testing"
)

func TestInspectCmd_IsRegistered(t *testing.T) {
	if inspectCmd == nil {
		t.Error("inspectCmd should not be nil")
	}

	if inspectCmd.Use != "inspect" {
		t.Errorf("Expected Use='inspect', got '%s'", inspectCmd.Use)
	}

	if inspectCmd.Short == "" {
		t.Error("inspectCmd should have a Short description")
	}

	if inspectCmd.Run == nil {
		t.Error("inspectCmd should have a Run function")
	}
}

func TestInspectCmd_Properties(t *testing.T) {
	if inspectCmd == nil {
		t.Fatal("inspectCmd is nil")
	}

	tests := map[string]string{
		"Use":   "inspect",
		"Short": "",
	}

	if inspectCmd.Use != tests["Use"] {
		t.Errorf("Expected Use='%s', got '%s'", tests["Use"], inspectCmd.Use)
	}

	if inspectCmd.Short != "" && len(inspectCmd.Short) == 0 {
		t.Error("Short description should not be empty")
	}
}

func TestInspectCmd_HasRunFunction(t *testing.T) {
	if inspectCmd == nil {
		t.Fatal("inspectCmd is nil")
	}

	if inspectCmd.Run == nil && inspectCmd.RunE == nil {
		t.Error("inspectCmd should have Run or RunE function")
	}
}
