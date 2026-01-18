# Neev Phase 4: Quick Reference Card

## 🚀 For Developers: Quick Integration Guide

### Initialize Logger
```go
import "github.com/neev-kit/neev/core/logger"

func main() {
    logger.Init()  // Call once at startup
}
```

### Use Logging
```go
logger.Info("Operation started", "module", "auth")
logger.Warn("Cache miss detected")
logger.Error("Failed to write file", "err", err)
```

### Create Custom Errors
```go
import "github.com/neev-kit/neev/core/errors"

// Quick constructors
err := errors.ErrBlueprintNotFound("my-blueprint")
err := errors.ErrFoundationMissing()
err := errors.ErrInvalidConfig("missing project_name")

// Generic custom error
err := errors.NewNeevError(
    errors.ErrTypeValidation,
    "invalid input data",
    nil,
)
```

### Handle Errors in CLI
```go
if neevErr, ok := err.(*errors.NeevError); ok {
    fmt.Printf("Error: %v\n", neevErr)
    fmt.Printf("💡 %s\n", neevErr.GetSolutionHint())
} else {
    fmt.Printf("Error: %v\n", err)
}
```

### Load Configuration
```go
import "github.com/neev-kit/neev/core/config"

cfg, err := config.LoadConfig(".")
if err != nil {
    return err
}

// Access config
logger.Info("Project", "name", cfg.ProjectName)
logger.Info("Foundation", "path", cfg.FoundationPath)
```

### Use Config with Inspect
```go
import "github.com/neev-kit/neev/core/foundation"

cfg, _ := config.LoadConfig(".")
warnings, err := foundation.InspectWithConfig(".", cfg)
if err != nil {
    return errors.ErrFoundationMissing()
}

for _, w := range warnings {
    logger.Warn(w)
}
```

---

## 📝 Environment Variables

| Variable | Values | Usage |
|----------|--------|-------|
| `NEEV_LOG` | `json` or empty | `json` for CI/CD, empty for colored output |
| `GITHUB_TOKEN` | Token string | Required for GoReleaser releases |

**Example:**
```bash
NEEV_LOG=json neev init  # JSON logging
neev init                # Colored logging (default)
```

---

## 🎯 Error Types Reference

| Error Type | Constructor | Solution Hint |
|-----------|-------------|--------------|
| `ErrTypeBlueprintNotFound` | `ErrBlueprintNotFound(name)` | Run `neev draft` to create one |
| `ErrTypeFoundation` | `ErrFoundationMissing()` | Run `neev init` first |
| `ErrTypeInvalidConfig` | `ErrInvalidConfig(reason)` | Check neev.yaml format |
| `ErrTypeIO` | `NewNeevError(ErrTypeIO, ...)` | Check file permissions |
| `ErrTypeValidation` | `NewNeevError(ErrTypeValidation, ...)` | Check input parameters |

---

## 📦 Building and Releasing

### Quick Snapshot (Local Testing)
```bash
goreleaser release --snapshot --clean
# Creates: dist/neev_* directories with binaries
```

### Full Release (To GitHub)
```bash
# 1. Tag the release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 2. Set GitHub token
export GITHUB_TOKEN=ghp_xxxxx

# 3. Build and release
goreleaser release --clean
```

### Test Downloaded Binary
```bash
# After snapshot build
ls -la dist/neev_darwin_arm64_v1/
./dist/neev_darwin_arm64_v1/neev --help
```

---

## 🧪 Testing

### Run All Tests
```bash
cd /path/to/neev/cli && go test ./...
cd /path/to/neev/core && go test ./...
```

### Run Specific Tests
```bash
go test -v ./core/foundation -run TestInspect
go test -v ./core/config
```

### Test Coverage
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 🔧 Configuration File Format

**File:** `neev.yaml` (in project root)

```yaml
project_name: "My Amazing App"
foundation_path: ".neev"
ignore_dirs:
  - node_modules
  - dist
  - build
  - vendor
  - .git
  - target
```

**Validation Rules:**
- `project_name` - Required, non-empty string
- `foundation_path` - Required, must be relative path
- `ignore_dirs` - Optional, default list provided

---

## ✨ Example Workflow

### Scenario: Add New Feature with Logging

```go
package myfeature

import (
    "github.com/neev-kit/neev/core/logger"
    "github.com/neev-kit/neev/core/errors"
)

func ProcessBlueprint(name string) error {
    logger.Info("Processing blueprint", "name", name)
    
    if name == "" {
        return errors.ErrInvalidConfig("blueprint name cannot be empty")
    }
    
    logger.Info("Blueprint processed successfully")
    return nil
}
```

**Output (default):**
```
ℹ️  Processing blueprint name=my-blueprint
ℹ️  Blueprint processed successfully
```

**Output (NEEV_LOG=json):**
```json
{"time":"2024-01-18T10:30:45Z","level":"INFO","msg":"Processing blueprint","name":"my-blueprint"}
{"time":"2024-01-18T10:30:45Z","level":"INFO","msg":"Blueprint processed successfully"}
```

---

## 🎓 Key Takeaways

1. **Always initialize logger:** `logger.Init()` in main()
2. **Use custom errors:** Provides user-friendly hints
3. **Config-aware code:** Use `InspectWithConfig()` instead of `Inspect()`
4. **Environment variables:** Set `NEEV_LOG=json` for CI/CD
5. **Test thoroughly:** Use `t.TempDir()` for test isolation
6. **Document errors:** Each error type has a solution hint

---

## 📚 Learn More

- **Full Guide:** See [PHASE4_PRODUCTION_HARDENING.md](PHASE4_PRODUCTION_HARDENING.md)
- **Implementation Details:** See [PHASE4_IMPLEMENTATION_SUMMARY.md](PHASE4_IMPLEMENTATION_SUMMARY.md)
- **Error Package:** [core/errors/errors.go](core/errors/errors.go)
- **Config Package:** [core/config/loader.go](core/config/loader.go)
- **Logger Package:** [core/logger/logger.go](core/logger/logger.go)

---

## 💡 Pro Tips

1. **Chain errors for debugging:**
   ```go
   return errors.NewNeevError(
       errors.ErrTypeIO,
       "failed to save blueprint",
       originalErr,  // Wrapped for debugging
   )
   ```

2. **Use structured logging for debugging:**
   ```bash
   NEEV_LOG=json neev init | jq '.msg' # Extract messages
   ```

3. **Test with temporary directories:**
   ```go
   tmpDir := t.TempDir()  // Auto-cleanup after test
   ```

4. **Snapshot before real release:**
   ```bash
   goreleaser release --snapshot --clean
   # Test binaries before pushing to GitHub
   ```
