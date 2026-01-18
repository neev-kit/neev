# ❓ Neev FAQ

Frequently asked questions and comprehensive troubleshooting for Neev.

## Table of Contents

1. [General Questions](#general-questions)
2. [Installation & Setup](#installation--setup)
3. [Usage Questions](#usage-questions)
4. [AI Integration](#ai-integration)
5. [Team & Collaboration](#team--collaboration)
6. [Technical Questions](#technical-questions)
7. [Troubleshooting](#troubleshooting)
8. [Best Practices](#best-practices)
9. [Advanced Topics](#advanced-topics)

---

## General Questions

### What is Neev?

Neev is a lightweight CLI tool that helps you document project specifications and generate AI-ready context. It acts as a bridge between your project's intent (what you want to build) and AI coding assistants (Claude, Copilot, Cursor, etc.).

### Why should I use Neev?

**Problems Neev solves:**
- ❌ Scattered project documentation
- ❌ Lost architectural decisions
- ❌ Poor AI code suggestions (lack of context)
- ❌ Team misalignment on features
- ❌ Specs getting out of sync with code

**Benefits:**
- ✅ Centralized, version-controlled specs
- ✅ AI-ready context aggregation
- ✅ Better team alignment
- ✅ Drift detection between specs and code
- ✅ Zero external dependencies

### How is Neev different from traditional documentation?

| Traditional Docs | Neev |
|------------------|------|
| Separate from code | Lives alongside code in `.neev/` |
| Often gets outdated | Drift detection keeps it current |
| Not AI-optimized | Designed for AI consumption |
| Unstructured | Consistent structure (blueprints + foundation) |
| Hard to aggregate | One command to aggregate all context |

### Do I need to use AI to benefit from Neev?

No! Neev is valuable even without AI:
- **Documentation** — Structured project specs
- **Team alignment** — Clear feature definitions
- **Onboarding** — New developers get complete context
- **Planning** — Document before building

AI integration is a bonus feature, not a requirement.

### What languages/frameworks does Neev support?

**All of them!** Neev is language-agnostic. It works with:
- JavaScript/TypeScript, Python, Go, Java, C#, Ruby, PHP, Rust...
- React, Vue, Angular, Django, Rails, Express, Spring Boot...
- Any tech stack or framework

Neev documents **what** you're building, not **how** you're building it.

### Is Neev open source?

Yes! Neev is MIT licensed and available on [GitHub](https://github.com/neev-kit/neev).

---

## Installation & Setup

### How do I install Neev?

**Three options:**

1. **Download binary** (recommended):
   ```bash
   # macOS
   curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_darwin_amd64.tar.gz | tar xz
   sudo mv neev /usr/local/bin/
   
   # Linux
   curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_linux_amd64.tar.gz | tar xz
   sudo mv neev /usr/local/bin/
   ```

2. **Go install**:
   ```bash
   go install github.com/neev-kit/neev/cli@latest
   ```

3. **Build from source**:
   ```bash
   git clone https://github.com/neev-kit/neev.git
   cd neev
   go build -o neev ./cli
   ```

See [GETTING_STARTED.md](GETTING_STARTED.md#installation) for details.

### I get "command not found: neev"

**Cause:** Neev binary is not in your PATH.

**Solutions:**

1. **Check if installed:**
   ```bash
   which neev
   # If empty, neev is not installed
   ```

2. **If using Go install, add $GOPATH/bin to PATH:**
   ```bash
   echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
   source ~/.bashrc
   ```

3. **If downloaded binary, move to PATH location:**
   ```bash
   sudo mv neev /usr/local/bin/
   ```

4. **Verify:**
   ```bash
   neev --version
   ```

### Can I install Neev without Go?

Yes! Download the pre-built binary for your platform from [GitHub Releases](https://github.com/neev-kit/neev/releases). No Go required.

### How do I update Neev?

```bash
# If installed via binary
# Download latest release and replace

# If installed via Go
go install github.com/neev-kit/neev/cli@latest

# Verify new version
neev --version
```

### Can I uninstall Neev?

```bash
# Remove binary
sudo rm /usr/local/bin/neev

# Or if installed via Go
rm $GOPATH/bin/neev

# Remove project .neev directories (optional)
# rm -rf /path/to/project/.neev
```

---

## Usage Questions

### Do I need to commit `.neev/` to version control?

**Yes!** The `.neev/` directory should be version controlled because:
- It contains your project specifications
- Team members need access to blueprints
- Specs should evolve with code
- History of architectural decisions is valuable

```bash
git add .neev/
git commit -m "docs: add project specifications"
```

### Can I use Neev in an existing project?

Absolutely! Run `neev init` in any project directory:

```bash
cd /path/to/existing/project
neev init
```

This creates `.neev/` alongside your existing code. Then document your current architecture and features.

### What should go in foundation vs blueprints?

**Foundation** — Project-wide, stable, shared across features:
- Tech stack choices
- Coding conventions
- Architectural principles
- Development workflow
- Team standards

**Blueprints** — Feature-specific, changes more often:
- Individual feature specs
- API endpoint designs
- Component architectures
- Module designs
- Specific implementations

**Rule of thumb:** If it applies to the whole project, it's foundation. If it's specific to one feature/component, it's a blueprint.

### How many blueprints should I create?

**It depends on your project size:**

- **Small project (1-3 developers):** 5-10 blueprints
- **Medium project (5-10 developers):** 15-30 blueprints
- **Large project (10+ developers):** 30-100 blueprints

**Guidelines:**
- One blueprint per major feature/component
- Split large features into multiple blueprints
- Don't create blueprints for trivial changes
- Quality over quantity

### Can I delete blueprints?

Yes, but consider archiving instead:

```bash
# Archive (recommended)
neev lay "blueprint-name"

# Or manually delete
rm -rf .neev/blueprints/blueprint-name
```

Archiving preserves history and updates changelog.

### Can I have multiple `.neev/` directories?

Not in the same project. Neev looks for `.neev/` in the current directory or parent directories.

For mono-repos with multiple projects:
- Option 1: One `.neev/` at root with all projects
- Option 2: Separate `.neev/` per sub-project (use remotes to share)

### What file formats does Neev support?

Neev uses **Markdown (`.md`)** exclusively. This provides:
- Human-readable documentation
- Version control friendly (text files)
- Tool compatibility (GitHub, IDEs, etc.)
- No vendor lock-in

### Can I customize blueprint templates?

Currently, blueprints always include `intent.md` and `architecture.md`. You can:
- Add additional `.md` files manually
- Customize the generated templates (requires code changes)
- Use your own naming conventions within files

### How do I organize large blueprints?

**Option 1: Multiple files in blueprint directory**
```
.neev/blueprints/api-gateway/
├── intent.md
├── architecture.md
├── security.md       # Additional file
└── deployment.md     # Additional file
```

**Option 2: Split into multiple blueprints**
```
.neev/blueprints/
├── api-gateway-core/
├── api-gateway-auth/
└── api-gateway-rate-limiting/
```

---

## AI Integration

### Which AI tools work with Neev?

Neev works with any AI tool that accepts text input:
- ✅ Claude (Anthropic)
- ✅ ChatGPT (OpenAI)
- ✅ GitHub Copilot
- ✅ Cursor
- ✅ Codeium
- ✅ Tabnine
- ✅ Any custom AI agents

### How do I use Neev with Claude?

```bash
# Generate Claude-optimized context
neev bridge --claude > context.md

# Copy contents
cat context.md | pbcopy  # macOS
cat context.md | xclip   # Linux

# Paste into Claude with your request:
# "Here's my project context: [paste]
#  Now implement the user authentication feature."
```

### How do I use Neev with GitHub Copilot?

```bash
# Generate Copilot instructions
neev instructions

# This creates/updates .github/copilot-instructions.md
# Commit it to your repo
git add .github/copilot-instructions.md
git commit -m "docs: update Copilot instructions"

# Copilot automatically reads this file
```

### Do I need to regenerate context every time?

**When to regenerate:**
- ✅ After creating new blueprints
- ✅ After updating existing blueprints
- ✅ Before starting a new feature
- ✅ When switching between major features

**Not needed for:**
- ❌ Small code changes
- ❌ Bug fixes (unless spec-related)
- ❌ Refactoring (unless architecture changes)

### The AI ignores my specifications. Why?

**Possible causes:**

1. **Too much context** — AI gets overwhelmed
   ```bash
   # Solution: Use --focus
   neev bridge --focus "specific-feature"
   ```

2. **Vague specifications** — Not clear enough
   ```markdown
   # Bad: "Build a good auth system"
   # Good: "Implement JWT-based auth with 1-hour expiry"
   ```

3. **Context not included** — Forgot to paste
   ```bash
   # Always verify you pasted the context
   ```

4. **Conflicting instructions** — Multiple contradictory specs
   ```bash
   # Review blueprints for conflicts
   grep -r "authentication" .neev/
   ```

### Can Neev improve my AI suggestions automatically?

Yes, with GitHub Copilot:

```bash
# Run this after updating specs
neev instructions
git commit -m "docs: update Copilot instructions"
```

Copilot reads `.github/copilot-instructions.md` automatically and improves suggestions based on your project context.

---

## Team & Collaboration

### How do multiple developers use Neev?

1. **One person initializes:**
   ```bash
   neev init
   git add .neev/
   git commit -m "docs: initialize Neev"
   git push
   ```

2. **Others pull and use:**
   ```bash
   git pull
   neev bridge  # Works immediately
   ```

3. **Team members create blueprints:**
   ```bash
   neev draft "My Feature"
   # Edit files
   git add .neev/blueprints/my-feature/
   git commit -m "docs: add my feature spec"
   git push
   ```

### Should blueprints be reviewed like code?

**Yes!** Blueprint reviews prevent:
- Architectural mistakes
- Conflicting designs
- Scope creep
- Implementation issues

**Process:**
1. Create feature branch
2. Add/update blueprint
3. Create pull request
4. Team reviews **spec** before code
5. Merge when approved
6. Then implement

### How do we keep specs in sync with code?

**Use drift detection:**

```bash
# Locally
neev inspect

# In CI/CD
neev inspect --json --strict
```

**Best practices:**
- Update specs when code changes significantly
- Review specs during code reviews
- Run `neev inspect` in CI/CD
- Make spec updates part of your definition of done

### Can different teams have different Neev setups?

Yes, using remotes:

**Backend team's `neev.yaml`:**
```yaml
remotes:
  - name: frontend-contracts
    path: "../frontend/.neev/foundation"
```

**Frontend team's `neev.yaml`:**
```yaml
remotes:
  - name: backend-api
    path: "../backend/.neev/foundation"
```

Both teams maintain their own specs but can sync shared knowledge.

### How do we handle conflicting blueprints?

**Prevention:**
- Regular blueprint reviews
- Clear ownership (who owns which blueprints)
- Foundation docs for shared standards

**Resolution:**
1. Identify conflict
2. Team discussion
3. Update one or both blueprints
4. Document decision in foundation if architectural

```markdown
# .neev/foundation/decisions.md

## Resolution: Authentication Approach
After conflict between blueprints, we've decided to use JWT.
All services must follow this standard.
```

---

## Technical Questions

### Does Neev require internet access?

No! Neev is completely offline:
- No API calls
- No external services
- No cloud dependencies
- Pure local file operations

### Does Neev store my data anywhere?

No. Everything stays in your `.neev/` directory on your machine. Neev never sends data to external servers.

### What's the performance impact?

**Minimal:**
- CLI operations are instant (< 100ms)
- File operations are local
- No build step required
- No runtime overhead

### Can I use Neev in CI/CD?

Absolutely! Common use cases:

```yaml
# GitHub Actions example
- name: Check specification drift
  run: |
    curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_linux_amd64.tar.gz | tar xz
    sudo mv neev /usr/local/bin/
    neev inspect --json --strict
```

See [TUTORIALS.md](TUTORIALS.md#tutorial-4-cicd-integration) for full example.

### What are the system requirements?

**Minimal:**
- Any OS (macOS, Linux, Windows)
- No dependencies (standalone binary)
- ~10 MB disk space for binary
- Any modern terminal

**No requirements for:**
- Go installation (if using binary)
- Node.js
- Python
- Docker

### Can I extend Neev with plugins?

Not currently. Neev is designed to be simple and focused. 

**Workarounds:**
- Custom scripts to process `.neev/` files
- Git hooks for automation
- CI/CD workflows for validation

### Is there an API or library version?

Neev is CLI-only currently. You can:
- Call it from scripts: `neev bridge > output.md`
- Parse its JSON output: `neev inspect --json`
- Process `.neev/` files directly (they're markdown)

---

## Troubleshooting

### Error: ".neev directory already exists"

**Cause:** Trying to run `neev init` in an already initialized project.

**Solution:**

```bash
# Check existing setup
ls -la .neev/

# If you want to reinitialize (BE CAREFUL - backs up first)
mv .neev .neev.backup
neev init

# Or just use existing setup
# No need to run init again
```

### Error: "blueprint already exists: <name>"

**Cause:** A blueprint with that name already exists.

**Solution:**

```bash
# Check existing blueprints
ls .neev/blueprints/

# Use a different name
neev draft "authentication-v2"

# Or remove existing (if intended)
rm -rf .neev/blueprints/authentication
neev draft "authentication"
```

### Error: "failed to read blueprints directory"

**Cause:** Missing or corrupted `.neev/blueprints/` directory.

**Solution:**

```bash
# Recreate structure
mkdir -p .neev/blueprints
mkdir -p .neev/foundation

# Or reinitialize
neev init
```

### `neev bridge` produces no output

**Cause:** No markdown files in `.neev/`.

**Solution:**

```bash
# Check for files
find .neev -name "*.md"

# If empty, create content
neev draft "Sample Blueprint"
echo "# Stack\n\nNode.js" > .neev/foundation/stack.md

# Try again
neev bridge
```

### `neev bridge --focus` returns nothing

**Cause:** Keyword doesn't match any file content.

**Solution:**

```bash
# Search to verify keyword exists
grep -ri "your-keyword" .neev/

# Try broader keyword
neev bridge --focus "api"

# Or use full context
neev bridge
```

### `neev inspect` shows drift but I disagree

**Causes:**
- False positives (code exists but not detected)
- Different understanding of "implementation"
- Module descriptors not configured

**Solutions:**

```bash
# Use module descriptors for accurate detection
# Create .neev/blueprints/feature/.module.yaml
cat > .neev/blueprints/feature/.module.yaml << 'EOF'
expected_files:
  - src/feature.js
  - tests/feature.test.js
EOF

neev inspect --use-descriptors

# Or update blueprint if code doesn't exist
# Or implement the feature if spec is correct
```

### Neev is slow

**This should not happen.** Neev operations are designed to be instant.

**Possible causes:**

1. **Extremely large `.neev/` directory**
   ```bash
   du -sh .neev/
   # If > 100MB, you may have non-markdown files
   ```

2. **Network drive or slow disk**
   ```bash
   # Move project to local SSD
   ```

3. **Antivirus scanning**
   ```bash
   # Add .neev/ to antivirus exclusions
   ```

If still slow, [open an issue](https://github.com/neev-kit/neev/issues).

### Permission denied errors

**Cause:** Insufficient permissions to create/modify files.

**Solution:**

```bash
# Check directory permissions
ls -la

# Fix ownership if needed
sudo chown -R $USER:$USER .neev/

# Or run with sudo (not recommended)
sudo neev init
```

### `neev sync-remotes` fails

**Cause:** Remote path not found or inaccessible.

**Solution:**

```bash
# Check neev.yaml configuration
cat .neev/neev.yaml

# Verify remote path exists
ls ../backend/.neev/foundation

# Fix path in neev.yaml
# Use absolute paths if relative paths don't work
```

### Binary doesn't run on Windows

**Cause:** Windows requires `.exe` extension.

**Solution:**

```bash
# Download Windows binary from releases
# https://github.com/neev-kit/neev/releases

# Rename if needed
mv neev neev.exe

# Add to PATH
# Control Panel > System > Environment Variables
```

### Bridge output is too large

**Cause:** Many blueprints or large files.

**Solution:**

```bash
# Use --focus to reduce size
neev bridge --focus "specific-area"

# Or split into multiple aggregations
neev bridge --focus "frontend" > frontend.md
neev bridge --focus "backend" > backend.md

# Or review blueprints and archive completed ones
neev lay "old-completed-feature"
```

---

## Best Practices

### How often should I run `neev inspect`?

**Recommendations:**
- **Daily:** If actively developing
- **Per PR:** In code review process
- **Weekly:** For maintenance
- **CI/CD:** On every push/PR

### When should I archive blueprints?

Archive when:
- ✅ Feature is fully implemented
- ✅ Feature is deployed to production
- ✅ Feature is stable (no major changes expected)
- ✅ Blueprint is no longer actively referenced

Don't archive:
- ❌ Features still in development
- ❌ Features with known issues
- ❌ Features being actively modified

### How detailed should blueprints be?

**Right level of detail:**
```markdown
# Good
## Authentication
Use JWT tokens with 1-hour expiry.
Store refresh tokens in database.
Implement rate limiting: 5 attempts per hour.
```

**Too vague:**
```markdown
# Too vague
## Authentication
Use secure authentication methods.
```

**Too detailed:**
```markdown
# Too detailed (implementation, not architecture)
## Authentication
\`\`\`javascript
const jwt = require('jsonwebtoken');
const secret = process.env.JWT_SECRET;
function generateToken(userId) {
  return jwt.sign({ userId }, secret, { expiresIn: '1h' });
}
\`\`\`
```

**Rule:** Describe **what** and **why**, not step-by-step **how**.

### Should I use Neev for documentation or specifications?

**Both!** Neev is designed for:

**Specifications (planning phase):**
- Document what you're going to build
- Get team alignment
- Generate AI context for implementation

**Documentation (maintenance phase):**
- Keep specs updated as code evolves
- Onboard new team members
- Reference during development

---

## Advanced Topics

### Can I use Neev with monorepos?

Yes! Two approaches:

**Approach 1: Single `.neev/` at root**
```
monorepo/
├── .neev/
│   ├── blueprints/
│   │   ├── frontend-feature/
│   │   └── backend-feature/
│   └── foundation/
├── frontend/
└── backend/
```

**Approach 2: Separate `.neev/` per project + remotes**
```
monorepo/
├── frontend/
│   └── .neev/
├── backend/
│   └── .neev/
└── shared-standards/
    └── .neev/
```

Use remotes to share across projects.

### How do I handle secrets in blueprints?

**Never commit secrets!**

**Bad:**
```markdown
## Database Connection
Host: prod-db.example.com
Password: super_secret_123
```

**Good:**
```markdown
## Database Connection
Use environment variables:
- `DB_HOST` — Database host
- `DB_PASSWORD` — Database password

See `.env.example` for required variables.
```

### Can I generate documentation websites from Neev?

Yes! Neev outputs markdown which can be converted:

```bash
# Generate markdown
neev bridge > docs/architecture.md

# Convert to HTML with Pandoc
pandoc docs/architecture.md -o docs/architecture.html

# Or use static site generators
# MkDocs, Docusaurus, Hugo, etc.
```

### How do I handle API versioning?

Create separate blueprints:

```bash
neev draft "API v1"
neev draft "API v2"

# Or within single blueprint
# .neev/blueprints/api/architecture.md:
## Version 1 (Deprecated)
...

## Version 2 (Current)
...

## Version 3 (Planned)
...
```

### Can I use Neev for database migrations?

Yes! Document schema changes:

```bash
neev draft "Database Schema v2"
# Document changes, migrations, rollbacks

# Use neev bridge to share with team
neev bridge --focus "database"
```

See [TUTORIALS.md](TUTORIALS.md#tutorial-8-database-schema-evolution) for detailed example.

---

## Getting More Help

### Where can I find examples?

- [Tutorials](TUTORIALS.md) — 8 step-by-step guides
- [Getting Started](GETTING_STARTED.md) — Walkthrough examples
- [Best Practices](BEST_PRACTICES.md) — Patterns and anti-patterns

### How do I report bugs?

[Open an issue on GitHub](https://github.com/neev-kit/neev/issues) with:
- Neev version (`neev --version`)
- Operating system
- Command that failed
- Expected vs actual behavior
- Error messages

### How do I request features?

[Start a discussion on GitHub](https://github.com/neev-kit/neev/discussions) to:
- Propose new features
- Discuss ideas
- Get community feedback

### How do I contribute?

See [CONTRIBUTING.md](CONTRIBUTING.md) for:
- Development setup
- Code standards
- Submission process

### Is there a community?

Yes!
- [GitHub Discussions](https://github.com/neev-kit/neev/discussions) — Ask questions
- [GitHub Issues](https://github.com/neev-kit/neev/issues) — Report bugs
- [Twitter/X](#) — Follow updates (coming soon)

---

## Quick Reference

### Common Commands

```bash
# Initialize project
neev init

# Create blueprint
neev draft "Feature Name"

# Generate full context
neev bridge

# Generate focused context
neev bridge --focus "keyword"

# Check for drift
neev inspect

# Archive completed work
neev lay "blueprint-name"

# Generate Copilot instructions
neev instructions

# Sync remotes
neev sync-remotes
```

### Common Workflows

```bash
# Planning a new feature
neev draft "New Feature"
# Edit intent.md and architecture.md
git commit -m "docs: add new feature spec"

# Getting AI help
neev bridge --focus "feature" > context.md
# Paste context.md to AI

# Checking project health
neev inspect --json --strict

# Team onboarding
neev bridge > ONBOARDING.md
```

### File Locations

```
.neev/
├── neev.yaml           # Configuration
├── blueprints/         # Feature specifications
│   └── feature-name/
│       ├── intent.md
│       └── architecture.md
└── foundation/         # Project-wide docs
    ├── stack.md
    ├── principles.md
    └── conventions.md
```

---

**Still have questions?** [Open a discussion](https://github.com/neev-kit/neev/discussions) or check our [documentation](README.md)!
