package paths

const (
	// ClaudexDir is the root directory for all Claudex artifacts
	ClaudexDir = ".claudex"

	// SessionsDir is the directory for session data
	SessionsDir = ".claudex/sessions"

	// LogsDir is the directory for log files
	LogsDir = ".claudex/logs"

	// ConfigFile is the configuration file path
	ConfigFile = ".claudex/config.toml"

	// PreferencesFile is the user preferences file path
	PreferencesFile = ".claudex/preferences.json"

	// Legacy paths (for migration detection)
	LegacySessionsDir = "sessions"
	LegacyLogsDir     = "logs"
	LegacyConfigFile  = ".claudex.toml"
)
