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
	// WarningMissingEndpoint indicates an API endpoint is documented but not implemented
	WarningMissingEndpoint WarningType = "MISSING_ENDPOINT"
	// WarningUndocumentedEndpoint indicates an endpoint exists but is not documented
	WarningUndocumentedEndpoint WarningType = "UNDOCUMENTED_ENDPOINT"
	// WarningSignatureMismatch indicates function signature doesn't match spec
	WarningSignatureMismatch WarningType = "SIGNATURE_MISMATCH"
	// WarningMissingFunction indicates an expected function is not found
	WarningMissingFunction WarningType = "MISSING_FUNCTION"
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
	TotalModules       int                `json:"total_modules"`
	MatchingModules    int                `json:"matching_modules"`
	MissingModules     int                `json:"missing_modules"`
	ExtraCodeDirs      int                `json:"extra_code_dirs"`
	TotalWarnings      int                `json:"total_warnings"`
	ErrorCount         int                `json:"error_count"`
	WarningCount       int                `json:"warning_count"`
	Languages          map[string]int     `json:"languages,omitempty"`           // Language name -> file count
	MissingEndpoints   int                `json:"missing_endpoints,omitempty"`   // Level 2
	UndocumentedEnds   int                `json:"undocumented_endpoints,omitempty"` // Level 2
	SignatureMismatches int               `json:"signature_mismatches,omitempty"` // Level 3
}

// ModuleDescriptor defines the expected structure of a module
type ModuleDescriptor struct {
	Name              string             `yaml:"name"`
	Description       string             `yaml:"description"`
	ExpectedFiles     []string           `yaml:"expected_files"`
	ExpectedDirs      []string           `yaml:"expected_dirs"`
	Patterns          []string           `yaml:"patterns"` // Glob patterns for files
	ExpectedFunctions []FunctionSpec     `yaml:"expected_functions,omitempty"` // Level 3
}

// FunctionSpec defines expected function/method signatures
type FunctionSpec struct {
	Name        string          `yaml:"name"`
	Language    string          `yaml:"language"` // go, python, javascript, java, csharp, ruby
	FilePattern string          `yaml:"file_pattern,omitempty"` // Where to find it
	Parameters  []ParameterSpec `yaml:"parameters,omitempty"`
	Returns     []ReturnSpec    `yaml:"returns,omitempty"`
	Visibility  string          `yaml:"visibility,omitempty"` // public, private, protected
}

// ParameterSpec defines a function parameter
type ParameterSpec struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

// ReturnSpec defines a return type
type ReturnSpec struct {
	Type string `yaml:"type"`
	Name string `yaml:"name,omitempty"` // For named returns
}

// Endpoint represents an API endpoint (HTTP handler)
type Endpoint struct {
	Method      string   `json:"method"`
	Path        string   `json:"path"`
	Handler     string   `json:"handler,omitempty"`
	File        string   `json:"file,omitempty"`
	Line        int      `json:"line,omitempty"`
	Language    string   `json:"language,omitempty"`
}

// Language represents a detected programming language
type Language string

const (
	LangGo         Language = "go"
	LangPython     Language = "python"
	LangJavaScript Language = "javascript"
	LangTypeScript Language = "typescript"
	LangJava       Language = "java"
	LangCSharp     Language = "csharp"
	LangRuby       Language = "ruby"
)
