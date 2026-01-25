# 🎯 Neev Complete Command Catalog

**This is the authoritative reference for ALL Neev commands.** Use this as the source of truth for documentation.

---

## Commands Overview

Neev provides 14 commands organized into 4 categories:

| Category | Commands | Purpose |
|----------|----------|---------|
| **Foundation** | init, lay | Set up and manage project structure |
| **Blueprints** | draft, bridge, inspect | Create and organize blueprints with polyglot drift detection |
| **Generation** | openapi, cucumber, handoff, instructions | Generate specifications and outputs |
| **Integration** | slash-commands, migrate, sync-remotes | AI tool integration and migration |
| **System** | completion, help | Shell integration and help |

---

## 1. Foundation Commands

### neev init

**Initialize a new Neev project foundation**

```bash
neev init
```

**Description:**
Creates the `.neev/` directory structure and supporting files needed for spec-driven development.

**Files Created:**
- `.neev/neev.yaml` - Project configuration
- `.neev/blueprints/` - Blueprint storage directory
- `.neev/foundation/` - Foundation specifications:
  - `stack.md` - Technology stack template
  - `principles.md` - Design principles template
  - `patterns.md` - Patterns & practices template
- `.neev/commands/registry.yaml` - Command registry for all AI tools
- `.github/copilot-instructions.md` - GitHub Copilot Chat instructions
- `.github/slash-commands.json` - Copilot slash command manifest
- `.cursor/commands.json` - Cursor IDE command integration
- `.vscode/commands.json` - VS Code command palette integration
- `AGENTS.md` - AI tool definitions and fallback instructions

**Flags:**
None

**Examples:**
```bash
# Initialize in current directory
neev init

# Output
🏗️  Laying foundation in /path/to/project
📋 Created AGENTS.md with slash command definitions
✅ Foundation laid successfully!
```

**Use Case:**
Run once per project to set up the foundation. Cannot reinitialize if `.neev/` already exists.

---

### neev lay

**Archive a completed blueprint into the foundation**

```bash
neev lay <blueprint_name> [flags]
```

**Description:**
Moves a blueprint from active development to the foundation archive. Updates the changelog and marks the blueprint as complete.

**Parameters:**
- `<blueprint_name>` - Name of the blueprint directory to archive

**Flags:**
None

**Examples:**
```bash
# Archive the user-authentication blueprint
neev lay user-authentication

# Output
📦 Archiving blueprint: user-authentication
✅ Blueprint moved to foundation archive
📝 Updated changelog
```

**Use Case:**
When a blueprint is complete and you want to preserve it in the foundation for reference.

---

## 2. Blueprint Commands

### neev draft

**Create a draft blueprint**

```bash
neev draft <title> [flags]
```

**Description:**
Creates a new blueprint with template files for capturing project specifications and architectural decisions.

**Parameters:**
- `<title>` - Blueprint title (can contain spaces)

**Flags:**
None

**Files Created:**
- `.neev/blueprints/<sanitized-name>/intent.md` - Purpose and goals
- `.neev/blueprints/<sanitized-name>/architecture.md` - Technical design

**Examples:**
```bash
# Create a user authentication blueprint
neev draft "user authentication"

# Create an API design blueprint
neev draft "REST API"

# Output
📋 Creating blueprint: user-authentication
✅ Blueprint created successfully!
```

**Output Structure:**
```
.neev/blueprints/user-authentication/
├── intent.md          # What this blueprint is for, goals, scope
└── architecture.md    # Technical architecture, design decisions, API specs
```

**Use Case:**
Before building a feature, create a blueprint to document your intent and design. Fill in `intent.md` with what you're building, and `architecture.md` with how you're building it.

---

### neev bridge

**Aggregate project context for AI agents**

```bash
neev bridge [flags]
```

**Description:**
Generates aggregated context from foundation and all blueprints, perfect for passing to AI coding assistants.

**Flags:**
- `-f, --focus <string>` - Filter context by keyword (e.g., "auth", "api")
- `--claude` - Format output optimized for Claude AI
- `--slash` - Format output for IDE slash commands
- `--with-remotes` - Include synced remote foundations in context

