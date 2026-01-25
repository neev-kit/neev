package inspect

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ValidateFunctionSignatures checks if function signatures match specs
func ValidateFunctionSignatures(opts InspectOptions, analyzer *PolyglotAnalyzer) ([]Warning, error) {
	var warnings []Warning
	
	// Get foundation modules with descriptors
	_, descriptors, err := getFoundationModules(opts.FoundationPath, true)
	if err != nil {
		return warnings, err
	}
	
	// If no descriptors with function specs, nothing to validate
	hasSpecs := false
	for _, desc := range descriptors {
		if len(desc.ExpectedFunctions) > 0 {
			hasSpecs = true
			break
		}
	}
	
	if !hasSpecs {
		return warnings, nil
	}
	
	// Extract all functions from code
	allFunctions, err := analyzer.ExtractAllFunctions(opts.RootDir, opts.IgnoreDirs)
	if err != nil {
		return warnings, fmt.Errorf("failed to extract functions from code: %w", err)
	}
	
	// Create lookup map by function name
	funcMap := make(map[string][]FunctionSignature)
	for _, fn := range allFunctions {
		funcMap[fn.Name] = append(funcMap[fn.Name], fn)
	}
	
	// Validate each expected function
	for moduleName, descriptor := range descriptors {
		for _, expectedFunc := range descriptor.ExpectedFunctions {
			// Find matching functions
			actualFuncs, found := funcMap[expectedFunc.Name]
			
			if !found || len(actualFuncs) == 0 {
				warning := Warning{
					Type:     WarningMissingFunction,
					Module:   moduleName,
					Message:  fmt.Sprintf("Expected function '%s' not found in code", expectedFunc.Name),
					Severity: "error",
					Remediation: fmt.Sprintf("Implement function '%s' in %s files or update module descriptor", 
						expectedFunc.Name, expectedFunc.Language),
				}
				warnings = append(warnings, warning)
				continue
			}
			
			// Check if any actual function matches the signature
			matchFound := false
			for _, actualFunc := range actualFuncs {
				// Check if language matches (if specified)
				if expectedFunc.Language != "" {
					actualLang := string(DetectLanguageByExtension(actualFunc.File))
					if actualLang != expectedFunc.Language {
						continue
					}
				}
				
				// Check if file pattern matches (if specified)
				if expectedFunc.FilePattern != "" {
					matched, _ := filepath.Match(expectedFunc.FilePattern, filepath.Base(actualFunc.File))
					if !matched {
						continue
					}
				}
				
				// Compare signatures
				sigWarnings := compareSignatures(moduleName, expectedFunc, actualFunc)
				if len(sigWarnings) == 0 {
					matchFound = true
					break
				}
				
				// If this is the only match, report the differences
				if len(actualFuncs) == 1 {
					warnings = append(warnings, sigWarnings...)
				}
			}
			
			// If no exact match found but we have similar functions
			if !matchFound && len(actualFuncs) > 0 {
				// Report the closest match
				actualFunc := actualFuncs[0]
				sigWarnings := compareSignatures(moduleName, expectedFunc, actualFunc)
				warnings = append(warnings, sigWarnings...)
			}
		}
	}
	
	return warnings, nil
}

