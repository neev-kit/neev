# 🚀 Getting Started with Neev

**Welcome to Neev!** This guide will walk you through everything you need to know to start using Neev effectively in your projects.

## 📋 Table of Contents

1. [What is Neev?](#what-is-neev)
2. [Installation](#installation)
3. [Your First Project](#your-first-project)
4. [Understanding the Basics](#understanding-the-basics)
5. [Common Workflows](#common-workflows)
6. [Next Steps](#next-steps)

## What is Neev?

Neev is a lightweight CLI framework that helps you bridge your project's intent with AI coding assistants. Think of it as a structured way to:

- **Document what you want to build** before writing code
- **Organize project knowledge** in version-controlled markdown files
- **Generate AI-ready context** for better code suggestions
- **Maintain alignment** between specs and implementation

### Why Use Neev?

✅ **Zero Setup** — No API keys, no external services, just local files  
✅ **AI-Optimized** — Output designed for Claude, Copilot, Cursor, etc.  
✅ **Version Controlled** — All specs live in `.neev/` alongside your code  
✅ **Framework Agnostic** — Works with any language or tech stack  

## Installation

### Option 1: Download Pre-built Binary (Recommended)

Choose your platform and run the appropriate commands:

**macOS (Intel):**
```bash
curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_darwin_amd64.tar.gz | tar xz
sudo mv neev /usr/local/bin/
```

**macOS (Apple Silicon):**
```bash
curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_darwin_arm64.tar.gz | tar xz
sudo mv neev /usr/local/bin/
```

**Linux (x86_64):**
```bash
curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_linux_amd64.tar.gz | tar xz
sudo mv neev /usr/local/bin/
```

**Windows:**
1. Visit [GitHub Releases](https://github.com/neev-kit/neev/releases)
2. Download the Windows `.zip` file
3. Extract and add to your PATH

### Option 2: Using Go Install

If you have Go installed (1.23+):

```bash
go install github.com/neev-kit/neev/cli@latest
```

### Option 3: Build from Source

```bash
git clone https://github.com/neev-kit/neev.git
cd neev
go mod download
go build -o neev ./cli
sudo mv neev /usr/local/bin/
```

### Verify Installation

```bash
neev --version
```

You should see the version number. If you get "command not found", ensure the binary is in your PATH.

## Your First Project

Let's walk through setting up Neev in a sample project. You can follow along with your own project or create a test directory.

### Step 1: Initialize Neev

Navigate to your project directory and initialize Neev:

```bash
cd /path/to/your/project
neev init
```

**What this does:**
- Creates `.neev/` directory in your project root
- Generates `neev.yaml` configuration file
- Sets up `blueprints/` and `foundation/` directories

**Directory structure after init:**
```
your-project/
├── .neev/
│   ├── neev.yaml          # Configuration
│   ├── blueprints/        # Feature specifications
│   └── foundation/        # Project-wide docs
└── [your existing code]
```

### Step 2: Create Your First Blueprint

A blueprint represents a feature or component you want to build. Let's create one for user authentication:

```bash
neev draft "User Authentication"
```

**What this creates:**
```
.neev/blueprints/user-authentication/
├── intent.md          # What and why
└── architecture.md    # How it works
```

### Step 3: Document Your Blueprint

Edit the generated files to describe your feature:

**`.neev/blueprints/user-authentication/intent.md`:**
```markdown
# User Authentication Intent

## Purpose
Implement secure user authentication system for the application.

## Goals
- Support email/password login
- Implement JWT-based sessions
- Add password reset functionality
- Include rate limiting for security

## Success Criteria
- Users can register and login securely
- Sessions persist across page refreshes
- Failed login attempts are limited
```

**`.neev/blueprints/user-authentication/architecture.md`:**
```markdown
# User Authentication Architecture

## Components
- Auth Service: Handle login/logout/register
- JWT Middleware: Validate tokens
- Database Schema: Users table with hashed passwords

## Technology Stack
- bcrypt for password hashing
- JWT for tokens
- PostgreSQL for user storage

## Security Considerations
- Never store plain text passwords
- Use secure random token generation
- Implement CSRF protection
```

### Step 4: Create Foundation Documents

Foundation documents describe project-wide conventions and principles:

```bash
# Create foundation files manually
echo "# Tech Stack\n\nWe use Node.js with Express, PostgreSQL, and React." > .neev/foundation/stack.md
echo "# Principles\n\n- Security first\n- Simple over clever\n- Test everything" > .neev/foundation/principles.md
```

### Step 5: Generate Context for AI

Now aggregate all your documentation into AI-ready context:

```bash
neev bridge
```

**Output:**
```markdown
# Project Foundation
## File: stack.md
# Tech Stack

We use Node.js with Express, PostgreSQL, and React.

## File: principles.md
# Principles

- Security first
- Simple over clever
- Test everything

## File: intent.md
# User Authentication Intent
...
```

### Step 6: Use with AI Assistants

Copy the output and share it with your AI coding assistant:

```bash
# Copy to clipboard (macOS)
neev bridge | pbcopy

# Copy to clipboard (Linux)
neev bridge | xclip -selection clipboard

# Save to file
neev bridge > context.md
```

Then, in your AI assistant (Claude, Cursor, GitHub Copilot, etc.):

```
Here's my project context:
[paste the output]

Now implement the user authentication service according to these specifications.
```

## Understanding the Basics

### Core Concepts

#### 1. Blueprints

**What:** Specifications for individual features or components  
**Where:** `.neev/blueprints/<blueprint-name>/`  
**Contains:** `intent.md` (what/why) and `architecture.md` (how)

**Example use cases:**
- Feature specifications (User Auth, Payment Processing)
- Component designs (API Gateway, Database Schema)
- Module documentation (Email Service, Logging System)

#### 2. Foundation

**What:** Project-wide knowledge shared across all blueprints  
**Where:** `.neev/foundation/`  
**Contains:** Any `.md` files you create (conventions, stack, principles, etc.)

**Example files:**
- `stack.md` — Technology choices
- `principles.md` — Development philosophy
- `conventions.md` — Coding standards
- `architecture.md` — System overview

#### 3. Context Bridging

**What:** Aggregating all specs into a single AI-ready document  
**How:** `neev bridge` command  
**Output:** Markdown with all foundation + blueprint content

**Key features:**
- Filter by keyword: `neev bridge --focus "auth"`
- Include remotes: `neev bridge --with-remotes`
- Claude-optimized: `neev bridge --claude`

### Configuration File

The `neev.yaml` file controls Neev's behavior:

```yaml
project_name: "My Project"
foundation_path: ".neev"
ignore_dirs:
  - node_modules
  - .git
  - __pycache__
  - dist
  - vendor
```

**Key options:**
- `project_name` — Display name for your project
- `foundation_path` — Where `.neev/` lives (default: `.neev`)
- `ignore_dirs` — Directories to skip during inspection

## Common Workflows

### Workflow 1: Planning a New Feature

```bash
# 1. Create a blueprint
neev draft "Shopping Cart"

# 2. Document intent and architecture
# Edit .neev/blueprints/shopping-cart/intent.md
# Edit .neev/blueprints/shopping-cart/architecture.md

# 3. Get focused context
neev bridge --focus "cart" > cart-spec.md

# 4. Share with AI to generate implementation
# Paste cart-spec.md content to your AI assistant
```

### Workflow 2: Onboarding New Team Members

```bash
# Generate comprehensive project overview
neev bridge > TEAM_ONBOARDING.md

# Share with new developers
# They now have complete context about the project
```

### Workflow 3: Code Review with Context

```bash
# Get context for specific area being reviewed
neev bridge --focus "payment" > payment-context.md

# Reviewers can understand the intent behind the code
```

### Workflow 4: Checking Project Drift

```bash
# Inspect if code matches specs
neev inspect

# Get structured report for CI/CD
neev inspect --json

# Strict mode (exit 1 if drift detected)
neev inspect --strict
```

### Workflow 5: GitHub Copilot Integration

```bash
# Generate Copilot instructions from your specs
neev instructions

# Creates/updates .github/copilot-instructions.md
# Copilot now understands your project better
```

### Workflow 6: Multi-Repo Projects

```yaml
# In neev.yaml, reference other repos
remotes:
  - name: api
    path: "../backend/.neev/foundation"
    public_only: true
```

```bash
# Sync remote foundations
neev sync-remotes

# Include in bridge
neev bridge --with-remotes
```

### Workflow 7: Archiving Completed Work

```bash
# Move completed blueprint to foundation
neev lay "user-authentication"

# Blueprint is archived, changelog is updated
```

## Next Steps

Congratulations! You now know the basics of Neev. Here's what to explore next:

### 📖 **Detailed Guides**

- **[API Reference](API_REFERENCE.md)** — Complete command documentation
- **[Tutorials](TUTORIALS.md)** — Step-by-step walkthroughs for specific use cases
- **[Best Practices](BEST_PRACTICES.md)** — Patterns and anti-patterns
- **[Usage Guide](USAGE.md)** — Comprehensive feature documentation

### 🛠️ **Advanced Topics**

- **Production Features** — See [BEST_PRACTICES.md](BEST_PRACTICES.md) for production patterns
  - Module descriptors with `.module.yaml`
  - Advanced drift detection
  - Polyrepo support
  - AI assistant integrations

### 💡 **Real Examples**

Try these example scenarios:

1. **Microservices Architecture**
   ```bash
   neev draft "User Service"
   neev draft "Product Service"
   neev draft "API Gateway"
   neev bridge
   ```

2. **API Documentation**
   ```bash
   neev draft "REST API Endpoints"
   neev draft "WebSocket Events"
   neev draft "Authentication Flow"
   neev bridge --focus "API"
   ```

3. **Team Conventions**
   ```bash
   # Create foundation docs
   echo "# Code Style\n\nFollow PEP 8 for Python..." > .neev/foundation/code-style.md
   echo "# Testing\n\nAll features need unit tests..." > .neev/foundation/testing.md
   neev bridge > CONVENTIONS.md
   ```

### 🤝 **Get Help**

- **Troubleshooting** — See [FAQ.md](FAQ.md)
- **Issues** — [Open an issue](https://github.com/neev-kit/neev/issues)
- **Discussions** — [Join discussions](https://github.com/neev-kit/neev/discussions)
- **Contributing** — See [CONTRIBUTING.md](CONTRIBUTING.md)

### 📚 **Learn More**

- **Architecture** — Understand how Neev works: [ARCHITECTURE.md](ARCHITECTURE.md)
- **Development** — Contribute to Neev: [DEVELOPMENT.md](DEVELOPMENT.md)
- **Philosophy** — Why we built Neev this way: [README.md](README.md#why-neev)

---

**Ready to build better software with AI?** Start documenting your next feature in Neev! 🚀