**Examples:**
```bash
# Get all project context
neev bridge

# Get context filtered to authentication-related items
neev bridge --focus auth

# Get context optimized for Claude AI
neev bridge --claude

# Save full context to file
neev bridge > project-context.md

# Get context for API-related components
neev bridge -f api > api-context.md

# Include remote foundations in context
neev bridge --with-remotes
```

**Output:**
Markdown document containing:
- All foundation files (stack, principles, patterns)
- All blueprint intents and architectures
- Structured for AI consumption

**Use Case:**
Run before asking AI agents for implementation. Pass the output to Claude, ChatGPT, or other AI tools to provide full project context.

---

### neev inspect

**Check project structure against foundation specifications with multi-language drift detection**

```bash
neev inspect [flags]
```

**Description:**
Analyzes the project structure and verifies it matches the foundation specifications. Detects drift between specifications and implementation across multiple programming languages. Supports three levels of analysis:

- **Level 1**: Directory/file structure checking (default)
- **Level 2**: OpenAPI contract validation (checks API endpoints)
- **Level 3**: Function signature validation (checks method signatures)

**Supported Languages:**
- Go (.go)
- Python (.py) 
- JavaScript/TypeScript (.js, .ts, .jsx, .tsx)
- Java (.java)
- C# (.cs)
- Ruby (.rb)

**Flags:**
- `--json` - Output results in JSON format
- `--strict` - Exit with code 1 if any drift is detected (for CI pipelines)
- `--use-descriptors` - Use `.module.yaml` files for detailed inspection
- `--depth int` - Depth of analysis: 1=structure, 2=+API, 3=+signatures (default: 1)
- `--check-api` - Validate OpenAPI specs (enables Level 2)
- `--check-signatures` - Validate function signatures (enables Level 3)
- `--check-tests` - Validate BDD test coverage (not yet implemented)

**Examples:**
```bash
# Basic structure check (Level 1)
neev inspect

# Output
🔍 POLYGLOT DRIFT DETECTION

═════ LANGUAGES DETECTED =====
Go (23 files) | Python (5 files) | JavaScript (12 files)

═════ LEVEL 1: STRUCTURE =====
✅ Modules: 12/12 matching
⚠️  Extra: utils/ (undocumented)

# Check API contracts (Level 2)
neev inspect --check-api
# or
neev inspect --depth 2

# Output includes
═════ LEVEL 2: API CONTRACTS =====
[Go]
  ✅ GET /api/users → ListUsers handler found
  🔴 DELETE /api/users/{id} → NO HANDLER (documented)

[Python]
  ✅ POST /api/users → create_user handler found

# Full validation including signatures (Level 3)
neev inspect --depth 3
# or
neev inspect --check-signatures

# Output includes
═════ LEVEL 3: SIGNATURES =====
[Go]
  ✅ ListUsers(w http.ResponseWriter, r *http.Request) → Matches
  ⚠️  CreateUser(..., Logger logger) → Extra parameter (undocumented)

# Get JSON output for parsing
neev inspect --json --depth 3

# Fail if any drift is detected (useful in CI/CD)
neev inspect --strict --check-api

# Use detailed module descriptors
neev inspect --use-descriptors
```

**Output Structure:**
- **Human-readable**: 
  - Language breakdown showing file counts
  - Grouped warnings by severity (error/warning/info)
  - Detailed remediation suggestions
  - Summary statistics by analysis level

- **JSON**: 
```json
{
  "success": false,
  "warnings": [
    {
      "type": "MISSING_ENDPOINT",
      "module": "api",
      "message": "API endpoint DELETE /api/users/{id} is documented but not implemented",
      "severity": "error",
      "remediation": "Implement handler for DELETE /api/users/{id}"
    }
  ],
  "summary": {
    "total_modules": 12,
    "matching_modules": 11,
    "missing_modules": 0,
    "extra_code_dirs": 1,
    "languages": {
      "go": 23,
      "python": 5,
      "javascript": 12
    },
    "missing_endpoints": 1,
    "undocumented_endpoints": 0,
    "signature_mismatches": 1,
    "total_warnings": 3,
    "error_count": 1,
    "warning_count": 2
  }
}
```

