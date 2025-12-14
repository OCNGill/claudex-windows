# Migrate Usecase

## Overview
The migrate package handles the initialization and migration of Claudex artifacts to the new `.claudex/` directory structure.

## Purpose
- Create `.claudex/` directory on first run
- Auto-create `config.toml` with sensible defaults
- Migrate legacy artifacts from old locations
- Ensure idempotent operations (safe to run multiple times)

## Architecture

### Components
- **Migrator**: Main migration orchestrator

### Migration Process
1. **Create `.claudex/` directory** if it doesn't exist
2. **Create default `config.toml`** if it doesn't exist
3. **Migrate legacy sessions/** → `.claudex/sessions/` (if exists)
4. **Migrate legacy logs/** → `.claudex/logs/` (if exists)
5. **Migrate legacy `.claudex.toml`** → `.claudex/config.toml` (overwrites default if exists)

## Key Features

### Atomic Operations
- Uses `Rename()` for atomic directory migration when possible
- Falls back to copy + delete for cross-filesystem migrations
- Ensures data integrity during migration

### Idempotency
- Safe to run multiple times
- Skips operations if destination already exists
- No side effects on repeated execution

### Error Handling
- Returns errors only for critical failures (directory creation, default config)
- Logs warnings for non-critical migration failures
- Continues execution even if legacy migrations fail

### Default Configuration
Creates `config.toml` with defaults:
```toml
[features]
autodoc_session_progress = true
autodoc_session_end = true
autodoc_frequency = 5
```

## Usage

```go
import (
    "github.com/spf13/afero"
    "github.com/maikelderhaeg/claudex/src/internal/usecases/migrate"
)

func main() {
    fs := afero.NewOsFs()
    migrator := migrate.New(fs)

    if err := migrator.Run(); err != nil {
        log.Fatalf("Migration failed: %v", err)
    }
}
```

## Migration Scenarios

### Scenario 1: Fresh Installation
- Creates `.claudex/` directory
- Creates `config.toml` with defaults
- No legacy artifacts to migrate

### Scenario 2: Existing Legacy Setup
- Creates `.claudex/` directory
- Creates default `config.toml`
- Migrates `sessions/` → `.claudex/sessions/`
- Migrates `logs/` → `.claudex/logs/`
- Migrates `.claudex.toml` → `.claudex/config.toml` (overwrites default)
- Removes legacy files after successful migration

### Scenario 3: Partial Legacy Setup
- Creates `.claudex/` directory
- Creates default `config.toml`
- Migrates only the legacy artifacts that exist
- Skips missing legacy artifacts (no errors)

### Scenario 4: Already Migrated
- Detects existing `.claudex/` directory
- Detects existing `config.toml`
- Skips all operations (no changes)
- Returns success immediately

## Dependencies
- `github.com/spf13/afero` - Filesystem abstraction
- `github.com/maikelderhaeg/claudex/src/internal/services/paths` - Path constants

## Files
- `migrate.go` - Main migration implementation
- `index.md` - This documentation file

## Future Enhancements
- [ ] Add migration version tracking
- [ ] Add rollback capability for failed migrations
- [ ] Add dry-run mode for testing migrations
- [ ] Add migration progress reporting for large datasets
