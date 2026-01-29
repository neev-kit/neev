# Skills System Implementation Guide

This document describes the technical implementation of Neev's AI tool integration and skill generation system.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     Skills System                           │
└─────────────────────────────────────────────────────────────┘

1. DETECTION LAYER
   ├─ Tool Detection (detect.go)
   │  ├─ detectCursor()
   │  ├─ detectClaude()
   │  ├─ detectCopilot()
   │  ├─ detectCodeium()
   │  ├─ detectSupabase()
   │  └─ detectPerplexity()
   │
   └─ Stores detected tools in Tool struct with:
      ├─ Type (ToolType enum)
      ├─ Name (friendly name)
      ├─ Path (installation location)
      └─ Config (ToolConfig with directories)

2. ADAPTER LAYER
   ├─ Adapter Interface
   │  ├─ Name()
   │  ├─ GenerateSkill()
   │  ├─ GenerateConfigFile()
   │  └─ GetMetadata()
   │
   ├─ CursorAdapter (JSON format)
   ├─ ClaudeAdapter (Markdown format)
   ├─ CopilotAdapter (Markdown format)
   ├─ CodeiumAdapter (JSON format)
   └─ FallbackAdapter (Markdown with docs)

3. GENERATOR LAYER
   ├─ SkillsGenerator
   │  ├─ GenerateSkills()
   │  ├─ generateFallbackDocumentation()
   │  └─ GenerateSummaryReport()
   │
   └─ WriteSkillToFile()
      └─ Uses adapter-specific formatting

4. CLI LAYER
   ├─ sync-skills command
   ├─ detect-tools command
   └─ skills-status command
```

## Core Components

### 1. Tool Detection (`core/tools/detect.go`)

**Purpose**: Identify installed AI tools on the system

**Key Types**:
```go
type ToolType string // Enum: "claude", "cursor", "copilot", etc.

type Tool struct {
    Type      ToolType
    Name      string
    Path      string
    Installed bool
    Config    ToolConfig
}

type ToolConfig struct {
    SkillsDir   string // Native skills directory
    ConfigDir   string // Tool config directory
    Native      bool   // Is native format
    CommandName string // CLI command name
}
```

**Detection Strategy**:
- Checks environment variables
- Scans standard OS-specific directories
- Verifies installation with `os.Stat()`
- Returns all detected tools with full metadata

**Example Detection Paths**:
```
macOS:
  - ~/Applications/Cursor.app/Contents/MacOS/Cursor
  - ~/.claude/ (Claude app)
  - ~/.vscode/extensions/ (VS Code extensions)

Linux:
  - ~/.cursor/
  - ~/.claude/
  - ~/.vscode/extensions/

Windows:
  - ~/AppData/Local/Cursor/Cursor.exe
  - ~/AppData/Local/Claude/
  - ~/AppData/Roaming/.vscode/extensions/
