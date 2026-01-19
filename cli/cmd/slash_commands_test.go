package cmd

import (
	"testing"
)

func TestSlashCommandsCmd_Properties(t *testing.T) {
	if slashCommandsCmd.Use == "" {
		t.Error("Expected non-empty Use")
	}

	if slashCommandsCmd.Short == "" {
		t.Error("Expected non-empty Short")
	}

	if slashCommandsCmd.Long == "" {
		t.Error("Expected non-empty Long")
	}
}

func TestSlashCommandsCmd_IsRegisteredWithRoot(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "slash-commands" {
			found = true
			break
		}
	}

	if !found {
		t.Error("slash-commands command not registered with root")
	}
}

func TestSlashCommandsCmd_HasListFlag(t *testing.T) {
	if flag := slashCommandsCmd.Flag("list"); flag == nil {
		t.Error("Expected --list flag")
	}
}

func TestSlashCommandsCmd_HasUpdateFlag(t *testing.T) {
	if flag := slashCommandsCmd.Flag("update"); flag == nil {
		t.Error("Expected --update flag")
	}
}

func TestSlashCommandsCmd_HasRegisterFlag(t *testing.T) {
	if flag := slashCommandsCmd.Flag("register"); flag == nil {
		t.Error("Expected --register flag")
	}
}

func TestSlashCommandsCmd_HasToolFlag(t *testing.T) {
	if flag := slashCommandsCmd.Flag("tool"); flag == nil {
		t.Error("Expected --tool flag")
	}
}

func TestListSlashCommands(t *testing.T) {
	// This is an output function, just verify it doesn't panic
	listSlashCommands()
}

func TestShowToolCommands(t *testing.T) {
	// Verify it handles tool names
	showToolCommands("claude-code")
	showToolCommands("cursor")
}
