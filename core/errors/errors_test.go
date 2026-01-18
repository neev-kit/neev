package errors

import (
	"errors"
	"testing"
)

func TestErrBlueprintNotFound(t *testing.T) {
	err := ErrBlueprintNotFound("test")
	if err == nil {
		t.Fatal("ErrBlueprintNotFound returned nil")
	}
	if err.Type != ErrTypeBlueprintNotFound {
		t.Errorf("Wrong error type: %s", err.Type)
	}
	if err.Message != "blueprint 'test' not found" {
		t.Errorf("Wrong message: %s", err.Message)
	}
	if err.Err != nil {
		t.Error("Expected nil underlying error")
	}
}

func TestErrFoundationMissing(t *testing.T) {
	err := ErrFoundationMissing()
	if err == nil {
		t.Fatal("ErrFoundationMissing returned nil")
	}
	if err.Type != ErrTypeFoundation {
		t.Errorf("Wrong error type: %s", err.Type)
	}
}

func TestErrInvalidConfig(t *testing.T) {
	err := ErrInvalidConfig("bad format")
	if err == nil {
		t.Fatal("ErrInvalidConfig returned nil")
	}
	if err.Type != ErrTypeInvalidConfig {
		t.Errorf("Wrong error type: %s", err.Type)
	}
	if err.Message != "invalid configuration: bad format" {
		t.Errorf("Wrong message: %s", err.Message)
	}
}

func TestNeevError_Error(t *testing.T) {
	err := ErrInvalidConfig("test")
	if err.Error() == "" {
		t.Error("Error() returned empty string")
	}
	if err.Error() != "invalid configuration: test" {
		t.Errorf("Error() returned wrong message: %s", err.Error())
	}
}

func TestNeevError_ErrorWithUnderlying(t *testing.T) {
	underlying := errors.New("underlying error")
	err := NewNeevError(ErrTypeIO, "io failed", underlying)
	errorMsg := err.Error()

	if errorMsg != "io failed: underlying error" {
		t.Errorf("Expected 'io failed: underlying error', got '%s'", errorMsg)
	}
}

func TestGetSolutionHint_BlueprintNotFound(t *testing.T) {
	err := ErrBlueprintNotFound("myblueprint")
	hint := err.GetSolutionHint()
	if hint == "" {
		t.Error("GetSolutionHint returned empty string")
	}
	if !contains(hint, "blueprint") {
		t.Errorf("Expected 'blueprint' in hint: %s", hint)
	}
}

func TestGetSolutionHint_FoundationMissing(t *testing.T) {
	err := ErrFoundationMissing()
	hint := err.GetSolutionHint()
	if hint == "" {
		t.Error("GetSolutionHint returned empty string")
	}
	// Check for common words related to foundation
	if !contains(hint, "Foundation") && !contains(hint, "foundation") {
		t.Errorf("Expected 'Foundation' or 'foundation' in hint: %s", hint)
	}
}

func TestGetSolutionHint_InvalidConfig(t *testing.T) {
	err := ErrInvalidConfig("bad yaml")
	hint := err.GetSolutionHint()
	if hint == "" {
		t.Error("GetSolutionHint returned empty string")
	}
	if !contains(hint, "configuration") {
		t.Errorf("Expected 'configuration' in hint: %s", hint)
	}
}

func TestGetSolutionHint_IOError(t *testing.T) {
	err := NewNeevError(ErrTypeIO, "file error", nil)
	hint := err.GetSolutionHint()
	if hint == "" {
		t.Error("GetSolutionHint returned empty string")
	}
	if !contains(hint, "files") {
		t.Errorf("Expected 'files' in hint: %s", hint)
	}
}

func TestGetSolutionHint_ValidationError(t *testing.T) {
	err := NewNeevError(ErrTypeValidation, "invalid input", nil)
	hint := err.GetSolutionHint()
	if hint == "" {
		t.Error("GetSolutionHint returned empty string")
	}
	if !contains(hint, "input") {
		t.Errorf("Expected 'input' in hint: %s", hint)
	}
}

func TestGetSolutionHint_UnknownError(t *testing.T) {
	err := NewNeevError(ErrTypeUnknown, "something broke", nil)
	hint := err.GetSolutionHint()
	if hint == "" {
		t.Error("GetSolutionHint returned empty string")
	}
	if !contains(hint, "unknown") {
		t.Errorf("Expected 'unknown' in hint: %s", hint)
	}
}

func TestGetSolutionHint_InvalidErrorType(t *testing.T) {
	err := NewNeevError(ErrorType("invalid_type"), "message", nil)
	hint := err.GetSolutionHint()
	if hint == "" {
		t.Error("GetSolutionHint returned empty string for invalid type")
	}
}

func TestNewNeevError(t *testing.T) {
	underlying := errors.New("underlying")
	err := NewNeevError(ErrTypeIO, "io message", underlying)

	if err.Type != ErrTypeIO {
		t.Errorf("Expected ErrTypeIO, got %s", err.Type)
	}
	if err.Message != "io message" {
		t.Errorf("Expected 'io message', got '%s'", err.Message)
	}
	if err.Err != underlying {
		t.Error("Expected underlying error to be preserved")
	}
}

func TestNeevError_Unwrap(t *testing.T) {
	underlying := errors.New("underlying error")
	err := NewNeevError(ErrTypeIO, "io failed", underlying)

	unwrapped := err.Unwrap()
	if unwrapped != underlying {
		t.Error("Unwrap() did not return the underlying error")
	}

	if !errors.Is(err, underlying) {
		t.Error("errors.Is() should work with unwrapped error")
	}
}

func TestNeevError_UnwrapNil(t *testing.T) {
	err := NewNeevError(ErrTypeIO, "io failed", nil)
	unwrapped := err.Unwrap()
	if unwrapped != nil {
		t.Error("Unwrap() should return nil when no underlying error")
	}
}

func TestErrorTypesExist(t *testing.T) {
	errorTypes := []ErrorType{
		ErrTypeBlueprintNotFound,
		ErrTypeFoundation,
		ErrTypeInvalidConfig,
		ErrTypeIO,
		ErrTypeValidation,
		ErrTypeUnknown,
	}

	for _, et := range errorTypes {
		if et == "" {
			t.Error("Error type should not be empty")
		}
	}
}

func TestBlueprintNotFoundMessage(t *testing.T) {
	testCases := []string{
		"feature-x",
		"my-feature",
		"test123",
		"",
	}

	for _, tc := range testCases {
		err := ErrBlueprintNotFound(tc)
		expectedMsg := "blueprint '" + tc + "' not found"
		if err.Message != expectedMsg {
			t.Errorf("For blueprint '%s': expected '%s', got '%s'", tc, expectedMsg, err.Message)
		}
	}
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
