# AI Skills & Tool Integration

Neev automatically detects installed AI tools and generates skills in their native formats.

## Table of Contents

1. [Supported Tools](#supported-tools)
2. [Automatic Detection](#automatic-detection)
3. [Native Directory Structure](#native-directory-structure)
4. [Tool-Specific Adapters](#tool-specific-adapters)
5. [Natural Language Fallback](#natural-language-fallback)
6. [Commands](#commands)
7. [Workflow](#workflow)

## Supported Tools

### Native Integration (Automatic Skills Generation)

| Tool | Format | Directory | Status |
|------|--------|-----------|--------|
| **Cursor** | JSON | `~/.cursor/skills/` | ✅ Full Support |
| **Claude** | Markdown | `~/.claude/skills/` | ✅ Full Support |
| **GitHub Copilot** | Markdown | `~/.copilot/skills/` | ✅ Full Support |
| **Codeium** | JSON | `~/.codeium/skills/` | ✅ Full Support |
| **Supabase** | JSON | `~/.supabase/skills/` | ✅ Full Support |
| **Perplexity AI** | Markdown | `~/.perplexity/skills/` | ✅ Full Support |

### Natural Language Fallback

Tools without native adapters have their skills generated in markdown format in `.neev/skills/` with documentation for manual integration.

## Automatic Detection

Neev automatically detects which AI tools you have installed:

```bash
# See all detected tools
neev detect-tools
```

Detection checks for:
- **Cursor**: `~/.cursor` directory or Cursor.app
- **Claude**: VS Code extension or Claude.app
- **GitHub Copilot**: VS Code `github.copilot` extension
- **Codeium**: VS Code `codeium` extension
- **Supabase**: `~/.supabase` directory
- **Perplexity**: `~/.perplexity` directory

## Native Directory Structure

Skills are generated directly in each tool's native directory:

```
~/.cursor/
├── skills/
│   ├── blueprint-name-1.json
│   ├── blueprint-name-2.json
│   ├── README.md
│   └── INDEX.md

~/.claude/
├── skills/
│   ├── blueprint-name-1.md
│   ├── blueprint-name-2.md
│   ├── README.md
│   └── INDEX.md

~/.copilot/
├── skills/
│   ├── blueprint-name-1.md
│   ├── blueprint-name-2.md
│   ├── README.md
│   └── INDEX.md
```

## Tool-Specific Adapters

### Cursor Adapter

- **Format**: JSON
- **Features**: Native support for skill loading
- **Configuration**: Auto-discovered from `~/.cursor`

```json
{
  "name": "Feature Name",
  "description": "Feature description",
  "type": "command",
  "content": "... skill implementation ...",
  "version": "1.0"
}
```

### Claude Adapter

- **Format**: Markdown
- **Features**: Rich markdown formatting, easy to read and modify
- **Configuration**: Auto-discovered from `~/.claude`

```markdown
# Skill: Feature Name

**Description:** Feature description
**Type:** command
**Version:** 1.0

## Implementation

... skill implementation ...
```

### GitHub Copilot Adapter

- **Format**: Markdown with code blocks
- **Features**: Supports code syntax highlighting
- **Configuration**: Auto-discovered from VS Code extensions

```markdown
# Feature Name

Feature description

## Type: command
## Version: 1.0

### Content

\`\`\`language
... implementation ...
\`\`\`
```

### Codeium Adapter

- **Format**: JSON
- **Features**: Similar to Cursor but Codeium-specific structure
- **Configuration**: Auto-discovered from VS Code extensions

### Fallback Adapters (Supabase, Perplexity)

- **Format**: JSON or Markdown
- **Features**: Automatic translation to tool's native format
- **Fallback**: Natural language markdown if tool not recognized

## Natural Language Fallback

If you're using an AI tool that doesn't have a native adapter:

1. **Skills are generated in `.neev/skills/`**
   ```
   .neev/
   └── skills/
       ├── feature-1.md
       ├── feature-2.md
       ├── README.md
       └── INDEX.md
   ```

2. **Documentation included**
   - Each skill has integration instructions
   - README.md explains how to import manually
   - Full support for natural language integration

3. **Manual Integration Steps**
   ```bash
   # 1. View generated skills
   cat .neev/skills/feature-name.md
   
   # 2. Copy content into your AI tool
   # 3. Test the skill in your environment
   ```

## Commands

### Sync Skills (Main Command)

Regenerate and synchronize skills across all detected tools:

```bash
# Regenerate all skills
neev sync-skills

# Show detailed output
neev sync-skills --verbose
```

**What it does:**
1. ✅ Detects installed AI tools
2. ✅ Loads blueprints from `.neev/blueprints/`
3. ✅ Generates skills in tool-native formats
4. ✅ Creates `.claude/skills/`, `.cursor/skills/`, etc.
5. ✅ Generates fallback documentation for unsupported tools
6. ✅ Creates README and INDEX files

### Detect Tools

View all detected AI tools:

```bash
# List detected tools
neev detect-tools

# Shows:
# - Tool names
# - Installation status
# - Native skill directories
# - Configuration paths
```

### Skills Status

Check the status of generated skills:

```bash
# Show skills status
neev skills-status

# Shows:
# - Which tools have skills
# - Number of available skills
# - Last generation time
# - Missing skills warnings
```

## Workflow

### 1. Create Blueprints

```bash
# Write your feature specification
cat > .neev/blueprints/auth-feature.md << 'EOF'
# Authentication Feature

Implement OAuth2 authentication with...
EOF
```

### 2. Initialize Skills

```bash
# First time setup - Neev auto-detects tools
neev init

# See what tools were detected
neev detect-tools
```

### 3. Generate Skills

```bash
# Generate skills from blueprints
neev sync-skills
```

### 4. Use in Your AI Tool

**For Cursor:**
- Skills automatically appear in `.cursor/skills/`
- Load directly in Cursor IDE

**For Claude:**
- Skills appear in `~/.claude/skills/`
- Claude automatically discovers them

**For GitHub Copilot:**
- Skills in `.copilot/skills/`
- Load via Copilot prompts

**For Unsupported Tools:**
- Open `.neev/skills/` 
- Copy content into your tool
- Follow integration documentation

### 5. Keep Skills Updated

When you modify blueprints, regenerate:

```bash
# Update blueprints
cat >> .neev/blueprints/auth-feature.md << 'EOF'
## Updated Requirements
EOF

# Regenerate skills
neev sync-skills
```

## Configuration

### Tool Detection

Neev uses environment variables and standard tool locations:

```bash
# Check detection
export HOME=$HOME  # Uses your home directory

# Standard paths checked:
# - ~/.cursor/
# - ~/.claude/
# - ~/.vscode/extensions/
# - ~/.codeium/
# - ~/.supabase/
# - ~/.perplexity/
```

### Custom Directories

If your tools are installed in custom locations, set environment variables:

```bash
export CLAUDE_HOME=/custom/path/.claude
export CURSOR_HOME=/custom/path/.cursor
neev sync-skills
```

## Troubleshooting

### Tools Not Detected

```bash
# Check detection directly
neev detect-tools

# Verify tool installation
ls -la ~/.cursor/    # For Cursor
ls -la ~/.claude/    # For Claude
ls -la ~/.vscode/extensions/ | grep copilot  # For Copilot
```

### Skills Not Generating

```bash
# Check blueprint directory
ls -la .neev/blueprints/

# Verify blueprint format (must be .md files)
file .neev/blueprints/*.md

# Regenerate with verbose output
neev sync-skills --verbose
```

### Missing Native Directory

```bash
# If ~/.cursor/skills/ doesn't exist:
mkdir -p ~/.cursor/skills/
mkdir -p ~/.claude/skills/
neev sync-skills
```

## Best Practices

1. **Keep Blueprints Updated**
   - Change blueprints → run `neev sync-skills`
   - Skills stay in sync with specifications

2. **Use Natural Language Fallback**
   - All skills available in `.neev/skills/` as markdown
   - Copy-paste friendly for any AI tool

3. **Version Control Skills**
   - `.neev/` is typically committed to git
   - Tool-specific directories are usually ignored
   - Skills are regenerated locally per tool

4. **Share Blueprints, Not Generated Skills**
   - Commit `.neev/blueprints/` to git
   - Each developer runs `neev sync-skills`
   - Each tool gets its own format automatically

5. **Review Generated Skills**
   - Check README.md in each tool's skills directory
   - Verify format matches tool's expectations
   - Customize if needed for your environment

## Examples

### Example 1: Multi-Tool Setup

```bash
# Initial setup
neev init
# ✓ Found Claude, Cursor, GitHub Copilot

# Generate for all tools
neev sync-skills
# ✓ Generated for Claude (~/.claude/skills/)
# ✓ Generated for Cursor (~/.cursor/skills/)
# ✓ Generated for Copilot (~/.copilot/skills/)

# Each tool gets format:
# - Claude: Markdown
# - Cursor: JSON
# - Copilot: Markdown
```

### Example 2: Using Fallback

```bash
# Using Supabase (fallback adapter)
neev sync-skills
# ✓ Generated for Supabase (~/.supabase/skills/)
# ✓ Generated fallback docs (.neev/skills/)

# For manual integration:
cat .neev/skills/my-feature.md
# Copy content into Supabase
```

### Example 3: Update Workflow

```bash
# Modify blueprint
vi .neev/blueprints/auth.md

# Regenerate all tools at once
neev sync-skills
# ✓ Updated Claude skills
# ✓ Updated Cursor skills
# ✓ Updated Copilot skills
```

## Related Commands

- `neev init` - Initialize and detect tools
- `neev bridge` - Generate project context
- `neev inspect` - Check blueprint drift
- `neev draft` - Create new blueprints

## Support

For issues or feature requests:
- 📖 See [USAGE.md](../USAGE.md#skill-generation)
- 🐛 Report at [GitHub Issues](https://github.com/neev-kit/neev/issues)
- 💬 Discuss in [Discussions](https://github.com/neev-kit/neev/discussions)
