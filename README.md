# 🏗️ Neev - Spec-Driven Development for Teams

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![Tests](https://img.shields.io/github/actions/workflow/status/neev-kit/neev/tests.yml?branch=main&style=flat-square&label=Tests)](https://github.com/neev-kit/neev/actions)
[![Release](https://img.shields.io/github/v/release/neev-kit/neev?style=flat-square&label=Release)](https://github.com/neev-kit/neev/releases)

**Document your architecture before you build it. Ship aligned code on first implementation.**

Neev is a spec-driven development framework that ensures **what you intend to build matches what you actually build**. It works by:

1. **Write blueprints** (markdown specs of features/components)
2. **Aggregate context** (run `neev bridge` to get full project context)
3. **Build with AI** (pass context to Claude, Cursor, Copilot for implementation)
4. **Verify alignment** (run `neev inspect` to catch drift from specs)

**No external APIs. No dependencies. All files versioned in your repository.**

## The Problem Neev Solves

❌ **Without specs**: Features get built differently than intended → rewrites → shipping late

✅ **With Neev**: Specs are clear → implementation matches intent → ship on first attempt

**Real benefit**: Teams ship aligned code faster because **thinking comes before typing**.

## 5-Minute Setup

### 1️⃣ Install

```bash
# macOS
curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_darwin_arm64.tar.gz | tar xz
sudo mv neev /usr/local/bin/

# Linux
curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_linux_amd64.tar.gz | tar xz
sudo mv neev /usr/local/bin/

# Or build from source
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
└── foundation/
    ├── stack.md                 # Tech stack
    ├── principles.md            # Design principles
    └── patterns.md              # Architecture patterns
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

- 🐛 [Report bugs](https://github.com/neev-kit/neev/issues)
- 💡 [Suggest features](https://github.com/neev-kit/neev/discussions)
- 🔧 [Submit PRs](https://github.com/neev-kit/neev/pulls)

---

## Next Steps

1. **Try it now:** `neev init && neev draft "My Feature"`
2. **Learn more:** [GETTING_STARTED.md](GETTING_STARTED.md) or [COMMAND_CATALOG.md](COMMAND_CATALOG.md)
3. **Build with AI:** `neev bridge` → copy to Claude/Copilot → ship aligned code
