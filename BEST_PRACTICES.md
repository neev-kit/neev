# 💡 Neev Best Practices

Proven patterns, recommendations, and anti-patterns for using Neev effectively in your projects.

## Table of Contents

1. [Blueprint Organization](#blueprint-organization)
2. [Foundation Structure](#foundation-structure)
3. [Naming Conventions](#naming-conventions)
4. [Documentation Writing](#documentation-writing)
5. [Version Control](#version-control)
6. [Team Workflows](#team-workflows)
7. [AI Integration](#ai-integration)
8. [Maintenance](#maintenance)
9. [Common Patterns](#common-patterns)
10. [Anti-Patterns](#anti-patterns)

---

## Blueprint Organization

### ✅ Do: Organize by Domain or Feature

```
.neev/blueprints/
├── user-authentication/
├── payment-processing/
├── notification-system/
└── analytics-dashboard/
```

**Why:** Clear boundaries, easy to find, aligns with feature development.

### ✅ Do: Keep Blueprints Focused

Each blueprint should represent **one cohesive component or feature**.

```bash
# Good: Focused blueprints
neev draft "User Authentication"
neev draft "Password Reset"
neev draft "Email Verification"

# Avoid: Too broad
neev draft "User System"  # Too vague
```

### ✅ Do: Split Large Systems

For complex systems, create multiple related blueprints:

```
.neev/blueprints/
├── api-gateway/              # Entry point
├── user-service/             # Domain service
├── product-service/          # Domain service
└── shared-event-bus/         # Infrastructure
```

### ❌ Don't: Create Too Many Small Blueprints

**Bad:**
```bash
neev draft "Login Button"
neev draft "Logout Button"
neev draft "Forgot Password Link"
```

**Good:**
```bash
neev draft "Authentication UI Components"
```

### ✅ Do: Use Consistent Depth

```
# Good: Consistent structure
.neev/blueprints/
├── user-auth/
│   ├── intent.md
│   └── architecture.md
├── payment/
│   ├── intent.md
│   └── architecture.md

# Avoid: Inconsistent depth
.neev/blueprints/
├── user-auth/
│   └── everything.md        # Non-standard
├── payment/
    ├── intent.md
    ├── architecture.md
    └── security.md           # Additional file (ok if needed)
```

---

## Foundation Structure

### ✅ Do: Keep Foundation Minimal and Stable

Foundation documents should be **project-wide** and **change rarely**.

**Good foundation files:**
```
.neev/foundation/
├── stack.md          # Tech choices
├── principles.md     # Core values
├── conventions.md    # Coding standards
└── architecture.md   # System overview
```

**Not for foundation:**
- Feature-specific details (use blueprints)
- Implementation details (use code comments)
- Temporary decisions (use ADRs in blueprints)

### ✅ Do: Document Architectural Decisions

Create `decisions.md` or `adr.md` in foundation:

```markdown
# Architectural Decision Records

## ADR-001: Use PostgreSQL for Primary Database

**Date:** 2024-01-15  
**Status:** Accepted  
**Context:** Need reliable, ACID-compliant database  
**Decision:** Use PostgreSQL 14+  
**Consequences:** Team needs PostgreSQL expertise  
```

### ✅ Do: Define Clear Conventions

Be specific in your conventions:

```markdown
# Coding Conventions

## API Endpoints
✅ Do: Use plural nouns
- `/users` not `/user`
- `/products/:id` not `/product/:id`

✅ Do: Use kebab-case for multi-word resources
- `/order-items` not `/orderItems` or `/OrderItems`

❌ Don't: Use verbs in URLs
- `/users/123` not `/getUser/123`
```

### ❌ Don't: Duplicate Information

**Bad:**
```
.neev/foundation/tech-stack.md    # Lists: Node.js, React
.neev/foundation/technologies.md   # Lists: Node.js, React
```

**Good:**
```
.neev/foundation/stack.md          # Single source of truth
```

---

## Naming Conventions

### ✅ Do: Use Clear, Descriptive Names

```bash
# Good
neev draft "User Authentication System"
neev draft "Payment Gateway Integration"
neev draft "Real-time Notifications"

# Avoid
neev draft "Feature 1"
neev draft "New Thing"
neev draft "System"
```

### ✅ Do: Be Consistent with Terminology

Pick terms and stick with them:

```bash
# Consistent
neev draft "User Service"
neev draft "Product Service"
neev draft "Order Service"

# Inconsistent
neev draft "User Service"
neev draft "Product Module"      # Module vs Service
neev draft "Order Component"     # Component vs Service
```

### ✅ Do: Use Business Language

Use terms your business/product team uses:

```bash
# Good (matches business language)
neev draft "Shopping Cart"
neev draft "Checkout Flow"
neev draft "Order Fulfillment"

# Technical but unclear to business
neev draft "Cart State Manager"
neev draft "Transaction Processor"
```

### ✅ Do: Version Major Changes

For significant updates:

```bash
neev draft "API v2"
neev draft "Authentication System v2"
neev draft "Database Schema v3"
```

---

## Documentation Writing

### ✅ Do: Write for Your Future Self

Assume you'll forget everything in 6 months:

```markdown
# Intent

## Purpose
Implement user authentication using JWT tokens.

## Why JWT?
- Stateless (no server-side sessions)
- Works across microservices
- Industry standard

## Alternatives Considered
- Session-based auth (rejected: requires Redis)
- OAuth 2.0 (rejected: overkill for our use case)
```

### ✅ Do: Include Examples

Always show examples:

```markdown
# Architecture

## API Endpoint

**Request:**
\`\`\`json
{
  "email": "user@example.com",
  "password": "secure123"
}
\`\`\`

**Response:**
\`\`\`json
{
  "token": "eyJhbG...",
  "expiresIn": 3600
}
\`\`\`
```

### ✅ Do: Document Constraints and Limitations

Be honest about limitations:

```markdown
## Limitations

- Maximum 1000 requests per minute per API key
- Files up to 10MB only
- No support for Internet Explorer 11
- Requires PostgreSQL 12+
```

### ✅ Do: Use Diagrams

Visual aids help:

```markdown
## Authentication Flow

\`\`\`
User → Frontend → API Gateway → Auth Service
                                    ↓
                                Database
                                    ↓
                    JWT Token ← Auth Service
\`\`\`
```

### ❌ Don't: Write Implementation Details

**Bad (too detailed):**
```markdown
Use bcrypt with 12 rounds. Initialize with 
`const bcrypt = require('bcrypt')` then call 
`bcrypt.hash(password, 12)`.
```

**Good (architectural):**
```markdown
Hash passwords using bcrypt with 12 rounds minimum.
Verify with constant-time comparison to prevent timing attacks.
```

### ❌ Don't: Write Prose

Keep it scannable:

**Bad:**
```markdown
So we decided to use React because it's really popular and 
everyone on the team knows it and it has a great ecosystem...
```

**Good:**
```markdown
## Tech Stack

### Frontend Framework: React 18

**Reasons:**
- Team expertise
- Large ecosystem
- Component reusability
- Strong TypeScript support
```

---

## Version Control

### ✅ Do: Commit Specifications Early

```bash
# Create blueprint first
neev draft "User Dashboard"

# Document intent and architecture
# (edit files)

# Commit BEFORE implementation
git add .neev/blueprints/user-dashboard/
git commit -m "docs: add user dashboard specification"

# Now implement
# (write code)
```

**Why:** Specs should drive implementation, not document it after the fact.

### ✅ Do: Use Meaningful Commit Messages

```bash
# Good
git commit -m "docs: add payment processing blueprint"
git commit -m "docs: update API authentication spec"
git commit -m "docs: archive completed user-auth blueprint"

# Avoid
git commit -m "update docs"
git commit -m "changes"
```

### ✅ Do: Review Specification Changes

Treat blueprint changes like code changes:

```bash
# Create feature branch
git checkout -b feature/add-payment-spec

# Make changes
neev draft "Payment Processing"
# (edit files)

# Submit for review
git add .neev/
git commit -m "docs: add payment processing specification"
git push origin feature/add-payment-spec

# Create Pull Request
# Get team feedback before implementation
```

### ✅ Do: Keep Specifications in Sync

```bash
# When code changes significantly
# Update blueprints to match
vim .neev/blueprints/user-auth/architecture.md

git commit -m "docs: update auth spec to reflect OAuth changes"
```

### ❌ Don't: Ignore Drift

```bash
# Bad: Ignoring drift warnings
$ neev inspect
⚠️  Drift detected: user-auth has no implementation

# Fix it or update spec
```

---

## Team Workflows

### ✅ Do: Define Roles

Establish who maintains what:

```markdown
# .neev/foundation/ownership.md

## Blueprint Ownership

- **Backend blueprints** → Backend team lead reviews
- **Frontend blueprints** → Frontend team lead reviews
- **Infrastructure blueprints** → DevOps team reviews
- **Foundation docs** → Technical architect maintains
```

### ✅ Do: Review Blueprints Before Implementation

**Process:**
1. Developer creates blueprint
2. Team reviews spec (not code yet)
3. Iterate on spec until agreed
4. Developer implements
5. Code review references spec

**Benefits:**
- Catch design issues early
- Team alignment before coding
- Better code reviews (check vs spec)

### ✅ Do: Use Blueprints in Planning

```markdown
# Sprint Planning

## Stories Ready for Development
1. User Authentication ✅ (blueprint approved)
2. Payment Processing ✅ (blueprint approved)
3. Email Notifications ⏳ (blueprint in review)

## Blocked
4. Push Notifications ❌ (no blueprint yet)
```

### ✅ Do: Generate Onboarding Docs

```bash
# For new team members
neev bridge > TEAM_ONBOARDING.md

# Include in welcome email:
# "Read TEAM_ONBOARDING.md for complete project context"
```

### ❌ Don't: Skip Blueprint for "Small" Changes

**Bad:**
```
"It's just a small feature, I'll skip the blueprint."
```

**Good:**
```bash
# Even for small features
neev draft "Add User Avatar Upload"
# Document intent and approach
# Then implement
```

**Why:** Small features accumulate. Without specs, you lose architectural coherence.

---

## AI Integration

### ✅ Do: Structure for AI Consumption

Write clear, scannable content:

```markdown
# Good: Clear structure
## Endpoint: POST /users
**Purpose:** Create new user  
**Auth:** Required  
**Rate Limit:** 10/minute  

# Avoid: Prose
This endpoint is used to create a new user. You need 
to be authenticated. We also have a rate limit...
```

### ✅ Do: Use `neev bridge --focus` for Specific Context

```bash
# When asking AI to implement specific feature
neev bridge --focus "authentication" > auth-context.md

# Share only relevant context
# "Implement authentication according to this spec:"
# [paste auth-context.md]
```

### ✅ Do: Include Examples in Blueprints

AI performs better with examples:

```markdown
## User Registration

**Example Request:**
\`\`\`json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "name": "John Doe"
}
\`\`\`

**Example Response:**
\`\`\`json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe"
}
\`\`\`
```

### ✅ Do: Update Copilot Instructions Regularly

```bash
# After blueprint changes
neev instructions

# Commit updated instructions
git add .github/copilot-instructions.md
git commit -m "docs: update Copilot instructions"
```

### ✅ Do: Use Claude-Optimized Output

When using Claude:

```bash
# Claude performs better with its format
neev bridge --claude > context-for-claude.md
```

### ❌ Don't: Overwhelm AI with Context

**Bad:**
```bash
# Too much context
neev bridge > huge-context.md  # 50,000 lines
# "Here's everything, implement this one endpoint"
```

**Good:**
```bash
# Focused context
neev bridge --focus "user-api" > focused.md
# "Implement user registration endpoint per this spec"
```

---

## Maintenance

### ✅ Do: Regular Drift Checks

```bash
# Weekly or in CI/CD
neev inspect

# Fix drift immediately
```

### ✅ Do: Archive Completed Work

```bash
# When feature is done and deployed
neev lay "user-authentication"

# Keeps blueprints directory clean
```

### ✅ Do: Update Specifications When Code Changes

```bash
# Code refactored?
# Update blueprint too

git commit -m "refactor: change auth to OAuth 2.0

- Updated .neev/blueprints/authentication/architecture.md
- Migrated code from JWT to OAuth 2.0"
```

### ✅ Do: Quarterly Documentation Review

**Process:**
1. Run `neev inspect`
2. Review all blueprints
3. Update outdated specs
4. Archive completed blueprints
5. Regenerate `neev bridge > DOCS.md`

### ❌ Don't: Let Specs Become Stale

**Warning signs:**
- Specs describe features that don't exist
- Code has features not in specs
- Team stops referencing blueprints
- Drift warnings ignored

**Fix:**
- Dedicate time to update specs
- Make spec updates part of PR requirements
- Use `neev inspect --strict` in CI/CD

---

## Common Patterns

### Pattern 1: Feature Planning

```bash
# 1. Create blueprint
neev draft "New Feature"

# 2. Document thoroughly
# (edit intent.md and architecture.md)

# 3. Review with team
git add .neev/ && git commit -m "docs: add feature spec"
git push origin feature/new-feature-spec

# 4. Get approval on spec

# 5. Generate context
neev bridge --focus "new-feature" > context.md

# 6. Implement with AI assistance
# (paste context.md to AI)

# 7. Code review against spec
# (reviewers check: does code match blueprint?)

# 8. After deployment, archive
neev lay "new-feature"
```

### Pattern 2: API Design First

```bash
# 1. Create API blueprint
neev draft "Payments API"

# 2. Document all endpoints
# (detailed request/response examples)

# 3. Generate API docs
neev bridge --focus "Payments API" > API_SPEC.md

# 4. Share with frontend team
# (they can start integration before backend is done)

# 5. Backend implements to spec
# (API contract is already defined)
```

### Pattern 3: Documentation-Driven Onboarding

```bash
# 1. New developer joins
# 2. Read generated documentation
neev bridge > FULL_CONTEXT.md

# 3. Questions? Search blueprints
grep -r "authentication" .neev/blueprints/

# 4. Deep dive on specific area
neev bridge --focus "authentication" > auth-guide.md
```

### Pattern 4: Multi-Repo Coordination

```yaml
# In frontend repo neev.yaml
remotes:
  - name: backend-api
    path: "../backend/.neev/foundation"

# In backend repo neev.yaml
remotes:
  - name: frontend-contracts
    path: "../frontend/.neev/foundation"
```

```bash
# Both repos sync
neev sync-remotes

# Generate context with shared knowledge
neev bridge --with-remotes
```

---

## Anti-Patterns

### ❌ Anti-Pattern 1: Documentation Dump

**Problem:**
```markdown
# intent.md
Lorem ipsum dolor sit amet, consectetur adipiscing elit...
(5000 words of prose)
```

**Solution:**
```markdown
# Intent

## Purpose
[One clear sentence]

## Goals
- Goal 1
- Goal 2

## Success Criteria
- Criterion 1
- Criterion 2
```

### ❌ Anti-Pattern 2: Copy-Paste Duplication

**Problem:**
```
.neev/blueprints/user-api/architecture.md      # Includes auth details
.neev/blueprints/product-api/architecture.md   # Duplicates auth details
.neev/blueprints/order-api/architecture.md     # Duplicates auth details
```

**Solution:**
```
.neev/foundation/authentication.md    # Auth details here
.neev/blueprints/*/architecture.md    # Reference foundation
```

```markdown
# User API Architecture

## Authentication
See [Authentication](../../foundation/authentication.md) for details.

This API requires valid JWT token in Authorization header.
```

### ❌ Anti-Pattern 3: Implementation Details in Blueprints

**Problem:**
```markdown
## Implementation

\`\`\`typescript
// Import bcrypt
const bcrypt = require('bcrypt');

// Hash password function
async function hashPassword(password: string) {
  const salt = await bcrypt.genSalt(12);
  return await bcrypt.hash(password, salt);
}
\`\`\`
```

**Solution:**
```markdown
## Security

Passwords must be hashed using bcrypt with 12+ rounds.
Never store plain text passwords.
Use constant-time comparison for verification.
```

### ❌ Anti-Pattern 4: Specs After Implementation

**Problem:**
```bash
# Write code first
git commit -m "implement user auth"

# Then document
neev draft "User Auth"
git commit -m "add docs for auth (already built)"
```

**Solution:**
```bash
# Document first
neev draft "User Auth"
git commit -m "docs: user auth specification"

# Then implement
git commit -m "feat: implement user auth per spec"
```

### ❌ Anti-Pattern 5: Ignoring Drift

**Problem:**
```bash
$ neev inspect
⚠️  Drift detected
$ # Ignore and continue coding
```

**Solution:**
```bash
$ neev inspect --strict  # In CI/CD
# Fix drift before merging
```

### ❌ Anti-Pattern 6: Vague Blueprints

**Problem:**
```markdown
# Intent
Build a good authentication system that works well.

# Architecture
Use industry standard practices and secure methods.
```

**Solution:**
```markdown
# Intent

## Purpose
Implement JWT-based authentication for the REST API.

## Goals
- Support email/password login
- Token expiry: 1 hour
- Refresh token: 30 days
- Rate limit: 5 attempts per hour

# Architecture

## Token Format
JWT with HS256 algorithm.

**Claims:**
\`\`\`json
{
  "sub": "user-id",
  "role": "user|admin",
  "exp": 1234567890
}
\`\`\`

## Database Schema
\`\`\`sql
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE,
  password_hash VARCHAR(255),
  created_at TIMESTAMP
);
\`\`\`
```

### ❌ Anti-Pattern 7: Over-Engineering Blueprints

**Problem:**
```
.neev/blueprints/user-auth/
├── intent.md
├── architecture.md
├── security.md
├── performance.md
├── monitoring.md
├── deployment.md
├── testing.md
├── api-spec.md
├── database.md
└── diagrams/
    ├── sequence.md
    ├── component.md
    └── deployment.md
```

**Solution:**
```
.neev/blueprints/user-auth/
├── intent.md          # What, why, goals
└── architecture.md    # How, including security, performance notes
```

Keep it simple. Add more files only when necessary.

---

## Checklist for Quality Blueprints

Before finalizing a blueprint, check:

- [ ] **Clear purpose** — One sentence answers "what is this?"
- [ ] **Defined goals** — 3-5 specific, measurable goals
- [ ] **Success criteria** — How do we know it's done?
- [ ] **Scope defined** — What's included and excluded?
- [ ] **Examples included** — Sample requests, responses, code
- [ ] **Constraints documented** — Limitations, requirements
- [ ] **Diagrams where helpful** — Visualize complex flows
- [ ] **No implementation details** — Architecture, not code
- [ ] **Links to related specs** — Cross-reference other blueprints
- [ ] **Reviewed by team** — At least one other person read it

---

## Quick Reference

### When to Use Blueprints vs Foundation

| Use Blueprint For | Use Foundation For |
|-------------------|-------------------|
| Specific features | Project-wide standards |
| Single components | Tech stack decisions |
| Time-bound work | Lasting principles |
| Detailed designs | High-level architecture |
| What you're building now | How you build everything |

### Command Quick Reference

```bash
# Common workflows
neev draft "Feature"              # Plan new work
neev bridge --focus "Feature"     # Get focused context
neev inspect                      # Check drift
neev lay "Feature"                # Archive when done
neev instructions                 # Update Copilot

# Team workflows
neev bridge > DOCS.md             # Generate docs
neev sync-remotes                 # Sync multi-repo
neev inspect --json --strict      # CI/CD check

# AI workflows
neev bridge --claude              # Claude format
neev bridge --focus "api"         # Focused context
```

---

## Additional Resources

- [Getting Started Guide](GETTING_STARTED.md) — Learn the basics
- [API Reference](API_REFERENCE.md) — Complete command docs
- [Tutorials](TUTORIALS.md) — Step-by-step examples
- [FAQ](FAQ.md) — Common questions

---

**Have a best practice to share?** [Open a discussion](https://github.com/neev-kit/neev/discussions) or [contribute](CONTRIBUTING.md)!
