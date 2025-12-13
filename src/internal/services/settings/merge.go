package settings

import (
	"encoding/json"
	"fmt"
)

// MergeSettings performs a smart merge of template and existing settings.
// It adds missing hooks from the template while preserving all existing hooks
// and user permissions. The merge is idempotent and identifies hooks by the
// (hook_type, command_path) tuple for deduplication.
//
// Parameters:
//   - template: The embedded template settings JSON
//   - existing: The current settings.local.json content (may be nil/empty)
//
// Returns:
//   - Merged settings JSON bytes
//   - Error if JSON parsing/marshaling fails
//
// Behavior:
//   - If existing is nil/empty, returns template directly
//   - Preserves all user permissions unchanged
//   - Adds missing hooks from template to each hook type
//   - Preserves all existing hooks (user customizations)
//   - Deduplicates by command path within each hook type
func MergeSettings(template, existing []byte) ([]byte, error) {
	// Parse template settings first (always required for validation)
	var templateSettings Settings
	if err := json.Unmarshal(template, &templateSettings); err != nil {
		return nil, fmt.Errorf("failed to parse template settings: %w", err)
	}

	// If no existing settings, return template directly
	if len(existing) == 0 {
		return template, nil
	}

	// Parse existing settings
	var existingSettings Settings
	if err := json.Unmarshal(existing, &existingSettings); err != nil {
		return nil, fmt.Errorf("failed to parse existing settings: %w", err)
	}

	// Start with existing settings as base (preserves permissions and all existing hooks)
	result := existingSettings

	// Initialize hooks map if nil
	if result.Hooks == nil {
		result.Hooks = make(map[string][]HookEntry)
	}

	// For each hook type in template, add missing hooks
	for hookType, templateEntries := range templateSettings.Hooks {
		// Get existing hook entries for this type (may be empty)
		existingEntries := result.Hooks[hookType]

		// Build set of existing command paths for this hook type
		existingCommands := make(map[string]bool)
		for _, entry := range existingEntries {
			for _, hook := range entry.Hooks {
				existingCommands[hook.Command] = true
			}
		}

		// Add missing hooks from template
		for _, templateEntry := range templateEntries {
			for _, templateHook := range templateEntry.Hooks {
				// If this command path doesn't exist, add it
				if !existingCommands[templateHook.Command] {
					// Add to first entry if it exists, otherwise create new entry
					if len(existingEntries) > 0 {
						// Check if first entry has same matcher as template entry
						// If so, add to that entry; otherwise create new entry
						if existingEntries[0].Matcher == templateEntry.Matcher {
							existingEntries[0].Hooks = append(existingEntries[0].Hooks, templateHook)
						} else {
							// Create new entry with template's matcher
							newEntry := HookEntry{
								Matcher: templateEntry.Matcher,
								Hooks:   []Hook{templateHook},
							}
							existingEntries = append(existingEntries, newEntry)
						}
					} else {
						// No existing entries, create new one
						newEntry := HookEntry{
							Matcher: templateEntry.Matcher,
							Hooks:   []Hook{templateHook},
						}
						existingEntries = append(existingEntries, newEntry)
					}

					// Mark as added to prevent duplicates in this merge
					existingCommands[templateHook.Command] = true
				}
			}
		}

		// Update result with merged entries
		result.Hooks[hookType] = existingEntries
	}

	// Marshal result to JSON
	mergedJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal merged settings: %w", err)
	}

	return mergedJSON, nil
}