```

### 2. Adapters (`core/tools/adapters.go`)

**Purpose**: Transform skills into tool-specific formats

**Adapter Interface**:
```go
type Adapter interface {
    Name() string
    GenerateSkill(skill SkillContent) (string, error)
    GenerateConfigFile(projectName string, skills []SkillContent) (string, error)
    GetMetadata() map[string]interface{}
}
```

**Built-in Adapters**:

| Adapter | Format | Extension | Native |
|---------|--------|-----------|--------|
| CursorAdapter | JSON | .json | Yes |
| ClaudeAdapter | Markdown | .md | Yes |
| CopilotAdapter | Markdown | .md | Yes |
| CodeiumAdapter | JSON | .json | Yes |
| FallbackAdapter | Markdown | .md | No |

**SkillContent Structure**:
```go
type SkillContent struct {
    Name        string // Unique skill identifier
    Description string // Human-readable description
    Content     string // Actual skill/blueprint content
    Type        string // "command", "snippet", "function", etc.
    Language    string // "go", "python", "markdown", etc.
    Version     string // "1.0", "2.0", etc.
}
```

### 3. Generator (`core/tools/generator.go`)

**Purpose**: Orchestrate skill generation across all detected tools

**SkillsGenerator Type**:
```go
type SkillsGenerator struct {
    projectName string    // Project identifier
    projectRoot string    // Root directory
    tools       []Tool    // Detected tools
    adapters    []Adapter // Tool-specific adapters
}
```

**Generation Workflow**:

1. **Detection Phase**
   - Call `DetectInstalledTools()`
   - Filter for installed tools only

2. **Loading Phase**
   - Load blueprints from `.neev/blueprints/`
   - Convert to `SkillContent` structs

3. **Generation Phase**
   - For each tool:
     - Create native skills directory
     - For each skill:
       - Use adapter to format
       - Write to native directory
     - Generate tool-specific README

4. **Fallback Phase**
   - If no tools: generate to `.neev/skills/`
   - Use FallbackAdapter for markdown docs

5. **Reporting Phase**
   - Generate summary report
   - Display generation statistics

### 4. CLI Commands (`cli/cmd/sync_skills.go`)

**Command: `neev sync-skills`**
- Detects tools
- Loads blueprints
- Generates skills
- Creates README/INDEX files
- Prints summary report

**Command: `neev detect-tools`**
- Lists all detected tools
- Shows installation paths
- Displays configuration locations
- Indicates native support

**Command: `neev skills-status`**
- Shows generation status
- Counts available skills
- Identifies missing skills
- Suggests next steps

## Data Flow

### Skill Generation Flow

```
User runs 'neev sync-skills'
         ↓
   Detect Tools
    ↙        ↘
Found    No Tools
Tools       ↓
 ↓       Create
Get      .neev/skills
Adapters  Using Fallback
 ↓        Adapter
For Each  ↓
Tool      Write
 ↓        Markdown
Create    Files
Native    ↓
Dir       Done
 ↓
For Each
Skill
 ↓
Adapter.
Generate
Skill()
 ↓
Write to
Native
Dir
 ↓
Generate
README
 ↓
Done
```

### Adapter Selection

```
Tool Type → GetAdapter()
                ↓
            ┌───┴───────────────────────────────┐
            ↓                                    ↓
      Native Adapter              FallbackAdapter
      ├─ CursorAdapter               │
      ├─ ClaudeAdapter              Uses
      ├─ CopilotAdapter            Markdown
      └─ CodeiumAdapter            format
                                   + docs
```

## Integration Points

### 1. During Project Init

```go
// cli/cmd/init.go
foundation.Initialize(cwd)
// ↓
detectedTools := tools.DetectInstalledTools()
// ↓
Print tool detection summary
```

### 2. During Skills Sync

```go
// cli/cmd/sync_skills.go
tools.DetectInstalledTools()
↓
loadBlueprintsAsSkills()
↓
SkillsGenerator.GenerateSkills()
↓
Print summary report
```

### 3. File Structure Created

```
Project Root/
├── .neev/
│   ├── blueprints/           (source)
│   └── skills/               (fallback)
│       ├── skill-1.md
│       ├── skill-2.md
│       └── README.md
│
├── ~/.claude/                (Claude home)
│   └── skills/               (native)
│       ├── skill-1.md
│       ├── skill-2.md
│       └── README.md
│
├── ~/.cursor/                (Cursor home)
│   └── skills/               (native)
│       ├── skill-1.json
│       ├── skill-2.json
│       └── README.md
│
└── ~/.copilot/               (Copilot home)
    └── skills/               (native)
        ├── skill-1.md
        ├── skill-2.md
        └── README.md
```

## Adding New Tool Support

To add support for a new AI tool:

### Step 1: Create Detection Function

```go
// core/tools/detect.go

const ToolMyTool ToolType = "mytool"

func detectMyTool() *Tool {
    // Check for installation
    if installed {
        return &Tool{
            Type:      ToolMyTool,
            Name:      "My Tool",
            Path:      path,
            Installed: true,
            Config: ToolConfig{
                SkillsDir:   skillsDir,
                ConfigDir:   configDir,
                Native:      true,
                CommandName: "mytool",
            },
        }
    }
    return nil
}

