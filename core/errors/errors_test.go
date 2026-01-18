package errors

import "testing"

func TestErrBlueprintNotFound(t *testing.T) {
	err := ErrBlueprintNotFound("test")
	if err == nil {
		t.Fatal("ErrBlueprintNotFound returned nil")
	}
	if err.Type != ErrTypeBlueprintNotFound {
		t.Errorf("Wrong error type: %s", err.Type)
	}
}

func TestErrFoundationMissing(t *testing.T) {
	err := ErrFoundationMissing()
	if err == nil {
		t.Fatal("ErrFoundationMissing returned nil")
	}
}

func TestNeevError_Error(t *testing.T) {
	err := ErrInvalidConfig("test")
	if err.Error() == "" {
		t.Error("Error() returned empty string")
	}
}

func TestGetSolutionHint(t *testing.T) {
	err := ErrBlueprintNotFound("myblueprint")
	hint := err.GetSolutionHint()
	if hint == "" {
		t.Error("GetSolutionHint returned empty string")
	}
}
