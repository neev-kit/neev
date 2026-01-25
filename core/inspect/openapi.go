package inspect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/neev-kit/neev/core/openapi"
	"gopkg.in/yaml.v3"
)

// OpenAPISpec represents a parsed OpenAPI specification
type OpenAPISpec struct {
	Endpoints []openapi.Endpoint
}

// ValidateOpenAPIContracts checks if code implements documented API endpoints
func ValidateOpenAPIContracts(opts InspectOptions, analyzer *PolyglotAnalyzer) ([]Warning, error) {
	var warnings []Warning
	
	// Try to find openapi.yaml or architecture.md in foundation/blueprints
	var specEndpoints []openapi.Endpoint
	var err error
	
	// Check for openapi.yaml in blueprints
	blueprintsPath := filepath.Join(opts.RootDir, ".neev", "blueprints")
	if stat, statErr := os.Stat(blueprintsPath); statErr == nil && stat.IsDir() {
		// Search for openapi.yaml files in blueprint directories
		filepath.Walk(blueprintsPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() && (info.Name() == "openapi.yaml" || info.Name() == "openapi.yml") {
				endpoints, parseErr := parseOpenAPIFile(path)
				if parseErr == nil {
					specEndpoints = append(specEndpoints, endpoints...)
				}
			}
			return nil
		})
	}
	
	// Also check for ARCHITECTURE.md in blueprints
	filepath.Walk(blueprintsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.ToLower(info.Name()) == "architecture.md" {
			endpoints, parseErr := openapi.ParseArchitecture(path)
			if parseErr == nil {
				specEndpoints = append(specEndpoints, endpoints...)
			}
		}
		return nil
	})
	
	// If no spec found, check foundation directory
	if len(specEndpoints) == 0 {
		archPath := filepath.Join(opts.FoundationPath, "ARCHITECTURE.md")
		if _, statErr := os.Stat(archPath); statErr == nil {
			specEndpoints, err = openapi.ParseArchitecture(archPath)
			if err != nil {
				// Not fatal, just no spec to validate against
				return warnings, nil
			}
		}
	}
	
	// If still no spec, nothing to validate
	if len(specEndpoints) == 0 {
		return warnings, nil
	}
	
	// Extract implemented endpoints from code
	implementedEndpoints, err := analyzer.ExtractAllEndpoints(opts.RootDir, opts.IgnoreDirs)
	if err != nil {
		return warnings, fmt.Errorf("failed to extract endpoints from code: %w", err)
	}
	
	// Compare documented vs implemented
	warnings = append(warnings, compareEndpoints(specEndpoints, implementedEndpoints)...)
	
	return warnings, nil
}

// parseOpenAPIFile parses an OpenAPI YAML file and extracts endpoints
func parseOpenAPIFile(filePath string) ([]openapi.Endpoint, error) {
	var endpoints []openapi.Endpoint
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	
	var spec map[string]interface{}
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, err
	}
	
	// Extract paths from OpenAPI spec
	paths, ok := spec["paths"].(map[string]interface{})
	if !ok {
		return endpoints, nil
	}
	
	for path, pathItem := range paths {
		pathMap, ok := pathItem.(map[string]interface{})
		if !ok {
			continue
		}
		
		// Check each HTTP method
		for method, operation := range pathMap {
			methodUpper := strings.ToUpper(method)
			if methodUpper == "GET" || methodUpper == "POST" || methodUpper == "PUT" || 
			   methodUpper == "DELETE" || methodUpper == "PATCH" || methodUpper == "OPTIONS" {
				
				endpoint := openapi.Endpoint{
					Method: methodUpper,
					Path:   path,
				}
				
				// Try to extract description
				if opMap, ok := operation.(map[string]interface{}); ok {
					if desc, ok := opMap["summary"].(string); ok {
						endpoint.Description = desc
					} else if desc, ok := opMap["description"].(string); ok {
						endpoint.Description = desc
					}
				}
				
				endpoints = append(endpoints, endpoint)
			}
		}
	}
	
	return endpoints, nil
}

// compareEndpoints compares documented and implemented endpoints
func compareEndpoints(specEndpoints []openapi.Endpoint, implEndpoints []Endpoint) []Warning {
	var warnings []Warning
	
	// Create lookup maps
	specMap := make(map[string]openapi.Endpoint)
	for _, ep := range specEndpoints {
		key := fmt.Sprintf("%s %s", ep.Method, normalizePath(ep.Path))
		specMap[key] = ep
	}
	
	implMap := make(map[string]Endpoint)
	for _, ep := range implEndpoints {
		key := fmt.Sprintf("%s %s", ep.Method, normalizePath(ep.Path))
		implMap[key] = ep
	}
	
	// Check for missing implementations (documented but not implemented)
	for key, specEp := range specMap {
		if _, exists := implMap[key]; !exists {
			warning := Warning{
				Type:     WarningMissingEndpoint,
				Module:   "api",
				Message:  fmt.Sprintf("API endpoint %s %s is documented but not implemented", specEp.Method, specEp.Path),
				Severity: "error",
				Remediation: fmt.Sprintf("Implement handler for %s %s or remove from API documentation", 
					specEp.Method, specEp.Path),
			}
			warnings = append(warnings, warning)
		}
	}
	
	// Check for undocumented implementations (implemented but not documented)
	for key, implEp := range implMap {
		if _, exists := specMap[key]; !exists {
			warning := Warning{
				Type:     WarningUndocumentedEndpoint,
				Module:   "api",
				Message:  fmt.Sprintf("API endpoint %s %s is implemented but not documented (found in %s:%d)", 
					implEp.Method, implEp.Path, filepath.Base(implEp.File), implEp.Line),
				Severity: "warning",
				Remediation: fmt.Sprintf("Add %s %s to API documentation (OpenAPI/ARCHITECTURE.md)", 
					implEp.Method, implEp.Path),
			}
			warnings = append(warnings, warning)
		}
	}
	
	return warnings
}

// normalizePath normalizes API paths for comparison
// Converts different parameter styles to a common format
func normalizePath(path string) string {
	// Normalize path parameters
	// /api/users/:id -> /api/users/{id}
	// /api/users/<id> -> /api/users/{id}
	path = strings.ReplaceAll(path, ":", "")
	path = strings.ReplaceAll(path, "<", "{")
	path = strings.ReplaceAll(path, ">", "}")
	
	// Remove trailing slashes
	path = strings.TrimSuffix(path, "/")
	
	// Ensure leading slash
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	
	return path
}
