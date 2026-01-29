package tools

import (
	"os"
	"path/filepath"
	"runtime"
)

// ToolType represents the type of AI tool
type ToolType int

const (
	ToolUnknown ToolType = iota
	ToolCursor
	ToolClaude
	ToolCopilot
	ToolCodeium
	ToolSupabase
	ToolPerplexity
)

// ToolConfig holds configuration for a tool
type ToolConfig struct {
	Native    bool
	SkillsDir string
	ConfigDir string
	Paths     []string
}

// Tool represents an installed AI tool
type Tool struct {
	Type      ToolType
	Name      string
	Installed bool
	Path      string
	Config    ToolConfig
}

// DetectInstalledTools detects all installed AI tools on the system
func DetectInstalledTools() []Tool {
	tools := make([]Tool, 0)

	// Detect each tool
	if tool := detectCursor(); tool.Installed {
		tools = append(tools, tool)
	}
	if tool := detectClaude(); tool.Installed {
		tools = append(tools, tool)
	}
	if tool := detectCopilot(); tool.Installed {
		tools = append(tools, tool)
	}
	if tool := detectCodeium(); tool.Installed {
		tools = append(tools, tool)
	}
	if tool := detectSupabase(); tool.Installed {
		tools = append(tools, tool)
	}
	if tool := detectPerplexity(); tool.Installed {
		tools = append(tools, tool)
	}

	return tools
}

// detectCursor detects Cursor IDE
func detectCursor() Tool {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".cursor")
	skillsDir := filepath.Join(configDir, "skills")

	_, err := os.Stat(configDir)
	installed := err == nil

	return Tool{
		Type:      ToolCursor,
		Name:      "Cursor",
		Installed: installed,
		Path:      configDir,
		Config: ToolConfig{
			Native:    installed,
			SkillsDir: skillsDir,
			ConfigDir: configDir,
			Paths:     []string{configDir},
		},
	}
}

// detectClaude detects Claude
func detectClaude() Tool {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".claude")
	skillsDir := filepath.Join(configDir, "skills")

	_, err := os.Stat(configDir)
	installed := err == nil

	return Tool{
		Type:      ToolClaude,
		Name:      "Claude",
		Installed: installed,
		Path:      configDir,
		Config: ToolConfig{
			Native:    installed,
			SkillsDir: skillsDir,
			ConfigDir: configDir,
			Paths:     []string{configDir},
		},
	}
}

// detectCopilot detects GitHub Copilot
func detectCopilot() Tool {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".copilot")
	skillsDir := filepath.Join(configDir, "skills")

	_, err := os.Stat(configDir)
	installed := err == nil

	return Tool{
		Type:      ToolCopilot,
		Name:      "GitHub Copilot",
		Installed: installed,
		Path:      configDir,
		Config: ToolConfig{
			Native:    installed,
			SkillsDir: skillsDir,
			ConfigDir: configDir,
			Paths:     []string{configDir},
		},
	}
}

// detectCodeium detects Codeium
func detectCodeium() Tool {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".codeium")
	skillsDir := filepath.Join(configDir, "skills")

	_, err := os.Stat(configDir)
	installed := err == nil

	return Tool{
		Type:      ToolCodeium,
		Name:      "Codeium",
		Installed: installed,
		Path:      configDir,
		Config: ToolConfig{
			Native:    installed,
			SkillsDir: skillsDir,
			ConfigDir: configDir,
			Paths:     []string{configDir},
		},
	}
}

// detectSupabase detects Supabase
func detectSupabase() Tool {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".supabase")
	skillsDir := filepath.Join(configDir, "skills")

	_, err := os.Stat(configDir)
	installed := err == nil

	return Tool{
		Type:      ToolSupabase,
		Name:      "Supabase",
		Installed: installed,
		Path:      configDir,
		Config: ToolConfig{
			Native:    installed,
			SkillsDir: skillsDir,
			ConfigDir: configDir,
			Paths:     []string{configDir},
		},
	}
}

// detectPerplexity detects Perplexity
func detectPerplexity() Tool {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".perplexity")
	skillsDir := filepath.Join(configDir, "skills")

	_, err := os.Stat(configDir)
	installed := err == nil

	return Tool{
		Type:      ToolPerplexity,
		Name:      "Perplexity",
		Installed: installed,
		Path:      configDir,
		Config: ToolConfig{
			Native:    installed,
			SkillsDir: skillsDir,
			ConfigDir: configDir,
			Paths:     []string{configDir},
		},
	}
}

// FindTool finds a tool by type
func FindTool(toolType ToolType, tools []Tool) *Tool {
	for i := range tools {
		if tools[i].Type == toolType {
			return &tools[i]
		}
	}
	return nil
}

// GetInstalledToolsNames returns names of all installed tools
func GetInstalledToolsNames(tools []Tool) []string {
	var names []string
	for _, tool := range tools {
		if tool.Installed {
			names = append(names, tool.Name)
		}
	}
	return names
}

// HasAnyTool checks if any tools are installed
func HasAnyTool(tools []Tool) bool {
	for _, tool := range tools {
		if tool.Installed {
			return true
		}
	}
	return false
}

// GetPlatform returns the current platform
func GetPlatform() string {
	return runtime.GOOS
}
