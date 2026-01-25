package inspect

import (
	"regexp"
	"strings"
)

// CSharpDetector handles C# language analysis
type CSharpDetector struct{}

// NewCSharpDetector creates a new C# language detector
func NewCSharpDetector() *CSharpDetector {
	return &CSharpDetector{}
}

// Detect returns true if this is a C# file
func (d *CSharpDetector) Detect(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".cs")
}

// Language returns the language identifier
func (d *CSharpDetector) Language() Language {
	return LangCSharp
}

// ExtractEndpoints finds HTTP endpoints in C# code
// Supports: ASP.NET Core attributes
func (d *CSharpDetector) ExtractEndpoints(filePath string, content []byte) ([]Endpoint, error) {
	var endpoints []Endpoint
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	
	// ASP.NET Core attribute patterns
	annotationPattern := regexp.MustCompile(`\[(HttpGet|HttpPost|HttpPut|HttpDelete|HttpPatch)\s*(?:\(["']([^"']+)["']\))?`)
	methodPattern := regexp.MustCompile(`^\s*(?:public|private|protected|internal)?\s+(?:async\s+)?(?:Task<)?(\w+)`)
	
	var currentMethod string
	var currentPath string
	
	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		// Check for HTTP attribute
		annotationMatch := annotationPattern.FindStringSubmatch(trimmedLine)
		if annotationMatch != nil {
			attr := annotationMatch[1]
			path := ""
			if len(annotationMatch) > 2 {
				path = annotationMatch[2]
			}
			
			// Convert attribute to HTTP method
			method := strings.ToUpper(strings.TrimPrefix(attr, "Http"))
			
			currentMethod = method
			currentPath = path
			continue
		}
		
		// If we have a pending attribute, look for the method definition
		if currentMethod != "" {
			methodMatch := methodPattern.FindStringSubmatch(trimmedLine)
			if methodMatch != nil && strings.Contains(trimmedLine, "(") {
				// Extract method name
				funcNamePattern := regexp.MustCompile(`\s+(\w+)\s*\(`)
				funcMatch := funcNamePattern.FindStringSubmatch(trimmedLine)
				var handler string
				if funcMatch != nil {
					handler = funcMatch[1]
				}
				
				endpoint := Endpoint{
					Method:   currentMethod,
					Path:     currentPath,
					Handler:  handler,
					File:     filePath,
					Line:     lineNum + 1,
					Language: string(LangCSharp),
				}
				endpoints = append(endpoints, endpoint)
				
				// Reset
				currentMethod = ""
				currentPath = ""
			}
		}
	}
	
	return endpoints, nil
}

// ExtractFunctions finds function signatures in C# code
func (d *CSharpDetector) ExtractFunctions(filePath string, content []byte) ([]FunctionSignature, error) {
	var functions []FunctionSignature
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	
	// Pattern: public ReturnType MethodName(params)
	funcPattern := regexp.MustCompile(`^\s*(public|private|protected|internal)?\s+(?:static\s+)?(?:async\s+)?(?:virtual\s+)?(?:override\s+)?(Task<\w+>|\w+(?:<[^>]+>)?)\s+([A-Z]\w*)\s*\((.*?)\)`)
	
	for lineNum, line := range lines {
		matches := funcPattern.FindStringSubmatch(line)
		if matches != nil {
			visibility := matches[1]
			if visibility == "" {
				visibility = "private" // default in C#
			}
			returnType := matches[2]
			funcName := matches[3]
			paramsStr := matches[4]
			
			// Parse parameters
			params := d.parseCSharpParameters(paramsStr)
			
			// Parse return type
			var returns []ReturnSpec
			if returnType != "void" {
				returns = append(returns, ReturnSpec{
					Type: returnType,
				})
			}
			
			function := FunctionSignature{
				Name:       funcName,
				Parameters: params,
				Returns:    returns,
				File:       filePath,
				Line:       lineNum + 1,
				Visibility: visibility,
			}
			functions = append(functions, function)
		}
	}
	
	return functions, nil
}

// parseCSharpParameters parses C# method parameters
func (d *CSharpDetector) parseCSharpParameters(paramsStr string) []ParameterSpec {
	var params []ParameterSpec
	
	if strings.TrimSpace(paramsStr) == "" {
		return params
	}
	
	// Split by comma
	parts := strings.Split(paramsStr, ",")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		// Remove parameter modifiers
		part = strings.TrimPrefix(part, "ref ")
		part = strings.TrimPrefix(part, "out ")
		part = strings.TrimPrefix(part, "params ")
		part = strings.TrimSpace(part)
		
		// Pattern: Type name or Type name = default
		tokens := strings.Fields(strings.Split(part, "=")[0])
		if len(tokens) >= 2 {
			// Last token is name, rest is type
			paramName := tokens[len(tokens)-1]
			paramType := strings.Join(tokens[:len(tokens)-1], " ")
			
			params = append(params, ParameterSpec{
				Name: paramName,
				Type: paramType,
			})
		}
	}
	
	return params
}
