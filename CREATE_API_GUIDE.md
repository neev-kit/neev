# Step-by-Step: Create an API for Users with Neev

This guide walks you through creating a complete API specification and generating it using Neev's blueprint system.

## Part 1: Installation (Local Setup)

### Step 1.1: Clone and Build Locally

```bash
# Navigate to the neev project
cd /Users/surajsrivastav/workspace/neev

# Build the binary
go build -o neev ./cli

# Verify installation
./neev --version
```

### Step 1.2: Install Globally (Optional)

```bash
# Copy to your PATH
sudo cp ./neev /usr/local/bin/neev

# Test from any directory
neev --version
```

---

## Part 2: Create an API Blueprint

### Step 2.1: Initialize a New Project

```bash
# Create a new project directory
mkdir my-api-project
cd my-api-project

# Initialize with Neev
neev init
```

This creates:
```
my-api-project/
├── .neev/
│   ├── AGENTS.md
│   └── blueprints/
└── README.md
```

### Step 2.2: Create the API Blueprint

```bash
# Create a blueprint for your User API
cat > .neev/blueprints/user-api.md << 'EOF'
# User Management API

## Intent
Provide a comprehensive REST API for managing user accounts with authentication, 
profile management, and role-based access control.

## Architecture

### API Endpoints

#### Authentication
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login user
- `POST /auth/logout` - Logout user
- `POST /auth/refresh` - Refresh token

#### Users
- `GET /users` - List all users (admin only)
- `GET /users/{id}` - Get user by ID
- `PUT /users/{id}` - Update user profile
- `DELETE /users/{id}` - Delete user account

#### Roles & Permissions
- `GET /roles` - List available roles
- `POST /users/{id}/roles` - Assign role to user
- `DELETE /users/{id}/roles/{roleId}` - Remove role from user

### Data Models

#### User
```json
{
  "id": "uuid",
  "email": "string (unique)",
  "name": "string",
  "passwordHash": "string (hashed)",
  "status": "active|inactive|suspended",
  "roles": ["string"],
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

#### AuthToken
```json
{
  "accessToken": "string",
  "refreshToken": "string",
  "expiresIn": "number (seconds)"
}
```

### API Spec

#### POST /auth/register
Request:
```json
{
  "email": "user@example.com",
  "name": "John Doe",
  "password": "secure_password"
}
```

Response (201):
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe",
  "status": "active",
  "roles": ["user"],
  "createdAt": "2024-01-28T10:00:00Z"
}
```

#### POST /auth/login
Request:
```json
{
  "email": "user@example.com",
  "password": "secure_password"
}
```

Response (200):
```json
{
  "accessToken": "eyJhbGc...",
  "refreshToken": "refresh_token_here",
  "expiresIn": 3600,
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "roles": ["user"]
  }
}
```

#### GET /users/{id}
Response (200):
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe",
  "status": "active",
  "roles": ["user"],
  "createdAt": "2024-01-28T10:00:00Z",
  "updatedAt": "2024-01-28T10:00:00Z"
}
```

### Security Considerations

1. **Authentication**: JWT tokens with 1-hour expiration
2. **Authorization**: Role-based access control (RBAC)
3. **Password**: bcrypt hashing with salt
4. **Rate Limiting**: 100 requests per minute per IP
5. **CORS**: Whitelist approved origins
6. **HTTPS**: All endpoints require TLS 1.2+
7. **Validation**: Input validation on all endpoints
8. **SQL Injection**: Use parameterized queries

### Error Handling

Standard HTTP status codes:
- 200: Success
- 201: Created
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 409: Conflict (duplicate email)
- 422: Unprocessable Entity (validation)
- 500: Internal Server Error

All errors return:
```json
{
  "error": "error_code",
  "message": "Human readable message",
  "details": {}
}
```

### Implementation Notes

- Language: Go / Node.js / Python (choose one)
- Database: PostgreSQL
- Cache: Redis for session management
- Framework: Echo / Express / FastAPI

---

## Part 3: Generate Skills for Your API

### Step 3.1: Detect Installed Tools

```bash
# See which AI tools are installed
./neev detect-tools

# Output:
# 🔍 Detecting AI Tools
# ═════════════════════════════════════════════════════════════
# 
# Detected Tools:
# 
# ✓ Claude
#   Config: /Users/surajsrivastav/.claude
#   Skills: /Users/surajsrivastav/.claude/skills
# 
# ✓ Cursor
#   Config: /Users/surajsrivastav/.cursor
#   Skills: /Users/surajsrivastav/.cursor/skills
# 
# Found 2 tool(s).
```

### Step 3.2: Generate Skills from Blueprint

```bash
# Generate skills for all detected tools
./neev sync-skills

# Output:
# ✅ Skills generated successfully!
# Skills Generation Report for my-api-project
# ═════════════════════════════════════════════════════════════
# Detected Tools: 2
#   - Claude: Installed
#   - Cursor: Installed
# 
# Blueprints Converted: 1
#   - user-api
# 
# Adapters Used: 2
#   - Claude (markdown)
#   - Cursor (json)
```

### Step 3.3: Check Skills Status

```bash
# See where skills were generated
./neev skills-status

# Output:
# 📊 Skills Status
# ═════════════════════════════════════════════════════════════
# 
# ✓ Claude: 1 skills
#   Directory: /Users/surajsrivastav/.claude/skills
# 
# ✓ Cursor: 1 skills
#   Directory: /Users/surajsrivastav/.cursor/skills
```

