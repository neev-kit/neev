package inspect

// WarningType categorizes different types of drift warnings
type WarningType string

const (
	// WarningMissingModule indicates a foundation spec exists but no matching code directory
	WarningMissingModule WarningType = "MISSING_MODULE"
	// WarningExtraCode indicates code exists without a foundation spec
	WarningExtraCode WarningType = "EXTRA_CODE"
	// WarningMismatchedName indicates a naming inconsistency
	WarningMismatchedName WarningType = "MISMATCHED_NAME"
	// WarningMissingFile indicates an expected file from descriptor is missing
	WarningMissingFile WarningType = "MISSING_FILE"
	// WarningUnexpectedFile indicates a file exists that wasn't expected
	WarningUnexpectedFile WarningType = "UNEXPECTED_FILE"
)

// Warning represents a single drift detection warning
type Warning struct {
	Type        WarningType `json:"type"`
	Module      string      `json:"module"`
	Message     string      `json:"message"`
	Severity    string      `json:"severity"`    // "error", "warning", "info"
	Remediation string      `json:"remediation"` // Suggested fix
}

// InspectResult contains the complete result of an inspection
type InspectResult struct {
	Success  bool      `json:"success"`
	Warnings []Warning `json:"warnings"`
	Summary  Summary   `json:"summary"`
}

// Summary provides high-level statistics about the inspection
type Summary struct {
	TotalModules       int `json:"total_modules"`
	MatchingModules    int `json:"matching_modules"`
	MissingModules     int `json:"missing_modules"`
	ExtraCodeDirs      int `json:"extra_code_dirs"`
	TotalWarnings      int `json:"total_warnings"`
	ErrorCount         int `json:"error_count"`
	WarningCount       int `json:"warning_count"`
}

// ModuleDescriptor defines the expected structure of a module
type ModuleDescriptor struct {
	Name          string   `yaml:"name"`
	Description   string   `yaml:"description"`
	ExpectedFiles []string `yaml:"expected_files"`
	ExpectedDirs  []string `yaml:"expected_dirs"`
	Patterns      []string `yaml:"patterns"` // Glob patterns for files
}
