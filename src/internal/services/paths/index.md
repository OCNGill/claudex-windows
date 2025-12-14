# Paths Package

## Purpose

The `paths` package provides centralized path constants for all Claudex artifacts. This ensures consistency across the codebase and supports the migration from legacy scattered paths to the new `.claudex/` folder structure.

## Overview

All Claudex artifacts (sessions, logs, config, preferences) are now stored under a single `.claudex/` directory in the working directory. This prevents conflicts with user directories and provides cleaner organization.

## Constants

### New Paths (Active)

- **ClaudexDir**: `.claudex` - Root directory for all Claudex artifacts
- **SessionsDir**: `.claudex/sessions` - Session data storage
- **LogsDir**: `.claudex/logs` - Log files
- **ConfigFile**: `.claudex/config.toml` - Configuration file
- **PreferencesFile**: `.claudex/preferences.json` - User preferences

### Legacy Paths (Migration Support)

- **LegacySessionsDir**: `sessions` - Old session directory location
- **LegacyLogsDir**: `logs` - Old logs directory location
- **LegacyConfigFile**: `.claudex.toml` - Old config file location

## Usage

```go
import "claudex/src/internal/services/paths"

// Use path constants instead of hardcoded strings
sessionPath := filepath.Join(workDir, paths.SessionsDir)
configPath := filepath.Join(workDir, paths.ConfigFile)
```

## Migration Support

The legacy path constants enable automatic detection and migration of existing artifacts:

1. Check if legacy paths exist
2. If found, move to new `.claudex/` locations
3. Continue normal operation with new paths

This migration is handled by the `migrate` usecase during application startup.

## Design Principles

- **Centralized**: Single source of truth for all path definitions
- **Explicit**: Clear distinction between new and legacy paths
- **Type-safe**: Constants prevent typos and string literal duplication
- **Migration-aware**: Supports smooth transition from old to new structure
