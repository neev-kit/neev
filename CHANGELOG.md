# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Windows path compatibility enhancements
- Comprehensive Windows path handling tests
- ACKNOWLEDGMENTS.md with proper attribution to Spec-Kit and OpenSpec
- CODE_OF_CONDUCT.md following Contributor Covenant
- SECURITY.md with vulnerability reporting guidelines
- GitHub issue templates (bug report, feature request, question)
- GitHub pull request template with comprehensive checklist
- OPEN_SOURCE_POLICY.md for compliance tracking

### Changed
- Updated README with Windows installation instructions (PowerShell and winget)
- Enhanced path handling to use `filepath.Join()` for cross-platform compatibility
- Improved COPILOT_SLASH_COMMANDS.md with better attribution

### Fixed
- Hardcoded path separators in `core/bridge/context.go`
- Path compatibility issues for Windows users

## [1.0.0] - 2026-01-19

### Added
- Initial stable release
- Core spec-driven development framework
- `neev init` - Project initialization
- `neev draft <title>` - Blueprint creation
- `neev bridge [--focus]` - Context aggregation for AI
- `neev inspect [--strict]` - Drift detection from specs
- `neev openapi <blueprint>` - OpenAPI spec generation
- `neev cucumber <blueprint>` - BDD test scaffolding
- `neev handoff <role>` - Handoff prompt generation
- `neev lay <blueprint>` - Blueprint archival
- `neev migrate --source <type>` - Migration from Spec-Kit/OpenSpec
- `neev sync-remotes` - Remote foundation synchronization
- `neev slash-commands` - GitHub Copilot slash command management
- GitHub Copilot slash commands integration
- Cursor IDE slash commands support
- VS Code Copilot Chat integration
- Comprehensive CLI documentation
- Example blueprints and use cases

### Features
- Spec-first development methodology
- AI-native workflow (Claude, Copilot, Cursor compatible)
- Offline-first architecture (no external APIs)
- Git-versioned specifications in `.neev/`
- Cross-platform support (macOS, Linux, Windows)
- Minimal dependencies
- Comprehensive test coverage (99%+)

### Documentation
- README.md - Quick start and overview
- GETTING_STARTED.md - Onboarding guide
- COMMAND_CATALOG.md - Full command reference
- API_REFERENCE.md - Developer API documentation
- DEVELOPMENT.md - Development environment setup
- CONTRIBUTING.md - Contribution guidelines
- MAINTAINERS.md - Team information
- FAQ.md - Common questions
- LICENSE - MIT license

### Technical Details
- Written in Go 1.25+
- Build system: GoReleaser (automated cross-platform builds)
- CI/CD: GitHub Actions
- Testing: 99%+ code coverage
- Linting: golangci-lint
- Supported platforms:
  - macOS (Intel: amd64, Apple Silicon: arm64)
  - Linux (x86_64: amd64)
  - Windows (x86_64: amd64)

## [0.1.0] - Pre-release

### Initial Development
- Project scaffolding
- Core foundation module
- Basic CLI framework
- Initial documentation

---

## How to Use This Changelog

- **Added** — New features
- **Changed** — Changes in existing functionality
- **Deprecated** — Soon-to-be removed features
- **Removed** — Now removed features
- **Fixed** — Bug fixes
- **Security** — Security vulnerability fixes

## Release Process

1. Update this CHANGELOG.md with changes for the release
2. Use semantic versioning: MAJOR.MINOR.PATCH
3. Tag release: `git tag v1.0.0`
4. Push tag: `git push origin v1.0.0`
5. GitHub Actions builds and publishes release
6. Create GitHub Release with CHANGELOG excerpt

## Versioning

- **MAJOR** version for incompatible API changes
- **MINOR** version for backwards-compatible functionality additions
- **PATCH** version for bug fixes

## Links

- [Releases](https://github.com/neev-kit/neev/releases)
- [GitHub](https://github.com/neev-kit/neev)
- [Issues](https://github.com/neev-kit/neev/issues)
- [Discussions](https://github.com/neev-kit/neev/discussions)
