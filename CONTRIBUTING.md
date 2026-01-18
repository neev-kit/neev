# Contributing to Neev

Thank you for your interest in contributing to Neev! This guide provides the essential information for working on the project.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Development Workflow](#development-workflow)
- [Code Standards](#code-standards)
- [Testing](#testing)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Reporting Issues](#reporting-issues)

## Getting Started

1. **Fork the Repository** - Click "Fork" on the GitHub repository page
2. **Clone Your Fork** - `git clone https://github.com/YOUR_USERNAME/neev.git`
3. **Add Upstream Remote** - `git remote add upstream https://github.com/neev-kit/neev.git`
4. **Create a Branch** - `git checkout -b feature/your-feature-name`

## Development Setup

### Prerequisites

- Go 1.23 or higher
- Git
- Make (optional but recommended)

### Local Environment

```bash
# Clone and navigate
git clone https://github.com/neev-kit/neev.git
cd neev

# Install dependencies
go mod download

# Build the CLI
go build -o neev ./cli

# Verify installation
./neev --version

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Project Structure

```
neev/
├── cli/                    # CLI application (Cobra)
│   ├── main.go            # Entry point
│   ├── go.mod             # Module definition
│   └── cmd/               # Commands
│       ├── root.go        # Root command with logger init
│       ├── init.go        # Project initialization
│       ├── draft.go       # Blueprint drafting
│       ├── bridge.go      # Context bridging
│       ├── lay.go         # Lay foundation (internal)
│       └── inspect.go     # Inspect projects (internal)
│
├── core/                  # Core business logic
│   ├── go.mod             # Module definition
│   ├── config/            # Configuration management
│   │   └── loader.go      # YAML config loading
│   ├── errors/            # Custom error handling
│   │   └── errors.go      # NeevError types
│   ├── logger/            # Structured logging
│   │   └── logger.go      # slog integration
│   ├── foundation/        # Project foundation
│   │   ├── init.go        # Init logic
│   │   ├── inspect.go     # Project inspection
│   │   └── paths.go       # Constants & paths
│   ├── blueprint/         # Blueprint management
│   │   ├── draft.go       # Draft creation
│   │   └── lay.go         # Blueprint laying
│   └── bridge/            # Context aggregation
│       └── context.go     # Context building
│
├── .github/workflows/     # CI/CD pipelines
│   ├── tests.yml          # Run tests on push/PR
│   └── release.yml        # Build & release on tag
│
└── go.work                # Workspace file for multi-module
```

## Development Workflow

### 1. **Create a Feature Branch**

```bash
# Update main branch first
git fetch upstream
git checkout main
git merge upstream/main

# Create feature branch
git checkout -b feature/add-new-command

# Or for fixes
git checkout -b fix/issue-123
```

### 2. **Make Your Changes**

- Write code following the code standards (see below)
- Add or update tests for new functionality
- Update documentation if behavior changes

### 3. **Run Tests Locally**

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestInspect_FindsMissingCodeDirectories ./core/foundation

# Run tests with coverage
go test -cover ./... | grep -E "ok|FAIL"

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 4. **Commit Your Work**

Follow the [Commit Guidelines](#commit-guidelines) below.

### 5. **Push and Create PR**

```bash
git push origin feature/add-new-command
# Then create PR via GitHub UI
```

## Code Standards

### Go Style Guidelines

- **Formatting**: Use `gofmt` (or `go fmt ./...`)
- **Linting**: Follow Go best practices
- **Naming**:
  - Functions: `CamelCase` (exported), `camelCase` (unexported)
  - Constants: `UPPER_CASE` for constants
  - Variables: Descriptive names, avoid single letters except loop counters
- **Comments**:
  - Exported functions/types must have doc comments
  - Start with function/type name: `// Init initializes...`
  - Add meaningful comments for complex logic

### Example Code Structure

```go
package foundation

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/neev-kit/neev/core/errors"
)

// Init initializes the project foundation with configuration.
func Init(path string) (*Project, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, errors.ErrFoundationMissing(path)
	}
	
	// Implementation...
	return &Project{}, nil
}
```

### Error Handling

Always use custom error types from `core/errors`:

```go
// ❌ Don't do this
return fmt.Errorf("blueprint not found")

// ✅ Do this
return errors.ErrBlueprintNotFound(name)
```

### Logging

Use the structured logger for all output:

```go
import "github.com/neev-kit/neev/core/logger"

// Initialize in main command
logger.Init()

// Use throughout
logger.Info("Starting process", "step", 1)
logger.Debug("Configuration loaded", "path", configPath)
logger.Warn("Deprecated flag used", "flag", "--old-flag")
logger.Error("Failed to initialize", "error", err)
```

### Configuration

Load configuration with validation:

```go
import "github.com/neev-kit/neev/core/config"

// Load with defaults
cfg := config.LoadConfig(configPath)

// Validate before use
if err := cfg.Validate(); err != nil {
	return errors.ErrInvalidConfig(err.Error())
}

// Access properties
dirs := cfg.GetIgnoreDirs()
```

## Testing

### Testing Guidelines

1. **Unit Tests**: Create `*_test.go` files alongside source files
2. **Test Names**: Describe what is being tested: `TestFunctionName_Scenario_Expected`
3. **Use `t.TempDir()`**: For file system operations
4. **Mock External Calls**: Don't make real API/network calls
5. **Table-Driven Tests**: Use for multiple scenarios

### Example Test

```go
func TestInit_CreatesFoundation(t *testing.T) {
	tmpDir := t.TempDir()
	
	project, err := Init(tmpDir)
	
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	
	if project == nil {
		t.Error("Expected project, got nil")
	}
}

func TestInspect_WithCustomConfig(t *testing.T) {
	tests := []struct {
		name        string
		configPath  string
		expectError bool
	}{
		{"valid config", "config.yaml", false},
		{"missing config", "", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test logic
		})
	}
}
```

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./core/foundation

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

## Commit Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/) for clear commit history:

### Format

```
type(scope): description

Optional longer explanation.
Optional footer with references: Fixes #123
```

### Types

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **test**: Test additions/updates
- **refactor**: Code refactoring without behavior changes
- **perf**: Performance improvements
- **chore**: Build, deps, or tooling changes

### Examples

```bash
git commit -m "feat(blueprint): add draft command for creating blueprints"

git commit -m "fix(bridge): handle empty project context gracefully

Fixes issue with context bridge returning nil when no blueprints exist.

Fixes #42"

git commit -m "docs: update README with installation instructions"

git commit -m "test(foundation): add 5 new test cases for Init function"
```

## Pull Request Process

### Before Creating a PR

1. **Ensure tests pass**: `go test ./...`
2. **Format code**: `go fmt ./...`
3. **Update documentation** if behavior changed
4. **Rebase on main**: `git rebase upstream/main`

### PR Template

```markdown
## Description
Brief description of changes.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update

## Related Issues
Fixes #123

## Testing
How to test these changes:
```

### PR Guidelines

- Keep PRs focused on a single feature/fix
- Include tests for new code
- Update README/documentation if needed
- Use descriptive PR titles
- Link related issues

### Review Process

1. **Automated Checks**: Tests and linting must pass
2. **Code Review**: Maintainers review changes
3. **Approval**: One approval required
4. **Merge**: Squash and merge to main

## Reporting Issues

### Bug Reports

Include:
- Go version (`go version`)
- OS and architecture
- Steps to reproduce
- Expected vs actual behavior
- Error messages/logs

### Feature Requests

Include:
- Use case and motivation
- Proposed solution (if any)
- Alternative approaches considered

### Template

```markdown
## Issue Type
[ ] Bug [ ] Feature [ ] Question

## Description
Clear description of the issue.

## Environment
- Go Version: x.xx
- OS: macOS/Linux/Windows
- Neev Version: v1.0.0

## Steps to Reproduce
1. 
2.

## Expected Behavior
...

## Actual Behavior
...
```

## Questions?

- **Discussions**: Use GitHub Discussions for questions
- **Issues**: Report bugs with detailed information
- **Discord**: Join our community Discord (if available)

Thank you for contributing to Neev! 🙏
