package inspect

import (
	"regexp"
	"strings"
)

// JavaDetector handles Java language analysis
type JavaDetector struct{}

// NewJavaDetector creates a new Java language detector
func NewJavaDetector() *JavaDetector {
	return &JavaDetector{}
}

// Detect returns true if this is a Java file
func (d *JavaDetector) Detect(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".java")
}

// Language returns the language identifier
func (d *JavaDetector) Language() Language {
	return LangJava
}

// ExtractEndpoints finds HTTP endpoints in Java code
// Supports: Spring Boot annotations
func (d *JavaDetector) ExtractEndpoints(filePath string, content []byte) ([]Endpoint, error) {
	var endpoints []Endpoint
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	
	// Spring Boot annotation patterns
	annotationPattern := regexp.MustCompile(`@(GetMapping|PostMapping|PutMapping|DeleteMapping|PatchMapping|RequestMapping)\s*(?:\((?:value\s*=\s*)?["']([^"']+)["'])?`)
	methodPattern := regexp.MustCompile(`^\s*(?:public|private|protected)?\s+(?:\w+\s+)*(\w+)\s*\(`)
	
	var currentAnnotation string
	var currentPath string
	var currentMethod string
	
	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		// Check for Spring annotation
		annotationMatch := annotationPattern.FindStringSubmatch(trimmedLine)
		if annotationMatch != nil {
			annotation := annotationMatch[1]
			path := ""
			if len(annotationMatch) > 2 {
				path = annotationMatch[2]
			}
			
			// Convert annotation to HTTP method
			method := "GET"
			switch annotation {
			case "PostMapping":
				method = "POST"
			case "PutMapping":
				method = "PUT"
			case "DeleteMapping":
				method = "DELETE"
			case "PatchMapping":
				method = "PATCH"
			case "RequestMapping":
				// Try to find method in annotation
				if strings.Contains(trimmedLine, "method") {
					if strings.Contains(trimmedLine, "POST") {
						method = "POST"
					} else if strings.Contains(trimmedLine, "PUT") {
						method = "PUT"
					} else if strings.Contains(trimmedLine, "DELETE") {
						method = "DELETE"
					}
				}
			}
			
			currentAnnotation = annotation
			currentPath = path
			currentMethod = method
			continue
		}
		
		// If we have a pending annotation, look for the method definition
		if currentAnnotation != "" {
			methodMatch := methodPattern.FindStringSubmatch(trimmedLine)
			if methodMatch != nil {
				handler := methodMatch[1]
				
				endpoint := Endpoint{
					Method:   currentMethod,
					Path:     currentPath,
					Handler:  handler,
					File:     filePath,
					Line:     lineNum + 1,
					Language: string(LangJava),
				}
				endpoints = append(endpoints, endpoint)
				
				// Reset
				currentAnnotation = ""
				currentPath = ""
				currentMethod = ""
			}
		}
	}
	
	return endpoints, nil
}

// ExtractFunctions finds function signatures in Java code
func (d *JavaDetector) ExtractFunctions(filePath string, content []byte) ([]FunctionSignature, error) {
	var functions []FunctionSignature
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	
	// Pattern: public ReturnType methodName(params)
	funcPattern := regexp.MustCompile(`^\s*(public|private|protected)?\s+(?:static\s+)?(?:final\s+)?(?:abstract\s+)?(\w+(?:<[^>]+>)?)\s+([a-zA-Z_]\w*)\s*\((.*?)\)`)
	
	for lineNum, line := range lines {
		matches := funcPattern.FindStringSubmatch(line)
		if matches != nil {
			visibility := matches[1]
			if visibility == "" {
				visibility = "package" // default in Java
			}
			returnType := matches[2]
			funcName := matches[3]
			paramsStr := matches[4]
			
			// Skip constructors
			if returnType == funcName {
				continue
			}
			
			// Parse parameters
			params := d.parseJavaParameters(paramsStr)
			
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

// parseJavaParameters parses Java method parameters
func (d *JavaDetector) parseJavaParameters(paramsStr string) []ParameterSpec {
	var params []ParameterSpec
	
	if strings.TrimSpace(paramsStr) == "" {
		return params
	}
	
	// Split by comma, handling generics
	parts := d.splitJavaParams(paramsStr)
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		// Remove annotations like @RequestBody
		part = regexp.MustCompile(`@\w+\s+`).ReplaceAllString(part, "")
		part = strings.TrimSpace(part)
		
		// Pattern: Type name or final Type name
		tokens := strings.Fields(part)
		if len(tokens) >= 2 {
			// Last token is name, rest is type
			paramName := tokens[len(tokens)-1]
			paramType := strings.Join(tokens[:len(tokens)-1], " ")
			// Remove final keyword
			paramType = strings.TrimPrefix(paramType, "final ")
			
			params = append(params, ParameterSpec{
				Name: paramName,
				Type: paramType,
			})
		}
	}
	
	return params
}

// splitJavaParams splits parameters considering generics
func (d *JavaDetector) splitJavaParams(paramsStr string) []string {
	var parts []string
	var current strings.Builder
	depth := 0
	
	for _, ch := range paramsStr {
		switch ch {
		case '<':
			depth++
			current.WriteRune(ch)
		case '>':
			depth--
			current.WriteRune(ch)
		case ',':
			if depth == 0 {
				parts = append(parts, current.String())
				current.Reset()
			} else {
				current.WriteRune(ch)
			}
		default:
			current.WriteRune(ch)
		}
	}
	
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	
	return parts
}
