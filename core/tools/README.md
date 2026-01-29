package tools

/*
Tools Module - AI Tool Detection and Skills Generation

This package provides automatic detection of installed AI tools (Claude, Cursor,
GitHub Copilot, etc.) and generates skills in their native formats.

MAIN COMPONENTS

1. DETECTION (detect.go)
   - DetectInstalledTools(): Find all installed AI tools
   - Tool type constants: ToolClaude, ToolCursor, ToolCopilot, etc.
   - Tool struct: Stores detection results with configuration
   - Individual detection functions: detectCursor(), detectClaude(), etc.

   Usage:
   ```go
   tools := tools.DetectInstalledTools()
   for _, tool := range tools {
       if tool.Installed {
           fmt.Printf("Found: %s at %s\n", tool.Name, tool.Config.SkillsDir)
       }
   }
   ```

2. ADAPTERS (adapters.go)
   - Adapter interface: Defines format conversion contract
   - CursorAdapter: Generates JSON skills
   - ClaudeAdapter: Generates Markdown skills
   - CopilotAdapter: Generates Markdown skills
   - CodeiumAdapter: Generates JSON skills
   - FallbackAdapter: Generates Markdown with documentation

   Usage:
   ```go
   tool := tools.FindTool(detectedTools, tools.ToolCursor)
   adapter := tools.GetAdapter(tool)
   skillContent, _ := adapter.GenerateSkill(skillData)
   ```

3. GENERATOR (generator.go)
   - SkillsGenerator: Orchestrates skill generation
   - GenerateSkills(): Create skills for all tools
   - GenerateSummaryReport(): Create human-readable report
   - Handles fallback when no tools detected

   Usage:
   ```go
   gen := tools.NewSkillsGenerator(projectName, projectRoot, tools)
   err := gen.GenerateSkills(blueprints)
   report := gen.GenerateSummaryReport(blueprints)
   ```

SUPPORTED TOOLS

Native Adapters (Full Support):
- Cursor IDE (JSON format in ~/.cursor/skills/)
- Claude App/VS Code (Markdown in ~/.claude/skills/)
- GitHub Copilot (Markdown in ~/.copilot/skills/)
- Codeium (JSON in ~/.codeium/skills/)

Fallback Support:
- Supabase (JSON with fallback)
- Perplexity AI (Markdown with fallback)
- Any other AI tool (Natural language markdown)

DIRECTORY STRUCTURE

Generated skills are placed in tool-native directories:

~/.cursor/skills/
  ├── feature-1.json
  ├── feature-2.json
  ├── README.md
  └── INDEX.md

~/.claude/skills/
  ├── feature-1.md
  ├── feature-2.md
  ├── README.md
  └── INDEX.md

~/.copilot/skills/
  ├── feature-1.md
  ├── feature-2.md
  ├── README.md
  └── INDEX.md

.neev/skills/ (fallback)
  ├── feature-1.md
  ├── feature-2.md
  ├── README.md
  └── INDEX.md

WORKFLOW

1. Detect Tools
   tools := tools.DetectInstalledTools()

2. Create Generator
   gen := tools.NewSkillsGenerator(name, root, tools)

3. Generate Skills
   gen.GenerateSkills(blueprints)

4. Report Results
   fmt.Println(gen.GenerateSummaryReport(blueprints))

KEY TYPES

ToolType
  - String enum for tool identification
  - Values: "claude", "cursor", "copilot", "codeium", "supabase", "perplexity"

Tool
  - Type: ToolType
  - Name: Friendly name
  - Path: Installation location
  - Installed: bool indicating if found
  - Config: ToolConfig with directories

ToolConfig
  - SkillsDir: Where skills are stored
  - ConfigDir: Tool's config directory
  - Native: Is native adapter available
  - CommandName: CLI command name

SkillContent
  - Name: Skill identifier
  - Description: Human-readable description
  - Content: Skill/blueprint content
  - Type: "command", "snippet", "function", etc.
  - Language: "go", "python", "markdown", etc.
  - Version: Version string

Adapter Interface
  - Name() string
  - GenerateSkill(skill SkillContent) (string, error)
  - GenerateConfigFile(projectName string, skills []SkillContent) (string, error)
  - GetMetadata() map[string]interface{}

EXAMPLES

Example 1: Detect Tools
```go
package main

import (
	"fmt"
	"github.com/neev-kit/neev/core/tools"
)

func main() {
	detected := tools.DetectInstalledTools()
	tools.PrintDetectionSummary(detected)
	
	// Output:
	// Detected AI Tools:
	//   ✓ Cursor (~/.cursor/skills)
	//   ✓ Claude App (~/.claude/skills)
	//   ✓ GitHub Copilot (~/.copilot/skills)
}
```

Example 2: Generate Skills
```go
func main() {
	// Detect tools
	detected := tools.DetectInstalledTools()
	
	// Load blueprints
	blueprints := []tools.SkillContent{
		{
			Name: "auth-feature",
			Description: "OAuth2 authentication",
			Content: "Implementation...",
			Type: "command",
			Language: "go",
			Version: "1.0",
		},
	}
	
	// Generate
	gen := tools.NewSkillsGenerator("myproject", "/path/to/project", detected)
	if err := gen.GenerateSkills(blueprints); err != nil {
		panic(err)
	}
	
	// Report
	fmt.Println(gen.GenerateSummaryReport(blueprints))
}
```

Example 3: Custom Tool Detection
```go
func main() {
	detected := tools.DetectInstalledTools()
	
	// Find specific tool
	cursor := tools.FindTool(detected, tools.ToolCursor)
	if cursor != nil && cursor.Installed {
		fmt.Printf("Cursor found at: %s\n", cursor.Config.SkillsDir)
		
		// Get adapter for that tool
		adapter := tools.GetAdapter(cursor)
		fmt.Printf("Using adapter: %s\n", adapter.Name())
	}
}
```

NATURAL LANGUAGE FALLBACK

When tools aren't detected or don't have native adapters:

1. Skills are generated in .neev/skills/
2. FallbackAdapter creates markdown with:
   - Skill description
   - Implementation code
   - Integration instructions
   - Manual copy-paste friendly

3. User can:
   - Read the markdown skill
   - Copy content to their AI tool
   - Follow integration guide
   - No tool installation required

ERROR HANDLING

Detection
- Returns empty slice if no tools found
- Graceful handling of missing directories
- No errors thrown, just empty results

Generation
- File I/O errors returned immediately
- Per-skill errors logged, generation continues
- Summary shows what succeeded/failed

Tool Not Detected
- Falls back to .neev/skills/
- Generates markdown documentation
- Clear instructions for manual integration

PERFORMANCE

Detection: O(n) where n = number of check paths (~20)
Generation: O(m × t) where m = blueprints, t = tools (~5)
File I/O: Sequential writes with standard buffering

Typical times:
- Tool detection: <10ms
- Skill generation (10 blueprints × 3 tools): <100ms
- Summary report: <10ms

TESTING

See detect_test.go, adapters_test.go, generator_test.go for:
- Unit tests for each detector
- Adapter format validation
- Generator workflow tests
- Integration tests

Run tests:
```bash
go test ./core/tools/... -v
```

ADDING NEW TOOLS

To add support for a new AI tool:

1. Add constant in detect.go:
   const ToolMyTool ToolType = "mytool"

2. Create detection function:
   func detectMyTool() *Tool { ... }

3. Register in DetectInstalledTools()

4. Create adapter in adapters.go:
   type MyToolAdapter struct { ... }

5. Add to GetAdapter() switch statement

6. Write tests

See SKILLS_IMPLEMENTATION.md for detailed guide.

RELATED COMMANDS

neev init
  - Detects tools during initialization
  - Suggests running neev sync-skills

neev sync-skills
  - Main command for skill generation
  - Uses this package's detection and generation

neev detect-tools
  - Lists all detected tools
  - Uses DetectInstalledTools()

neev skills-status
  - Checks generated skill status
  - Counts files in skill directories

DOCUMENTATION

- SKILLS.md: User guide
- SKILLS_IMPLEMENTATION.md: Developer guide
- This file (tools_README.txt): Package documentation

ISSUES & SUPPORT

For issues:
- Check SKILLS.md troubleshooting section
- Open GitHub issue with tool name and OS
- Include output of: neev detect-tools

*/
