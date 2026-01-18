# Neev Phase 4: Production Hardening & Release Guide

This guide covers the production hardening updates to the Neev CLI and instructions for building and releasing binaries.

## Overview of Phase 4 Changes

### 1. Structured Logging (`core/logger`)

We've integrated `log/slog` (Go's standard structured logging package) for professional, human-readable logs with JSON support for CI/CD pipelines.

**Features:**
- **Human-Readable Logs** (default): Colored output with emojis for quick visual scanning
  - 🔍 Debug messages
  - ℹ️  Info messages
  - ⚠️  Warnings
  - ❌ Errors

- **JSON Logs** (for CI/CD): Set `NEEV_LOG=json` environment variable

**Usage in Code:**
```go
import "github.com/neev-kit/neev/core/logger"

// In main/Execute()
logger.Init()

// Logging calls
logger.Info("Starting deployment")
logger.Warn("Config not found, using defaults")
logger.Error("Failed to create blueprint", "err", err)
```

**Environment Variables:**
- `NEEV_LOG=json` - Use JSON output (useful for log aggregation services)
- Default (no env var) - Human-readable colored output

### 2. Custom Error Handling (`core/errors`)

We've created a custom error system with solution hints for better user experience.

**Error Types:**
- `ErrTypeBlueprintNotFound` - Blueprint doesn't exist
- `ErrTypeFoundation` - Foundation/`.neev` directory missing
- `ErrTypeInvalidConfig` - Configuration file is invalid
- `ErrTypeIO` - File read/write errors
- `ErrTypeValidation` - Invalid input data
- `ErrTypeUnknown` - Unexpected errors

**Usage:**
```go
import "github.com/neev-kit/neev/core/errors"

// Creating errors
err := errors.ErrBlueprintNotFound("my-blueprint")

// In CLI error handling
if neevErr, ok := err.(*errors.NeevError); ok {
    fmt.Printf("Error: %v\n", neevErr)
    fmt.Printf("💡 %s\n", neevErr.GetSolutionHint())
}

// Solution hints are provided automatically for users:
// "Error: blueprint 'my-blueprint' not found
//  💡 Make sure the blueprint exists in the .neev/blueprints/ directory. Run `neev draft` to create one."
```

### 3. Configuration Management (`core/config`)

Formalized `neev.yaml` configuration with a loader that handles defaults gracefully.

**Configuration Structure:**
```yaml
project_name: "My App"
foundation_path: ".neev"
ignore_dirs:
  - node_modules
  - dist
  - build
  - vendor
```

**Usage:**
```go
import "github.com/neev-kit/neev/core/config"

// Load config (auto-defaults if missing)
cfg, err := config.LoadConfig(".")

// Use with inspection
warnings, err := foundation.InspectWithConfig(cwd, cfg)

// Save modified config
err := config.SaveConfig(".", cfg)
```

### 4. GoReleaser Configuration (`.goreleaser.yaml`)

Pre-configured for multi-platform builds and releases to GitHub.

**Supported Platforms:**
- `darwin/amd64` - Intel Macs
- `darwin/arm64` - Apple Silicon Macs
- `linux/amd64` - Linux x86_64
- `windows/amd64` - Windows

**Features:**
- Automatic tarball/zip creation per platform
- SHA256 checksums
- GitHub release automation
- Homebrew tap integration (commented, ready to enable)

## Building and Releasing

### Prerequisites

1. **Install GoReleaser:**
   ```bash
   # On macOS with Homebrew
   brew install goreleaser

   # Or download from https://goreleaser.com/install
   ```

2. **Verify Go version:**
   ```bash
   go version  # Should be 1.21+ for best compatibility
   ```

### Local Snapshot Build (No Release)

Test the build process locally without creating a GitHub release:

```bash
# Create a snapshot build (won't be tagged)
goreleaser release --snapshot --clean

# This will:
# - Create binaries in ./dist/
# - Skip GitHub release
# - Useful for testing build configuration
# - Output example:
#   dist/neev_linux_amd64_v1/neev
#   dist/neev_darwin_amd64_v1/neev
#   dist/neev_darwin_arm64_v1/neev
#   dist/neev_windows_amd64_v1/neev.exe
```

**Verify the snapshot build:**
```bash
# Test the Intel Mac binary (example)
./dist/neev_darwin_amd64_v1/neev --help

# Check all architectures built
ls -la dist/
```

### Full Release Process

When ready to release to GitHub:

