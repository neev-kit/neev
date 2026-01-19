package cucumber

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/neev-kit/neev/core/openapi"
)

// GenerateFeatureFile generates a Cucumber/Gherkin feature file from parsed endpoints
func GenerateFeatureFile(endpoints []openapi.Endpoint, blueprintName string) (string, error) {
	var builder strings.Builder
	
	// Feature header
	featureName := formatFeatureName(blueprintName)
	builder.WriteString(fmt.Sprintf("Feature: %s\n", featureName))
	builder.WriteString(fmt.Sprintf("  As an API consumer\n"))
	builder.WriteString(fmt.Sprintf("  I want to interact with %s endpoints\n", blueprintName))
	builder.WriteString(fmt.Sprintf("  So that I can perform operations on the system\n\n"))
	
	// Generate scenarios for each endpoint
	for _, endpoint := range endpoints {
		scenario := generateScenario(endpoint)
		builder.WriteString(scenario)
		builder.WriteString("\n")
	}
	
	return builder.String(), nil
}

func generateScenario(endpoint openapi.Endpoint) string {
	var builder strings.Builder
	
	// Scenario title
	scenarioName := fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path)
	builder.WriteString(fmt.Sprintf("  Scenario: %s\n", scenarioName))
	
	// Given step
	builder.WriteString("    Given the API is available\n")
	
	// Add authentication if needed (for POST, PUT, DELETE, PATCH)
	if endpoint.Method != "GET" {
		builder.WriteString("    And I am authenticated\n")
	}
	
	// Add path parameters
	if hasPathParams(endpoint.Path) {
		paramNames := extractPathParams(endpoint.Path)
		for _, param := range paramNames {
			builder.WriteString(fmt.Sprintf("    And I have a valid %s\n", param))
		}
	}
	
	// When step
	methodAction := getMethodAction(endpoint.Method)
	path := endpoint.Path
	builder.WriteString(fmt.Sprintf("    When I %s to \"%s\"\n", methodAction, path))
	
	// Add request body for POST, PUT, PATCH
	if endpoint.Method == "POST" || endpoint.Method == "PUT" || endpoint.Method == "PATCH" {
		if endpoint.Request != "" {
			builder.WriteString("    And I send the following JSON payload:\n")
			builder.WriteString("      \"\"\"\n")
			for _, line := range strings.Split(strings.TrimSpace(endpoint.Request), "\n") {
				builder.WriteString(fmt.Sprintf("      %s\n", line))
			}
			builder.WriteString("      \"\"\"\n")
		}
	}
	
	// Add query parameters
	for _, param := range endpoint.Parameters {
		if param.In == "query" {
			builder.WriteString(fmt.Sprintf("    And I include query parameter \"%s\" with value \"<value>\"\n", param.Name))
		}
	}
	
	// Then step
	expectedStatus := getExpectedStatus(endpoint.Method)
	builder.WriteString(fmt.Sprintf("    Then the response status should be %s\n", expectedStatus))
	
	// Add response validation for GET endpoints
	if endpoint.Method == "GET" {
		builder.WriteString("    And the response should contain valid data\n")
	} else if endpoint.Method == "POST" {
		builder.WriteString("    And the response should contain the created resource\n")
	} else if endpoint.Method == "DELETE" {
		builder.WriteString("    And the resource should be deleted\n")
	}
	
	return builder.String()
}

// GenerateStepDefinitions generates skeleton step definition file
func GenerateStepDefinitions(language string) (string, error) {
	switch language {
	case "go":
		return generateGoStepDefinitions(), nil
	case "javascript", "js":
		return generateJavaScriptStepDefinitions(), nil
	case "python":
		return generatePythonStepDefinitions(), nil
	default:
		return "", fmt.Errorf("unsupported language: %s (supported: go, javascript, python)", language)
	}
}

func generateGoStepDefinitions() string {
	return `package steps

import (
	"github.com/cucumber/godog"
)

// APIContext holds the test context
type APIContext struct {
	response     interface{}
	statusCode   int
	requestBody  string
	queryParams  map[string]string
}

func (ctx *APIContext) theAPIIsAvailable() error {
	// TODO: Implement API availability check
	return nil
}

func (ctx *APIContext) iAmAuthenticated() error {
	// TODO: Implement authentication setup
	return nil
}

func (ctx *APIContext) iHaveAValidParameter(param string) error {
	// TODO: Implement parameter setup for param
	return nil
}

func (ctx *APIContext) iSendRequestTo(method, path string) error {
	// TODO: Implement HTTP request for method and path
	return nil
}

func (ctx *APIContext) iSendJSONPayload(payload string) error {
	ctx.requestBody = payload
	// TODO: Attach payload to request
	return nil
}

func (ctx *APIContext) iIncludeQueryParameter(name, value string) error {
	if ctx.queryParams == nil {
		ctx.queryParams = make(map[string]string)
	}
	ctx.queryParams[name] = value
	return nil
}

func (ctx *APIContext) theResponseStatusShouldBe(expectedStatus int) error {
	// TODO: Verify response status matches expectedStatus
	return nil
}

func (ctx *APIContext) theResponseShouldContainValidData() error {
	// TODO: Validate response data
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	apiCtx := &APIContext{}
	
	ctx.Step(` + "`^the API is available$`" + `, apiCtx.theAPIIsAvailable)
	ctx.Step(` + "`^I am authenticated$`" + `, apiCtx.iAmAuthenticated)
	ctx.Step(` + "`^I have a valid (.+)$`" + `, apiCtx.iHaveAValidParameter)
	ctx.Step(` + "`^I (GET|POST|PUT|PATCH|DELETE) to \"([^\"]+)\"$`" + `, apiCtx.iSendRequestTo)
	ctx.Step(` + "`^I send the following JSON payload:$`" + `, apiCtx.iSendJSONPayload)
	ctx.Step(` + "`^I include query parameter \"([^\"]+)\" with value \"([^\"]+)\"$`" + `, apiCtx.iIncludeQueryParameter)
	ctx.Step(` + "`^the response status should be (\\d+)$`" + `, apiCtx.theResponseStatusShouldBe)
	ctx.Step(` + "`^the response should contain valid data$`" + `, apiCtx.theResponseShouldContainValidData)
}
`
}

