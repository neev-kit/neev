# Development Guide for Neev

Complete guide for setting up your development environment and building/testing Neev locally.

## Table of Contents

- [System Requirements](#system-requirements)
- [Installation & Setup](#installation--setup)
- [Building](#building)
- [Testing](#testing)
- [Code Organization](#code-organization)
- [Adding Features](#adding-features)
- [Debugging](#debugging)
- [Performance Profiling](#performance-profiling)

## System Requirements

### Minimum Requirements

| Component | Version | Notes |
|-----------|---------|-------|
| Go | 1.23+ | Required for building |
| Git | 2.30+ | For version control |
| Make | 4.0+ | Optional, for convenience |
| OS | macOS/Linux/Windows | Multi-platform support |

### Optional Tools

- **golangci-lint**: Code linting (`brew install golangci-lint`)
- **delve**: Debugging (`go install github.com/go-delve/delve/cmd/dlv@latest`)
- **goprofiler**: Performance analysis (built-in with Go)

## Installation & Setup

### Step 1: Clone Repository

```bash
git clone https://github.com/neev-kit/neev.git
cd neev
```

### Step 2: Initialize Go Modules

```bash
# Download dependencies
go mod download

# Verify dependencies
go mod verify
```

### Step 3: Verify Setup

```bash
# Check Go version
go version

# List dependencies
go mod graph

# Build binary
go build -o neev ./cli

# Test binary
./neev --version
```

### Step 4: Optional - Setup IDE

#### VSCode

```bash
# Install Go extension
# Recommended: golang.go by Google

# Create .vscode/settings.json
cat > .vscode/settings.json << 'EOF'
{
  "go.lintOnSave": "package",
  "go.lintTool": "golangci-lint",
  "go.coverageDecorator": "gutter",
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": "explicit"
    }
  }
}
EOF
```

#### GoLand / IntelliJ IDEA

- Built-in Go support
- Project structure auto-detected
- Run configurations pre-configured

## Building

### Development Build

```bash
# Build with debug symbols (slower)
go build -o neev ./cli

# Run directly (no build needed for quick iteration)
go run ./cli/main.go --help
```

### Release Build

```bash
# Build optimized binary
go build -ldflags="-s -w" -o neev ./cli

# Build for specific OS/ARCH
GOOS=linux GOARCH=amd64 go build -o neev-linux-amd64 ./cli
GOOS=darwin GOARCH=arm64 go build -o neev-darwin-arm64 ./cli
GOOS=windows GOARCH=amd64 go build -o neev-windows-amd64.exe ./cli
```

### Build Flags

Add version information to builds:

```bash
VERSION=$(git describe --tags --always)
go build -ldflags="-X main.Version=${VERSION}" -o neev ./cli
```

## Testing

### Run All Tests

```bash
# Basic test run
go test ./...

# Verbose output
go test -v ./...

# Show test coverage
go test -cover ./...
```

### Run Specific Tests

```bash
# Single package
go test ./core/foundation

# Single test
go test -run TestInit ./core/foundation

# Tests matching pattern
go test -run TestInspect ./core/foundation

# Single file (multiple packages)
go test -run TestDraft ./core/blueprint ./core/foundation
```

### Coverage Analysis

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Coverage for specific package
go test -cover ./core/foundation
```

### Test Output Formats

```bash
# JSON output for CI systems
go test -json ./... > test-results.json

# Verbose with coverage
go test -v -coverprofile=coverage.out ./...

# Show which tests ran
go test -v ./... 2>&1 | grep -E "RUN|PASS|FAIL"
```

### Benchmark Tests

```bash
# Run benchmarks (if defined)
go test -bench=. -benchmem ./...

# Run specific benchmark
go test -bench=BenchmarkInit -benchmem ./core/foundation

# Profile benchmark
go test -bench=BenchmarkBridge -cpuprofile=cpu.prof ./core/bridge
```

## Code Organization

### Module Structure

```
neev/
├── cli/                    # CLI application
│   └── cmd/               # Command definitions (Cobra)
│       ├── root.go        # Root command
│       ├── init.go        # Init command
│       ├── draft.go       # Draft command
│       ├── bridge.go      # Bridge command
│       └── *_test.go      # Command tests
│
└── core/                  # Core business logic
    ├── config/            # Configuration
    │   ├── loader.go      # YAML loading
    │   └── loader_test.go
    ├── errors/            # Error handling
    │   ├── errors.go      # Custom error types
    │   └── errors_test.go
    ├── logger/            # Logging
    │   ├── logger.go      # Structured logging
    │   └── logger_test.go
    ├── foundation/        # Foundation/base operations
    │   ├── init.go        # Initialization
    │   ├── inspect.go     # Project inspection
    │   ├── paths.go       # Constants
    │   └── *_test.go
    ├── blueprint/         # Blueprint operations
    │   ├── draft.go       # Draft creation
    │   ├── lay.go         # Blueprint laying
    │   └── *_test.go
    └── bridge/            # Context bridging
        ├── context.go     # Context aggregation
        └── context_test.go
```

### Dependency Flow

```
cli/cmd/
  ↓ imports
core/(foundation|blueprint|bridge|config|logger|errors)
  ↓
core/errors ← Used by all packages
core/logger ← Used by all packages
core/config ← Used where needed
```

## Adding Features

### Adding a New Command

1. **Create command file** in `cli/cmd/`:

```go
// File: cli/cmd/mycommand.go
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/neev-kit/neev/core/logger"
)

var myCmd = &cobra.Command{
	Use:   "mycommand",
	Short: "Description of my command",
	Long:  "Longer description with examples...",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Running mycommand")
		// Implementation
	},
}

func init() {
	rootCmd.AddCommand(myCmd)
	myCmd.Flags().StringVar(&flagVar, "flag", "", "Help text")
}
```

2. **Register in root.go**:

```go
// Already done if you call rootCmd.AddCommand() in init()
```

3. **Add tests** in `cli/cmd/mycommand_test.go`:

```go
func TestMyCommand_Success(t *testing.T) {
	cmd := &cobra.Command{}
	// Test logic
}
```

### Adding Core Functionality

1. **Create package** in `core/`:

```
core/myfeature/
├── implementation.go      # Main logic
├── implementation_test.go # Tests
└── types.go              # Type definitions (if needed)
```

2. **Follow patterns**:

```go
package myfeature

import (
	"github.com/neev-kit/neev/core/errors"
	"github.com/neev-kit/neev/core/logger"
)

// DoSomething performs an operation.
func DoSomething(input string) (string, error) {
	if input == "" {
		return "", errors.ErrInvalidConfig("input required")
	}
	
	logger.Debug("Processing input", "input", input)
	
	// Implementation
	result := "output"
	logger.Info("Completed operation", "result", result)
	
	return result, nil
}
```

3. **Export in appropriate module**:
   - Add function to public API
   - Document with comments
   - Add tests

## Debugging

### Using Printf-style Debugging

```go
logger.Debug("Value at step X", "value", myVar, "type", fmt.Sprintf("%T", myVar))
```

### Using Delve Debugger

```bash
# Start debug session
dlv debug ./cli

# Set breakpoint
(dlv) break main.main

# Continue execution
(dlv) continue

# Print variables
(dlv) print myVariable

# Step through code
(dlv) next
(dlv) step
```

### VS Code Debugging

Create `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Connect to Delve",
      "type": "go",
      "mode": "local",
      "remotePath": "",
      "port": 2345,
      "host": "127.0.0.1",
      "showLog": true,
      "trace": "verbose"
    }
  ]
}
```

### Common Debugging Patterns

```go
// Log function entry/exit
func MyFunction(input string) error {
	logger.Debug("Entering MyFunction", "input", input)
	defer logger.Debug("Exiting MyFunction")
	
	// Implementation
	return nil
}

// Log conditions
if condition {
	logger.Debug("Condition met", "expected", true, "actual", condition)
}

// Log state changes
before := myVar
myVar = newValue
logger.Debug("State changed", "before", before, "after", myVar)
```

## Performance Profiling

### CPU Profiling

```bash
# Add to benchmark test
go test -bench=. -cpuprofile=cpu.prof ./core/foundation

# Analyze
go tool pprof cpu.prof
(pprof) top
(pprof) list function_name
```

### Memory Profiling

```bash
# Generate memory profile
go test -memprofile=mem.prof ./core/foundation

# Analyze
go tool pprof mem.prof
(pprof) top
(pprof) alloc_space  # Total allocations
(pprof) inuse_space  # Current memory use
```

### Profiling Commands

```bash
# View top 20 functions by CPU time
go tool pprof -top cpu.prof | head -25

# Generate graph (requires graphviz)
go tool pprof -graph cpu.prof > cpu.graph

# Interactive analysis
go tool pprof cpu.prof
# Then use: top, list, web, etc.
```

### Example Benchmark

```go
func BenchmarkBridge(b *testing.B) {
	ctx := setupTestContext()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_ = ctx.Aggregate()
	}
}
```

## Continuous Integration

### Local CI Simulation

```bash
# Run what GitHub Actions runs
go test -v ./...
go build -o neev ./cli

# Check test coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### GitHub Actions

- **tests.yml**: Runs on push/PR
- **release.yml**: Runs on version tags

See `.github/workflows/` for details.

## Troubleshooting

### Common Issues

| Issue | Solution |
|-------|----------|
| `go: cannot find main module` | Run `go mod download` |
| `imports cycle not allowed` | Check for circular dependencies |
| Tests fail locally but pass in CI | Check file permissions, paths, temp directories |
| Build fails with version error | Ensure Go 1.23+ is installed |

### Getting Help

- Check existing issues on GitHub
- Review test cases for examples
- Read code comments and documentation
- Ask in discussions or issues