1. **Create and push a git tag:**
   ```bash
   # Tag should follow semantic versioning
   git tag -a v1.0.0 -m "Release v1.0.0: Production hardening and logging"
   git push origin v1.0.0
   ```

2. **Build and release:**
   ```bash
   # Set GitHub token (required for release)
   export GITHUB_TOKEN=<your-github-token>

   # Create the release
   goreleaser release --clean
   ```

3. **What happens:**
   - Builds binaries for all platforms
   - Creates tarballs/zips with checksums
   - Pushes to GitHub Releases
   - Generates CHANGELOG from commits

4. **Result:** Users can download from https://github.com/neev-kit/neev/releases

### Enabling Homebrew Distribution (Future)

When ready to distribute via Homebrew:

1. Create a `homebrew-tap` repository (if not exists)
2. Uncomment the `brews` section in `.goreleaser.yaml`
3. Set environment variable:
   ```bash
   export HOMEBREW_TAP_TOKEN=<your-github-token>
   goreleaser release --clean
   ```
4. Users can then install:
   ```bash
   brew tap neev-kit/neev
   brew install neev
   ```

## Running Tests

### Run All Tests

```bash
# From project root
go test ./...

# With verbose output
go test -v ./...

# With coverage
go test -cover ./...
```

### Run Specific Test Suite

```bash
# Test inspect functionality
go test -v ./core/foundation -run TestInspect

# Test config loading
go test -v ./core/config

# Test logger
go test -v ./core/logger
```

### Test Coverage Report

```bash
# Generate coverage report
go test -v -coverprofile=coverage.out ./...

# View in HTML
go tool cover -html=coverage.out
# Opens in your browser
```

## Example Workflow

### Scenario: Release v1.0.0 with Production Hardening

```bash
# 1. Ensure all changes are committed
git status

# 2. Test locally first
go test ./...

# 3. Create snapshot build to verify
goreleaser release --snapshot --clean

# 4. Test the snapshot binary
./dist/neev_darwin_arm64_v1/neev init --help

# 5. Tag the release
git tag -a v1.0.0 -m "Release v1.0.0: Production hardening & logging"

# 6. Push tag
git push origin v1.0.0

# 7. Build and release
export GITHUB_TOKEN=your-token-here
goreleaser release --clean

# 8. Verify release on GitHub
# https://github.com/neev-kit/neev/releases/tag/v1.0.0
```

## Troubleshooting

### Build fails with "go mod tidy" error
```bash
# Ensure all dependencies are updated
go mod tidy
go mod download
```

### Snapshot build shows "no GoReleaser config found"
```bash
# Ensure .goreleaser.yaml is in the project root
# Check file exists and is valid YAML
cat .goreleaser.yaml
```

### Release fails with GitHub token error
```bash
# Verify token has repo write permissions
export GITHUB_TOKEN=your-token
goreleaser check
```

### Binary signature/entitlements on macOS
```bash
# For future macOS notarization, add to .goreleaser.yaml:
# macos:
#   - enabled: true
#     signing:
#       - identity_file: /path/to/cert
```

## Environment Variables Reference

| Variable | Default | Purpose |
|----------|---------|---------|
| `NEEV_LOG` | (empty) | Set to `json` for JSON output; empty for colored output |
| `GITHUB_TOKEN` | (none) | Required for GoReleaser to push releases to GitHub |
| `HOMEBREW_TAP_TOKEN` | (none) | Token for Homebrew tap updates (future use) |

## Production Checklist

Before releasing to production:

- [ ] All tests pass: `go test ./...`
- [ ] Snapshot build succeeds: `goreleaser release --snapshot --clean`
- [ ] Binaries can execute: `./dist/neev_*/neev --help`
- [ ] Error handling shows solution hints
- [ ] Logger outputs correctly (test with `NEEV_LOG=json`)
- [ ] Config loads defaults when missing
- [ ] Changelog is up-to-date
- [ ] Version number is bumped in code/docs
- [ ] Git tag is created and pushed
- [ ] GitHub token is available
- [ ] Release notes are prepared

## Additional Resources

- GoReleaser Docs: https://goreleaser.com
- Go log/slog Package: https://pkg.go.dev/log/slog
- Semantic Versioning: https://semver.org
- GitHub Releases API: https://docs.github.com/en/rest/releases

## Support

For issues with Phase 4 implementation:
- Check the error solution hints: `neev <command> 2>&1`
- Enable debug logging: `NEEV_LOG=debug neev <command>`
- Review `.goreleaser.yaml` configuration
- Test snapshot build locally first