func generateJavaScriptStepDefinitions() string {
	return `const { Given, When, Then } = require('@cucumber/cucumber');
const assert = require('assert');

// Test context
let response;
let statusCode;
let requestBody;
let queryParams = {};

Given('the API is available', async function () {
  // TODO: Implement API availability check
});

Given('I am authenticated', async function () {
  // TODO: Implement authentication setup
});

Given('I have a valid {word}', async function (param) {
  // TODO: Implement parameter setup for param
});

When('I {word} to {string}', async function (method, path) {
  // TODO: Implement HTTP request for method and path
});

When('I send the following JSON payload:', async function (payload) {
  requestBody = payload;
  // TODO: Attach payload to request
});

When('I include query parameter {string} with value {string}', async function (name, value) {
  queryParams[name] = value;
});

Then('the response status should be {int}', async function (expectedStatus) {
  // TODO: Verify response status matches expectedStatus
  assert.strictEqual(statusCode, expectedStatus);
});

Then('the response should contain valid data', async function () {
  // TODO: Validate response data
});
`
}

func generatePythonStepDefinitions() string {
	return `from behave import given, when, then
import requests

@given('the API is available')
def step_api_available(context):
    # TODO: Implement API availability check
    pass

@given('I am authenticated')
def step_authenticated(context):
    # TODO: Implement authentication setup
    pass

@given('I have a valid {param}')
def step_valid_param(context, param):
    # TODO: Implement parameter setup for param
    pass

@when('I {method} to "{path}"')
def step_send_request(context, method, path):
    # TODO: Implement HTTP request for method and path
    pass

@when('I send the following JSON payload')
def step_json_payload(context):
    context.request_body = context.text
    # TODO: Attach payload to request
    pass

@when('I include query parameter "{name}" with value "{value}"')
def step_query_param(context, name, value):
    if not hasattr(context, 'query_params'):
        context.query_params = {}
    context.query_params[name] = value

@then('the response status should be {status:d}')
def step_response_status(context, status):
    # TODO: Verify response status matches status
    assert context.status_code == status

@then('the response should contain valid data')
def step_valid_data(context):
    # TODO: Validate response data
    pass
`
}

// GenerateCucumber is the main entry point for generating Cucumber tests
func GenerateCucumber(architecturePath, blueprintName, outputPath, language string) error {
	// Parse the architecture file
	endpoints, err := openapi.ParseArchitecture(architecturePath)
	if err != nil {
		return fmt.Errorf("failed to parse architecture: %w", err)
	}

	if len(endpoints) == 0 {
		return fmt.Errorf("no API endpoints found in architecture file")
	}

	// Generate feature file
	featureContent, err := GenerateFeatureFile(endpoints, blueprintName)
	if err != nil {
		return fmt.Errorf("failed to generate feature file: %w", err)
	}

	// Write feature file
	featurePath := filepath.Join(outputPath, "api.feature")
	if err := os.WriteFile(featurePath, []byte(featureContent), 0644); err != nil {
		return fmt.Errorf("failed to write feature file: %w", err)
	}

	// Generate step definitions if language is specified
	if language != "" {
		stepDefs, err := GenerateStepDefinitions(language)
		if err != nil {
			return fmt.Errorf("failed to generate step definitions: %w", err)
		}

		// Determine file extension and name
		var stepsFileName string
		switch language {
		case "go":
			stepsFileName = "steps.go"
		case "javascript", "js":
			stepsFileName = "steps.js"
		case "python":
			stepsFileName = "steps.py"
		}

		stepsPath := filepath.Join(outputPath, stepsFileName)
		if err := os.WriteFile(stepsPath, []byte(stepDefs), 0644); err != nil {
			return fmt.Errorf("failed to write step definitions: %w", err)
		}
	}

	return nil
}

// Helper functions

func hasPathParams(path string) bool {
	return strings.Contains(path, ":") || strings.Contains(path, "{")
}

func extractPathParams(path string) []string {
	var params []string
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			params = append(params, strings.TrimPrefix(part, ":"))
		} else if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			params = append(params, strings.Trim(part, "{}"))
		}
	}
	return params
}

func getMethodAction(method string) string {
	switch method {
	case "GET":
		return "GET"
	case "POST":
		return "POST"
	case "PUT":
		return "PUT"
	case "DELETE":
		return "DELETE"
	case "PATCH":
		return "PATCH"
	default:
		return method
	}
}

func getExpectedStatus(method string) string {
	switch method {
	case "POST":
		return "201"
	case "DELETE":
		return "204"
	default:
		return "200"
	}
}

func formatFeatureName(blueprintName string) string {
	// Convert kebab-case to Title Case
	words := strings.Split(blueprintName, "-")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ") + " API"
}
