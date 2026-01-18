# 🏗️ Neev - AI-Ready Blueprint Framework

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![Tests](https://img.shields.io/github/actions/workflow/status/neev-kit/neev/tests.yml?branch=main&style=flat-square&label=Tests)](https://github.com/neev-kit/neev/actions)
[![Release](https://img.shields.io/github/v/release/neev-kit/neev?style=flat-square&label=Release)](https://github.com/neev-kit/neev/releases)

**Build better software by bridging project intent with AI coding assistants.**

Neev is a lightweight CLI framework that helps you capture project blueprints, aggregate context, and seamlessly hand off to AI agents. No dependencies on external APIs or complex setup — just structured markdown files versioned in your repository.

## Why Neev?

Traditional AI coding assistants work best with clear project context. Neev solves this by:

- **🎯 Explicit Intent**: Capture what you want to build before implementation
- **📚 Organized Context**: Structure project knowledge in `.neev/` (version controlled)
- **🤖 AI-Ready**: Generate context aggregations perfect for LLM consumption
- **🔧 Zero Friction**: Works with any AI tool — no API keys or configuration
- **⚡ Fast Setup**: Initialize and start drafting in seconds

## Quick Start

### 1. Installation

**Option A: Download Pre-built Binary (Recommended)**

Download the latest stable release from [GitHub Releases](https://github.com/neev-kit/neev/releases):

```bash
# macOS (Intel)
curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_darwin_amd64.tar.gz | tar xz
sudo mv neev /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_darwin_arm64.tar.gz | tar xz
sudo mv neev /usr/local/bin/

# Linux
curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_linux_amd64.tar.gz | tar xz
sudo mv neev /usr/local/bin/

# Windows (PowerShell)
# Download from https://github.com/neev-kit/neev/releases and extract
```

**Option B: Build from Source**

```bash
# Clone and build
git clone https://github.com/neev-kit/neev.git
cd neev
go mod download
go build -o neev ./cli
sudo mv neev /usr/local/bin/

# Or install directly
go install github.com/neev-kit/neev/cli@latest
```

**Verify Installation**

```bash
neev --version
neev --help
```

### 2. Initialize Your Project

```bash
cd /path/to/your/project
neev init
```

Creates:
```
.neev/
├── neev.yaml              # Project configuration
├── blueprints/            # Your blueprint collection
└── foundation/            # Project foundations & principles
```

### 3. Create Your First Blueprint

```bash
neev draft "user-authentication"
neev draft "Database Schema"
```

### 4. Aggregate Context for AI

```bash
# Get full project context
neev bridge

# Filter by keyword
neev bridge --focus "authentication"

# Save to file
neev bridge > context.md
```

## Core Concepts

### Blueprint
A markdown-based specification of a feature or component you want to build. Each blueprint is self-contained and can reference others.

**Example**:
```
.neev/blueprints/user-auth/
├── intent.md          # What and why
├── architecture.md    # How it works
├── api-spec.md        # API contracts
└── security.md        # Security considerations
```

### Foundation
Project-wide principles, conventions, and architectural decisions. Shared across all blueprints.

**Example**:
```
.neev/foundation/
├── principles.md      # Project values
├── stack.md           # Technology choices
└── conventions.md     # Coding standards
```

### Context
Aggregated, searchable project information ready for AI consumption. Generated via `neev bridge`.

## How It Works

```
┌─────────────────────────────┐
│  Write Blueprints & Docs    │
│  (Markdown in .neev/)       │
└────────────┬────────────────┘
             │
             ▼
┌─────────────────────────────┐
│  Run: neev bridge           │
│  (Aggregate & optionally    │
│   filter by keywords)       │
└────────────┬────────────────┘
             │
             ▼
┌─────────────────────────────┐
│  Get Context Output         │
│  (Ready for AI agents)      │
└────────────┬────────────────┘
             │
             ▼
┌─────────────────────────────┐
│  Share with AI Coding       │
│  Assistant (Claude, Cursor, │
│  GitHub Copilot, etc.)      │
└─────────────────────────────┘
```

## Commands

### `neev init`

Initialize Neev in your project.

```bash
neev init
```

**Creates:**
- `.neev/` directory structure
- `neev.yaml` configuration file
- Empty `blueprints/` and `foundation/` directories
- Prevents accidental overwrites

### `neev draft <title>`

Create a new blueprint with template files.

```bash
neev draft "user-authentication"
neev draft "API Gateway"
```

**Creates:**
- Blueprint directory with sanitized name
- `intent.md` — Purpose and goals
- `architecture.md` — Technical design

### `neev bridge [flags]`

Aggregate and output project context.

```bash
neev bridge                    # Full context
neev bridge --focus auth       # Filter by keyword
neev bridge -f db > context.md # Save to file
```

**Flags:**
- `--focus, -f` — Filter by keyword

**Output:** Markdown with all foundation + blueprint content

### `neev inspect` (internal)

Analyze project structure and find missing blueprints.

```bash
neev inspect                    # Human-readable output
neev inspect --json             # Machine-readable JSON for CI/CD
neev inspect --use-descriptors  # File-level validation with .module.yaml
```

### `neev sync-remotes`

Synchronize remote foundation sources from other repositories.

```bash
neev sync-remotes               # Sync all remotes from neev.yaml
neev sync-remotes --json        # JSON output
```

### `neev instructions`

Generate GitHub Copilot instructions from your foundation and blueprints.

```bash
neev instructions               # Creates .github/copilot-instructions.md
```

## Production Features

Neev includes production-grade features for enterprise use:

### 🔍 Advanced Drift Detection
- **Structured Warnings**: Categorized drift detection (MISSING_MODULE, EXTRA_CODE, etc.)
- **Module Descriptors**: Define expected files and patterns in `.module.yaml` files
- **CI/CD Integration**: JSON output for automated checks

```bash
neev inspect --json             # Get structured drift report
neev inspect --use-descriptors  # Validate against module descriptors
```

### 🌐 Polyrepo Support
- **Remote Foundations**: Reference foundations from other repositories
- **Cross-Repo Context**: Include external specs in bridge context
- **Public/Private Control**: Filter what gets shared with `public_only`

```yaml
# neev.yaml
remotes:
  - name: api
    path: "../backend/.neev/foundation"
    public_only: true
```

```bash
neev sync-remotes               # Sync all remotes
neev bridge --with-remotes      # Include remotes in context
```

### 🤖 AI Assistant Integration
- **GitHub Copilot**: Auto-generate instructions from your specs
- **Claude Optimization**: Special formatting for Claude AI
- **Context Management**: Smart aggregation for better AI suggestions

```bash
neev instructions               # Generate Copilot instructions
neev bridge --claude            # Claude-optimized output
neev bridge --claude --with-remotes  # Full context for Claude
```

See [PRODUCTION_ENHANCEMENTS.md](PRODUCTION_ENHANCEMENTS.md) for detailed documentation.

## Configuration

The `neev.yaml` file controls Neev behavior:

```yaml
project_name: My Project
foundation_path: .neev
ignore_dirs:
  - node_modules
  - .git
  - __pycache__
  - vendor
```

**Options:**
- `project_name` — Display name for your project
- `foundation_path` — Where `.neev/` directory lives (default: `.neev`)
- `ignore_dirs` — Directories to skip during inspection

## Project Structure (Neev Repository)

```
neev/
├── cli/                     # CLI commands
│   └── cmd/
│       ├── root.go          # Root command (logger init)
│       ├── init.go          # Initialize foundation
│       ├── draft.go         # Create blueprints
│       ├── bridge.go        # Aggregate context
│       └── *_test.go        # Command tests
│
├── core/                    # Business logic
│   ├── foundation/          # Init & inspect projects
│   ├── blueprint/           # Blueprint management
│   ├── bridge/              # Context aggregation
│   ├── config/              # Configuration loading
│   ├── errors/              # Custom error types
│   └── logger/              # Structured logging
│
├── .github/workflows/       # CI/CD
│   ├── tests.yml            # Run tests
│   └── release.yml          # Build & release
│
└── Documentation
    ├── README.md            # This file
    ├── CONTRIBUTING.md      # Development guidelines
    ├── DEVELOPMENT.md       # Setup & debugging
    ├── ARCHITECTURE.md      # System design
    └── USAGE.md             # Detailed usage guide
```

## Getting Started

### Step 1: Install Neev

```bash
# Build from source
git clone https://github.com/neev-kit/neev.git
cd neev
go build -o neev ./cli

# Or use go install
go install github.com/neev-kit/neev/cli@latest
```

### Step 2: Initialize Your Project

```bash
cd /path/to/your/project
neev init
```

### Step 3: Write Blueprints

```bash
# Create blueprints for features you want to build
neev draft "User Authentication"
neev draft "Database Layer"
neev draft "API Gateway"

# Edit the generated files with details
# .neev/blueprints/user-authentication/intent.md
# .neev/blueprints/user-authentication/architecture.md
```

### Step 4: Use with AI Assistants

```bash
# Generate context
neev bridge > context.md

# Copy to Claude, Cursor, GitHub Copilot, etc.
# Or pipe directly to your AI tool
neev bridge | pbcopy  # macOS
neev bridge | xclip   # Linux
```

## Examples

### Example: Building a Payment System

**1. Create blueprints:**
```bash
neev draft "Payment Processing"
neev draft "Webhook Management"
neev draft "Error Handling"
```

**2. Add foundation:**
```
.neev/foundation/
├── stack.md         # "We use Go, PostgreSQL, Redis"
├── principles.md    # "Security first, simplicity second"
└── patterns.md      # "Repository pattern, dependency injection"
```

**3. Aggregate context:**
```bash
neev bridge --focus payment
```

**4. Share with AI:**
Paste the output into your AI coding assistant with your implementation request.

### Example: Team Onboarding

**1. Document architecture:**
```bash
neev draft "System Overview"
neev draft "Authentication Flow"
neev draft "Database Schema"
```

**2. Create foundation:**
```
.neev/foundation/contributing.md   # How to contribute
.neev/foundation/conventions.md    # Code style & patterns
```

**3. Share with new team members:**
```bash
neev bridge > ONBOARDING.md
```

## Use Cases

| Use Case | How Neev Helps |
|----------|---|
| **AI Pair Programming** | Context-aware coding with structured project knowledge |
| **Onboarding** | New team members get structured project overview |
| **Architecture Decisions** | Document and share technical choices |
| **API Documentation** | Maintain API specs alongside implementation |
| **Feature Planning** | Capture requirements before implementation |
| **Code Review** | Reviewers understand intent + architecture |

## Development

### Building & Testing

```bash
# Build
go build -o neev ./cli

# Run tests
go test ./...

# Test coverage
go test -cover ./...

# Run specific command
go run ./cli init
```

### Setup Local Development

See [DEVELOPMENT.md](DEVELOPMENT.md) for:
- Detailed setup instructions
- IDE configuration (VSCode, GoLand)
- Debugging with Delve
- Performance profiling

### Understanding the Codebase

See [ARCHITECTURE.md](ARCHITECTURE.md) for:
- System design overview
- Component interactions
- Data flow diagrams
- Extension points

### Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for:
- Development workflow
- Code standards
- Testing guidelines
- Commit conventions
- PR process

## Key Features

- **🚀 Zero Setup** — No API keys, no dependencies, no cloud login
- **🔒 Local First** — All files stored in `.neev/`, fully version controlled
- **🎯 AI-Optimized** — Output structured for LLM consumption
- **🎨 Markdown Native** — Work in the format you already know
- **⚡ Multi-Module** — Go-based with modular architecture
- **📦 Production Ready** — Structured logging, error handling, configuration management
- **✅ Well Tested** — 50+ comprehensive tests across all packages
- **🔄 CI/CD Ready** — GitHub Actions workflows included

## Standards & Practices

Neev follows industry best practices:

- **Structured Logging** — `slog` for consistent, parseable logs
- **Error Handling** — Custom error types with solution hints
- **Configuration** — YAML-based with validation and defaults
- **Testing** — Unit tests, integration tests, table-driven tests
- **Conventional Commits** — Clear git history
- **Clean Architecture** — Separation of concerns, testable design

See [ARCHITECTURE.md](ARCHITECTURE.md) for details.

## Comparison with Similar Tools

| Feature | Neev | Spec-Kit | OpenSpec |
|---------|------|----------|----------|
| **Local Files** | ✅ | ✅ | ✅ |
| **No External API** | ✅ | ✅ | ✅ |
| **Greenfield (0→1)** | ✅ | ⭐ | ✅ |
| **Brownfield (1→n)** | ✅ | ✅ | ⭐ |
| **Blueprint Templates** | ✅ | ✅ | ✅ |
| **Context Aggregation** | ✅ | Limited | ✅ |
| **Written in** | Go | Python | TypeScript |
| **CLI First** | ✅ | ✅ | ✅ |

**Best For:**
- **Neev** — Go projects, CLI tools, fast setup
- **Spec-Kit** — Comprehensive spec-driven workflow
- **OpenSpec** — Teams with complex change management

## Real-World Example

```bash
# 1. Initialize in your Go project
$ neev init

# 2. Plan features as blueprints
$ neev draft "User API"
$ neev draft "Authentication"

# 3. Document in .neev/foundation/
# Edited: .neev/foundation/principles.md
#   - Security-first design
#   - RESTful APIs
#   - PostgreSQL for persistence

# 4. Get AI-ready context
$ neev bridge --focus user > user-context.md

# 5. Share with Claude/Copilot for implementation
```

The output is ready for: "Build this according to the context above"

## Troubleshooting

| Issue | Solution |
|-------|----------|
| `command not found: neev` | Ensure `$GOPATH/bin` is in `$PATH` or build locally |
| `.neev` already exists | Use `neev init` only once per project |
| No blueprints generated | Run `neev draft "name"` to create blueprints |
| Bridge output is empty | Check that `.neev/foundation/` and `.neev/blueprints/` have `.md` files |

See [USAGE.md](USAGE.md) for detailed troubleshooting.

## FAQ

**Q: Do I need to commit `.neev/` to git?**  
A: Yes! `.neev/` contains your project knowledge and should be versioned.

**Q: Can I use Neev with non-Go projects?**  
A: Absolutely. Neev works with any project type.

**Q: How does Neev compare to writing prompts manually?**  
A: Neev structures your knowledge so AI gets context automatically, reducing manual copy/paste and keeping things in sync.

**Q: Can teams share blueprints?**  
A: Yes. Common patterns can be captured in `.neev/foundation/` and reused across projects.

**Q: What about large projects?**  
A: Use `--focus` flag to filter context by keywords. Blueprints can reference each other.

## Status

- ✅ **Phase 1**: Blueprint drafting & context bridging
- ✅ **Phase 2**: Test coverage & CLI hardening
- ✅ **Phase 3**: Comprehensive test suite
- ✅ **Phase 4**: Production hardening (logging, errors, config, CI/CD)
- 🚀 **v1.0.0**: Ready for production use

## License

MIT License — See [LICENSE](LICENSE) file

## Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md)

- 🐛 Found a bug? [Open an issue](https://github.com/neev-kit/neev/issues)
- ✨ Have an idea? [Start a discussion](https://github.com/neev-kit/neev/discussions)
- 🔧 Want to contribute? [See CONTRIBUTING.md](CONTRIBUTING.md)

## Maintainers

See [MAINTAINERS.md](MAINTAINERS.md) for core team and advisors.

---

**Ready to build better software with AI?** 

## 📚 Documentation

- **🚀 [Getting Started](GETTING_STARTED.md)** — Complete beginner's guide
- **📖 [Usage Guide](USAGE.md)** — Detailed command reference
- **📋 [API Reference](API_REFERENCE.md)** — Complete command documentation
- **🎓 [Tutorials](TUTORIALS.md)** — 8 step-by-step walkthroughs
- **💡 [Best Practices](BEST_PRACTICES.md)** — Patterns and anti-patterns
- **❓ [FAQ](FAQ.md)** — Common questions and troubleshooting
- **🏗️ [Architecture](ARCHITECTURE.md)** — System design
- **💻 [Development](DEVELOPMENT.md)** — Contributing guide

👉 **New to Neev?** Start with [Getting Started](GETTING_STARTED.md)