// compareSignatures compares expected and actual function signatures
func compareSignatures(moduleName string, expected FunctionSpec, actual FunctionSignature) []Warning {
	var warnings []Warning
	
	// Compare parameters
	if len(expected.Parameters) != len(actual.Parameters) {
		warning := Warning{
			Type:     WarningSignatureMismatch,
			Module:   moduleName,
			Message:  fmt.Sprintf("Function '%s' has %d parameters but expected %d (in %s:%d)", 
				actual.Name, len(actual.Parameters), len(expected.Parameters), 
				filepath.Base(actual.File), actual.Line),
			Severity: "warning",
			Remediation: fmt.Sprintf("Update function signature to match spec: %s", 
				formatExpectedSignature(expected)),
		}
		warnings = append(warnings, warning)
	} else {
		// Check each parameter
		for i, expectedParam := range expected.Parameters {
			if i >= len(actual.Parameters) {
				break
			}
			actualParam := actual.Parameters[i]
			
			// Compare parameter names
			if expectedParam.Name != "" && actualParam.Name != expectedParam.Name {
				warning := Warning{
					Type:     WarningSignatureMismatch,
					Module:   moduleName,
					Message:  fmt.Sprintf("Function '%s' parameter %d: name '%s' doesn't match expected '%s' (in %s:%d)", 
						actual.Name, i+1, actualParam.Name, expectedParam.Name, 
						filepath.Base(actual.File), actual.Line),
					Severity: "info",
					Remediation: fmt.Sprintf("Consider renaming parameter to '%s' for consistency", expectedParam.Name),
				}
				warnings = append(warnings, warning)
			}
			
			// Compare parameter types (normalized)
			if expectedParam.Type != "" && !typesMatch(expectedParam.Type, actualParam.Type) {
				warning := Warning{
					Type:     WarningSignatureMismatch,
					Module:   moduleName,
					Message:  fmt.Sprintf("Function '%s' parameter '%s': type '%s' doesn't match expected '%s' (in %s:%d)", 
						actual.Name, actualParam.Name, actualParam.Type, expectedParam.Type, 
						filepath.Base(actual.File), actual.Line),
					Severity: "warning",
					Remediation: fmt.Sprintf("Update parameter type to '%s'", expectedParam.Type),
				}
				warnings = append(warnings, warning)
			}
		}
	}
	
	// Compare return types
	if len(expected.Returns) != len(actual.Returns) {
		warning := Warning{
			Type:     WarningSignatureMismatch,
			Module:   moduleName,
			Message:  fmt.Sprintf("Function '%s' has %d return values but expected %d (in %s:%d)", 
				actual.Name, len(actual.Returns), len(expected.Returns), 
				filepath.Base(actual.File), actual.Line),
			Severity: "warning",
			Remediation: fmt.Sprintf("Update return types to match spec: %s", 
				formatExpectedSignature(expected)),
		}
		warnings = append(warnings, warning)
	} else {
		// Check each return type
		for i, expectedReturn := range expected.Returns {
			if i >= len(actual.Returns) {
				break
			}
			actualReturn := actual.Returns[i]
			
			if expectedReturn.Type != "" && !typesMatch(expectedReturn.Type, actualReturn.Type) {
				warning := Warning{
					Type:     WarningSignatureMismatch,
					Module:   moduleName,
					Message:  fmt.Sprintf("Function '%s' return type %d: '%s' doesn't match expected '%s' (in %s:%d)", 
						actual.Name, i+1, actualReturn.Type, expectedReturn.Type, 
						filepath.Base(actual.File), actual.Line),
					Severity: "warning",
					Remediation: fmt.Sprintf("Update return type to '%s'", expectedReturn.Type),
				}
				warnings = append(warnings, warning)
			}
		}
	}
	
	// Compare visibility
	if expected.Visibility != "" && expected.Visibility != actual.Visibility {
		warning := Warning{
			Type:     WarningSignatureMismatch,
			Module:   moduleName,
			Message:  fmt.Sprintf("Function '%s' has visibility '%s' but expected '%s' (in %s:%d)", 
				actual.Name, actual.Visibility, expected.Visibility, 
				filepath.Base(actual.File), actual.Line),
			Severity: "info",
			Remediation: fmt.Sprintf("Change visibility to '%s' if intended for external use", expected.Visibility),
		}
		warnings = append(warnings, warning)
	}
	
	return warnings
}

// typesMatch checks if two type strings match (with normalization)
func typesMatch(expected, actual string) bool {
	// Normalize types
	expected = normalizeType(expected)
	actual = normalizeType(actual)
	
	// Exact match
	if expected == actual {
		return true
	}
	
	// Handle common variations
	// e.g., "string" vs "String", "int" vs "Integer", etc.
	expectedLower := strings.ToLower(expected)
	actualLower := strings.ToLower(actual)
	
	if expectedLower == actualLower {
		return true
	}
	
	// Handle pointer vs non-pointer in Go
	if strings.HasPrefix(actual, "*") && expected == strings.TrimPrefix(actual, "*") {
		return true
	}
	if strings.HasPrefix(expected, "*") && actual == strings.TrimPrefix(expected, "*") {
		return true
	}
	
	return false
}

// normalizeType normalizes type strings for comparison
func normalizeType(typ string) string {
	typ = strings.TrimSpace(typ)
	
	// Remove package prefixes for common types
	// e.g., "http.ResponseWriter" -> "ResponseWriter"
	parts := strings.Split(typ, ".")
	if len(parts) > 1 {
		// Keep the last part only for simple comparison
		// This is a simplification but works for many cases
		return parts[len(parts)-1]
	}
	
	return typ
}

// formatExpectedSignature formats expected function signature for display
func formatExpectedSignature(spec FunctionSpec) string {
	var parts []string
	
	// Add visibility
	if spec.Visibility != "" {
		parts = append(parts, spec.Visibility)
	}
	
	// Add function name
	parts = append(parts, spec.Name)
	
	// Add parameters
	var params []string
	for _, param := range spec.Parameters {
		if param.Name != "" {
			params = append(params, fmt.Sprintf("%s %s", param.Name, param.Type))
		} else {
			params = append(params, param.Type)
		}
	}
	signature := strings.Join(parts, " ") + "(" + strings.Join(params, ", ") + ")"
	
	// Add return types
	if len(spec.Returns) > 0 {
		var returns []string
		for _, ret := range spec.Returns {
			returns = append(returns, ret.Type)
		}
		signature += " -> (" + strings.Join(returns, ", ") + ")"
	}
	
	return signature
}
