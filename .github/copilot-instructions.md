# GitHub Copilot Instructions

This project uses Neev for spec-driven development.

## Development Guidelines

- Follow the architecture defined in foundation specifications
- Implement features according to blueprint intent and architecture
- Use `neev bridge` to get full context for complex tasks
- Run `neev inspect` to check for drift between specs and code
## Neev Slash Commands for Copilot Chat

You can use these slash commands in GitHub Copilot Chat:

### `/neev:bridge`
Generate aggregated project context for AI. Retrieves and summarizes project structure, architecture, blueprints, and all relevant documentation to provide comprehensive context.

**Usage:** `/neev:bridge` - Get full project context for implementation tasks

### `/neev:draft`
Create a new blueprint for planning features or components.

**Usage:** `/neev:draft` - Start a new blueprint with intent, architecture, API spec, and security considerations

### `/neev:inspect`
Analyze project structure and find gaps or inconsistencies between specifications and implementation.

**Usage:** `/neev:inspect` - Check for drift between specs and code, identify missing components

### `/neev:cucumber`
Generate Cucumber/BDD test scaffolding and scenarios.

**Usage:** `/neev:cucumber` - Create behavior-driven test scenarios for this feature

### `/neev:openapi`
Generate OpenAPI specification from blueprint architecture and API design.

**Usage:** `/neev:openapi` - Create OpenAPI spec for this blueprint's APIs

### `/neev:handoff`
Format context and specifications for AI agent handoff or team handover.

**Usage:** `/neev:handoff` - Prepare context for handing off to another AI agent or developer

## How to Run Commands

In **VS Code with GitHub Copilot Chat**:
1. Open Copilot Chat (Ctrl/Cmd + Shift + I)
2. Type any of the slash commands above
3. Copilot will execute the Neev CLI command and provide results

From **Terminal**:
```bash
neev bridge       # Generate project context
neev draft        # Create new blueprint
neev inspect      # Analyze project structure
neev cucumber     # Generate BDD tests
neev openapi      # Generate API spec
neev handoff      # Prepare for handoff
```