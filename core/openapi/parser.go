package openapi

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Endpoint represents a parsed API endpoint from architecture.md
type Endpoint struct {
	Method      string
	Path        string
	Description string
	Request     string
	Response    string
	Parameters  []Parameter
}

// Parameter represents a query or path parameter
type Parameter struct {
	Name        string
	In          string // "query", "path", "header"
	Description string
	Required    bool
	Schema      string
}

// ParseArchitecture parses architecture.md and extracts API endpoints
func ParseArchitecture(filePath string) ([]Endpoint, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open architecture file: %w", err)
	}
	defer file.Close()

	var endpoints []Endpoint
	scanner := bufio.NewScanner(file)
	
	// Regular expressions for parsing
	endpointRe := regexp.MustCompile(`^###\s+(GET|POST|PUT|DELETE|PATCH)\s+(.+)$`)
	paramRe := regexp.MustCompile(`^-\s+` + "`" + `(\w+)` + "`" + `\s*(?:\(([^)]+)\))?\s*:?\s*(.*)$`)
	
	var currentEndpoint *Endpoint
	var inCodeBlock bool
	var codeBlockType string
	var codeBlockContent strings.Builder
	
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		
		// Detect code blocks
		if strings.HasPrefix(trimmedLine, "```") {
			if !inCodeBlock {
				inCodeBlock = true
				// Extract language identifier (e.g., "json", "sql")
				codeBlockType = strings.TrimPrefix(trimmedLine, "```")
				codeBlockContent.Reset()
			} else {
				// End of code block
				inCodeBlock = false
				content := codeBlockContent.String()
				
				// Assign code block to current endpoint
				if currentEndpoint != nil {
					if codeBlockType == "json" {
						if strings.Contains(strings.ToLower(content), "request") || 
						   currentEndpoint.Request == "" {
							currentEndpoint.Request = content
						} else {
							currentEndpoint.Response = content
						}
					}
				}
			}
			continue
		}
		
		if inCodeBlock {
			codeBlockContent.WriteString(line + "\n")
			continue
		}
		
		// Check for endpoint definition
		if matches := endpointRe.FindStringSubmatch(line); matches != nil {
			// Save previous endpoint if exists
			if currentEndpoint != nil {
				// Detect path parameters before saving
				detectPathParameters(currentEndpoint)
				endpoints = append(endpoints, *currentEndpoint)
			}
			
			// Create new endpoint
			currentEndpoint = &Endpoint{
				Method: matches[1],
				Path:   strings.TrimSpace(matches[2]),
			}
			continue
		}
		
		// Parse description (first line after endpoint)
		if currentEndpoint != nil && currentEndpoint.Description == "" && trimmedLine != "" && !strings.HasPrefix(trimmedLine, "**") {
			currentEndpoint.Description = trimmedLine
			continue
		}
		
		// Parse query parameters
		if currentEndpoint != nil && strings.HasPrefix(trimmedLine, "**Query Parameters:**") {
			// Next lines are parameters
			continue
		}
		
		// Parse parameter line
		if currentEndpoint != nil && strings.HasPrefix(trimmedLine, "-") {
			if matches := paramRe.FindStringSubmatch(trimmedLine); matches != nil {
				param := Parameter{
					Name:        matches[1],
					In:          "query",
					Description: matches[3],
					Required:    false,
				}
				if matches[2] != "" {
					param.Description = matches[2] + ": " + matches[3]
				}
				currentEndpoint.Parameters = append(currentEndpoint.Parameters, param)
			}
		}
	}
	
	// Save last endpoint
	if currentEndpoint != nil {
		// Detect path parameters before saving
		detectPathParameters(currentEndpoint)
		endpoints = append(endpoints, *currentEndpoint)
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading architecture file: %w", err)
	}
	
	return endpoints, nil
}

// detectPathParameters detects and adds path parameters from the endpoint path
func detectPathParameters(endpoint *Endpoint) {
	if !strings.Contains(endpoint.Path, ":") {
		return
	}
	
	pathParamRe := regexp.MustCompile(`:(\w+)`)
	pathParams := pathParamRe.FindAllStringSubmatch(endpoint.Path, -1)
	for _, match := range pathParams {
		paramName := match[1]
		// Check if not already added
		exists := false
		for _, p := range endpoint.Parameters {
			if p.Name == paramName && p.In == "path" {
				exists = true
				break
			}
		}
		if !exists {
			endpoint.Parameters = append(endpoint.Parameters, Parameter{
				Name:        paramName,
				In:          "path",
				Description: fmt.Sprintf("%s identifier", paramName),
				Required:    true,
				Schema:      "string",
			})
		}
	}
}
