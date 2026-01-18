# 📖 Neev Tutorials

Step-by-step tutorials for common Neev use cases. Each tutorial is self-contained and includes working examples.

## Table of Contents

1. [Tutorial 1: Building a REST API with Neev](#tutorial-1-building-a-rest-api-with-neev)
2. [Tutorial 2: Microservices Architecture Planning](#tutorial-2-microservices-architecture-planning)
3. [Tutorial 3: Team Onboarding Documentation](#tutorial-3-team-onboarding-documentation)
4. [Tutorial 4: CI/CD Integration](#tutorial-4-cicd-integration)
5. [Tutorial 5: Multi-Repo Projects with Remotes](#tutorial-5-multi-repo-projects-with-remotes)
6. [Tutorial 6: Using Neev with GitHub Copilot](#tutorial-6-using-neev-with-github-copilot)
7. [Tutorial 7: Migration from Existing Documentation](#tutorial-7-migration-from-existing-documentation)
8. [Tutorial 8: Database Schema Evolution](#tutorial-8-database-schema-evolution)

---

## Tutorial 1: Building a REST API with Neev

**Goal:** Plan and document a REST API for a blog platform using Neev, then use AI to implement it.

**Time:** 20 minutes

### Prerequisites

- Neev installed ([Installation guide](GETTING_STARTED.md#installation))
- A project directory
- Access to an AI coding assistant (Claude, Cursor, GitHub Copilot, etc.)

### Step 1: Initialize Your Project

```bash
mkdir blog-api
cd blog-api
neev init
```

**Expected output:**
```
🏗️  Laying foundation in /Users/you/blog-api
✅ Foundation laid successfully!
```

### Step 2: Define Your Tech Stack

Create foundation documents to establish project-wide decisions:

```bash
cat > .neev/foundation/stack.md << 'EOF'
# Technology Stack

## Backend
- **Language:** Node.js (v18+)
- **Framework:** Express.js
- **Database:** PostgreSQL 14
- **ORM:** Prisma

## API Design
- RESTful architecture
- JSON responses
- JWT authentication
- API versioning (/v1/)

## Development Tools
- ESLint for linting
- Jest for testing
- Prettier for formatting
EOF
```

### Step 3: Define Development Principles

```bash
cat > .neev/foundation/principles.md << 'EOF'
# Development Principles

## Code Quality
- Write tests for all endpoints
- Use TypeScript for type safety
- Follow REST conventions
- Document all public APIs

## Security
- Never expose sensitive data
- Validate all input
- Use parameterized queries
- Implement rate limiting

## Performance
- Paginate list endpoints
- Use database indexes
- Cache where appropriate
- Monitor query performance
EOF
```

### Step 4: Create API Blueprints

Create blueprints for each major API component:

```bash
# User management
neev draft "User Management API"

# Posts/Articles
neev draft "Posts API"

# Comments
neev draft "Comments API"

# Authentication
neev draft "Authentication API"
```

### Step 5: Document User Management API

Edit `.neev/blueprints/user-management-api/intent.md`:

```markdown
# User Management API Intent

## Purpose
Provide CRUD operations for user accounts in the blog platform.

## Goals
- Allow user registration
- Support profile updates
- Enable user deletion (soft delete)
- List all users (admin only)

## Out of Scope
- Authentication (handled by Authentication API)
- Password reset (phase 2)
- Email verification (phase 2)

## Success Criteria
- Users can create accounts with email and password
- Users can update their profiles (name, bio, avatar)
- Admins can list all users with pagination
- Soft deletion preserves data integrity
```

Edit `.neev/blueprints/user-management-api/architecture.md`:

```markdown
# User Management API Architecture

## Endpoints

### POST /v1/users
Create a new user account.

**Request:**
\`\`\`json
{
  "email": "user@example.com",
  "password": "secure_password",
  "name": "John Doe"
}
\`\`\`

**Response (201):**
\`\`\`json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe",
  "createdAt": "2024-01-01T00:00:00Z"
}
\`\`\`

### GET /v1/users/:id
Retrieve user by ID.

**Response (200):**
\`\`\`json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe",
  "bio": "Software developer",
  "avatarUrl": "https://...",
  "createdAt": "2024-01-01T00:00:00Z"
}
\`\`\`

### PUT /v1/users/:id
Update user profile (authenticated user only).

**Request:**
\`\`\`json
{
  "name": "Jane Doe",
  "bio": "Updated bio"
}
\`\`\`

### DELETE /v1/users/:id
Soft delete user account (admin or owner only).

**Response (204):** No content

## Database Schema

\`\`\`sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  bio TEXT,
  avatar_url TEXT,
  deleted_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
\`\`\`

## Security Considerations
- Hash passwords with bcrypt (12 rounds)
- Validate email format
- Enforce password strength (min 8 chars)
- Never return password hashes in responses
- Rate limit registration (5 per hour per IP)

## Error Handling
- 400: Invalid input
- 401: Unauthorized
- 403: Forbidden
- 404: User not found
- 409: Email already exists
- 500: Server error
```

### Step 6: Document Posts API

Edit `.neev/blueprints/posts-api/intent.md`:

```markdown
# Posts API Intent

## Purpose
Enable users to create, read, update, and delete blog posts.

## Goals
- Support markdown content
- Allow draft and published states
- Enable post categorization with tags
- Implement search functionality
- Support pagination

## Success Criteria
- Users can create and edit their posts
- Posts can be saved as drafts
- Published posts are visible to all users
- Posts can be filtered by tags
- Search works across title and content
```

Edit `.neev/blueprints/posts-api/architecture.md`:

```markdown
# Posts API Architecture

## Endpoints

### POST /v1/posts
Create a new post (authenticated users).

**Request:**
\`\`\`json
{
  "title": "My First Post",
  "content": "# Introduction\n\nContent here...",
  "tags": ["tutorial", "nodejs"],
  "status": "draft"
}
\`\`\`

### GET /v1/posts
List all published posts with pagination.

**Query Parameters:**
- `page` (default: 1)
- `limit` (default: 20, max: 100)
- `tag` (filter by tag)
- `search` (search in title/content)

**Response (200):**
\`\`\`json
{
  "posts": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "totalPages": 8
  }
}
\`\`\`

### GET /v1/posts/:id
Retrieve a single post by ID.

### PUT /v1/posts/:id
Update post (author only).

### DELETE /v1/posts/:id
Delete post (author or admin only).

## Database Schema

\`\`\`sql
CREATE TABLE posts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  status VARCHAR(20) DEFAULT 'draft',
  published_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE tags (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE post_tags (
  post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
  tag_id UUID REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (post_id, tag_id)
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_posts_published_at ON posts(published_at);
\`\`\`
```

### Step 7: Generate Context for AI

Now aggregate all your specifications:

```bash
# Full context
neev bridge > api-context.md

# Or focus on specific area
neev bridge --focus "User Management" > user-api-context.md
```

### Step 8: Use with AI Assistant

Open your AI coding assistant (Claude, Cursor, etc.) and provide the context:

```
I have documented a REST API specification using Neev. Here's the complete context:

[paste contents of api-context.md]

Please implement the User Management API according to these specifications. 
Include:
1. Express.js routes
2. Prisma schema
3. Input validation with Joi
4. Unit tests with Jest
5. Error handling middleware
```

### Step 9: Implement and Iterate

As your AI generates code:

1. **Review** — Check if it matches your specs
2. **Test** — Run the generated tests
3. **Refine** — Update blueprints if needed
4. **Regenerate** — Run `neev bridge` again and ask AI to adjust

### Step 10: Check for Drift

As you build, check if implementation matches specs:

```bash
neev inspect
```

### Result

You now have:
- ✅ Complete API specification
- ✅ AI-generated implementation
- ✅ Documentation in version control
- ✅ Tests and validation
- ✅ Team alignment on design

---

## Tutorial 2: Microservices Architecture Planning

**Goal:** Plan a microservices architecture for an e-commerce platform.

**Time:** 30 minutes

### Overview

We'll design:
- User Service
- Product Service
- Order Service
- Payment Service
- API Gateway
- Message Queue integration

### Step 1: Initialize and Create Foundation

```bash
mkdir ecommerce-platform
cd ecommerce-platform
neev init
```

Create architecture principles:

```bash
cat > .neev/foundation/architecture.md << 'EOF'
# Microservices Architecture Principles

## Service Design
- Each service owns its database
- Services communicate via events
- RESTful APIs for synchronous calls
- Message queue for async operations

## Technology Stack
- **Services:** Node.js with Express
- **Gateway:** Kong or custom Node.js
- **Message Queue:** RabbitMQ
- **Databases:** PostgreSQL per service
- **Cache:** Redis
- **Container:** Docker
- **Orchestration:** Kubernetes

## Communication Patterns
- Synchronous: REST APIs
- Asynchronous: Event-driven via RabbitMQ
- Service discovery: Kubernetes DNS
- Load balancing: Kubernetes Services

## Data Consistency
- Eventual consistency where possible
- Saga pattern for distributed transactions
- Event sourcing for audit trails
EOF
```

### Step 2: Create Service Blueprints

```bash
neev draft "User Service"
neev draft "Product Service"
neev draft "Order Service"
neev draft "Payment Service"
neev draft "API Gateway"
neev draft "Event Bus"
```

### Step 3: Document User Service

`.neev/blueprints/user-service/intent.md`:

```markdown
# User Service Intent

## Responsibility
Manage user accounts, authentication, and profiles.

## Owns
- User database
- Authentication tokens
- User preferences
- Profile information

## APIs Provided
- `POST /users` — Register
- `GET /users/:id` — Get profile
- `PUT /users/:id` — Update profile
- `POST /auth/login` — Login
- `POST /auth/logout` — Logout

## Events Published
- `user.created`
- `user.updated`
- `user.deleted`

## Events Consumed
- None (independent service)
```

`.neev/blueprints/user-service/architecture.md`:

```markdown
# User Service Architecture

## Database Schema
\`\`\`sql
CREATE DATABASE user_service;

CREATE TABLE users (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  name VARCHAR(255),
  status VARCHAR(20) DEFAULT 'active',
  created_at TIMESTAMP DEFAULT NOW()
);
\`\`\`

## Event Publishing
When a user is created:
\`\`\`json
{
  "event": "user.created",
  "data": {
    "userId": "uuid",
    "email": "user@example.com",
    "timestamp": "2024-01-01T00:00:00Z"
  }
}
\`\`\`

## API Endpoints

### POST /users
Create new user.

### GET /users/:id
Get user by ID.

### Authentication Flow
1. Client sends credentials to `/auth/login`
2. Service validates against database
3. Issues JWT token (1 hour expiry)
4. Client includes token in `Authorization` header
```

### Step 4: Document Order Service

`.neev/blueprints/order-service/intent.md`:

```markdown
# Order Service Intent

## Responsibility
Manage customer orders and order lifecycle.

## Owns
- Order database
- Order state management
- Order history

## APIs Provided
- `POST /orders` — Create order
- `GET /orders/:id` — Get order
- `GET /orders` — List user orders
- `PUT /orders/:id/cancel` — Cancel order

## Events Published
- `order.created`
- `order.confirmed`
- `order.cancelled`
- `order.completed`

## Events Consumed
- `payment.completed` — Confirm order
- `payment.failed` — Cancel order
- `product.out_of_stock` — Update order status
```

### Step 5: Document Cross-Service Communication

`.neev/foundation/integration.md`:

```markdown
# Service Integration Patterns

## Saga: Order Processing

### Happy Path
1. User creates order → Order Service
2. Order Service publishes `order.created`
3. Payment Service consumes event
4. Payment Service processes payment
5. Payment Service publishes `payment.completed`
6. Order Service consumes event
7. Order Service publishes `order.confirmed`
8. Product Service updates inventory

### Failure Path
1. Payment fails
2. Payment Service publishes `payment.failed`
3. Order Service cancels order
4. User is notified

## API Gateway Routes
\`\`\`yaml
/users/** → User Service
/products/** → Product Service
/orders/** → Order Service
/payments/** → Payment Service
\`\`\`

## Service Discovery
Services register with Kubernetes:
\`\`\`yaml
apiVersion: v1
kind: Service
metadata:
  name: user-service
spec:
  selector:
    app: user-service
  ports:
  - port: 3000
\`\`\`
```

### Step 6: Generate Architecture Document

```bash
neev bridge > ARCHITECTURE_SPEC.md
```

### Step 7: Create Service Diagrams

Add to `.neev/foundation/diagrams.md`:

```markdown
# System Diagrams

## Service Dependencies
\`\`\`
API Gateway
├── User Service
├── Product Service
├── Order Service
│   ├─> Product Service (inventory check)
│   └─> Payment Service (payment processing)
└── Payment Service
    └─> Order Service (payment confirmation)
\`\`\`

## Event Flow: Order Creation
\`\`\`
Client → API Gateway → Order Service
                       ├─> order.created event
                       │
Payment Service <──────┘
       │
       ├─> payment.completed event
       │
Order Service <────────┘
       │
       └─> order.confirmed event
\`\`\`
```

### Step 8: Use for Implementation

Share with your team or AI:

```bash
neev bridge > microservices-spec.md
# Share microservices-spec.md with development team
```

---

## Tutorial 3: Team Onboarding Documentation

**Goal:** Create comprehensive onboarding documentation for new developers.

**Time:** 25 minutes

### Step 1: Initialize in Existing Project

```bash
cd /path/to/existing/project
neev init
```

### Step 2: Document Project Overview

```bash
cat > .neev/foundation/overview.md << 'EOF'
# Project Overview

## What We Build
A SaaS platform for project management with real-time collaboration.

## Our Mission
Help distributed teams work together seamlessly.

## Target Users
- Small to medium businesses
- Remote teams
- Project managers
- Developers

## Key Features
- Real-time task updates
- Team chat
- File sharing
- Time tracking
- Reporting dashboards
EOF
```

### Step 3: Document Tech Stack

```bash
cat > .neev/foundation/stack.md << 'EOF'
# Technology Stack

## Frontend
- React 18 with TypeScript
- Next.js for SSR
- TailwindCSS for styling
- Zustand for state management
- React Query for data fetching

## Backend
- Node.js with Express
- PostgreSQL database
- Prisma ORM
- Redis for caching
- WebSocket for real-time features

## Infrastructure
- AWS ECS for hosting
- RDS for database
- S3 for file storage
- CloudFront for CDN
- GitHub Actions for CI/CD

## Development Tools
- ESLint + Prettier
- Jest + React Testing Library
- Playwright for E2E tests
- Docker for local development
EOF
```

### Step 4: Document Development Workflow

```bash
cat > .neev/foundation/workflow.md << 'EOF'
# Development Workflow

## Getting Started

### 1. Clone Repository
\`\`\`bash
git clone https://github.com/company/project
cd project
\`\`\`

### 2. Install Dependencies
\`\`\`bash
npm install
\`\`\`

### 3. Setup Environment
\`\`\`bash
cp .env.example .env
# Edit .env with your credentials
\`\`\`

### 4. Run Database Migrations
\`\`\`bash
npm run db:migrate
npm run db:seed
\`\`\`

### 5. Start Development Server
\`\`\`bash
npm run dev
# Frontend: http://localhost:3000
# Backend: http://localhost:3001
\`\`\`

## Daily Workflow

### Create Feature Branch
\`\`\`bash
git checkout main
git pull
git checkout -b feature/your-feature-name
\`\`\`

### Make Changes
1. Write code
2. Write tests
3. Run tests: `npm test`
4. Run linter: `npm run lint`

### Submit Pull Request
1. Push branch: `git push origin feature/your-feature-name`
2. Create PR on GitHub
3. Request review from team
4. Address feedback
5. Merge when approved

## Code Review Checklist
- [ ] Tests pass
- [ ] Linter passes
- [ ] Code follows style guide
- [ ] Documentation updated
- [ ] No console.logs left
- [ ] Error handling included
EOF
```

### Step 5: Document Conventions

```bash
cat > .neev/foundation/conventions.md << 'EOF'
# Coding Conventions

## File Naming
- Components: `PascalCase.tsx` (e.g., `UserProfile.tsx`)
- Utilities: `camelCase.ts` (e.g., `formatDate.ts`)
- Tests: `*.test.ts` or `*.test.tsx`
- Styles: `ComponentName.module.css`

## Component Structure
\`\`\`typescript
// Imports
import React from 'react';
import { useQuery } from 'react-query';

// Types
interface UserProfileProps {
  userId: string;
}

// Component
export const UserProfile: React.FC<UserProfileProps> = ({ userId }) => {
  // Hooks
  const { data, isLoading } = useQuery(['user', userId], fetchUser);

  // Early returns
  if (isLoading) return <Loading />;
  
  // Main render
  return (
    <div>
      {/* Component JSX */}
    </div>
  );
};
\`\`\`

## API Naming
- Use RESTful conventions
- Plural nouns: `/users`, `/posts`
- ID parameters: `/users/:id`
- Actions: `/users/:id/activate`

## Database Naming
- Tables: `snake_case` plural (e.g., `user_profiles`)
- Columns: `snake_case` (e.g., `created_at`)
- Foreign keys: `table_id` (e.g., `user_id`)

## Git Commit Messages
Follow conventional commits:
- `feat: add user profile page`
- `fix: resolve login redirect issue`
- `docs: update API documentation`
- `test: add tests for auth flow`
- `refactor: simplify data fetching logic`
EOF
```

### Step 6: Create Feature Documentation Blueprints

```bash
neev draft "Authentication System"
neev draft "Real-time Updates"
neev draft "File Upload System"
```

Document each with current implementation details.

### Step 7: Generate Onboarding Guide

```bash
neev bridge > ONBOARDING.md
```

### Step 8: Add to Repository

```bash
git add .neev/ ONBOARDING.md
git commit -m "docs: add Neev-based onboarding documentation"
git push
```

### Result

New team members can now:
- Read `ONBOARDING.md` for complete context
- Understand project architecture
- Follow established conventions
- See how features are designed

---

## Tutorial 4: CI/CD Integration

**Goal:** Integrate Neev drift detection into your CI/CD pipeline.

**Time:** 15 minutes

### Prerequisites

- Neev initialized in your project
- GitHub repository with Actions enabled

### Step 1: Create Drift Detection Workflow

Create `.github/workflows/spec-drift.yml`:

```yaml
name: Specification Drift Check

on:
  pull_request:
    branches: [main, develop]
  push:
    branches: [main]

jobs:
  drift-detection:
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
      
      - name: Download Neev
        run: |
          curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_linux_amd64.tar.gz | tar xz
          chmod +x neev
          sudo mv neev /usr/local/bin/
      
      - name: Verify Neev installation
        run: neev --version
      
      - name: Check for specification drift
        run: |
          neev inspect --json --strict > drift-report.json
          cat drift-report.json
        continue-on-error: true
        id: drift-check
      
      - name: Upload drift report
        uses: actions/upload-artifact@v3
        with:
          name: drift-report
          path: drift-report.json
      
      - name: Comment on PR
        if: github.event_name == 'pull_request' && steps.drift-check.outcome == 'failure'
        uses: actions/github-script@v6
        with:
          script: |
            const fs = require('fs');
            const drift = JSON.parse(fs.readFileSync('drift-report.json', 'utf8'));
            
            const comment = `## ⚠️ Specification Drift Detected
            
            The code implementation has drifted from specifications:
            
            **Missing implementations:**
            ${drift.drift.missing.map(m => `- ${m}`).join('\n') || 'None'}
            
            **Extra code (no specs):**
            ${drift.drift.extra.map(e => `- ${e}`).join('\n') || 'None'}
            
            Please update specifications or code to align.`;
            
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: comment
            });
      
      - name: Fail if drift detected
        if: steps.drift-check.outcome == 'failure'
        run: exit 1
```

### Step 2: Create Documentation Generation Workflow

Create `.github/workflows/update-docs.yml`:

```yaml
name: Update Documentation

on:
  push:
    branches: [main]
    paths:
      - '.neev/**'

jobs:
  generate-docs:
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Download Neev
        run: |
          curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_linux_amd64.tar.gz | tar xz
          chmod +x neev
          sudo mv neev /usr/local/bin/
      
      - name: Generate architecture documentation
        run: |
          neev bridge > docs/ARCHITECTURE.md
          neev bridge --claude > docs/ARCHITECTURE_CLAUDE.md
      
      - name: Generate Copilot instructions
        run: neev instructions
      
      - name: Commit updated documentation
        run: |
          git config --local user.email "github-actions[bot]@users.noreply.github.com"
          git config --local user.name "github-actions[bot]"
          git add docs/ .github/copilot-instructions.md
          git diff --staged --quiet || git commit -m "docs: auto-update from Neev specs"
          git push
```

### Step 3: Add Status Badge to README

Add to your `README.md`:

```markdown
[![Spec Drift](https://github.com/your-org/your-repo/actions/workflows/spec-drift.yml/badge.svg)](https://github.com/your-org/your-repo/actions/workflows/spec-drift.yml)
```

### Step 4: Test the Workflow

```bash
# Commit and push
git add .github/workflows/
git commit -m "ci: add Neev drift detection"
git push

# Check Actions tab in GitHub
```

### Result

Now your CI/CD:
- ✅ Checks for drift on every PR
- ✅ Comments on PRs when drift is detected
- ✅ Auto-generates documentation
- ✅ Updates Copilot instructions
- ✅ Enforces spec-code alignment

---

## Tutorial 5: Multi-Repo Projects with Remotes

**Goal:** Share foundation documents across multiple repositories.

**Time:** 20 minutes

### Scenario

You have:
- `shared-standards` repo — Company-wide standards
- `backend-api` repo — Backend service
- `frontend-app` repo — Frontend application

### Step 1: Setup Shared Standards Repo

```bash
mkdir shared-standards
cd shared-standards
neev init

cat > .neev/foundation/coding-standards.md << 'EOF'
# Company Coding Standards

## Code Review
- All code must be reviewed
- Tests required for all features
- Documentation required for APIs

## Security
- Never commit secrets
- Use environment variables
- Validate all inputs
- Follow OWASP guidelines
EOF

cat > .neev/foundation/git-workflow.md << 'EOF'
# Git Workflow

## Branch Naming
- Features: `feature/description`
- Bugs: `fix/description`
- Hotfixes: `hotfix/description`

## Commit Messages
Use conventional commits:
- `feat:` new features
- `fix:` bug fixes
- `docs:` documentation
- `test:` tests
- `refactor:` refactoring
EOF

git add .neev/
git commit -m "docs: add shared standards"
git push
```

### Step 2: Setup Backend API with Remote

```bash
cd ../backend-api
neev init

# Configure remote in neev.yaml
cat > .neev/neev.yaml << 'EOF'
project_name: "Backend API"
foundation_path: ".neev"

remotes:
  - name: company-standards
    path: "../shared-standards/.neev/foundation"
    public_only: true

ignore_dirs:
  - node_modules
  - dist
  - .git
EOF
```

### Step 3: Sync Remote Foundations

```bash
neev sync-remotes
```

**Output:**
```
🔄 Syncing remote foundations...
✅ Synced company-standards (2 files)
```

**Check synced files:**
```bash
tree .neev/remotes/
# .neev/remotes/
# └── company-standards/
#     ├── coding-standards.md
#     └── git-workflow.md
```

### Step 4: Create Local Blueprints

```bash
neev draft "Authentication API"
neev draft "Database Schema"

# Document your backend-specific specs
```

### Step 5: Generate Context with Remotes

```bash
# Include remote standards
neev bridge --with-remotes > full-context.md

# The output includes:
# 1. Local foundation docs
# 2. Remote foundation docs (from shared-standards)
# 3. Local blueprint docs
```

### Step 6: Setup Frontend with Same Remote

```bash
cd ../frontend-app
neev init

cat > .neev/neev.yaml << 'EOF'
project_name: "Frontend App"
foundation_path: ".neev"

remotes:
  - name: company-standards
    path: "../shared-standards/.neev/foundation"
    public_only: true
  - name: backend-api-docs
    path: "../backend-api/.neev/foundation"
    public_only: true

ignore_dirs:
  - node_modules
  - .next
  - .git
EOF

neev sync-remotes
```

### Step 7: Cross-Repo Context

Now frontend can include backend context:

```bash
neev bridge --with-remotes > frontend-context.md

# Includes:
# - Company standards
# - Backend API foundation (if any)
# - Frontend blueprints
```

### Step 8: Automate Sync in CI/CD

Add to `.github/workflows/sync-specs.yml`:

```yaml
name: Sync Remote Specifications

on:
  schedule:
    - cron: '0 0 * * *'  # Daily
  workflow_dispatch:

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          token: ${{ secrets.PAT_TOKEN }}
          submodules: recursive
      
      - name: Checkout shared standards
        uses: actions/checkout@v3
        with:
          repository: your-org/shared-standards
          path: ./shared-standards
          token: ${{ secrets.PAT_TOKEN }}
      
      - name: Install Neev
        run: |
          curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_linux_amd64.tar.gz | tar xz
          sudo mv neev /usr/local/bin/
      
      - name: Sync remotes
        run: neev sync-remotes
      
      - name: Commit if changed
        run: |
          git config --local user.email "github-actions[bot]@users.noreply.github.com"
          git config --local user.name "github-actions[bot]"
          git add .neev/remotes/
          git diff --staged --quiet || git commit -m "chore: sync remote specifications"
          git push
```

### Result

- ✅ Shared standards across repos
- ✅ Automatic syncing
- ✅ Consistent documentation
- ✅ Cross-repo context aggregation

---

## Tutorial 6: Using Neev with GitHub Copilot

**Goal:** Enhance GitHub Copilot suggestions using Neev specifications.

**Time:** 15 minutes

### Step 1: Initialize Neev

```bash
cd your-project
neev init
```

### Step 2: Create Detailed Foundation

```bash
cat > .neev/foundation/coding-patterns.md << 'EOF'
# Preferred Coding Patterns

## Error Handling
Always use custom error classes:

\`\`\`typescript
class ValidationError extends Error {
  constructor(message: string, public field: string) {
    super(message);
    this.name = 'ValidationError';
  }
}

throw new ValidationError('Invalid email', 'email');
\`\`\`

## Async/Await
Always use try-catch with async:

\`\`\`typescript
async function fetchUser(id: string) {
  try {
    const response = await api.get(\`/users/\${id}\`);
    return response.data;
  } catch (error) {
    logger.error('Failed to fetch user', { id, error });
    throw new FetchError('User not found');
  }
}
\`\`\`

## React Components
Use functional components with TypeScript:

\`\`\`typescript
interface UserCardProps {
  user: User;
  onSelect: (user: User) => void;
}

export const UserCard: React.FC<UserCardProps> = ({ user, onSelect }) => {
  return (
    <div onClick={() => onSelect(user)}>
      <h3>{user.name}</h3>
      <p>{user.email}</p>
    </div>
  );
};
\`\`\`
EOF
```

### Step 3: Create API Blueprints

```bash
neev draft "User API Client"

cat > .neev/blueprints/user-api-client/architecture.md << 'EOF'
# User API Client Architecture

## Module Structure
\`\`\`typescript
// src/api/users.ts
import { apiClient } from './client';

export interface User {
  id: string;
  email: string;
  name: string;
  role: 'user' | 'admin';
}

export const usersApi = {
  async getUser(id: string): Promise<User> {
    const response = await apiClient.get(\`/users/\${id}\`);
    return response.data;
  },
  
  async listUsers(): Promise<User[]> {
    const response = await apiClient.get('/users');
    return response.data;
  },
  
  async createUser(data: Omit<User, 'id'>): Promise<User> {
    const response = await apiClient.post('/users', data);
    return response.data;
  }
};
\`\`\`

## Error Handling
All API methods should catch and transform errors:

\`\`\`typescript
try {
  return await apiClient.get(url);
} catch (error) {
  if (error.response?.status === 404) {
    throw new NotFoundError('User not found');
  }
  throw new ApiError('Failed to fetch user');
}
\`\`\`
EOF
```

### Step 4: Generate Copilot Instructions

```bash
neev instructions
```

This creates `.github/copilot-instructions.md`:

```markdown
# GitHub Copilot Instructions

## Project Context
This project uses TypeScript, React, and RESTful APIs.

## Foundation Modules

### coding-patterns.md
Preferred coding patterns for error handling, async/await, and React components.

## Active Blueprints

### user-api-client
User API client with TypeScript interfaces and error handling.

## Development Guidelines
- Use custom error classes
- Always use try-catch with async functions
- Functional React components with TypeScript
- Follow the patterns shown in foundation documents
```

### Step 5: Commit to Version Control

```bash
git add .neev/ .github/copilot-instructions.md
git commit -m "docs: add Neev specifications for Copilot"
git push
```

### Step 6: Use in Your IDE

GitHub Copilot will now automatically:
- Read `.github/copilot-instructions.md`
- Suggest code following your patterns
- Use your preferred error handling
- Match your component structure

### Step 7: Test Copilot Understanding

Open a new file in your IDE and start typing:

```typescript
// Type this comment:
// Create a user API client method to update a user

// Copilot should suggest something like:
async updateUser(id: string, data: Partial<User>): Promise<User> {
  try {
    const response = await apiClient.put(`/users/${id}`, data);
    return response.data;
  } catch (error) {
    logger.error('Failed to update user', { id, error });
    throw new ApiError('Failed to update user');
  }
}
```

### Step 8: Keep Instructions Updated

Whenever you update blueprints:

```bash
# Update blueprint
vim .neev/blueprints/user-api-client/architecture.md

# Regenerate instructions
neev instructions

# Commit
git add .neev/ .github/copilot-instructions.md
git commit -m "docs: update API client specifications"
git push
```

### Result

GitHub Copilot now:
- ✅ Suggests code matching your patterns
- ✅ Uses your error handling approach
- ✅ Follows your TypeScript conventions
- ✅ Understands your project structure

---

## Tutorial 7: Migration from Existing Documentation

**Goal:** Convert existing README/Wiki documentation to Neev format.

**Time:** 30 minutes

### Scenario

You have a project with scattered documentation:
- `README.md` — Setup instructions
- `docs/api.md` — API documentation
- `docs/architecture.md` — System design
- `docs/deployment.md` — Deployment guide

### Step 1: Initialize Neev

```bash
cd your-project
neev init
```

### Step 2: Analyze Existing Documentation

```bash
# See what you have
find docs -name "*.md"
```

### Step 3: Migrate General Documentation to Foundation

```bash
# Copy architecture docs
cp docs/architecture.md .neev/foundation/

# Copy deployment docs
cp docs/deployment.md .neev/foundation/

# Extract tech stack from README
# (manually create this)
cat > .neev/foundation/stack.md << 'EOF'
# Technology Stack

[Extract relevant sections from README.md]
EOF
```

### Step 4: Convert API Documentation to Blueprints

Your `docs/api.md` might have sections for different endpoints. Split them:

```bash
# Create blueprint for User API
neev draft "User API"

# Extract user-related endpoints from docs/api.md
# Copy to .neev/blueprints/user-api/architecture.md

# Create blueprint for Products API
neev draft "Products API"

# Extract product-related endpoints
# Copy to .neev/blueprints/products-api/architecture.md
```

### Step 5: Add Intent to Blueprints

For each blueprint, add context to `intent.md`:

```bash
cat > .neev/blueprints/user-api/intent.md << 'EOF'
# User API Intent

## Purpose
Provide CRUD operations for user management.

## Goals
- Support user registration and login
- Enable profile management
- Provide admin user listing

## Scope
Covers all `/api/users` endpoints.
EOF
```

### Step 6: Update README

Replace detailed sections in README with references:

```markdown
# Your Project

Brief description...

## Documentation

- **Getting Started**: See [GETTING_STARTED.md](GETTING_STARTED.md)
- **Architecture**: See `.neev/foundation/architecture.md`
- **API Reference**: See blueprints in `.neev/blueprints/`
- **Deployment**: See `.neev/foundation/deployment.md`

## Generate Full Context

\`\`\`bash
neev bridge > FULL_DOCUMENTATION.md
\`\`\`

## For Developers

\`\`\`bash
# Get focused context
neev bridge --focus "API" > api-docs.md
neev bridge --focus "deployment" > deployment-guide.md
\`\`\`
```

### Step 7: Cleanup Old Documentation

```bash
# Create backup
mkdir docs/archive
mv docs/*.md docs/archive/

# Keep only necessary docs
# Most content now lives in .neev/
```

### Step 8: Generate Aggregated Documentation

```bash
# Create comprehensive docs
neev bridge > DOCUMENTATION.md

# Add to git
git add .neev/ DOCUMENTATION.md README.md
git commit -m "docs: migrate to Neev structure"
```

### Step 9: Setup Auto-Generation

Add to `.github/workflows/docs.yml`:

```yaml
name: Generate Documentation

on:
  push:
    branches: [main]
    paths:
      - '.neev/**'

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Install Neev
        run: |
          curl -L https://github.com/neev-kit/neev/releases/latest/download/neev_linux_amd64.tar.gz | tar xz
          sudo mv neev /usr/local/bin/
      - name: Generate docs
        run: |
          neev bridge > DOCUMENTATION.md
          neev instructions
      - name: Commit
        run: |
          git config --local user.email "github-actions[bot]@users.noreply.github.com"
          git config --local user.name "github-actions[bot]"
          git add DOCUMENTATION.md .github/copilot-instructions.md
          git diff --staged --quiet || git commit -m "docs: auto-generate from Neev"
          git push
```

### Result

- ✅ Centralized documentation in `.neev/`
- ✅ Structured blueprints per feature
- ✅ Auto-generated comprehensive docs
- ✅ Version-controlled specifications
- ✅ AI-ready context

---

## Tutorial 8: Database Schema Evolution

**Goal:** Document and track database schema changes using Neev.

**Time:** 20 minutes

### Step 1: Initialize and Create Blueprint

```bash
cd your-project
neev init
neev draft "Database Schema"
```

### Step 2: Document Initial Schema

`.neev/blueprints/database-schema/intent.md`:

```markdown
# Database Schema Intent

## Purpose
Define the database structure for the application.

## Goals
- Support user management
- Track orders and products
- Enable reporting
- Ensure data integrity

## Design Principles
- Normalize where appropriate
- Use foreign keys for relationships
- Add indexes for query performance
- Include audit timestamps
```

`.neev/blueprints/database-schema/architecture.md`:

```markdown
# Database Schema Architecture

## Schema Version: 1.0.0

### Users Table
\`\`\`sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  role VARCHAR(50) DEFAULT 'user',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
\`\`\`

### Products Table
\`\`\`sql
CREATE TABLE products (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  description TEXT,
  price DECIMAL(10, 2) NOT NULL,
  stock_quantity INT DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_products_name ON products(name);
\`\`\`

### Orders Table
\`\`\`sql
CREATE TABLE orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  status VARCHAR(50) DEFAULT 'pending',
  total_amount DECIMAL(10, 2) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
\`\`\`

### Order Items Table
\`\`\`sql
CREATE TABLE order_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
  product_id UUID REFERENCES products(id),
  quantity INT NOT NULL,
  price DECIMAL(10, 2) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);
\`\`\`

## Relationships
- Users have many Orders (1:N)
- Orders have many Order Items (1:N)
- Products appear in many Order Items (1:N)

## Indexes Strategy
- Primary keys automatically indexed
- Foreign keys indexed for join performance
- Email indexed for login lookups
- Status indexed for filtering
- Timestamps indexed for date range queries

## Migration Strategy
- Use Prisma for schema migrations
- Always create migration files
- Test migrations on staging first
- Keep rollback scripts ready
```

### Step 3: Track Schema Changes

Create a changelog:

`.neev/blueprints/database-schema/changelog.md`:

```markdown
# Database Schema Changelog

## [1.0.0] - 2024-01-15
### Added
- Initial schema with users, products, orders, order_items tables
- Indexes on foreign keys and frequently queried fields

## [Planned] 2.0.0
### To Add
- Reviews table for product reviews
- Wishlists table for user favorites
- Audit log table for compliance
```

### Step 4: Document Schema Evolution

When making changes, create a new blueprint:

```bash
neev draft "Schema Migration v2"
```

`.neev/blueprints/schema-migration-v2/intent.md`:

```markdown
# Schema Migration v2 Intent

## Purpose
Add product reviews and wishlist functionality.

## Changes
- Add `reviews` table
- Add `wishlists` table
- Modify `products` table to track review count

## Backward Compatibility
- All existing queries continue to work
- New columns have defaults
- No breaking changes
```

`.neev/blueprints/schema-migration-v2/architecture.md`:

```markdown
# Schema Migration v2 Architecture

## New Tables

### Reviews Table
\`\`\`sql
CREATE TABLE reviews (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id UUID REFERENCES products(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  rating INT CHECK (rating >= 1 AND rating <= 5),
  comment TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_reviews_product_id ON reviews(product_id);
CREATE INDEX idx_reviews_user_id ON reviews(user_id);
CREATE INDEX idx_reviews_rating ON reviews(rating);

-- Constraint: One review per user per product
CREATE UNIQUE INDEX idx_reviews_unique ON reviews(product_id, user_id);
\`\`\`

### Wishlists Table
\`\`\`sql
CREATE TABLE wishlists (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  product_id UUID REFERENCES products(id) ON DELETE CASCADE,
  added_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_wishlists_user_id ON wishlists(user_id);
CREATE INDEX idx_wishlists_product_id ON wishlists(product_id);

-- Constraint: One wishlist entry per user per product
CREATE UNIQUE INDEX idx_wishlists_unique ON wishlists(user_id, product_id);
\`\`\`

## Modifications

### Products Table
\`\`\`sql
ALTER TABLE products 
ADD COLUMN review_count INT DEFAULT 0,
ADD COLUMN average_rating DECIMAL(3, 2) DEFAULT 0.0;

CREATE INDEX idx_products_average_rating ON products(average_rating);
\`\`\`

## Migration Script
\`\`\`sql
BEGIN;

-- Create new tables
CREATE TABLE reviews (...);
CREATE TABLE wishlists (...);

-- Modify products table
ALTER TABLE products ADD COLUMN review_count INT DEFAULT 0;
ALTER TABLE products ADD COLUMN average_rating DECIMAL(3, 2) DEFAULT 0.0;

-- Create indexes
CREATE INDEX idx_reviews_product_id ON reviews(product_id);
-- ... all other indexes

COMMIT;
\`\`\`

## Rollback Script
\`\`\`sql
BEGIN;

DROP TABLE IF EXISTS wishlists CASCADE;
DROP TABLE IF EXISTS reviews CASCADE;

ALTER TABLE products DROP COLUMN review_count;
ALTER TABLE products DROP COLUMN average_rating;

COMMIT;
\`\`\`
```

### Step 5: Generate Schema Documentation

```bash
# Get current schema docs
neev bridge --focus "Database Schema" > DATABASE_SCHEMA.md

# Get migration docs
neev bridge --focus "Migration" > MIGRATION_GUIDE.md
```

### Step 6: After Migration, Archive Old Version

```bash
# Archive the migration blueprint
neev lay "Schema Migration v2"

# Update main schema blueprint with new state
# Edit .neev/blueprints/database-schema/architecture.md
# Update version to 2.0.0 and include new tables
```

### Result

- ✅ Complete schema documentation
- ✅ Tracked schema evolution
- ✅ Migration and rollback scripts
- ✅ Version history
- ✅ AI-ready context for schema work

---

## Next Steps

Now that you've completed these tutorials, you can:

1. **Apply to your project** — Start using Neev in your real projects
2. **Explore advanced features** — Check [PRODUCTION_ENHANCEMENTS.md](PRODUCTION_ENHANCEMENTS.md)
3. **Share with your team** — Help team members get started
4. **Contribute** — Share your patterns and improvements

## Additional Resources

- [Getting Started Guide](GETTING_STARTED.md)
- [API Reference](API_REFERENCE.md)
- [Best Practices](BEST_PRACTICES.md)
- [FAQ](FAQ.md)
- [GitHub Discussions](https://github.com/neev-kit/neev/discussions)

---

**Questions or feedback?** [Open an issue](https://github.com/neev-kit/neev/issues) or [start a discussion](https://github.com/neev-kit/neev/discussions)!
