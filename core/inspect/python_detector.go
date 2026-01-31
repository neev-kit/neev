package inspect

import (
	"regexp"
	"strings"
)

// PythonDetector handles Python language analysis
type PythonDetector struct{}

// NewPythonDetector creates a new Python language detector
func NewPythonDetector() *PythonDetector {
	return &PythonDetector{}
}

// Detect returns true if this is a Python file
func (d *PythonDetector) Detect(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".py")
}

// Language returns the language identifier
func (d *PythonDetector) Language() Language {
	return LangPython
}

// ExtractEndpoints finds HTTP endpoints in Python code
// Supports: Flask, FastAPI, Django
func (d *PythonDetector) ExtractEndpoints(filePath string, content []byte) ([]Endpoint, error) {
	var endpoints []Endpoint
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	
	// Patterns for different Python web frameworks
	patterns := []*regexp.Regexp{
		// Flask: @app.route("/path", methods=["GET"])
		regexp.MustCompile(`@(app|router|bp|blueprint)\.(route|get|post|put|delete|patch)\s*\(\s*["']([^"']+)["'](?:.*methods\s*=\s*\[["'](\w+)["']\])?`),
		
		// FastAPI: @app.get("/path") or @router.post("/path")
		regexp.MustCompile(`@(app|router)\.(get|post|put|delete|patch|options|head)\s*\(\s*["']([^"']+)["']`),
		
		// Django URL patterns: path('api/users', views.list_users)
		regexp.MustCompile(`path\s*\(\s*["']([^"']+)["']\s*,\s*([^,)]+)`),
		
		// Django URL patterns: url(r'^api/users$', views.list_users)
		regexp.MustCompile(`url\s*\(\s*r?["']([^"']+)["']\s*,\s*([^,)]+)`),
	}
	
	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		for _, pattern := range patterns {
			matches := pattern.FindStringSubmatch(trimmedLine)
			if matches != nil {
				var method, path, handler string
				
				if strings.Contains(trimmedLine, "@app") || strings.Contains(trimmedLine, "@router") {
					// Flask/FastAPI decorator
					if len(matches) >= 3 {
						decoratorMethod := strings.ToUpper(matches[2])
						path = matches[3]
						
						// Try to extract method from decorator or subsequent methods=
						if decoratorMethod == "ROUTE" {
							// Look for methods parameter
							if len(matches) > 4 && matches[4] != "" {
								method = strings.ToUpper(matches[4])
							} else {
								method = "GET" // Default
							}
						} else {
							method = decoratorMethod
						}
						
						// Next line should have the function name
						if lineNum+1 < len(lines) {
							nextLine := strings.TrimSpace(lines[lineNum+1])
							funcMatch := regexp.MustCompile(`^def\s+(\w+)`).FindStringSubmatch(nextLine)
							if funcMatch != nil {
								handler = funcMatch[1]
							}
						}
					}
				} else if strings.Contains(trimmedLine, "path(") || strings.Contains(trimmedLine, "url(") {
					// Django URL pattern
					if len(matches) >= 3 {
						method = "GET" // Django doesn't specify method in URL patterns
						path = matches[1]
						handler = strings.TrimSpace(matches[2])
					}
				}
				
				if method != "" && path != "" {
					endpoint := Endpoint{
						Method:   method,
						Path:     path,
						Handler:  handler,
						File:     filePath,
						Line:     lineNum + 1,
						Language: string(LangPython),
					}
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}
	
	return endpoints, nil
}

// ExtractFunctions finds function signatures in Python code
func (d *PythonDetector) ExtractFunctions(filePath string, content []byte) ([]FunctionSignature, error) {
	var functions []FunctionSignature
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	
	// Pattern: def function_name(params) -> return_type:
	funcPattern := regexp.MustCompile(`^(async\s+)?def\s+([a-zA-Z_]\w+)\s*\((.*?)\)(?:\s*->\s*([^:]+))?:`)
	
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		
		matches := funcPattern.FindStringSubmatch(line)
		if matches != nil {
			funcName := matches[2]
			paramsStr := matches[3]
			returnType := strings.TrimSpace(matches[4])
			
			// Skip private functions (starting with _)
			if strings.HasPrefix(funcName, "_") && !strings.HasPrefix(funcName, "__") {
				continue
			}
			
			// Parse parameters
			params := d.parsePythonParameters(paramsStr)
			
			// Parse return type
			var returns []ReturnSpec
			if returnType != "" {
				returns = append(returns, ReturnSpec{
					Type: returnType,
				})
			}
			
			// Determine visibility
			visibility := "public"
			if strings.HasPrefix(funcName, "_") {
				visibility = "private"
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

// parsePythonParameters parses Python function parameters
func (d *PythonDetector) parsePythonParameters(paramsStr string) []ParameterSpec {
	var params []ParameterSpec
	
	if strings.TrimSpace(paramsStr) == "" {
		return params
	}
	
	// Split by comma, handling default values
	parts := strings.Split(paramsStr, ",")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "self" || part == "cls" {
			continue
		}
		
		// Remove default values
		part = strings.Split(part, "=")[0]
		part = strings.TrimSpace(part)
		
		// Pattern: name: type or just name
		if strings.Contains(part, ":") {
			tokens := strings.SplitN(part, ":", 2)
			params = append(params, ParameterSpec{
				Name: strings.TrimSpace(tokens[0]),
				Type: strings.TrimSpace(tokens[1]),
			})
		} else {
			params = append(params, ParameterSpec{
				Name: part,
				Type: "Any",
			})
		}
	}
	
	return params
}
