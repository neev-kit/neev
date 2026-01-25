package inspect

import (
	"regexp"
	"strings"
)

// GoDetector handles Go language analysis
type GoDetector struct{}

// NewGoDetector creates a new Go language detector
func NewGoDetector() *GoDetector {
	return &GoDetector{}
}

// Detect returns true if this is a Go file
func (d *GoDetector) Detect(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".go")
}

// Language returns the language identifier
func (d *GoDetector) Language() Language {
	return LangGo
}

// ExtractEndpoints finds HTTP endpoints in Go code
// Supports: Gin, Echo, Fiber, Chi, net/http
func (d *GoDetector) ExtractEndpoints(filePath string, content []byte) ([]Endpoint, error) {
	var endpoints []Endpoint
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	
	// Patterns for different Go web frameworks
	patterns := []*regexp.Regexp{
		// Gin: router.GET("/path", handler)
		regexp.MustCompile(`(?i)(router|r|app|engine|e|g)\.(GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD)\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`),
		
		// Echo: app.GET("/path", handler)
		regexp.MustCompile(`(?i)(e|echo|app)\.(GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD)\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`),
		
		// Chi: r.Get("/path", handler)
		regexp.MustCompile(`(?i)(r|router)\.(Get|Post|Put|Delete|Patch|Options|Head)\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`),
		
		// net/http: http.HandleFunc("/path", handler)
		regexp.MustCompile(`http\.HandleFunc\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`),
		
		// ServeMux: mux.HandleFunc("/path", handler)
		regexp.MustCompile(`(?i)(mux|serveMux|router)\.(HandleFunc|Handle)\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`),
	}
	
	for lineNum, line := range lines {
		for _, pattern := range patterns {
			matches := pattern.FindStringSubmatch(line)
			if matches != nil {
				var method, path, handler string
				
				// Handle different match group structures
				if len(matches) >= 4 {
					if strings.Contains(matches[0], "HandleFunc") || strings.Contains(matches[0], "http.") {
						// http.HandleFunc pattern
						if len(matches) == 3 {
							method = "GET" // Default for HandleFunc
							path = matches[1]
							handler = matches[2]
						} else {
							method = strings.ToUpper(matches[2])
							path = matches[3]
							handler = matches[4]
						}
					} else {
						// Framework pattern (Gin, Echo, Chi)
						method = strings.ToUpper(matches[2])
						path = matches[3]
						if len(matches) > 4 {
							handler = matches[4]
						} else {
							handler = ""
						}
					}
					
					// Clean up handler name
					handler = strings.TrimSpace(handler)
					handler = strings.Split(handler, ",")[0] // Remove additional args
					
					endpoint := Endpoint{
						Method:  method,
						Path:    path,
						Handler: handler,
						File:    filePath,
						Line:    lineNum + 1,
						Language: string(LangGo),
					}
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}
	
	return endpoints, nil
}

// ExtractFunctions finds function signatures in Go code
func (d *GoDetector) ExtractFunctions(filePath string, content []byte) ([]FunctionSignature, error) {
	var functions []FunctionSignature
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	
	// Pattern: func FunctionName(params) returns
	funcPattern := regexp.MustCompile(`^func\s+(\*?\w+\.)?\s*([A-Z]\w+)\s*\((.*?)\)\s*(.*)$`)
	
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip test functions for now
		if strings.Contains(line, "Test") && strings.Contains(line, "*testing.T") {
			continue
		}
		
		matches := funcPattern.FindStringSubmatch(line)
		if matches != nil {
			receiver := matches[1]
			funcName := matches[2]
			paramsStr := matches[3]
			returnsStr := strings.TrimSpace(matches[4])
			
			// Skip if method receiver is present (we want exported functions primarily)
			if receiver != "" {
				continue
			}
			
			// Parse parameters
			params := d.parseGoParameters(paramsStr)
			
			// Parse return types
			returns := d.parseGoReturns(returnsStr)
			
			// Determine visibility (exported = public in Go)
			visibility := "private"
			if len(funcName) > 0 && funcName[0] >= 'A' && funcName[0] <= 'Z' {
				visibility = "public"
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

// parseGoParameters parses Go function parameters
func (d *GoDetector) parseGoParameters(paramsStr string) []ParameterSpec {
	var params []ParameterSpec
	
	if strings.TrimSpace(paramsStr) == "" {
		return params
	}
	
	// Split by comma, handling potential function types
	parts := strings.Split(paramsStr, ",")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		// Simple pattern: name type or just type
		tokens := strings.Fields(part)
		if len(tokens) >= 2 {
			params = append(params, ParameterSpec{
				Name: tokens[0],
				Type: strings.Join(tokens[1:], " "),
			})
		} else if len(tokens) == 1 {
			// Just type, no name
			params = append(params, ParameterSpec{
				Name: "",
				Type: tokens[0],
			})
		}
	}
	
	return params
}

// parseGoReturns parses Go function return types
func (d *GoDetector) parseGoReturns(returnsStr string) []ReturnSpec {
	var returns []ReturnSpec
	
	returnsStr = strings.TrimSpace(returnsStr)
	if returnsStr == "" {
		return returns
	}
	
	// Remove parentheses for multiple returns
	returnsStr = strings.TrimPrefix(returnsStr, "(")
	returnsStr = strings.TrimSuffix(returnsStr, ")")
	returnsStr = strings.TrimSuffix(returnsStr, "{") // Remove function body start
	returnsStr = strings.TrimSpace(returnsStr)
	
	if returnsStr == "" {
		return returns
	}
	
	// Split by comma
	parts := strings.Split(returnsStr, ",")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		returns = append(returns, ReturnSpec{
			Type: part,
		})
	}
	
	return returns
}
