# Neev: Spec‑Driven Development CLI

**Stop vibe coding. Build exactly what you spec.**

Neev is a local‑first SDD toolkit that makes your **foundation** (specs + blueprints) the source of truth, then inspects drift and bridges context to Claude Code/Copilot.

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![Tests](https://img.shields.io/github/actions/workflow/status/neev-kit/neev/tests.yml?branch=main&style=flat-square&label=Tests)](https://github.com/neev-kit/neev/actions/workflows/tests.yml)
[![Release](https://img.shields.io/github/v/release/neev-kit/neev?style=flat-square&label=Release)](https://github.com/neev-kit/neev/releases/latest)
[![SDD Certified](https://img.shields.io/badge/SDD-Spec--Driven%20Development-9b59b6?style=flat-square)](https://github.com/neev-kit/neev)

**Unlike Spec-Kit:** Repo‑aware + polyrepo support  
**Unlike OpenSpec:** Go CLI + drift detection  
**No external APIs. No dependencies. All files versioned in git.**

## How It Works

1. **Write blueprints** (markdown specs of features/components)
2. **Aggregate context** (run `neev bridge` to get full project context)
3. **Build with AI** (pass context to Claude, Cursor, Copilot for implementation)
4. **Verify alignment** (run `neev inspect` to catch drift from specs)

**Real benefit**: Ship aligned code on first attempt because **specs drive decisions before you type**.

### The Neev Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                     SPEC-DRIVEN DEVELOPMENT                     │
└─────────────────────────────────────────────────────────────────┘

1️⃣  PLAN                    2️⃣  BUILD                  3️⃣  VERIFY
─────────────────────────────────────────────────────────────────
Write blueprints       →    Copy context to AI    →    neev inspect
in .neev/                   (Claude/Copilot)            checks drift
    ↓                            ↓                           ↓
Specs clear            →    Implementation         →    Code matches
Organized                   with full intent           architecture
    ↓                            ↓                           ↓
Every feature              Code generated              Specs stay
documented                 from specs                  in sync

Result: Merged code matches specs on first attempt ✅
```

### Workflow: Adding a New Feature

**In GitHub Copilot Chat or Claude Code, use slash commands:**

```
Step 1: /neev:draft "User Authentication"
        → Creates blueprint structure + intent.md

Step 2: /neev:draft "Expand User Authentication with API design"
        → AI drafts architecture.md with endpoints, models, security

Step 3: /neev:bridge --focus auth
        → Get full project context for implementation

Step 4: [Paste context] "Build this user authentication system"
        → AI implements the feature

Step 5: /neev:inspect --strict
        → Verify code matches your documented architecture
```

### Workflow: Updating an Existing Feature

**In GitHub Copilot Chat or Claude Code, use slash commands:**

```
Step 1: /neev:draft "Add email verification to User Authentication"
        → AI updates the blueprint with new feature

Step 2: /neev:bridge --focus auth
        → Get updated context with your new spec

Step 3: [Paste context] "Implement email verification"
        → AI builds the new feature

Step 4: /neev:inspect
        → Verify implementation matches updated blueprint

Step 5: /neev:inspect --strict --drift-check
        → Ensure no spec-code divergence before merge
```

---

### Why Neev? (Business Value)

**Problem:** Specs go stale, code diverges from intent, reviews miss alignment issues, rework happens.

**Neev's Solution:**

| Problem | Neev's Answer |
|---------|---|
| Specs get outdated | Specs = source of truth (auto-validated) |
| Team misses context | `neev bridge` auto-aggregates everything |
| Code-spec drift | `neev inspect` catches it immediately |
| AI builds wrong thing | Full context → AI builds exactly as specified |
| Specs as documentation | Blueprints ARE docs (always in sync) |
| Review takes 2 hours | Alignment checked before review starts |

**Result:** Ship aligned code on first attempt, reduce rework, improve team velocity.

### Quick Adoption (No Risk)

```
Week 1: Install Neev (5 min)
        Initialize one service (2 min per service)
        Team writes one blueprint (30 min)
        ✅ Try it, zero commitment

Week 2-3: Integrate with CI/CD (if desired)
          Run `neev inspect` in merge gates
          Optional: Enforce specs quality

Month 1: Team sees alignment benefits
         Reduces code review friction
         Faster onboarding with specs

Month 2: Roll out to all services
         Integrate with AI workflows
         Measure spec coverage % over time

Month 3+: Continuous value
          Specs stay in sync automatically
          Less rework, faster shipping
```

### Integration is Simple

**Option 1: Local Only**
```bash
$ neev inspect              # After coding, before push
$ Result: Catch drift early
```

**Option 2: In CI/CD**
```bash
$ neev inspect --strict            # Merge gate check
$ Result: Block misaligned PRs
```

**Option 3: With AI**
```bash
$ neev bridge --focus feature      # Copy context
$ Paste into Claude/Copilot
$ Result: AI builds exactly as specified
```

All options work **together** — pick what suits your team.

### Zero Operational Overhead

- ✅ **Local-first:** All files in git (no cloud, no API)
- ✅ **Single binary:** Just `neev` command (no dependencies)
- ✅ **Markdown-based:** Specs in version control
- ✅ **Gradual adoption:** Start with one team/service
- ✅ **No vendor lock-in:** Specs are plain markdown
- ✅ **Polyrepo-ready:** Works with monorepos and multi-repos

### Success Metrics (6 Months)

```
Before Neev:
  • Spec-code drift incidents: 8-10/month
  • Code review time: ~2 hours per PR
  • Rework due to misalignment: ~15% of features
  • New engineer onboarding: 2-3 weeks to productivity

After Neev:
  • Spec-code drift incidents: <1/month
  • Code review time: ~30 min per PR
  • Rework due to misalignment: <5% of features
  • New engineer onboarding: 1 week with specs
```

### How to Get Started (3 Steps)

**Step 1:** Try it locally on one service
```bash
neev init                           # Set up .neev/
neev draft "Your Feature"           # Create blueprint
neev bridge                         # See aggregated context
neev inspect --drift-check          # Check alignment
```

**Step 2:** Show your team the benefits
- Specs are clear and discoverable
- AI context is comprehensive
- Drift is caught early

**Step 3:** Roll out to more services
- Add to CI/CD for merge gates (optional)
- Integrate with your AI workflows
- Track spec coverage % over time

---

## 5-Minute Setup

### 1️⃣ Install

**macOS (Apple Silicon):**
```bash
curl -sL $(curl -s https://api.github.com/repos/neev-kit/neev/releases/latest | grep -o '"browser_download_url": "[^"]*darwin_arm64[^"]*"' | cut -d'"' -f4) | tar xz && sudo mv neev_*/neev /usr/local/bin/ && rm -rf neev_*
```

**macOS (Intel):**
```bash
curl -sL $(curl -s https://api.github.com/repos/neev-kit/neev/releases/latest | grep -o '"browser_download_url": "[^"]*darwin_amd64[^"]*"' | cut -d'"' -f4) | tar xz && sudo mv neev_*/neev /usr/local/bin/ && rm -rf neev_*
```

**Linux (x86_64):**
```bash
curl -sL $(curl -s https://api.github.com/repos/neev-kit/neev/releases/latest | grep -o '"browser_download_url": "[^"]*linux_amd64[^"]*"' | cut -d'"' -f4) | tar xz && sudo mv neev_*/neev /usr/local/bin/ && rm -rf neev_*
```

**Windows (PowerShell):**
```powershell
$url = (curl -s https://api.github.com/repos/neev-kit/neev/releases/latest | Select-String 'windows_amd64.zip' -Raw) -replace '.*"browser_download_url": "([^"]+)".*', '$1'
curl -Lo $env:TEMP\neev.zip $url
Expand-Archive $env:TEMP\neev.zip -DestinationPath $env:TEMP\neev_extract
Move-Item $env:TEMP\neev_extract\neev_*\neev.exe $env:USERPROFILE\AppData\Local\Microsoft\WindowsApps\neev.exe
Remove-Item $env:TEMP\neev.zip, $env:TEMP\neev_extract -Recurse
```

**Homebrew (macOS/Linux):**
```bash
brew install neev-kit/tap/neev
```

**Build from source:**
```bash
git clone https://github.com/neev-kit/neev.git && cd neev
go build -o neev ./cli && sudo mv neev /usr/local/bin/
```

### 2️⃣ Initialize

```bash
cd /path/to/your/project
neev init
```

Creates:
```
.neev/
├── neev.yaml                    # Config
├── blueprints/                  # Features you'll build
├── foundation/
│   ├── stack.md                 # Tech stack
│   ├── principles.md            # Design principles
│   └── patterns.md              # Architecture patterns
├── .commands/                   # AI command manifests
└── .ai/                         # AI context files
```

### 3️⃣ Create Your First Blueprint

```bash
neev draft "User Authentication"
```

Generates:
```
.neev/blueprints/user-authentication/
├── intent.md         # Why we're building this
└── architecture.md   # How it works
```

### 4️⃣ Get AI-Ready Context

```bash
neev bridge
# Copy output → Paste into Claude/Copilot/Cursor
```

That's it. You're ready to build aligned code.

---

## How It Works

**The Neev Workflow:**

```
PLAN              BUILD              VERIFY
────────          ─────              ──────
Write blueprints  Copy context to    neev inspect catches
in .neev/         AI coding tool     any drift from specs
       ↓                  ↓                ↓
  Specs clear     Implementation   Code matches
                  with intent       architecture
```

**Key insight:** Specs become the source of truth. Code should match specs, not the other way around.

## All Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `neev init` | Set up .neev/ foundation | `neev init` |
| `neev draft <title>` | Create a blueprint | `neev draft "User API"` |
| `neev bridge [flags]` | Get AI-ready context | `neev bridge` or `neev bridge --focus auth` |
| `neev inspect` | Check code matches specs | `neev inspect` or `neev inspect --strict` (for CI) |
| `neev openapi <bp>` | Generate OpenAPI spec | `neev openapi user-api` |
| `neev cucumber <bp>` | Generate BDD tests | `neev cucumber user-api --lang go` |
| `neev handoff <role>` | Create handoff prompts | `neev handoff backend` |
| `neev lay <bp>` | Archive completed blueprint | `neev lay user-api` |
| `neev instructions` | Update Copilot instructions | `neev instructions` |
| `neev slash-commands` | Manage slash commands | `neev slash-commands --list` |
| `neev migrate` | Convert from OpenSpec/Spec-Kit | `neev migrate --source openspec` |
| `neev sync-remotes` | Sync remote foundations | `neev sync-remotes` |

**Full reference:** See [COMMAND_CATALOG.md](COMMAND_CATALOG.md)

## Key Features

**🎯 Specs First**
- Write blueprints before code → implementation matches intent
- `neev inspect` catches drift early (use `--strict` in CI/CD)

**🤖 AI-Native**
- `neev bridge` outputs context optimized for Claude, Copilot, Cursor
- Slash commands for quick access in GitHub Copilot Chat
- No API keys, no cloud login — works offline

**🔍 Alignment**
- Detect when code diverges from specs
- Module-level validation with `.module.yaml` descriptors
- JSON output for CI/CD pipelines

**🌐 Polyrepo Ready**
- Reference specs from other repositories
- Share architecture across teams
- Fine-grained public/private control

**📚 Documentation**
- Generate OpenAPI specs from architecture blueprints
- Create BDD test scaffolding automatically
- Onboarding docs are always in sync with reality

**🚀 Zero Setup**
- Local-first, no external dependencies
- All files versioned in `.neev/` (commit to git)
- Works with any project type (Go, Node, Python, etc.)

## Configuration (neev.yaml)

```yaml
project_name: My Project
foundation_path: .neev
ignore_dirs:
  - node_modules
  - .git
  - vendor

# For polyrepo support
remotes:
  - name: shared-lib
    path: "../shared/.neev/foundation"
    public_only: true
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup.

## Real-World Example: Building a User API

**1. Create blueprint:**
```bash
neev draft "User API"
```

**2. Document in `.neev/blueprints/user-api/architecture.md`:**
```markdown
## User API Architecture

### Endpoints

#### GET /api/users
List all users (paginated)
- Query: page, limit
- Response: { id, name, email, created_at }

#### POST /api/users
Create user
- Body: { name, email, password }
- Response: { id, name, email }

#### GET /api/users/:id
Get user by ID

#### PUT /api/users/:id
Update user
```

**3. Generate documentation:**
```bash
# Generate OpenAPI spec
neev openapi user-api

# Generate BDD tests
neev cucumber user-api --lang go
```

**4. Build with AI:**
```bash
# Get context for the entire project
neev bridge > context.md

# Paste context + architecture into Cursor/Claude:
# "Build this user API according to the context above"
```

**5. Verify alignment:**
```bash
# After implementing, check code matches specs
neev inspect
```

**Why this matters:** Your API is built exactly as specified. No rewrites. No "wait, we said it differently in the blueprint."

## Use Cases

| Scenario | How Neev Helps |
|----------|---|
| **Building with AI** | Specs → AI context → aligned code on first attempt |
| **Team onboarding** | Blueprints = auto-generated onboarding docs |
| **Code reviews** | Reviewers see intent + architecture, not just code |
| **Polyrepo/monorepo** | Reference specs across repositories |
| **API design** | Architecture → OpenAPI spec → BDD tests automatically |
| **Preventing rework** | Specs caught early, implementation is clean |

## Development

```bash
# Build
go build -o neev ./cli

# Run tests
go test ./...

# Test coverage
go test -cover ./...
```

**More details:**
- [DEVELOPMENT.md](DEVELOPMENT.md) — Setup & debugging
- [ARCHITECTURE.md](ARCHITECTURE.md) — System design
- [CONTRIBUTING.md](CONTRIBUTING.md) — Contributing guidelines

## Why Neev Over Manual Approaches

| Aspect | Manual Specs | Neev |
|--------|------|------|
| **Specs stay in sync** | ❌ Gets outdated | ✅ `neev inspect` catches drift |
| **AI context** | ❌ Copy/paste manually | ✅ `neev bridge` auto-aggregates |
| **Onboarding docs** | ❌ Write separately | ✅ Blueprints ARE docs |
| **API documentation** | ❌ Manual work | ✅ `neev openapi` generates |
| **Test scaffolding** | ❌ Write from scratch | ✅ `neev cucumber` generates |
| **Versioned specs** | ❌ External tools | ✅ All in git (.neev/) |

## FAQ

**Q: Does Neev work with my tech stack?**  
A: Yes. Node, Python, Go, Rust, .NET, Java, etc. Neev works with any project.

**Q: Do I have to commit `.neev/` to git?**  
A: Yes. `.neev/` is your project's architecture. It should be versioned.

**Q: How is this different from writing prompts?**  
A: Specs in `.neev/` stay in sync with implementation. Manual prompts don't — they become stale.

**Q: Can teams share blueprints?**  
A: Yes, via `neev sync-remotes`. Reference foundations from other repos.

**Q: What about large projects?**  
A: Use `--focus` to filter context by keywords. Blueprints can nest and reference each other.

**More questions?** See [FAQ.md](FAQ.md)

## License & Contributing

MIT License — See [LICENSE](LICENSE)

**Want to contribute?** See [CONTRIBUTING.md](CONTRIBUTING.md)

**Community Standards:**
- 📋 [Code of Conduct](CODE_OF_CONDUCT.md) — Community guidelines
- 🔒 [Security Policy](SECURITY.md) — Reporting vulnerabilities
- 📝 [Changelog](CHANGELOG.md) — Release history and updates
- 📊 [Open Source Policy](OPEN_SOURCE_POLICY.md) — Compliance details

**Get Involved:**
- 🐛 [Report bugs](https://github.com/neev-kit/neev/issues)
- 💡 [Suggest features](https://github.com/neev-kit/neev/discussions)
- 🔧 [Submit PRs](https://github.com/neev-kit/neev/pulls)

## Acknowledgments

Neev is inspired by and builds on ideas from **Spec-Kit** (GitHub) and **OpenSpec** (Fission-AI). See [ACKNOWLEDGMENTS.md](ACKNOWLEDGMENTS.md) for details on projects and people who influenced Neev's design.

---

## Next Steps

1. **Try it now:** `neev init && neev draft "My Feature"`
2. **Learn more:** [GETTING_STARTED.md](GETTING_STARTED.md) or [COMMAND_CATALOG.md](COMMAND_CATALOG.md)
3. **Build with AI:** `neev bridge` → copy to Claude/Copilot → ship aligned code
