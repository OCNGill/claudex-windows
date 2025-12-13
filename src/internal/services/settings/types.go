package settings

// Settings represents the structure of settings.local.json configuration file.
// It contains permissions for tool usage and hooks for various lifecycle events.
type Settings struct {
	Permissions Permissions            `json:"permissions"`
	Hooks       map[string][]HookEntry `json:"hooks"`
}

// Permissions defines access control for tool invocations.
// It uses allow/deny/ask lists to control which tools can be executed.
type Permissions struct {
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
	Ask   []string `json:"ask"`
}

// HookEntry represents a hook configuration for a specific event type.
// Each entry can optionally have a matcher regex and contains a list of hooks to execute.
type HookEntry struct {
	Matcher string `json:"matcher,omitempty"`
	Hooks   []Hook `json:"hooks"`
}

// Hook defines a single hook command to execute.
// Currently supports "command" type hooks with shell script execution.
type Hook struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}
