# 📚 Neev API Reference

Complete reference documentation for all Neev commands, flags, and configuration options.

## Table of Contents

- [Command Overview](#command-overview)
- [Commands](#commands)
  - [neev init](#neev-init)
  - [neev draft](#neev-draft)
  - [neev bridge](#neev-bridge)
  - [neev inspect](#neev-inspect)
  - [neev instructions](#neev-instructions)
  - [neev sync-remotes](#neev-sync-remotes)
  - [neev lay](#neev-lay)
  - [neev migrate](#neev-migrate)
  - [neev handoff](#neev-handoff)
- [Configuration Reference](#configuration-reference)
- [Exit Codes](#exit-codes)
- [Environment Variables](#environment-variables)

## Command Overview

| Command | Purpose | Common Use Case |
|---------|---------|-----------------|
| `init` | Initialize Neev in a project | Starting fresh with Neev |
| `draft` | Create a new blueprint | Planning a feature |
| `bridge` | Aggregate context | Getting AI-ready specs |
| `inspect` | Check project drift | CI/CD validation |
| `instructions` | Generate Copilot instructions | Improving AI suggestions |
| `sync-remotes` | Sync remote foundations | Multi-repo projects |
| `lay` | Archive completed blueprint | Marking work done |
| `migrate` | Convert from other tools | Switching to Neev |
| `handoff` | Create agent prompts | AI agent workflows |

---

## Commands

### neev init

Initialize Neev foundation in your project.

#### Syntax

```bash
neev init [flags]
```

#### Description

Creates the `.neev/` directory structure with:
- `neev.yaml` — Configuration file
- `blueprints/` — Blueprint storage
- `foundation/` — Foundation documents

#### Flags

*None currently*

#### Examples

```bash
# Basic initialization
cd /path/to/project
neev init

# Check what was created
ls -la .neev/
```

#### Output

```
🏗️  Laying foundation in /path/to/project
✅ Foundation laid successfully!
```

#### Exit Codes

- `0` — Success
- `1` — Error (e.g., `.neev` already exists)

#### Errors

| Error Message | Cause | Solution |
|---------------|-------|----------|
| `.neev directory already exists` | Already initialized | Use existing setup or backup and reinit |
| `permission denied` | Insufficient permissions | Check directory permissions |
| `failed to create directory` | File system error | Check disk space and permissions |

#### Notes

- Safe to run multiple times (will error if already exists)
- Creates minimal structure — you add content
- All files are plain text markdown

---

### neev draft

Create a new blueprint with template files.

#### Syntax

```bash
neev draft <blueprint-name> [flags]
```

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `blueprint-name` | Yes | Name of the blueprint (can contain spaces) |

#### Flags

*None currently*

#### Name Sanitization

Blueprint names are automatically sanitized:
- Spaces → hyphens (`-`)
- Uppercase → lowercase
- Special chars preserved where valid

**Examples:**
- `"User Auth"` → `user-auth`
- `"API v2.0"` → `api-v2.0`
- `"Database Schema"` → `database-schema`

#### Generated Files

Creates two template files:

**`intent.md`:**
```markdown
# Template for intent.md
```

**`architecture.md`:**
```markdown
# Template for architecture.md
```

#### Examples

```bash
# Simple name
neev draft "Authentication"
# Creates: .neev/blueprints/authentication/

# Multi-word name
neev draft "Payment Processing"
# Creates: .neev/blueprints/payment-processing/

# With version
neev draft "API v2"
# Creates: .neev/blueprints/api-v2/
```

#### Output

```
✅ Created blueprint at .neev/blueprints/authentication
```

#### Exit Codes

- `0` — Success
- `1` — Error (e.g., blueprint already exists)

#### Errors

| Error Message | Cause | Solution |
|---------------|-------|----------|
| `blueprint already exists: <name>` | Duplicate name | Use different name or remove existing |
| `failed to create directory` | File system error | Check permissions |
| `missing blueprint name` | No name provided | Provide a name as argument |

---

### neev bridge

Aggregate context from all blueprints and foundation.

#### Syntax

```bash
neev bridge [flags]
```

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--focus` | `-f` | string | — | Filter by keyword |
| `--claude` | — | boolean | false | Claude-optimized output |
| `--slash` | — | boolean | false | Format for IDE slash commands |
| `--with-remotes` | — | boolean | false | Include synced remotes |

#### Description

Aggregates all markdown files from:
1. `.neev/foundation/` — Project-wide docs
2. `.neev/blueprints/*/` — All blueprints
3. `.neev/remotes/` — Synced remote foundations (if `--with-remotes`)

Output is written to stdout with section headers.

#### Examples

```bash
# Get all context
neev bridge

# Filter by keyword
neev bridge --focus "auth"
neev bridge -f "database"

# Claude-optimized format
neev bridge --claude

# Include remote foundations
neev bridge --with-remotes

# Save to file
neev bridge > context.md

# Combine flags
neev bridge --focus "api" --claude --with-remotes > api-context.md

# Pipe to clipboard (macOS)
neev bridge | pbcopy

# Pipe to clipboard (Linux)
neev bridge | xclip -selection clipboard
```

#### Output Format

```markdown
# Project Foundation
## File: stack.md
[content of stack.md]

## File: principles.md
[content of principles.md]

## File: intent.md
[content from blueprint 1]

## File: architecture.md
[content from blueprint 1]
...
```

#### Focus Filtering

When using `--focus <keyword>`:
- Only includes files containing the keyword
- Case-sensitive matching
- Partial word matches

```bash
# Matches: "authentication", "auth", "Authentication"
neev bridge --focus "auth"

# No matches if keyword not found
neev bridge --focus "xyz"  # May output only headers
```

#### Claude Optimization

When using `--claude`:
- Adds XML-style tags for better parsing
- Optimizes token usage
- Improves Claude's context understanding

#### Exit Codes

- `0` — Success (even if no content found)
- `1` — Error (e.g., cannot read directories)

#### Errors

| Error Message | Cause | Solution |
|---------------|-------|----------|
| `failed to read blueprints` | Missing directory | Run `neev init` first |
| `failed to read file` | Permission error | Check file permissions |

---

### neev inspect

Check if project structure matches foundation specifications.

#### Syntax

```bash
neev inspect [flags]
```

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--json` | boolean | false | Output in JSON format |
| `--strict` | boolean | false | Exit 1 if drift detected |
| `--use-descriptors` | boolean | false | Use `.module.yaml` files |

#### Description

Analyzes your project to detect:
- **Missing modules** — Blueprints without implementation
- **Extra code** — Code without blueprints
- **Drift** — Misalignment between specs and code

#### Examples

```bash
# Human-readable output
neev inspect

# JSON output (for CI/CD)
neev inspect --json

# Strict mode (fail on drift)
neev inspect --strict

# Use module descriptors
neev inspect --use-descriptors

# Combine flags
neev inspect --json --strict > drift-report.json
```

#### Output (Human-Readable)

```
🔍 Inspecting project structure...

✅ Found 3 blueprints:
  - user-authentication
  - api-gateway
  - database-schema

⚠️  Potential drift detected:
  - Missing implementation: user-authentication
  - Extra code: src/legacy/

💡 Tip: Use --use-descriptors for detailed validation
```

#### Output (JSON)

```json
{
  "blueprints": [
    {
      "name": "user-authentication",
      "path": ".neev/blueprints/user-authentication",
      "status": "missing_implementation"
    }
  ],
  "drift": {
    "missing": ["user-authentication"],
    "extra": ["src/legacy"]
  },
  "has_drift": true
}
```

#### Module Descriptors

Create `.module.yaml` in blueprint directories to define expectations:

```yaml
# .neev/blueprints/user-authentication/.module.yaml
expected_files:
  - src/auth/login.js
  - src/auth/register.js
  - tests/auth.test.js
expected_patterns:
  - src/auth/**/*.js
```

#### Exit Codes

- `0` — No drift (or drift found without `--strict`)
- `1` — Drift found with `--strict` flag

#### CI/CD Integration

```yaml
# GitHub Actions example
- name: Check specification drift
  run: |
    neev inspect --json --strict > drift-report.json
    cat drift-report.json
```

---

### neev instructions

Generate GitHub Copilot instructions from your foundation and blueprints.

#### Syntax

```bash
neev instructions [flags]
```

#### Flags

*None currently*

#### Description

Creates or updates `.github/copilot-instructions.md` with:
- Foundation module summaries
- Active blueprint intents
- Development guidelines

This helps GitHub Copilot provide better, context-aware suggestions.

#### Examples

```bash
# Generate instructions
neev instructions

# Check what was created
cat .github/copilot-instructions.md
```

#### Output

```
✅ GitHub Copilot instructions generated at .github/copilot-instructions.md
```

#### Generated File Structure

```markdown
# GitHub Copilot Instructions

## Project Context

[Summary of foundation modules]

## Active Blueprints

### user-authentication
[Intent summary]

### api-gateway
[Intent summary]

## Development Guidelines

[From foundation documents]
```

#### Exit Codes

- `0` — Success
- `1` — Error (e.g., cannot create `.github/` directory)

#### Notes

- Run after creating/updating blueprints
- Commit `.github/copilot-instructions.md` to version control
- Copilot automatically reads this file

---

### neev sync-remotes

Synchronize remote foundation sources from other repositories.

#### Syntax

```bash
neev sync-remotes [flags]
```

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--json` | boolean | false | Output in JSON format |

#### Description

Syncs remote foundations defined in `neev.yaml` into `.neev/remotes/`.

Useful for:
- Multi-repo projects
- Shared foundation documents
- Cross-team specifications

#### Configuration

Define remotes in `neev.yaml`:

```yaml
remotes:
  - name: backend-api
    path: "../backend/.neev/foundation"
    public_only: true
  - name: shared-services
    path: "../../shared/.neev/foundation"
    public_only: false
```

#### Examples

```bash
# Sync all remotes
neev sync-remotes

# JSON output
neev sync-remotes --json

# Check synced files
ls -la .neev/remotes/
```

#### Output

```
🔄 Syncing remote foundations...
✅ Synced backend-api (5 files)
✅ Synced shared-services (3 files)
```

#### Exit Codes

- `0` — Success
- `1` — Error (e.g., remote path not found)

#### Notes

- Synced files are cached locally
- Re-run to update from remotes
- Use `neev bridge --with-remotes` to include in context

---

### neev lay

Archive a completed blueprint into the foundation.

#### Syntax

```bash
neev lay <blueprint-name> [flags]
```

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `blueprint-name` | Yes | Name of the blueprint to archive |

#### Flags

*None currently*

#### Description

Moves a completed blueprint from `blueprints/` to foundation archive and updates changelog.

Useful for:
- Marking work as complete
- Keeping blueprints directory clean
- Maintaining project history

#### Examples

```bash
# Archive a blueprint
neev lay "user-authentication"

# Blueprint is moved and changelog updated
```

#### Output

```
✅ Blueprint 'user-authentication' laid into foundation
📝 Changelog updated
```

#### Exit Codes

- `0` — Success
- `1` — Error (e.g., blueprint not found)

---

### neev migrate

Convert existing projects from other specification formats to Neev.

#### Syntax

```bash
neev migrate [flags]
```

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--source` | `-s` | string | `auto` | Source type: `openspec`, `speckit`, or `auto` |
| `--dry-run` | — | boolean | false | Preview changes without applying |
| `--backup` | — | boolean | false | Backup existing `.neev` directory |

#### Description

Converts existing specification structures to Neev format.

Supports:
- **OpenSpec** — Converts OpenSpec structure
- **SpecKit** — Converts SpecKit structure
- **Auto** — Detects source type automatically

#### Examples

```bash
# Auto-detect and migrate
neev migrate

# Specify source type
neev migrate --source openspec
neev migrate -s speckit

# Dry run (preview only)
neev migrate --dry-run

# Create backup before migrating
neev migrate --backup

# Combine flags
neev migrate --source openspec --backup --dry-run
```

#### Output

```
🔄 Migrating from OpenSpec to Neev...
✅ Converted 5 specifications
✅ Migration complete
```

#### Exit Codes

- `0` — Success
- `1` — Error (e.g., source not detected)

---

### neev handoff

Create a structured handoff prompt with role-specific instructions.

#### Syntax

```bash
neev handoff <role> [flags]
```

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `role` | Yes | Role name (must have `.neev/agents/<role>.md`) |

#### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--markdown` | boolean | true | Wrap output in markdown fence |

#### Description

Generates a structured prompt for AI agents using role-specific instructions from `.neev/agents/<role>.md`.

#### Agent File Format

```markdown
# Role: Backend Developer

## Context
You are an expert backend developer...

## Responsibilities
- Implement API endpoints
- Write tests
- Review code

## Guidelines
- Follow REST conventions
- Use async/await
```

#### Examples

```bash
# Create handoff for backend developer
neev handoff backend

# Without markdown wrapping
neev handoff backend --markdown=false

# Save to file
neev handoff frontend > frontend-prompt.md
```

#### Output

```markdown
\`\`\`
# Handoff to: Backend Developer

## Role Context
You are an expert backend developer...

## Project Foundation
[Foundation content]

## Active Blueprints
[Blueprint content]

## Your Responsibilities
- Implement API endpoints
- Write tests
\`\`\`
```

#### Exit Codes

- `0` — Success
- `1` — Error (e.g., role file not found)

---

## Configuration Reference

### neev.yaml

The main configuration file for Neev projects.

#### Location

`.neev/neev.yaml`

#### Full Example

```yaml
# Project identification
project_name: "My Awesome Project"

# Directory configuration
foundation_path: ".neev"

# Directories to ignore during inspection
ignore_dirs:
  - node_modules
  - .git
  - __pycache__
  - dist
  - build
  - vendor
  - .next
  - target

# Remote foundation sources (optional)
remotes:
  - name: backend-api
    path: "../backend/.neev/foundation"
    public_only: true
  - name: shared-standards
    path: "../../shared/.neev/foundation"
    public_only: false

# Version (for future use)
version: "1.0"
```

#### Field Reference

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `project_name` | string | No | Directory name | Display name for project |
| `foundation_path` | string | No | `.neev` | Path to foundation directory |
| `ignore_dirs` | []string | No | Common dirs | Directories to skip |
| `remotes` | []Remote | No | `[]` | Remote foundation sources |
| `version` | string | No | `1.0` | Config version |

#### Remote Configuration

```yaml
remotes:
  - name: string          # Unique identifier
    path: string          # Relative or absolute path
    public_only: boolean  # Filter public content only
```

---

## Exit Codes

All Neev commands use consistent exit codes:

| Code | Meaning | When It Occurs |
|------|---------|----------------|
| `0` | Success | Command completed successfully |
| `1` | Error | Any error occurred during execution |

**Special cases:**
- `neev inspect --strict` — Exit 1 if drift detected
- All other commands — Exit 1 on any error

---

## Environment Variables

Neev currently does not use environment variables for configuration. All configuration is done via `neev.yaml`.

**Future considerations:**
- `NEEV_PATH` — Override default `.neev` location
- `NEEV_LOG_LEVEL` — Control logging verbosity
- `NEEV_NO_COLOR` — Disable colored output

---

## Global Flags

These flags work with all commands:

| Flag | Description | Example |
|------|-------------|---------|
| `-h, --help` | Show help for command | `neev bridge --help` |

---

## Version Information

Check Neev version:

```bash
neev --version
```

Output:
```
neev version X.Y.Z
```

---

## Getting Help

### Command-Specific Help

```bash
neev <command> --help
```

### General Help

```bash
neev --help
neev help
```

### Documentation

- [Getting Started](GETTING_STARTED.md)
- [Usage Guide](USAGE.md)
- [Tutorials](TUTORIALS.md)
- [Best Practices](BEST_PRACTICES.md)
- [FAQ](FAQ.md)

---

**Need more help?** [Open an issue](https://github.com/neev-kit/neev/issues) or [start a discussion](https://github.com/neev-kit/neev/discussions).
