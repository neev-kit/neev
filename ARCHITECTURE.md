# Neev Architecture Guide

Complete overview of Neev's design, components, and how they interact.

## Table of Contents

- [High-Level Overview](#high-level-overview)
- [Core Principles](#core-principles)
- [System Architecture](#system-architecture)
- [Package Organization](#package-organization)
- [Data Flow](#data-flow)
- [Error Handling](#error-handling)
- [Configuration System](#configuration-system)
- [Design Patterns](#design-patterns)
- [Extensibility](#extensibility)

## High-Level Overview

Neev is a CLI framework that bridges AI agents with project context. It operates in three main phases:

```
┌─────────────────────────────────────────────────┐
│           Neev CLI Application                   │
│  ┌───────────────────────────────────────────┐  │
│  │  Commands (init, draft, bridge, inspect)  │  │
│  └─────────────┬──────────────────────────────┘  │
│                │                                 │
│  ┌─────────────▼──────────────────────────────┐  │
│  │  Core Business Logic                       │  │
│  │  (foundation, blueprint, bridge modules)   │  │
│  └─────────────┬──────────────────────────────┘  │
│                │                                 │
│  ┌─────────────▼──────────────────────────────┐  │
│  │  Cross-cutting Concerns                    │  │
│  │  (logging, errors, config)                 │  │
│  └──────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
         │                              │
         ▼                              ▼
    File System              AI Agents / External Tools
    (.neev directory)        (context aggregation)
```

## Core Principles

### 1. **Separation of Concerns**

- **CLI Layer** (`cli/cmd/`): User interface and command routing
- **Business Logic** (`core/*/`): Domain functionality
- **Cross-cutting** (`core/logger`, `core/errors`): Shared utilities

### 2. **Error Handling**

- Custom error types with solution hints
- Errors flow up through call stack
- Root command handles NeevError display

### 3. **Configuration Management**

- YAML-based configuration in `.neev/neev.yaml`
- Type-safe configuration loader
- Defaults provided for all fields

### 4. **Structured Logging**

- Consistent logging via `slog` (Go standard)
- Human-readable + JSON output modes
- Emoji indicators for visual clarity

### 5. **File System as Database**

- No external dependencies (database, API)
- Local `.neev/` directory is single source of truth
- Easy to version control and audit

## System Architecture

### Component Diagram

```
┌────────────────────────────────────────────────────────┐
│                    CLI Commands                         │
│  ┌──────────┬──────────┬──────────┬──────────────────┐ │
│  │  init    │  draft   │  bridge  │  inspect/lay    │ │
│  └────┬─────┴────┬─────┴────┬─────┴────┬────────────┘ │
└───────┼──────────┼──────────┼─────────────────────────┘
        │          │          │          │
        ▼          ▼          ▼          ▼
┌─────────────────────────────────────────────────────────┐
│                Core Business Logic                      │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Foundation                                      │  │
│  │  - Initialize projects                          │  │
│  │  - Inspect project structure                    │  │
│  │  - Manage .neev directory                       │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Blueprint                                       │  │
│  │  - Create drafts                                │  │
│  │  - Manage blueprint templates                   │  │
│  │  - Generate blueprint structure                 │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Bridge                                          │  │
│  │  - Aggregate project context                    │  │
│  │  - Filter by keywords                           │  │
│  │  - Format for AI consumption                    │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
        │          │          │          │
        ▼          ▼          ▼          ▼
┌─────────────────────────────────────────────────────────┐
│           Cross-Cutting Concerns                        │
│  ┌────────────────┬────────────────┬─────────────────┐ │
│  │  Logger        │  Errors        │  Config         │ │
│  │  (structured   │  (custom types │  (YAML loader   │ │
│  │   logging)     │   & hints)      │   & validator)  │ │
│  └────────────────┴────────────────┴─────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

## Package Organization

### `cli/` - Command Line Interface

**Purpose**: User-facing CLI application

**Key Files**:
- `main.go` - Entry point
- `cmd/root.go` - Root command with logger initialization
- `cmd/init.go` - Initialize foundation
- `cmd/draft.go` - Create blueprints
- `cmd/bridge.go` - Aggregate context
- `cmd/lay.go` - Lay foundation (internal)
- `cmd/inspect.go` - Inspect projects (internal)

**Responsibilities**:
- Parse command-line arguments
- Call appropriate core functions
- Handle NeevError types and display hints
- Initialize logger on startup

**Dependencies**: Cobra, core packages

### `core/foundation/` - Project Foundation

**Purpose**: Initialize and manage project foundation

**Key Files**:
- `init.go` - Initialize `.neev/` directory structure
- `inspect.go` - Scan project structure and find missing parts
- `paths.go` - Constants and helper functions

**Responsibilities**:
- Create `.neev/` directory
- Generate `neev.yaml` template
- Validate existing foundations
- Detect project drift

**Key Functions**:
```go
func Init(path string) error           // Initialize foundation
func Inspect(path string) (Report, error) // Scan project
func GetFoundationPath() string         // Get .neev path
```

### `core/blueprint/` - Blueprint Management

**Purpose**: Create and manage project blueprints

**Key Files**:
- `draft.go` - Create blueprint drafts
- `lay.go` - Lay blueprint structure

**Responsibilities**:
- Generate blueprint templates
- Create blueprint directories
- Validate blueprint names
- Generate blueprint metadata

**Key Functions**:
```go
func Draft(title string) (string, error)   // Create blueprint
func Lay(name string) (string, error)      // Lay blueprint
```

### `core/bridge/` - Context Aggregation

**Purpose**: Aggregate project context for AI agents

**Key Files**:
- `context.go` - Build aggregated context

**Responsibilities**:
- Read foundation documents
- Read blueprint documents
- Aggregate into single context
- Filter by keywords
- Format for AI consumption

**Key Functions**:
```go
func BuildContext(foundationPath string) (string, error)
func FilterContext(context, keyword string) (string, error)
```

### `core/config/` - Configuration Management

**Purpose**: Load and validate project configuration

**Key Files**:
- `loader.go` - YAML configuration loader

**Responsibilities**:
- Parse `neev.yaml`
- Provide defaults
- Validate configuration
- Expose configuration to other packages

**Key Functions**:
```go
func LoadConfig(path string) (*Config, error)
func (c *Config) Validate() error
func (c *Config) GetIgnoreDirs() []string
```

### `core/errors/` - Error Handling

**Purpose**: Custom error types with helpful hints

**Key Files**:
- `errors.go` - Error type definitions

**Responsibilities**:
- Define error types
- Provide solution hints
- Enable rich error information

**Error Types**:
```go
type NeevError struct {
	ErrorType    string // e.g., "blueprint_not_found"
	Message      string
	SolutionHint string
	Details      map[string]string
}

// Constructors
ErrBlueprintNotFound(name string)
ErrFoundationMissing(path string)
ErrInvalidConfig(reason string)
ErrTypeIO(message string)
ErrTypeValidation(message string)
```

### `core/logger/` - Structured Logging

**Purpose**: Consistent, structured logging throughout app

**Key Files**:
- `logger.go` - Structured logging implementation

**Responsibilities**:
- Initialize logger on startup
- Provide logging functions
- Support human + JSON output
- Add emoji indicators

**Key Functions**:
```go
func Init()                           // Initialize logger
func Info(msg string, args ...any)    // Log info
func Debug(msg string, args ...any)   // Log debug
func Warn(msg string, args ...any)    // Log warning
func Error(msg string, args ...any)   // Log error
```

**Output Modes**:
- Human: Colored output with emoji (default)
- JSON: Machine-readable (set `NEEV_LOG=json`)

## Data Flow

### Command: `neev init`

```
User Input: neev init
    ↓
cmd/init.go Execute()
    ↓
foundation.Init(currentPath)
    ├─ Check if .neev exists
    ├─ Create directories
    │  ├─ .neev/
    │  ├─ .neev/blueprints/
    │  └─ .neev/foundation/
    ├─ Generate neev.yaml
    └─ Return success/error
    ↓
Logger outputs result
    ↓
Return to user
```

### Command: `neev bridge`

```
User Input: neev bridge [--focus keyword]
    ↓
cmd/bridge.go Execute()
    ↓
Load Config
    ├─ config.LoadConfig()
    └─ Get foundation path
    ↓
Build Context
    ├─ bridge.BuildContext(foundationPath)
    ├─ Read foundation/ markdown files
    ├─ Read blueprints/ markdown files
    ├─ Aggregate into single string
    └─ Filter by focus keyword (if provided)
    ↓
Format for output
    ├─ Add markdown headers
    └─ Structure sections
    ↓
Return to user (stdout or file)
```

### Error Flow

```
Core function detects error
    ↓
Returns NeevError
    ├─ Type: error category
    ├─ Message: what happened
    ├─ SolutionHint: how to fix
    └─ Details: additional info
    ↓
Propagates up call stack
    ↓
Root command catches error
    ├─ Type assert to *NeevError
    ├─ Log error
    ├─ Display solution hint
    └─ Exit with status 1
    ↓
User sees helpful message
```

## Error Handling

### Error Types

```go
// Blueprint not found
ErrBlueprintNotFound(name string)
// Hint: "Run 'neev list' to see available blueprints"

// Foundation missing
ErrFoundationMissing(path string)
// Hint: "Run 'neev init' to set up your project"

// Invalid configuration
ErrInvalidConfig(reason string)
// Hint: "Check neev.yaml configuration"

// IO errors
ErrTypeIO(message string)
// Hint: "Check file permissions and paths"

// Validation errors
ErrTypeValidation(message string)
// Hint: "Verify input according to constraints"
```

### Error Display

```
❌ Error: Blueprint not found
   Blueprint 'my-blueprint' does not exist

💡 Solution: Run 'neev list' to see available blueprints
```

## Configuration System

### Configuration File

Location: `.neev/neev.yaml`

```yaml
project_name: My Project
foundation_path: .neev

ignore_dirs:
  - node_modules
  - .git
  - __pycache__
  - vendor
```

### Configuration Loading

```go
cfg := config.LoadConfig(".neev/neev.yaml")

// Validate
if err := cfg.Validate(); err != nil {
    return errors.ErrInvalidConfig(err.Error())
}

// Use
dirs := cfg.GetIgnoreDirs()
name := cfg.ProjectName
```

## Design Patterns

### 1. **Dependency Injection**

Functions accept parameters rather than using globals:

```go
// ✅ Good
func Inspect(foundationPath, configPath string) error {
    cfg := config.LoadConfig(configPath)
    // ...
}

// ❌ Avoid
var GlobalConfig *Config
func Inspect() error {
    // ...
}
```

### 2. **Error as Values**

Use typed errors for control flow:

```go
result, err := bridge.BuildContext(path)
if err != nil {
    if neevErr, ok := err.(*errors.NeevError); ok {
        logger.Error(neevErr.Message, "hint", neevErr.SolutionHint)
    }
    return err
}
```

### 3. **Factory Functions**

Create instances with validation:

```go
func Init(path string) error {
    if err := validatePath(path); err != nil {
        return ErrInvalidConfig("path invalid")
    }
    // Create structure...
}
```

### 4. **Single Responsibility**

Each package does one thing well:

- `foundation/` - Initialize & inspect
- `blueprint/` - Create blueprints
- `bridge/` - Aggregate context
- `config/` - Load configuration
- `errors/` - Define error types
- `logger/` - Structured logging

## Extensibility

### Adding New Commands

1. Create new file in `cli/cmd/newcmd.go`
2. Define `cobra.Command`
3. Register with root command in `init()`
4. Call core functions from new package

### Adding New Error Types

1. Add constructor to `core/errors/errors.go`
2. Set appropriate ErrorType, Message, SolutionHint
3. Use in appropriate core package
4. Display in root command

### Adding New Configuration Options

1. Add field to `Config` struct
2. Add default in `DefaultConfig()`
3. Add validation in `Validate()`
4. Update `neev.yaml` template
5. Document in `CONFIGURATION.md`

### Adding New Core Functionality

1. Create new package in `core/`
2. Define clear public API
3. Use `core/errors` for error handling
4. Use `core/logger` for logging
5. Write comprehensive tests
6. Document in `ARCHITECTURE.md`

## Testing Strategy

### Test Organization

```
core/foundation/
├── init.go
├── init_test.go         # Tests for init.go
├── inspect.go
└── inspect_test.go      # Tests for inspect.go
```

### Test Types

1. **Unit Tests**: Single function behavior
2. **Integration Tests**: Multiple components together
3. **File System Tests**: Use `t.TempDir()`
4. **Error Tests**: Verify error types and messages

### Example

```go
func TestInit_CreatesFoundation(t *testing.T) {
    tmpDir := t.TempDir()
    
    err := Init(tmpDir)
    
    if err != nil {
        t.Fatalf("Init failed: %v", err)
    }
    
    if _, err := os.Stat(filepath.Join(tmpDir, ".neev")); err != nil {
        t.Error("Foundation directory not created")
    }
}
```

---

**Next Steps**: See [CONTRIBUTING.md](CONTRIBUTING.md) for development workflow and [DEVELOPMENT.md](DEVELOPMENT.md) for setup instructions.
