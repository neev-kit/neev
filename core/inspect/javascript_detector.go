package inspect

import (
	"regexp"
	"strings"
)

// JavaScriptDetector handles JavaScript and TypeScript language analysis
type JavaScriptDetector struct{}

// NewJavaScriptDetector creates a new JavaScript/TypeScript language detector
func NewJavaScriptDetector() *JavaScriptDetector {
	return &JavaScriptDetector{}
}

// Detect returns true if this is a JavaScript or TypeScript file
func (d *JavaScriptDetector) Detect(filePath string) bool {
	lower := strings.ToLower(filePath)
	return strings.HasSuffix(lower, ".js") || 
		   strings.HasSuffix(lower, ".ts") || 
		   strings.HasSuffix(lower, ".jsx") || 
		   strings.HasSuffix(lower, ".tsx")
}

// Language returns the language identifier
func (d *JavaScriptDetector) Language() Language {
	return LangJavaScript
}

// ExtractEndpoints finds HTTP endpoints in JavaScript/TypeScript code
// Supports: Express, Fastify, Koa
func (d *JavaScriptDetector) ExtractEndpoints(filePath string, content []byte) ([]Endpoint, error) {
	var endpoints []Endpoint
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	
	// Patterns for different JavaScript web frameworks
	patterns := []*regexp.Regexp{
		// Express: app.get("/path", handler) or router.post("/path", handler)
		regexp.MustCompile(`(app|router|express)\.(get|post|put|delete|patch|options|head|all)\s*\(\s*["'\x60]([^"'\x60]+)["'\x60]\s*,\s*([^)]+)\)`),
		
		// Express with arrow functions: app.get("/path", (req, res) => {...})
		regexp.MustCompile(`(app|router)\.(get|post|put|delete|patch)\s*\(\s*["'\x60]([^"'\x60]+)["'\x60]\s*,\s*(?:async\s*)?\([^)]*\)\s*=>`),
		
		// Fastify: fastify.get("/path", handler)
		regexp.MustCompile(`(fastify|server)\.(get|post|put|delete|patch|options|head)\s*\(\s*["'\x60]([^"'\x60]+)["'\x60]`),
		
		// Koa: router.get("/path", handler)
		regexp.MustCompile(`(router)\.(get|post|put|delete|patch|options|head)\s*\(\s*["'\x60]([^"'\x60]+)["'\x60]`),
	}
	
	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		for _, pattern := range patterns {
			matches := pattern.FindStringSubmatch(trimmedLine)
			if matches != nil {
				if len(matches) >= 4 {
					method := strings.ToUpper(matches[2])
					path := matches[3]
					handler := ""
					
					if len(matches) > 4 {
						handler = strings.TrimSpace(matches[4])
						// Clean up handler name
						handler = strings.Split(handler, ",")[0]
						handler = strings.TrimSpace(handler)
					}
					
					endpoint := Endpoint{
						Method:   method,
						Path:     path,
						Handler:  handler,
						File:     filePath,
						Line:     lineNum + 1,
						Language: string(LangJavaScript),
					}
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}
	
	return endpoints, nil
}

// ExtractFunctions finds function signatures in JavaScript/TypeScript code
func (d *JavaScriptDetector) ExtractFunctions(filePath string, content []byte) ([]FunctionSignature, error) {
	var functions []FunctionSignature
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	
	// Patterns for different function declarations
	patterns := []*regexp.Regexp{
		// function functionName(params): returnType
		regexp.MustCompile(`^(?:export\s+)?(?:async\s+)?function\s+([a-zA-Z_$][\w$]*)\s*\((.*?)\)(?:\s*:\s*([^{]+))?`),
		
		// const functionName = (params): returnType =>
		regexp.MustCompile(`^(?:export\s+)?const\s+([a-zA-Z_$][\w$]*)\s*=\s*(?:async\s*)?\((.*?)\)(?:\s*:\s*([^=]+))?\s*=>`),
		
		// TypeScript method: public methodName(params): returnType
		regexp.MustCompile(`^(public|private|protected)?\s*(?:async\s+)?([a-zA-Z_$][\w$]*)\s*\((.*?)\)(?:\s*:\s*([^{]+))?`),
	}
	
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		
		for _, pattern := range patterns {
			matches := pattern.FindStringSubmatch(line)
			if matches != nil {
				var funcName, paramsStr, returnType, visibility string
				
				// Determine which pattern matched and extract accordingly
				if strings.Contains(matches[0], "function") {
					funcName = matches[1]
					paramsStr = matches[2]
					if len(matches) > 3 {
						returnType = strings.TrimSpace(matches[3])
					}
					visibility = "public"
				} else if strings.Contains(matches[0], "const") {
					funcName = matches[1]
					paramsStr = matches[2]
					if len(matches) > 3 {
						returnType = strings.TrimSpace(matches[3])
					}
					visibility = "public"
				} else if len(matches) > 2 {
					// TypeScript method with visibility
					if matches[1] != "" {
						visibility = matches[1]
					} else {
						visibility = "public"
					}
					funcName = matches[2]
					paramsStr = matches[3]
					if len(matches) > 4 {
						returnType = strings.TrimSpace(matches[4])
					}
				}
				
				// Skip if no function name
				if funcName == "" {
					continue
				}
				
				// Parse parameters
				params := d.parseJSParameters(paramsStr)
				
				// Parse return type
				var returns []ReturnSpec
				if returnType != "" {
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
	}
	
	return functions, nil
}

// parseJSParameters parses JavaScript/TypeScript function parameters
func (d *JavaScriptDetector) parseJSParameters(paramsStr string) []ParameterSpec {
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
		
		// Remove default values
		part = strings.Split(part, "=")[0]
		part = strings.TrimSpace(part)
		
		// Pattern: name: type or just name
		if strings.Contains(part, ":") {
			tokens := strings.SplitN(part, ":", 2)
			paramName := strings.TrimSpace(tokens[0])
			paramType := strings.TrimSpace(tokens[1])
			
			// Remove optional marker (?)
			paramName = strings.TrimSuffix(paramName, "?")
			
			params = append(params, ParameterSpec{
				Name: paramName,
				Type: paramType,
			})
		} else {
			// Remove spread operator
			part = strings.TrimPrefix(part, "...")
			params = append(params, ParameterSpec{
				Name: part,
				Type: "any",
			})
		}
	}
	
	return params
}
