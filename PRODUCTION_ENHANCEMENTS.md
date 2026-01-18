# Neev Production Enhancements - Phase 1-4 Implementation

This document describes the production-grade enhancements added to Neev across four key axes.

## Overview

Neev has been enhanced from a functional prototype into a production-ready hybrid spec-driven development framework. The enhancements include:

1. **Richer Drift Detection** - Structured, machine-readable drift analysis
2. **Polyrepo/Remote Foundations** - Support for cross-repository foundation sharing
3. **AI Tooling Integration** - Native integration with GitHub Copilot and Claude

## 1. Richer Drift Detection

### What Changed
- New `core/inspect` package with structured types and categorized warnings
- Support for module descriptors (`.module.yaml`) for detailed file-level inspection
- JSON output for CI/CD integration

### New Commands
```bash
# Structured inspection with detailed output
neev inspect

# JSON output for programmatic consumption
neev inspect --json

# File-level inspection using module descriptors
neev inspect --use-descriptors
```

### Module Descriptors
Create `.neev/foundation/<module>.module.yaml` files to define expected structure:

```yaml
name: auth
description: Authentication module
expected_files:
  - handler.go
  - service.go
expected_dirs:
  - models
patterns:
  - "*.go"
  - "**/*_test.go"
```

### Warning Categories
- `MISSING_MODULE` - Foundation spec exists but no code directory
- `EXTRA_CODE` - Code exists without foundation spec
- `MISSING_FILE` - Expected file from descriptor is missing
- `MISMATCHED_NAME` - Naming inconsistency detected

## 2. Polyrepo/Remote Foundations

### What Changed
- Extended `neev.yaml` to support remote foundation sources
- New `neev sync-remotes` command to sync external foundations
- Bridge command can include remote contexts with `--with-remotes`

### Configuration
Add remotes to `neev.yaml`:

```yaml
project_name: My UI Project
foundation_path: .neev
remotes:
  - name: api
    path: "../my-api-repo/.neev/foundation"
    public_only: true
  - name: shared
    path: "../shared-lib/.neev/foundation"
    public_only: false
```

### New Commands
```bash
# Sync all configured remotes
neev sync-remotes

# Sync remotes with JSON output
neev sync-remotes --json

# Include remotes in bridge context
neev bridge --with-remotes
```

### Remote Sync Behavior
- Copies markdown files from remote foundations to `.neev/remotes/<name>/`
- Skips `archive/` subdirectories
- When `public_only: true`, skips files starting with `_`
- Synced remotes are included in bridge context when `--with-remotes` is used

## 3. AI Tooling Integration

### What Changed
- New `core/instructions` package for generating AI assistant instructions
- Copilot instructions generation command
- Claude-optimized bridge output format

### GitHub Copilot Integration
```bash
# Generate Copilot instructions
neev instructions
```

This creates `.github/copilot-instructions.md` with:
- Foundation module summary
- Active blueprint intents
- Development guidelines

### Claude Integration
```bash
# Get Claude-optimized context
neev bridge --claude

# With remotes included
neev bridge --claude --with-remotes
```

Claude mode adds:
- Explicit "RULES AND CONSTRAINTS" section
- Clear section markers with emojis
- "CURRENT TASK" summary
- Enhanced formatting for Claude's understanding

## Backwards Compatibility

All existing commands continue to work unchanged:
- `neev init` - Initialize project
- `neev draft` - Create blueprint
- `neev bridge` - Aggregate context (default mode unchanged)
- `neev inspect` - Drift detection (legacy mode as default)
- `neev lay` - Archive blueprint

## Usage Examples

### Example 1: Advanced Drift Detection
```bash
# Create a module descriptor
cat > .neev/foundation/auth.module.yaml << EOF
name: auth
description: Authentication module
expected_files:
  - handler.go
  - middleware.go
patterns:
  - "*.go"
EOF

# Run inspection with descriptors
neev inspect --use-descriptors

# Get JSON for CI pipeline
neev inspect --json > drift-report.json
```

### Example 2: Polyrepo Setup
```bash
# Configure remotes in neev.yaml
cat >> neev.yaml << EOF
remotes:
  - name: api
    path: "../backend/.neev/foundation"
    public_only: true
EOF

# Sync remotes
neev sync-remotes

# Use in bridge
neev bridge --with-remotes > full-context.md
```

### Example 3: AI Assistant Workflow
```bash
# Generate Copilot instructions
neev instructions

# Get Claude-formatted context
neev bridge --claude > claude-context.md

# Include remotes for cross-repo context
neev bridge --claude --with-remotes > full-claude-context.md
```

## Package Structure

### New Packages
- `core/inspect` - Structured drift detection
- `core/remotes` - Remote foundation synchronization
- `core/instructions` - AI assistant instruction generation

### Enhanced Packages
- `core/config` - Extended with Remote type
- `core/bridge` - Added remote context and Claude formatting
- `cli/cmd` - New commands: inspect (enhanced), sync-remotes, instructions

## Testing

All new functionality is covered by unit tests:
- `core/inspect/inspect_test.go` - Drift detection tests
- `core/remotes/sync_test.go` - Remote sync tests
- `core/instructions/copilot_test.go` - Instructions generation tests

Run tests:
```bash
go test ./core/... ./cli/... -cover
```

## Implementation Notes

### Design Decisions
1. **No breaking changes** - All existing commands work as before
2. **Opt-in features** - Advanced features require explicit flags
3. **Composable** - Features work independently or together
4. **Fast and local** - No external dependencies or API calls

### Performance
- Inspect: O(n) where n = number of modules
- Sync: O(m) where m = number of files in remotes
- Instructions: O(b) where b = number of blueprints

All operations are fast enough for interactive use.

## Future Enhancements

While Axis 2 (Spec-Driven Tests) was designed but not implemented in this phase, the foundation is ready for:
- `neev test-gen` command
- Test scenario scaffolding
- Language-agnostic test templates

This can be added in a future phase without affecting existing functionality.
