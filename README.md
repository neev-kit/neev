# Neev - AI-Ready Blueprint Framework

Neev is a Go-based CLI framework that helps developers create and manage project blueprints for AI integration. It provides tools to draft blueprints, aggregate project context, and bridge to external systems like AI agents.

## Features

- **🏗️ Project Initialization** (`neev init`) - Set up a `.neev` foundation directory with configuration
- **📝 Blueprint Drafting** (`neev draft`) - Create structured blueprint templates for your project
- **🌉 Context Bridging** (`neev bridge`) - Aggregate project context for AI agents and external tools

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/neev-kit/neev
cd neev

# Install dependencies
go mod download

# Build the CLI
go build -o neev ./cli
```

### Basic Usage

```bash
# Initialize a new project foundation
./neev init

# Create a new blueprint
./neev draft "my-blueprint"

# Aggregate project context
./neev bridge

# Get context with specific focus
./neev bridge --focus "database"
```

## Project Structure

```
neev/
├── cli/                         # CLI application
│   ├── main.go                 # Entry point
│   └── cmd/                    # Cobra commands
│       ├── root.go             # Root command
│       ├── init.go             # Init command
│       ├── draft.go            # Draft command
│       └── bridge.go           # Bridge command
├── core/                        # Core functionality
│   ├── foundation/             # Project foundation
│   │   ├── paths.go            # Constants and paths
│   │   └── init.go             # Initialization logic
│   ├── blueprint/              # Blueprint management
│   │   └── draft.go            # Draft creation
│   └── bridge/                 # Context aggregation
│       └── context.go          # Context building
└── .neev/                      # Generated project directory
    ├── neev.yaml               # Project configuration
    ├── blueprints/             # Blueprint storage
    └── foundation/             # Foundation specs
```

## Commands

### `neev init`

Initialize a new Neev foundation in your project.

**Usage:**
```bash
neev init
```

**What it does:**
- Creates `.neev/` directory structure
- Generates `neev.yaml` configuration file
- Sets up `blueprints/` and `foundation/` subdirectories
- Prevents accidental overwrites of existing projects

**Output:**
```
🏗️  Laying foundation in /path/to/project
✅ Foundation laid successfully!
```

### `neev draft <title>`

Create a new blueprint for your project.

**Usage:**
```bash
neev draft "my-blueprint"
neev draft "Authentication Service"
```

**What it does:**
- Sanitizes the blueprint name (converts to lowercase, replaces spaces with hyphens)
- Creates blueprint directory in `.neev/blueprints/`
- Generates template files (`intent.md`, `architecture.md`)
- Prevents duplicate blueprint names

**Output:**
```
✅ Created blueprint at .neev/blueprints/my-blueprint
✅ Created blueprint at .neev/blueprints/authentication-service
```

### `neev bridge [flags]`

Aggregate project context for AI agents and external systems.

**Usage:**
```bash
neev bridge
neev bridge --focus "database"
neev bridge -f "authentication"
```

**Flags:**
- `--focus, -f` (string) - Filter context by keyword

**What it does:**
- Reads all `.md` files from `.neev/foundation/`
- Reads all `.md` files from `.neev/blueprints/`
- Aggregates content into a single context string
- Optionally filters by focus keyword if provided
- Returns formatted context suitable for AI processing

**Output:**
```
# Project Foundation
## File: neev.yaml
...

## File: intent.md
...

## File: architecture.md
...
```

## Configuration

The `neev.yaml` configuration file is auto-generated during initialization:

```yaml
version: "1.0"
name: "my-project"
description: "Project description"
```

## Blueprint Structure

Each blueprint created with `neev draft` contains:

- `intent.md` - Purpose and goals of the blueprint
- `architecture.md` - Technical architecture details

You can extend blueprints by adding additional `.md` files that will be included when running `neev bridge`.

## Development

### Building from Source

```bash
go run ./cli init
go run ./cli draft "test"
go run ./cli bridge
```

### Running Tests

```bash
go test ./...
```

### Coverage

```bash
go test -cover ./...
```

## Architecture

### Core Module (`core/`)

- **foundation**: Project initialization and configuration
- **blueprint**: Blueprint creation and management
- **bridge**: Context aggregation logic

### CLI Module (`cli/`)

Built with [Cobra](https://github.com/spf13/cobra), provides user-facing commands with rich terminal styling via [Lipgloss](https://github.com/charmbracelet/lipgloss).

### Key Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `gopkg.in/yaml.v3` - YAML parsing

## Use Cases

### For Developers
- Document project architecture and intent
- Create reusable blueprint templates for common patterns
- Share project context with AI coding assistants

### For AI Integration
- Generate structured prompts from project blueprints
- Aggregate context for RAG (Retrieval-Augmented Generation) systems
- Maintain consistent project documentation

### For Teams
- Standardize project structure across repositories
- Ensure new contributors understand project intent
- Store architectural decisions in version control

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

See LICENSE file for details.

## Support

For issues and questions, please open an issue on GitHub.
