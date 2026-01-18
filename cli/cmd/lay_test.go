package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestLayCmd_IsRegistered(t *testing.T) {
	if layCmd == nil {
		t.Error("layCmd should not be nil")
	}

	if layCmd.Use != "lay <blueprint_name>" {
		t.Errorf("Expected Use='lay <blueprint_name>', got '%s'", layCmd.Use)
	}

	if layCmd.Short == "" {
		t.Error("layCmd should have a Short description")
	}

	if layCmd.Run == nil {
		t.Error("layCmd should have a Run function")
	}

	if layCmd.Args == nil {
		t.Error("layCmd should have Args validation")
	}
}

func TestLayCmd_ArgumentValidation(t *testing.T) {
	// Test with no arguments
	cmd := &cobra.Command{}
	err := layCmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Should fail with no arguments")
	}

	// Test with correct number of arguments
	err = layCmd.Args(cmd, []string{"blueprint-name"})
	if err != nil {
		t.Errorf("Should accept single argument: %v", err)
	}

	// Test with too many arguments
	err = layCmd.Args(cmd, []string{"blueprint", "extra"})
	if err == nil {
		t.Error("Should fail with multiple arguments")
	}
}

func TestLayCmd_Properties(t *testing.T) {
	if layCmd == nil {
		t.Fatal("layCmd is nil")
	}

	if layCmd.Use != "lay <blueprint_name>" {
		t.Errorf("Expected Use='lay <blueprint_name>', got '%s'", layCmd.Use)
	}

	if len(layCmd.Short) == 0 {
		t.Error("Short description should not be empty")
	}

	if layCmd.Run == nil {
		t.Error("layCmd should have Run function")
	}
}

func TestLayCmd_HasRunFunction(t *testing.T) {
	if layCmd == nil {
		t.Fatal("layCmd is nil")
	}

	if layCmd.Run == nil && layCmd.RunE == nil {
		t.Error("layCmd should have Run or RunE function")
	}
}