**Module Descriptor Example (`.module.yaml`):**
```yaml
name: api
description: REST API handlers
expected_functions:
  - name: ListUsers
    language: go
    file_pattern: "handlers/*.go"
    parameters:
      - name: w
        type: "http.ResponseWriter"
      - name: r
        type: "*http.Request"
    returns:
      - type: error
    visibility: public
    
  - name: CreateUser
    language: go
    parameters:
      - name: ctx
        type: context.Context
      - name: user
        type: "*User"
    returns:
      - type: "*User"
      - type: error
```

**Use Cases:**
1. **CI/CD Integration**: Run `neev inspect --strict --check-api` to fail builds on drift
2. **API Contract Testing**: Use `--check-api` to verify all documented endpoints are implemented
3. **Signature Validation**: Use `--check-signatures` to ensure function signatures match specs
4. **Polyglot Projects**: Automatically detects and validates code in 6+ languages
5. **Refactoring Safety**: Verify changes don't break documented interfaces

**Exit Codes:**
- `0` - No drift detected (or only warnings in non-strict mode)
- `1` - Drift detected with `--strict` flag

---

## 3. Generation Commands

### neev openapi

**Generate OpenAPI specification from a blueprint**

```bash
neev openapi <blueprint> [flags]
```

**Description:**
Parses the `architecture.md` from a blueprint and generates an OpenAPI 3.1 specification file.

**Parameters:**
- `<blueprint>` - Name of the blueprint directory containing architecture.md

**Flags:**
None

**Examples:**
```bash
# Generate OpenAPI spec from API blueprint
neev openapi api-design

# Output
📄 Parsing architecture.md from api-design blueprint
✅ Generated: api-design/openapi.yaml
```

**Output:**
- `<blueprint>/openapi.yaml` - Full OpenAPI 3.1 specification

**Expected Architecture Format:**
Your `architecture.md` should contain API endpoint definitions:
```
# API Endpoints

## GET /users
Returns list of users
Parameters: page, limit
Response: { id, name, email }

## POST /users
Create new user
Body: { name, email, password }
Response: { id, name, email, created_at }
```

**Use Case:**
Document your REST API design in a blueprint, then auto-generate the official OpenAPI specification.

---

### neev cucumber

**Generate Cucumber/BDD test scaffolding**

```bash
neev cucumber <blueprint> [flags]
```

**Description:**
Parses the `architecture.md` from a blueprint and generates Cucumber feature files and step definition scaffolds for API testing.

**Parameters:**
- `<blueprint>` - Name of the blueprint directory containing architecture.md

**Flags:**
- `-l, --lang <language>` - Language for step definitions: `go`, `javascript`, `python`

**Examples:**
```bash
# Generate Cucumber tests in Go
neev cucumber api-design --lang go

# Generate Cucumber tests in Python
neev cucumber api-design -l python

# Output
🧪 Generating Cucumber tests from api-design blueprint
✅ Generated: api-design/features/
✅ Generated: api-design/steps/
```

**Output:**
- `<blueprint>/features/` - Gherkin feature files
- `<blueprint>/steps/` - Step definition scaffolds in specified language

**Use Case:**
Generate BDD test structure from your API blueprint before implementing tests.

---

### neev handoff

**Create AI agent handoff prompts**

```bash
neev handoff <role> [flags]
```

**Description:**
Generates a structured handoff prompt with role-specific instructions from `.neev/agents/<role>.md`.

**Parameters:**
- `<role>` - Agent role (e.g., "backend", "frontend", "devops")

**Flags:**
- `--markdown` - Wrap output in markdown code fence for copy-paste (default: true)

**Examples:**
```bash
# Generate handoff for backend developer
neev handoff backend

# Generate handoff for frontend developer
neev handoff frontend

# Get raw output without markdown fence
neev handoff backend --markdown=false

# Output
👋 Generating handoff prompt for: backend

Role: backend
Context: [Full project context]
Guidelines: [Role-specific guidelines]
Active blueprints: [Relevant blueprints]
```

