package inspect

import (
	"regexp"
	"strings"
)

// RubyDetector handles Ruby language analysis
type RubyDetector struct{}

// NewRubyDetector creates a new Ruby language detector
func NewRubyDetector() *RubyDetector {
	return &RubyDetector{}
}

// Detect returns true if this is a Ruby file
func (d *RubyDetector) Detect(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".rb")
}

// Language returns the language identifier
func (d *RubyDetector) Language() Language {
	return LangRuby
}

// ExtractEndpoints finds HTTP endpoints in Ruby code
// Supports: Rails routes, Sinatra routes
func (d *RubyDetector) ExtractEndpoints(filePath string, content []byte) ([]Endpoint, error) {
	var endpoints []Endpoint
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	
	// Patterns for different Ruby web frameworks
	patterns := []*regexp.Regexp{
		// Rails routes: get '/api/users', to: 'users#index'
		regexp.MustCompile(`(get|post|put|delete|patch)\s+['"]([^'"]+)['"],?\s*(?:to:\s*['"]([^'"]+)['"])?`),
		
		// Sinatra: get '/api/users' do
		regexp.MustCompile(`^(get|post|put|delete|patch)\s+['"]([^'"]+)['"]\s+do`),
		
		// Rails resources: resources :users
		regexp.MustCompile(`resources\s+:(\w+)`),
	}
	
	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		for _, pattern := range patterns {
			matches := pattern.FindStringSubmatch(trimmedLine)
			if matches != nil {
				method := strings.ToUpper(matches[1])
				path := matches[2]
				handler := ""
				
				if len(matches) > 3 && matches[3] != "" {
					handler = matches[3]
				}
				
				// Handle resources shorthand
				if strings.Contains(trimmedLine, "resources") {
					// Rails resources create standard RESTful routes
					resourceMatch := regexp.MustCompile(`resources\s+:(\w+)`).FindStringSubmatch(trimmedLine)
					if resourceMatch != nil {
						resource := resourceMatch[1]
						// We'll just record the basic ones
						restMethods := []struct {
							method string
							path   string
							action string
						}{
							{"GET", "/" + resource, "index"},
							{"POST", "/" + resource, "create"},
							{"GET", "/" + resource + "/:id", "show"},
							{"PUT", "/" + resource + "/:id", "update"},
							{"DELETE", "/" + resource + "/:id", "destroy"},
						}
						
						for _, rm := range restMethods {
							endpoint := Endpoint{
								Method:   rm.method,
								Path:     rm.path,
								Handler:  resource + "#" + rm.action,
								File:     filePath,
								Line:     lineNum + 1,
								Language: string(LangRuby),
							}
							endpoints = append(endpoints, endpoint)
						}
						continue
					}
				}
				
				endpoint := Endpoint{
					Method:   method,
					Path:     path,
					Handler:  handler,
					File:     filePath,
					Line:     lineNum + 1,
					Language: string(LangRuby),
				}
				endpoints = append(endpoints, endpoint)
			}
		}
	}
	
	return endpoints, nil
}

// ExtractFunctions finds function signatures in Ruby code
func (d *RubyDetector) ExtractFunctions(filePath string, content []byte) ([]FunctionSignature, error) {
	var functions []FunctionSignature
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	
	// Pattern: def method_name(params)
	funcPattern := regexp.MustCompile(`^def\s+(self\.)?([a-zA-Z_]\w*[!?]?)\s*(?:\((.*?)\))?`)
	
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		
		matches := funcPattern.FindStringSubmatch(line)
		if matches != nil {
			isClassMethod := matches[1] != ""
			funcName := matches[2]
			paramsStr := ""
			if len(matches) > 3 {
				paramsStr = matches[3]
			}
			
			// Skip private methods (by convention, though detection is hard)
			// In Ruby, private is usually declared separately
			
			// Parse parameters
			params := d.parseRubyParameters(paramsStr)
			
			// Ruby doesn't have static typing by default, so return type is unknown
			var returns []ReturnSpec
			
			// Determine visibility
			visibility := "public"
			if isClassMethod {
				visibility = "public" // class methods are typically public
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

// parseRubyParameters parses Ruby method parameters
func (d *RubyDetector) parseRubyParameters(paramsStr string) []ParameterSpec {
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
		paramName := strings.Split(part, "=")[0]
		paramName = strings.TrimSpace(paramName)
		
		// Remove splat operator (*)
		paramName = strings.TrimPrefix(paramName, "*")
		paramName = strings.TrimPrefix(paramName, "**")
		
		// Remove keyword argument indicator (:)
		paramName = strings.TrimSuffix(paramName, ":")
		
		params = append(params, ParameterSpec{
			Name: paramName,
			Type: "Object", // Ruby is dynamically typed
		})
	}
	
	return params
}
