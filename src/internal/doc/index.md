# doc

Documentation generation and update services for Claudex sessions.

## Core Files

- `interface.go` - DocumentationUpdater interface definition
- `updater.go` - Background Claude invocation for documentation updates
- `transcript.go` - JSONL transcript parsing and formatting
- `prompts.go` - Prompt template loading and building

## Subdirectories

- `rangeupdater/` - Range-based documentation updates using Git commit ranges
  - `claude.go` - Background Claude invocation for index.md regeneration
  - `updater.go` - Core range-based documentation update logic
  - `resolver.go` - Commit range resolution and analysis
  - `types.go` - Type definitions for range updates
  - `skiprules.go` - Rules for skipping documentation updates
  - `fallback.go` - Fallback strategies for update failures

## Tests

- `transcript_test.go` - Tests for transcript parsing
- `prompts_test.go` - Tests for prompt template handling
- `updater_test.go` - Tests for the documentation updater