**Use Case:**
When handing off work to another developer or AI agent, provide role-specific context and instructions.

---

### neev instructions

**Generate GitHub Copilot Chat instructions**

```bash
neev instructions [flags]
```

**Description:**
Generates or updates `.github/copilot-instructions.md` with context about your project that helps GitHub Copilot provide better suggestions.

**Files:**
- `.github/copilot-instructions.md` - Updated with latest foundation and blueprint context

**Flags:**
None

**Examples:**
```bash
# Generate/update Copilot instructions
neev instructions

# Output
📝 Generating GitHub Copilot instructions...
✅ Created/updated: .github/copilot-instructions.md
```

**Includes:**
- Foundation module summary
- Active blueprint intents
- Development guidelines
- Slash command definitions for all 6 commands

**Use Case:**
Run after adding new blueprints or updating foundation to keep Copilot Chat informed about your project.

---

## 4. Integration Commands

### neev slash-commands

**Configure slash commands for AI tools**

```bash
neev slash-commands [flags]
```

**Description:**
Manages slash command configuration and integration for various AI coding assistants (Claude Code, Cursor, CodeBuddy, OpenCode, Qoder, Codex, RooCode, GitHub Copilot).

**Flags:**
- `--list` - List all available slash commands
- `--tool <string>` - Show commands for a specific AI tool
- `--update` - Update AGENTS.md with latest commands

**Examples:**
```bash
# List all available slash commands
neev slash-commands --list

# Show commands available for Cursor IDE
neev slash-commands --tool cursor

# Show commands for Claude Code
neev slash-commands --tool "claude-code"

# Update AGENTS.md with latest commands
neev slash-commands --update

# Output
✅ Slash commands registered:
  /neev:bridge       - Get full project context
  /neev:draft        - Create a new blueprint
  /neev:inspect      - Analyze project structure
  /neev:cucumber     - Generate BDD tests
  /neev:openapi      - Generate API spec
  /neev:handoff      - Prepare for AI handoff
```

**Supported AI Tools:**
- Claude Code
- Cursor IDE
- CodeBuddy
- OpenCode
- Qoder
- Codex
- RooCode
- GitHub Copilot

**Use Case:**
Check which slash commands are available and how to use them in your AI tool of choice.

---

### neev migrate

**Convert existing projects to Neev structure**

```bash
neev migrate [flags]
```

**Description:**
Migrates existing projects (OpenSpec, Spec-Kit) to Neev structure.

**Flags:**
- `-s, --source <string>` - Source type: `openspec`, `speckit`, or `auto` (default: auto)
- `--dry-run` - Preview changes without applying
- `--backup` - Create backup of existing `.neev` directory

**Examples:**
```bash
# Auto-detect and migrate
neev migrate

# Migrate from OpenSpec (auto-detect)
neev migrate --source openspec

# Preview migration without changes
neev migrate --dry-run

# Backup existing .neev before migrating
neev migrate --backup

# Output
🔄 Migrating from: openspec
📋 Detected 3 specs
✅ Migration preview complete (no changes made)
  - Convert: specs/api.yml → blueprints/api/architecture.md
  - Convert: specs/db.yml → blueprints/database/architecture.md
  - Convert: specs/auth.yml → blueprints/auth/architecture.md

# Run without --dry-run to apply
```

**Use Case:**
If you're already using OpenSpec or Spec-Kit, migrate to Neev to get all the new features.

---

### neev sync-remotes

**Synchronize remote foundations**

```bash
neev sync-remotes [flags]
```

**Description:**
Syncs remote foundation sources defined in `neev.yaml` to `.neev/remotes/`. Allows referencing external foundations from other repositories.

**Flags:**
- `--json` - Output results in JSON format

**Configuration in neev.yaml:**
```yaml
remotes:
  - name: api
    path: "../my-api-repo/.neev/foundation"
    public_only: true
  - name: shared
    path: "../shared-lib/.neev/foundation"
    public_only: false
```

**Examples:**
```bash
# Sync all configured remotes
neev sync-remotes

# Output as JSON
neev sync-remotes --json

# Output
🔄 Syncing remotes...
  ✅ api → .neev/remotes/api
  ✅ shared → .neev/remotes/shared
```

