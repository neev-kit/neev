# Quick Reference: Skills System

## Commands

### Detect Tools
```bash
neev detect-tools
```
Lists all detected AI tools and their configuration directories.

### Sync Skills
```bash
neev sync-skills
```
Regenerates skills from blueprints for all detected tools.

### Check Status
```bash
neev skills-status
```
Shows which tools have generated skills and count of available skills.

## Supported Tools

| Tool | Format | Location | Status |
|------|--------|----------|--------|
| Cursor | JSON | `~/.cursor/skills/` | ✅ Full |
| Claude | Markdown | `~/.claude/skills/` | ✅ Full |
| GitHub Copilot | Markdown | `~/.copilot/skills/` | ✅ Full |
| Codeium | JSON | `~/.codeium/skills/` | ✅ Full |
| Supabase | JSON | `~/.supabase/skills/` | ✅ Full |
| Perplexity | Markdown | `~/.perplexity/skills/` | ✅ Full |
| Any other tool | Markdown | `.neev/skills/` | ✅ Fallback |

## Workflow

```
1. Create blueprints
   mkdir -p .neev/blueprints/
   cat > .neev/blueprints/feature.md << EOF
   # Feature Name
   Feature description...
   EOF

2. Detect tools (automatic)
   neev init

3. Generate skills
   neev sync-skills

4. Use in your tool
   - Cursor: ~/.cursor/skills/
   - Claude: ~/.claude/skills/
   - Copilot: ~/.copilot/skills/
   - Other: .neev/skills/ (copy manually)

5. Update when blueprints change
   neev sync-skills
```

## Directory Structure

**Blueprints (source):**
```
.neev/blueprints/
├── auth-feature.md
├── api-design.md
└── database-schema.md
```

**Generated Skills:**
```
~/.cursor/skills/          (JSON format)
~/.claude/skills/          (Markdown format)
~/.copilot/skills/         (Markdown format)
~/.codeium/skills/         (JSON format)
.neev/skills/              (Markdown fallback)
```

## File Formats

### Cursor (JSON)
```json
{
  "name": "Feature Name",
  "description": "Description",
  "type": "command",
  "content": "...implementation...",
  "version": "1.0"
}
```

### Claude/Copilot (Markdown)
```markdown
# Feature Name

**Description:** Description here
**Type:** command
**Version:** 1.0

## Implementation

...implementation...
```

## Troubleshooting

**Tools not detected?**
```bash
neev detect-tools  # See what's found
echo $HOME         # Verify HOME is set
ls ~/.cursor/      # Check Cursor installation
ls ~/.claude/      # Check Claude installation
```

**Skills not generating?**
```bash
ls -la .neev/blueprints/  # Verify blueprints exist
neev sync-skills --verbose # See detailed output
mkdir -p ~/.cursor/skills/ # Ensure directory exists
```

**Using unsupported tool?**
```bash
cat .neev/skills/skill-name.md  # Read fallback
# Copy content into your tool manually
```

## Integration

**In VS Code:**
1. Claude Extension → Load from `~/.claude/skills/`
2. Copilot → Loads from `~/.copilot/skills/`

**In Cursor:**
1. Automatically loads from `~/.cursor/skills/`

**In Other Tools:**
1. Open `.neev/skills/`
2. Copy markdown content
3. Paste into tool's skill editor

## Examples

### Example 1: Multi-Tool Setup
```bash
neev init
# ✓ Found Claude, Cursor, Copilot

neev sync-skills
# ✓ Generated Claude skills
# ✓ Generated Cursor skills
# ✓ Generated Copilot skills
```

### Example 2: Fallback Usage
```bash
neev sync-skills
# No tools found, using fallback

ls .neev/skills/
# skill-1.md
# skill-2.md
# README.md
```

### Example 3: Update Workflow
```bash
# Modify blueprint
vim .neev/blueprints/auth.md

# Regenerate
neev sync-skills

# All tools updated automatically
```

## Common Tasks

**View all detected tools:**
```bash
neev detect-tools
```

**Regenerate all skills:**
```bash
neev sync-skills
```

**Check skill status:**
```bash
neev skills-status
```

**View a generated skill:**
```bash
cat ~/.cursor/skills/skill-name.json
cat ~/.claude/skills/skill-name.md
cat .neev/skills/skill-name.md
```

**Remove and regenerate:**
```bash
rm -rf ~/.cursor/skills
rm -rf ~/.claude/skills
rm -rf ~/.copilot/skills
neev sync-skills
```

## Environment Variables

**Custom tool locations:**
```bash
export CLAUDE_HOME=/custom/path/.claude
export CURSOR_HOME=/custom/path/.cursor
neev sync-skills
```

## Files Included

- `SKILLS.md` - Full user guide
- `SKILLS_IMPLEMENTATION.md` - Developer guide
- `core/tools/detect.go` - Tool detection
- `core/tools/adapters.go` - Format adapters
- `core/tools/generator.go` - Skills generation
- `cli/cmd/sync_skills.go` - CLI commands

## Support

- 📖 See [SKILLS.md](SKILLS.md) for detailed guide
- 🔧 See [SKILLS_IMPLEMENTATION.md](SKILLS_IMPLEMENTATION.md) for technical details
- 🐛 Open issue at [GitHub](https://github.com/neev-kit/neev/issues)
