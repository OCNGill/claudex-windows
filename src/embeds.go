package claudex

import "embed"

//go:embed profiles
var Profiles embed.FS

//go:embed .claude/settings.local.json
var SettingsTemplate []byte
