# Neev Usage Guide

A comprehensive guide to using Neev CLI for project blueprint management and AI context aggregation.

## Table of Contents

1. [Installation](#installation)
2. [Basic Workflow](#basic-workflow)
3. [Command Reference](#command-reference)
4. [Examples](#examples)
5. [Best Practices](#best-practices)
6. [Troubleshooting](#troubleshooting)

## Installation

### Option 1: Download Pre-built Binary (Recommended)

Download the latest stable release from [GitHub Releases](https://github.com/neev-kit/neev/releases):

```bash
# macOS (Intel)
curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_darwin_amd64.tar.gz | tar xz
sudo mv neev /usr/local/bin/

# macOS (Apple Silicon/M1/M2)
curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_darwin_arm64.tar.gz | tar xz
sudo mv neev /usr/local/bin/

# Linux (x86_64)
curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_linux_amd64.tar.gz | tar xz
sudo mv neev /usr/local/bin/

# Windows: Visit https://github.com/neev-kit/neev/releases and download the .zip file
```

### Option 2: Build from Source

```bash
git clone https://github.com/neev-kit/neev
cd neev
go mod download
go build -o neev ./cli
```

### Option 3: Using Go Install

```bash
go install github.com/neev-kit/neev/cli@latest
```

The binary will be installed to `$GOPATH/bin/neev`. Ensure `$GOPATH/bin` is in your `$PATH`.

### Verify Installation

```bash
neev --version
neev --help
```

## Basic Workflow

### Step 1: Initialize Your Project

Start by initializing Neev in your project directory:

```bash
cd /path/to/your/project
neev init
```

This creates:
- `.neev/` - Root configuration directory
- `.neev/blueprints/` - Blueprint storage
- `.neev/foundation/` - Foundation specifications
- `.neev/neev.yaml` - Configuration file

### Step 2: Create Blueprints

Create one or more blueprints for different aspects of your project:

```bash
# Create architecture blueprint
neev draft "Architecture"

# Create API blueprint
neev draft "API Design"

# Create database blueprint
neev draft "Database Schema"
```

Each blueprint includes:
- `intent.md` - Purpose of this component
- `architecture.md` - Technical details

### Step 3: Document Your Blueprints

Edit the generated blueprint files with your project details:

```
.neev/blueprints/architecture/
├── intent.md          # What this blueprint is for
└── architecture.md    # Technical architecture

.neev/blueprints/api-design/
├── intent.md          # API scope and goals
└── architecture.md    # API endpoints and specs

.neev/blueprints/database-schema/
├── intent.md          # Database purpose
└── architecture.md    # Schema design
```

### Step 4: Aggregate Context

Generate full project context for AI agents or documentation:

```bash
# Get all project context
neev bridge > context.md

# Get context focused on specific topic
neev bridge --focus "database" > database-context.md

# Get context for API-related components
neev bridge -f "API" > api-context.md
```

## Command Reference

### neev init

Initialize a Neev project foundation.

**Syntax:**
```bash
neev init [flags]
```

**Flags:**
- None (currently)

**Examples:**
```bash
neev init
```

**Output:**
```
🏗️  Laying foundation in /Users/username/project
✅ Foundation laid successfully!
```

**Errors:**
- `.neev directory already exists` - Cannot reinitialize existing projects
- `permission denied` - Insufficient permissions to create directory

---

### neev draft

Create a new blueprint with template files.

**Syntax:**
```bash
neev draft <blueprint-name> [flags]
```

**Arguments:**
- `<blueprint-name>` (required) - Name of the blueprint (spaces allowed, will be sanitized)

**Flags:**
- None currently, but reserved for future options

**Naming Rules:**
- Spaces are converted to hyphens
- Uppercase converted to lowercase
- Example: `"My API Service"` → `my-api-service`

**Examples:**
```bash
# Simple name
neev draft "Authentication"
# Result: .neev/blueprints/authentication/

# Multi-word name
neev draft "User Management Service"
# Result: .neev/blueprints/user-management-service/

# With special characters (sanitized)
neev draft "API v2.0"
# Result: .neev/blueprints/api-v2.0/
```

**Output:**
```
✅ Created blueprint at .neev/blueprints/authentication
```

**Errors:**
- `blueprint already exists: ...` - Blueprint name conflicts with existing blueprint
- `failed to create blueprint directory: ...` - File system error
- `failed to create file ...: ...` - Unable to create template files

**Generated Files:**

```markdown
# intent.md
Template for intent.md

# architecture.md
Template for architecture.md
```

---

### neev bridge

Aggregate context from all blueprints and foundation.

**Syntax:**
```bash
neev bridge [flags]
```

**Flags:**
- `--focus, -f` (string) - Filter context by keyword (optional)

**Examples:**
```bash
# Get all project context
neev bridge

# Filter for database-related content
neev bridge --focus "database"
neev bridge -f "database"

# Filter for authentication
neev bridge --focus "auth"

# Pipe to file
neev bridge > full-context.md

# Pipe to AI tool
neev bridge | some-ai-tool
```

**Output Format:**

The output includes:
1. Project Foundation header
2. All foundation files (from `.neev/foundation/`)
3. All blueprint files (from `.neev/blueprints/*/`)
4. Each file is prefixed with a section header

Example output:
```
# Project Foundation
## File: neev.yaml
version: "1.0"
name: "my-project"

## File: intent.md
...

## File: architecture.md
...
```

**Focus Filtering:**

When using `--focus`, only files containing that keyword are included:

```bash
# Only includes files mentioning "cache"
neev bridge --focus "cache"

# Case-sensitive matching
neev bridge --focus "Cache"  # Different from above
```

**Errors:**
- `failed to read blueprints directory: ...` - Missing `.neev/blueprints/`
- `failed to read directory ...: ...` - Permission or access error
- `failed to read file ...: ...` - Cannot read specific file

---

## Examples

### Example 1: E-Commerce Project Setup

```bash
# Initialize
cd ~/projects/ecommerce
neev init

# Create domain blueprints
neev draft "Product Catalog"
neev draft "Shopping Cart"
neev draft "Payment Processing"
neev draft "User Authentication"

# Edit blueprints with your specifications
# (Edit .neev/blueprints/*/intent.md and architecture.md)

# Aggregate for documentation
neev bridge > docs/system-architecture.md

# Get payment-related context
neev bridge --focus "Payment" > docs/payment-spec.md
```

### Example 2: Using with AI Tools

```bash
# Generate context for Claude/ChatGPT
neev bridge > ai-context.txt

# Use with LLM via piping
neev bridge | llm prompt "Based on this project structure, suggest improvements"

# Focus on specific area for AI review
neev bridge --focus "authentication" | ai-code-review
```

### Example 3: Team Documentation

```bash
# Core architecture
neev draft "System Architecture"
neev draft "Data Models"

# API specifications
neev draft "REST API"
neev draft "WebSocket Endpoints"

# Operations
neev draft "Deployment Process"
neev draft "Monitoring & Alerts"

# Generate handbook
neev bridge > TECHNICAL_HANDBOOK.md
```

### Example 4: Microservices Setup

```bash
neev init

# Service blueprints
neev draft "User Service"
neev draft "Product Service"
neev draft "Order Service"
neev draft "Payment Service"

# Cross-cutting concerns
neev draft "API Gateway"
neev draft "Message Queue"
neev draft "Shared Database"

# Get architecture overview
neev bridge

# Get service-specific context
neev bridge --focus "User Service"
```

## Best Practices

### 1. Naming Conventions

Use clear, descriptive names that reflect the blueprint's purpose:

```bash
# Good
neev draft "Authentication Service"
neev draft "Database Schema"
neev draft "API Specification"

# Avoid
neev draft "thing"
neev draft "stuff"
neev draft "temp"
```

### 2. Documentation Standards

Keep `intent.md` and `architecture.md` concise but comprehensive:

**intent.md Template:**
```markdown
# Intent

## Purpose
Describe what this blueprint is for.

## Goals
- Goal 1
- Goal 2
- Goal 3

## Out of Scope
What this blueprint does NOT cover.
```

**architecture.md Template:**
```markdown
# Architecture

## Overview
High-level diagram or description.

## Components
- Component A
- Component B

## Design Decisions
Explain key architectural choices.

## Dependencies
Internal and external dependencies.

## Data Flow
Describe how data moves through this component.
```

### 3. Organization Strategy

Group related blueprints logically:

```
Domain-driven:
- User Management
- Product Catalog
- Order Processing

Layer-based:
- Database Layer
- API Layer
- Service Layer

Feature-based:
- Authentication
- Payment Processing
- Notifications
```

### 4. Version Control

Commit your blueprints:

```bash
git add .neev/
git commit -m "Add initial project blueprints"
```

This ensures your architectural decisions are tracked historically.

### 5. Focus Keywords

Use consistent keywords in your blueprints for easier filtering:

```
# In intent.md and architecture.md, use keywords like:
- [DATABASE]
- [API]
- [SECURITY]
- [PERFORMANCE]
- [DEPLOYMENT]
```

Then use focus to filter:
```bash
neev bridge --focus "DATABASE"
neev bridge --focus "SECURITY"
```

### 6. Maintenance

Regularly update blueprints as your project evolves:

```bash
# Review and update quarterly
neev bridge > current-state.md
# Compare with last version to identify changes
```

## Troubleshooting

### Problem: `.neev directory already exists`

**Cause:** Attempting to run `neev init` on an already initialized project.

**Solution:**
```bash
# Check existing .neev/
ls -la .neev/

# If needed, backup then remove
mv .neev/ .neev.backup
neev init
```

### Problem: `blueprint already exists: ...`

**Cause:** Trying to create a blueprint with the same name.

**Solution:**
```bash
# List existing blueprints
ls -la .neev/blueprints/

# Use a different name
neev draft "Authentication v2"

# Or remove the old one
rm -rf .neev/blueprints/authentication
neev draft "Authentication"
```

### Problem: `failed to read blueprints directory`

**Cause:** Missing or corrupted `.neev/blueprints/` directory.

**Solution:**
```bash
# Recreate the structure
mkdir -p .neev/blueprints
mkdir -p .neev/foundation

# Or reinitialize (backup first)
mv .neev/ .neev.backup
neev init
```

### Problem: Empty output from `neev bridge`

**Cause:** No markdown files in blueprints or foundation directories.

**Solution:**
```bash
# Create a blueprint first
neev draft "Sample"

# Check the structure
find .neev -type f -name "*.md"

# Run bridge again
neev bridge
```

### Problem: Focus filter returns nothing

**Cause:** Keyword doesn't exist in any files.

**Solution:**
```bash
# Check what's in your files
neev bridge | grep -i "your-keyword"

# Use a different focus term
neev bridge --focus "different-keyword"
```

## Advanced Usage

### Integrating with Documentation Tools

```bash
# Generate with MkDocs
neev bridge > docs/generated-context.md

# Convert to different format
neev bridge | pandoc -f markdown -t pdf -o architecture.pdf

# HTML output
neev bridge | pandoc -f markdown -t html -o architecture.html
```

### CI/CD Integration

```bash
# GitHub Actions example
- name: Generate architecture context
  run: |
    neev bridge > ./docs/current-architecture.md
    git add ./docs/current-architecture.md
    git commit -m "Auto-update architecture documentation"
```

### Git Hooks

```bash
# Pre-commit: ensure blueprints are documented
#!/bin/bash
# .git/hooks/pre-commit

if [ -d ".neev/blueprints" ]; then
  neev bridge > /tmp/context.md
  echo "✓ Project context generated successfully"
fi
```

## Getting Help

```bash
# General help
neev --help

# Command-specific help
neev init --help
neev draft --help
neev bridge --help

# Verbose output (if supported)
neev -v init
```

---

For more information, visit: https://github.com/neev-kit/neev