---

## Part 4: Use Skills in Your AI Tool

### For Claude Users:

1. Open [Claude](https://claude.ai)
2. Navigate to your project: `/Users/surajsrivastav/.claude/skills`
3. Open `user-api.md`
4. Copy the content into Claude's instructions or custom context
5. Ask Claude to: "Using the User API blueprint, implement the authentication endpoints in [language]"

### For Cursor IDE Users:

1. Open Cursor IDE
2. Navigate to `.cursor/skills/` folder
3. Find `user-api.json`
4. The JSON contains structured skill data
5. Use in Cursor's AI context for better code generation

### For GitHub Copilot Users:

1. Open VS Code with GitHub Copilot
2. Check `.copilot/skills/` folder
3. Open `user-api.md`
4. Reference in comments: `// @skill user-api`
5. Copilot will use this context for suggestions

---

## Part 5: Implement the API

### Step 5.1: Create Project Structure

```bash
# Create directories
mkdir -p src/{handlers,models,middleware,database}
mkdir -p tests
mkdir -p config

# Create main files
touch src/main.go
touch go.mod
touch go.sum
```

### Step 5.2: Ask AI to Implement

**Prompt for Claude/Copilot:**

```
Using the User API blueprint in .neev/blueprints/user-api.md:

1. Implement the User model with validation
2. Create database migrations for PostgreSQL
3. Implement the POST /auth/register endpoint
4. Implement the POST /auth/login endpoint with JWT
5. Add middleware for authentication
6. Add role-based access control

Use Go with the Echo framework.
```

### Step 5.3: Regenerate Skills When Updating Blueprint

```bash
# Edit your blueprint
nano .neev/blueprints/user-api.md

# Regenerate skills
./neev sync-skills

# Now your AI tools have the updated spec!
```

---

## Part 6: Complete Implementation Examples

### Create Additional Blueprints for Your API

#### Product Catalog API

```bash
cat > .neev/blueprints/product-api.md << 'EOF'
# Product Catalog API

## Intent
Provide RESTful API for product management with inventory tracking.

## Endpoints
- `GET /products` - List products
- `GET /products/{id}` - Get product details
- `POST /products` - Create product (admin)
- `PUT /products/{id}` - Update product (admin)
- `DELETE /products/{id}` - Delete product (admin)
- `POST /products/{id}/inventory` - Update inventory

## Models
- Product (id, name, description, price, inventory)
- Category (id, name, description)
- Review (id, productId, rating, text)
EOF
```

#### Orders API

```bash
cat > .neev/blueprints/orders-api.md << 'EOF'
# Orders API

## Intent
Handle customer orders with payment processing and status tracking.

## Endpoints
- `POST /orders` - Create order
- `GET /orders/{id}` - Get order details
- `GET /orders` - List user's orders
- `PUT /orders/{id}` - Update order status (admin)
- `POST /orders/{id}/payment` - Process payment

## Models
- Order (id, userId, items, total, status, createdAt)
- OrderItem (productId, quantity, price)
EOF
```

### Regenerate All Skills

```bash
# Now sync all blueprints to your AI tools
./neev sync-skills

# Your AI assistant now has complete API documentation!
```

---

## Part 7: Workflow Summary

```
┌─────────────────────────────────────────────────────────────┐
│ Complete Workflow                                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ 1. Create blueprint (.neev/blueprints/api.md)              │
│    └─> Define API endpoints, models, security              │
│                                                             │
│ 2. Detect tools (neev detect-tools)                        │
│    └─> Finds Claude, Cursor, Copilot, etc.                │
│                                                             │
│ 3. Generate skills (neev sync-skills)                      │
│    └─> Converts blueprint to native tool formats           │
│       ├─> Claude: ~/.claude/skills/*.md                    │
│       ├─> Cursor: ~/.cursor/skills/*.json                  │
│       └─> Copilot: ~/.copilot/skills/*.md                 │
│                                                             │
│ 4. Use in AI tool                                          │
│    └─> Load skill file into your AI assistant              │
│                                                             │
│ 5. Implement with AI assistance                            │
│    └─> Ask AI to code the implementation                   │
│                                                             │
│ 6. Update blueprint (repeat step 1)                        │
│    └─> Evolve your API spec                                │
│                                                             │
│ 7. Regenerate skills (repeat step 3)                       │
│    └─> AI tools always have latest spec                    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Common Commands Reference

```bash
# Install
go build -o neev ./cli

# Initialize project
./neev init

# Detect tools
./neev detect-tools

# Generate/regenerate skills
./neev sync-skills

# Check skills status
./neev skills-status

# View help
./neev --help
./neev sync-skills --help
```

---

## Next Steps

1. ✅ Build neev locally: `go build -o neev ./cli`
2. ✅ Create project: `mkdir my-api && cd my-api`
3. ✅ Initialize: `../neev init`
4. ✅ Create blueprint: Create `.neev/blueprints/user-api.md`
5. ✅ Detect tools: `../neev detect-tools`
6. ✅ Generate skills: `../neev sync-skills`
7. ✅ Use in your AI tool: Load `.claude/skills/user-api.md` or `.cursor/skills/user-api.json`
8. ✅ Implement with AI: Ask Claude/Copilot to build endpoints
9. ✅ Update as needed: Edit blueprint, run `sync-skills` again

Happy building! 🚀