**Use Case:**
Reference foundation specifications from other projects in your team (monorepo, microservices, shared libraries).

---

## 5. Slash Commands Reference

All Neev commands can be accessed as slash commands in GitHub Copilot Chat and other AI tools:

| Slash Command | Equivalent CLI | Purpose |
|---------------|----------------|---------|
| `/neev:bridge` | `neev bridge` | Get full project context |
| `/neev:draft` | `neev draft` | Create a new blueprint |
| `/neev:inspect` | `neev inspect` | Analyze project structure |
| `/neev:cucumber` | `neev cucumber` | Generate BDD tests |
| `/neev:openapi` | `neev openapi` | Generate API spec |
| `/neev:handoff` | `neev handoff` | Prepare for AI handoff |

**Usage in GitHub Copilot Chat:**
```
@Copilot /neev:bridge
@Copilot /neev:draft Create authentication system
@Copilot /neev:inspect Check for drift
```

---

## 6. Common Workflows

### Workflow 1: Create New Feature

```bash
# 1. Create a blueprint for the feature
neev draft "Feature Name"

# 2. Edit .neev/blueprints/feature-name/
#    - Fill intent.md with what you're building
#    - Fill architecture.md with how you're building it

# 3. Get project context including your new blueprint
neev bridge > context.md

# 4. Pass context to AI agent for implementation
cat context.md | xclip  # Copy to clipboard, then paste in Copilot Chat

# 5. After implementation, verify it matches specifications
neev inspect
```

### Workflow 2: Generate API Documentation

```bash
# 1. Create an API blueprint
neev draft "REST API"

# 2. Document your endpoints in architecture.md
# 3. Generate OpenAPI spec
neev openapi "REST API"

# 4. Use generated openapi.yaml for API docs
```

### Workflow 3: Set Up BDD Testing

```bash
# 1. Create an API blueprint
neev draft "API Design"

# 2. Document endpoints in architecture.md
# 3. Generate Cucumber tests
neev cucumber "API Design" --lang go

# 4. Implement the generated test steps
```

### Workflow 4: Hand Off to Another Developer

```bash
# 1. Create a handoff role file: .neev/agents/backend.md
# 2. Generate handoff prompt
neev handoff backend

# 3. Share output with the developer
```

---

## 7. Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error / `--strict` flag detected drift |
| 2 | Invalid command or arguments |
| 3 | File not found or permission denied |

---

## 8. System Commands

### neev completion

Generate shell completion scripts for bash, zsh, etc.

```bash
neev completion bash
neev completion zsh
neev completion fish
```

### neev help

Get help on any command.

```bash
neev help
neev help bridge
neev help draft
```

---

## When to Use Each Command

| Task | Command | Frequency |
|------|---------|-----------|
| Initialize project | `neev init` | Once per project |
| Plan a feature | `neev draft` | Per feature |
| Get AI context | `neev bridge` | Before each AI interaction |
| Verify implementation | `neev inspect` | Before commits/PRs |
| Generate API docs | `neev openapi` | After API changes |
| Generate BDD tests | `neev cucumber` | After design decisions |
| Hand off work | `neev handoff` | At project milestones |
| Archive blueprint | `neev lay` | When feature complete |
| Update AI instructions | `neev instructions` | After adding blueprints |
| Check slash commands | `neev slash-commands` | Learning/troubleshooting |
| Migrate from other tools | `neev migrate` | One-time migration |
| Sync external specs | `neev sync-remotes` | In monorepos/shared libs |

---

## Related Documentation

- 🎯 **This file** - Complete command catalog
- 📚 [USAGE.md](USAGE.md) - Practical usage guide
- 🏗️ [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture
- 🤝 [CONTRIBUTING.md](CONTRIBUTING.md) - Development setup
- 💻 [COPILOT_SLASH_COMMANDS.md](COPILOT_SLASH_COMMANDS.md) - Slash command details

---

## Need Help?

```bash
# General help
neev help

# Help for a specific command
neev [command] --help

# Examples
neev draft --help
neev bridge --help
neev inspect --help
```
