// Package migrate provides migration functionality for Claudex.
// It handles creating the .claudex/ directory structure, auto-creating
// default configuration files, and migrating legacy artifacts from old locations.
package migrate

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/afero"

	"claudex/internal/services/paths"
)

const defaultConfigContent = `# Claudex Configuration
# See documentation for all available options

[features]
autodoc_session_progress = true
autodoc_session_end = true
autodoc_frequency = 5
`

// Migrator handles migration of legacy Claudex artifacts and initialization
// of the .claudex/ directory structure.
type Migrator struct {
	fs afero.Fs
}

// New creates a new Migrator instance with the provided filesystem.
func New(fs afero.Fs) *Migrator {
	return &Migrator{fs: fs}
}

// Run executes the migration process:
// 1. Creates .claudex/ directory if it doesn't exist
// 2. Creates config.toml with defaults if it doesn't exist
// 3. Migrates legacy sessions/ directory if it exists
// 4. Migrates legacy logs/ directory if it exists
// 5. Migrates legacy .claudex.toml config if it exists (overwrites default)
//
// Returns error only on critical failures. Non-critical issues are logged as warnings.
// This operation is idempotent and safe to run multiple times.
func (m *Migrator) Run() error {
	// Step 1: Create .claudex/ directory
	if err := m.ensureClaudexDir(); err != nil {
		return fmt.Errorf("failed to create .claudex directory: %w", err)
	}

	// Step 2: Create default config.toml if it doesn't exist
	if err := m.ensureDefaultConfig(); err != nil {
		return fmt.Errorf("failed to create default config: %w", err)
	}

	// Step 3-5: Migrate legacy artifacts
	// These are non-critical - we log warnings but don't fail the migration
	m.migrateLegacySessions()
	m.migrateLegacyLogs()
	m.migrateLegacyConfig()

	return nil
}

// ensureClaudexDir creates the .claudex/ directory if it doesn't exist.
func (m *Migrator) ensureClaudexDir() error {
	exists, err := afero.DirExists(m.fs, paths.ClaudexDir)
	if err != nil {
		return err
	}

	if !exists {
		if err := m.fs.MkdirAll(paths.ClaudexDir, 0755); err != nil {
			return err
		}
		log.Printf("Created %s directory", paths.ClaudexDir)
	}

	return nil
}

// ensureDefaultConfig creates config.toml with default values if it doesn't exist.
// If the file already exists, it does nothing (preserves user configuration).
func (m *Migrator) ensureDefaultConfig() error {
	exists, err := afero.Exists(m.fs, paths.ConfigFile)
	if err != nil {
		return err
	}

	if !exists {
		if err := afero.WriteFile(m.fs, paths.ConfigFile, []byte(defaultConfigContent), 0644); err != nil {
			return err
		}
		log.Printf("Created default config at %s", paths.ConfigFile)
	}

	return nil
}

// migrateLegacySessions migrates the legacy sessions/ directory to .claudex/sessions/
func (m *Migrator) migrateLegacySessions() {
	if err := m.migrateDirectory(paths.LegacySessionsDir, paths.SessionsDir); err != nil {
		log.Printf("Warning: Failed to migrate legacy sessions: %v", err)
	}
}

// migrateLegacyLogs migrates the legacy logs/ directory to .claudex/logs/
func (m *Migrator) migrateLegacyLogs() {
	if err := m.migrateDirectory(paths.LegacyLogsDir, paths.LogsDir); err != nil {
		log.Printf("Warning: Failed to migrate legacy logs: %v", err)
	}
}

// migrateLegacyConfig migrates the legacy .claudex.toml to .claudex/config.toml
// This overwrites the default config if a legacy config exists.
func (m *Migrator) migrateLegacyConfig() {
	exists, err := afero.Exists(m.fs, paths.LegacyConfigFile)
	if err != nil {
		log.Printf("Warning: Failed to check for legacy config: %v", err)
		return
	}

	if !exists {
		return
	}

	// Read legacy config
	content, err := afero.ReadFile(m.fs, paths.LegacyConfigFile)
	if err != nil {
		log.Printf("Warning: Failed to read legacy config: %v", err)
		return
	}

	// Write to new location (overwrites default)
	if err := afero.WriteFile(m.fs, paths.ConfigFile, content, 0644); err != nil {
		log.Printf("Warning: Failed to migrate legacy config: %v", err)
		return
	}

	// Remove legacy config file
	if err := m.fs.Remove(paths.LegacyConfigFile); err != nil {
		log.Printf("Warning: Failed to remove legacy config file: %v", err)
		return
	}

	log.Printf("Migrated legacy config from %s to %s", paths.LegacyConfigFile, paths.ConfigFile)
}

// migrateDirectory moves a directory from source to destination atomically.
// If destination already exists, it skips the migration.
// After successful migration, it removes the source directory.
func (m *Migrator) migrateDirectory(source, dest string) error {
	// Check if source exists
	sourceExists, err := afero.DirExists(m.fs, source)
	if err != nil {
		return fmt.Errorf("failed to check source directory: %w", err)
	}
	if !sourceExists {
		return nil // Nothing to migrate
	}

	// Check if destination already exists
	destExists, err := afero.DirExists(m.fs, dest)
	if err != nil {
		return fmt.Errorf("failed to check destination directory: %w", err)
	}
	if destExists {
		log.Printf("Destination %s already exists, skipping migration from %s", dest, source)
		return nil
	}

	// Use atomic rename if possible (same filesystem)
	// This is the most efficient and atomic operation
	if err := m.fs.Rename(source, dest); err != nil {
		// If rename fails, fall back to copy + delete
		// This can happen if source and dest are on different filesystems
		return m.copyAndRemoveDirectory(source, dest)
	}

	log.Printf("Migrated %s to %s", source, dest)
	return nil
}

// copyAndRemoveDirectory copies a directory recursively and then removes the source.
// This is a fallback when atomic rename is not possible.
func (m *Migrator) copyAndRemoveDirectory(source, dest string) error {
	// Ensure destination parent directory exists
	destParent := filepath.Dir(dest)
	if err := m.fs.MkdirAll(destParent, 0755); err != nil {
		return fmt.Errorf("failed to create destination parent directory: %w", err)
	}

	// Walk source directory and copy all files
	err := afero.Walk(m.fs, source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path from source
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		// Calculate destination path
		destPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			// Create directory at destination
			return m.fs.MkdirAll(destPath, info.Mode())
		}

		// Copy file
		content, err := afero.ReadFile(m.fs, path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		if err := afero.WriteFile(m.fs, destPath, content, info.Mode()); err != nil {
			return fmt.Errorf("failed to write %s: %w", destPath, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to copy directory: %w", err)
	}

	// Remove source directory after successful copy
	if err := m.fs.RemoveAll(source); err != nil {
		return fmt.Errorf("failed to remove source directory: %w", err)
	}

	log.Printf("Copied and removed %s to %s", source, dest)
	return nil
}
