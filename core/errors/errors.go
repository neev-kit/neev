package errors

import (
	"fmt"
)

// ErrorType represents the category of error
type ErrorType string

const (
	ErrTypeBlueprintNotFound ErrorType = "blueprint_not_found"
	ErrTypeFoundation        ErrorType = "foundation_missing"
	ErrTypeInvalidConfig     ErrorType = "invalid_config"
	ErrTypeIO                ErrorType = "io_error"
	ErrTypeValidation        ErrorType = "validation_error"
	ErrTypeUnknown           ErrorType = "unknown"
)

// NeevError is a custom error type for Neev-specific errors
type NeevError struct {
	Type    ErrorType
	Message string
	Err     error
}

// NewNeevError creates a new NeevError with the given type and message
func NewNeevError(errType ErrorType, message string, err error) *NeevError {
	return &NeevError{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// Error implements the error interface
func (e *NeevError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *NeevError) Unwrap() error {
	return e.Err
}

// GetSolutionHint provides a user-friendly solution hint based on the error type
func (e *NeevError) GetSolutionHint() string {
	switch e.Type {
	case ErrTypeBlueprintNotFound:
		return "Make sure the blueprint exists in the .neev/blueprints/ directory. Run `neev draft` to create one."
	case ErrTypeFoundation:
		return "Foundation is missing. Try running `neev init` first to set up your project."
	case ErrTypeInvalidConfig:
		return "Your neev.yaml configuration is invalid. Check the format and try again."
	case ErrTypeIO:
		return "An error occurred while reading or writing files. Check file permissions and disk space."
	case ErrTypeValidation:
		return "The provided input is invalid. Check the parameters and try again."
	case ErrTypeUnknown:
		fallthrough
	default:
		return "An unknown error occurred. Enable debug mode with NEEV_LOG=debug for more details."
	}
}

// ErrBlueprintNotFound returns a new ErrTypeBlueprintNotFound error
func ErrBlueprintNotFound(blueprintName string) *NeevError {
	return NewNeevError(
		ErrTypeBlueprintNotFound,
		fmt.Sprintf("blueprint '%s' not found", blueprintName),
		nil,
	)
}

// ErrFoundationMissing returns a new ErrTypeFoundation error
func ErrFoundationMissing() *NeevError {
	return NewNeevError(
		ErrTypeFoundation,
		"foundation directory or files not found",
		nil,
	)
}

// ErrInvalidConfig returns a new ErrTypeInvalidConfig error
func ErrInvalidConfig(reason string) *NeevError {
	return NewNeevError(
		ErrTypeInvalidConfig,
		fmt.Sprintf("invalid configuration: %s", reason),
		nil,
	)
}
