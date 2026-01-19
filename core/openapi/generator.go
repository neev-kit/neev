package openapi

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// OpenAPISpec represents an OpenAPI 3.1 specification
type OpenAPISpec struct {
	OpenAPI string                 `yaml:"openapi"`
	Info    Info                   `yaml:"info"`
	Paths   map[string]PathItem    `yaml:"paths"`
	Components *Components          `yaml:"components,omitempty"`
}

// Info contains API metadata
type Info struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description,omitempty"`
	Version     string `yaml:"version"`
}

// PathItem describes operations available on a single path
type PathItem struct {
	Get    *Operation `yaml:"get,omitempty"`
	Post   *Operation `yaml:"post,omitempty"`
	Put    *Operation `yaml:"put,omitempty"`
	Delete *Operation `yaml:"delete,omitempty"`
	Patch  *Operation `yaml:"patch,omitempty"`
}

// Operation describes a single API operation
type Operation struct {
	Summary     string                `yaml:"summary,omitempty"`
	Description string                `yaml:"description,omitempty"`
	Parameters  []ParameterSpec       `yaml:"parameters,omitempty"`
	RequestBody *RequestBody          `yaml:"requestBody,omitempty"`
	Responses   map[string]Response   `yaml:"responses"`
}

// ParameterSpec describes a single operation parameter
type ParameterSpec struct {
	Name        string      `yaml:"name"`
	In          string      `yaml:"in"`
	Description string      `yaml:"description,omitempty"`
	Required    bool        `yaml:"required,omitempty"`
	Schema      SchemaRef   `yaml:"schema"`
}

// RequestBody describes a request body
type RequestBody struct {
	Description string                `yaml:"description,omitempty"`
	Required    bool                  `yaml:"required,omitempty"`
	Content     map[string]MediaType  `yaml:"content"`
}

// Response describes a single response
type Response struct {
	Description string                `yaml:"description"`
	Content     map[string]MediaType  `yaml:"content,omitempty"`
}

// MediaType provides schema and examples for media type
type MediaType struct {
	Schema   SchemaRef              `yaml:"schema,omitempty"`
	Example  interface{}            `yaml:"example,omitempty"`
}

// SchemaRef references a schema
type SchemaRef struct {
	Type       string                 `yaml:"type,omitempty"`
	Format     string                 `yaml:"format,omitempty"`
	Properties map[string]SchemaRef   `yaml:"properties,omitempty"`
	Items      *SchemaRef             `yaml:"items,omitempty"`
	Ref        string                 `yaml:"$ref,omitempty"`
}

// Components holds reusable objects
type Components struct {
	Schemas map[string]SchemaRef `yaml:"schemas,omitempty"`
}

// GenerateOpenAPISpec generates an OpenAPI 3.1 specification from parsed endpoints
func GenerateOpenAPISpec(endpoints []Endpoint, blueprintName string) (*OpenAPISpec, error) {
	spec := &OpenAPISpec{
		OpenAPI: "3.1.0",
		Info: Info{
			Title:       formatTitle(blueprintName),
			Description: fmt.Sprintf("API specification for %s", blueprintName),
			Version:     "1.0.0",
		},
		Paths: make(map[string]PathItem),
	}

	for _, endpoint := range endpoints {
		// Convert :param to {param} format for OpenAPI
		path := convertPathParams(endpoint.Path)
		
		// Get or create path item
		pathItem, exists := spec.Paths[path]
		if !exists {
			pathItem = PathItem{}
		}

		// Create operation
		operation := &Operation{
			Summary:     endpoint.Description,
			Description: endpoint.Description,
			Responses:   make(map[string]Response),
		}

		// Add parameters
		for _, param := range endpoint.Parameters {
			operation.Parameters = append(operation.Parameters, ParameterSpec{
				Name:        param.Name,
				In:          param.In,
				Description: param.Description,
				Required:    param.Required,
				Schema: SchemaRef{
					Type: inferType(param.Schema),
				},
			})
		}

		// Add request body for POST, PUT, PATCH
		if endpoint.Method == "POST" || endpoint.Method == "PUT" || endpoint.Method == "PATCH" {
			if endpoint.Request != "" {
				operation.RequestBody = &RequestBody{
					Description: "Request body",
					Required:    true,
					Content: map[string]MediaType{
						"application/json": {
							Schema: SchemaRef{
								Type: "object",
							},
						},
					},
				}
			}
		}

		// Add default responses
		operation.Responses["200"] = Response{
			Description: "Successful response",
			Content: map[string]MediaType{
				"application/json": {
					Schema: SchemaRef{
						Type: "object",
					},
				},
			},
		}

		if endpoint.Method == "POST" {
			operation.Responses["201"] = Response{
				Description: "Resource created",
			}
		}

		operation.Responses["400"] = Response{
			Description: "Bad request",
		}

		operation.Responses["500"] = Response{
			Description: "Internal server error",
		}

		// Assign operation to the correct method
		switch endpoint.Method {
		case "GET":
			pathItem.Get = operation
		case "POST":
			pathItem.Post = operation
		case "PUT":
			pathItem.Put = operation
		case "DELETE":
			pathItem.Delete = operation
		case "PATCH":
			pathItem.Patch = operation
		}

		spec.Paths[path] = pathItem
	}

	return spec, nil
}

// GenerateYAML generates YAML output from OpenAPI spec
func GenerateYAML(spec *OpenAPISpec) ([]byte, error) {
	data, err := yaml.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OpenAPI spec: %w", err)
	}
	return data, nil
}

// GenerateOpenAPI is the main entry point for generating OpenAPI spec from architecture file
func GenerateOpenAPI(architecturePath, blueprintName string) ([]byte, error) {
	// Parse the architecture file
	endpoints, err := ParseArchitecture(architecturePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse architecture: %w", err)
	}

	if len(endpoints) == 0 {
		return nil, fmt.Errorf("no API endpoints found in architecture file")
	}

	// Generate OpenAPI spec
	spec, err := GenerateOpenAPISpec(endpoints, blueprintName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate OpenAPI spec: %w", err)
	}

	// Convert to YAML
	yamlData, err := GenerateYAML(spec)
	if err != nil {
		return nil, err
	}

	return yamlData, nil
}

// Helper functions

func convertPathParams(path string) string {
	// Convert :param to {param} format
	re := regexp.MustCompile(`:(\w+)`)
	return re.ReplaceAllString(path, "{$1}")
}

func inferType(schema string) string {
	if schema == "" {
		return "string"
	}
	return schema
}

func formatTitle(name string) string {
	// Convert kebab-case to Title Case
	words := strings.Split(name, "-")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}