// Register in DetectInstalledTools()
if myTool := detectMyTool(); myTool != nil {
    tools = append(tools, *myTool)
}
```

### Step 2: Create Adapter

```go
// core/tools/adapters.go

type MyToolAdapter struct {
    tool *Tool
}

func NewMyToolAdapter(tool *Tool) *MyToolAdapter {
    return &MyToolAdapter{tool: tool}
}

func (a *MyToolAdapter) Name() string {
    return "My Tool"
}

func (a *MyToolAdapter) GenerateSkill(skill SkillContent) (string, error) {
    // Format skill in tool's native format
    // Return formatted string or error
}

func (a *MyToolAdapter) GenerateConfigFile(projectName string, skills []SkillContent) (string, error) {
    // Generate tool-specific config/readme
}

func (a *MyToolAdapter) GetMetadata() map[string]interface{} {
    return map[string]interface{}{
        "adapter":    "MyTool",
        "native":     true,
        "formatType": "my-format",
        // ... other metadata
    }
}
```

### Step 3: Register Adapter

```go
// core/tools/adapters.go - GetAdapter function

func GetAdapter(tool *Tool) Adapter {
    switch tool.Type {
    case ToolMyTool:
        return NewMyToolAdapter(tool)
    // ... other cases
    default:
        return NewFallbackAdapter(tool)
    }
}
```

## Testing

### Unit Tests Structure

```
core/tools/
├── detect_test.go          # Detection tests
├── adapters_test.go        # Adapter tests
└── generator_test.go       # Generator tests

cli/cmd/
└── sync_skills_test.go     # CLI integration tests
```

### Test Examples

**Detection Test**:
```go
func TestDetectCursor(t *testing.T) {
    // Mock file system
    // Call detectCursor()
    // Verify Tool struct
}
```

**Adapter Test**:
```go
func TestCursorGenerateSkill(t *testing.T) {
    adapter := NewCursorAdapter(&tool)
    skill := SkillContent{...}
    
    result, err := adapter.GenerateSkill(skill)
    
    // Verify JSON format
    // Verify required fields
}
```

**Generator Test**:
```go
func TestGenerateSkills(t *testing.T) {
    gen := NewSkillsGenerator("test", "/tmp/test", tools)
    err := gen.GenerateSkills(blueprints)
    
    // Verify files created
    // Verify content format
    // Verify directory structure
}
```

## Performance Considerations

1. **Detection**: O(n) where n = number of standard paths checked
2. **Generation**: O(m × t) where m = blueprints, t = tools
3. **File I/O**: Sequential writes with buffering

**Optimizations**:
- Cache tool detection results
- Parallel adapter generation (future)
- Incremental skill updates (future)

## Error Handling

**Error Types**:
```go
- Tool detection failures: logged, continue
- File I/O errors: return formatted error
- Adapter errors: logged per skill
- Missing directories: auto-create with MkdirAll
```

**User-Friendly Messages**:
- Clear error messages with solutions
- Suggestions for next steps
- Links to documentation

## Future Enhancements

1. **Dynamic Tool Loading**
   - Plugin system for custom adapters
   - Third-party tool support

2. **Incremental Updates**
   - Track file hashes
   - Only regenerate changed skills

3. **Skill Templating**
   - Custom skill templates
   - Template variables substitution

4. **Conflict Resolution**
   - Detect manual modifications
   - Merge strategy for updates

5. **Skill Validation**
   - Schema validation per tool
   - Syntax checking before write

## Related Files

- [SKILLS.md](../SKILLS.md) - User guide
- [core/tools/detect.go](../core/tools/detect.go) - Detection implementation
- [core/tools/adapters.go](../core/tools/adapters.go) - Adapter implementations
- [core/tools/generator.go](../core/tools/generator.go) - Generator implementation
- [cli/cmd/sync_skills.go](../cli/cmd/sync_skills.go) - CLI commands
- [cli/cmd/init.go](../cli/cmd/init.go) - Initialization integration
