package config

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// TestLoad_EmptyConfig_ReturnsDefaults verifies that an empty config file returns default feature values
func TestLoad_EmptyConfig_ReturnsDefaults(t *testing.T) {
	fs := afero.NewMemMapFs()
	configPath := "/test/.claudex/config.toml"

	// Create empty config file
	err := afero.WriteFile(fs, configPath, []byte(""), 0644)
	require.NoError(t, err)

	// Load config
	cfg, err := Load(fs, configPath)
	require.NoError(t, err)

	// Assert defaults
	require.True(t, cfg.Features.AutodocSessionProgress, "AutodocSessionProgress should default to true")
	require.True(t, cfg.Features.AutodocSessionEnd, "AutodocSessionEnd should default to true")
	require.Equal(t, 5, cfg.Features.AutodocFrequency, "AutodocFrequency should default to 5")
}

// TestLoad_NoConfigFile_ReturnsDefaults verifies that missing config file returns defaults
func TestLoad_NoConfigFile_ReturnsDefaults(t *testing.T) {
	fs := afero.NewMemMapFs()
	configPath := "/test/.claudex/config.toml"

	// Don't create config file

	// Load config
	cfg, err := Load(fs, configPath)
	require.NoError(t, err)

	// Assert defaults
	require.True(t, cfg.Features.AutodocSessionProgress)
	require.True(t, cfg.Features.AutodocSessionEnd)
	require.Equal(t, 5, cfg.Features.AutodocFrequency)
	require.Empty(t, cfg.Doc)
	require.False(t, cfg.NoOverwrite)
}

// TestLoad_PartialFeatures_UsesDefaults verifies partial [features] section uses defaults for missing fields
func TestLoad_PartialFeatures_UsesDefaults(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected Features
	}{
		{
			name: "Only autodoc_session_progress",
			content: `[features]
autodoc_session_progress = false`,
			expected: Features{
				AutodocSessionProgress: false,
				AutodocSessionEnd:      true, // default
				AutodocFrequency:       5,    // default
			},
		},
		{
			name: "Only autodoc_session_end",
			content: `[features]
autodoc_session_end = false`,
			expected: Features{
				AutodocSessionProgress: true, // default
				AutodocSessionEnd:      false,
				AutodocFrequency:       5, // default
			},
		},
		{
			name: "Only autodoc_frequency",
			content: `[features]
autodoc_frequency = 10`,
			expected: Features{
				AutodocSessionProgress: true, // default
				AutodocSessionEnd:      true, // default
				AutodocFrequency:       10,
			},
		},
		{
			name: "Two fields set",
			content: `[features]
autodoc_session_progress = false
autodoc_frequency = 15`,
			expected: Features{
				AutodocSessionProgress: false,
				AutodocSessionEnd:      true, // default
				AutodocFrequency:       15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			configPath := "/test/.claudex/config.toml"

			err := afero.WriteFile(fs, configPath, []byte(tt.content), 0644)
			require.NoError(t, err)

			cfg, err := Load(fs, configPath)
			require.NoError(t, err)

			require.Equal(t, tt.expected.AutodocSessionProgress, cfg.Features.AutodocSessionProgress)
			require.Equal(t, tt.expected.AutodocSessionEnd, cfg.Features.AutodocSessionEnd)
			require.Equal(t, tt.expected.AutodocFrequency, cfg.Features.AutodocFrequency)
		})
	}
}

// TestLoad_FullFeatures_ParsesAllValues verifies full [features] section parses all values correctly
func TestLoad_FullFeatures_ParsesAllValues(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected Features
	}{
		{
			name: "All features disabled",
			content: `[features]
autodoc_session_progress = false
autodoc_session_end = false
autodoc_frequency = 1`,
			expected: Features{
				AutodocSessionProgress: false,
				AutodocSessionEnd:      false,
				AutodocFrequency:       1,
			},
		},
		{
			name: "All features enabled with custom frequency",
			content: `[features]
autodoc_session_progress = true
autodoc_session_end = true
autodoc_frequency = 20`,
			expected: Features{
				AutodocSessionProgress: true,
				AutodocSessionEnd:      true,
				AutodocFrequency:       20,
			},
		},
		{
			name: "Mixed settings",
			content: `[features]
autodoc_session_progress = false
autodoc_session_end = true
autodoc_frequency = 10`,
			expected: Features{
				AutodocSessionProgress: false,
				AutodocSessionEnd:      true,
				AutodocFrequency:       10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			configPath := "/test/.claudex/config.toml"

			err := afero.WriteFile(fs, configPath, []byte(tt.content), 0644)
			require.NoError(t, err)

			cfg, err := Load(fs, configPath)
			require.NoError(t, err)

			require.Equal(t, tt.expected, cfg.Features)
		})
	}
}

// TestLoad_InvalidFrequency_StillLoads verifies config loads even with problematic frequency values
func TestLoad_InvalidFrequency_StillLoads(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantValue int
	}{
		{
			name: "Negative frequency",
			content: `[features]
autodoc_frequency = -5`,
			wantValue: -5, // TOML will parse it, validation happens elsewhere
		},
		{
			name: "Zero frequency",
			content: `[features]
autodoc_frequency = 0`,
			wantValue: 0, // TOML will parse it, validation happens elsewhere
		},
		{
			name: "Large frequency",
			content: `[features]
autodoc_frequency = 1000`,
			wantValue: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			configPath := "/test/.claudex/config.toml"

			err := afero.WriteFile(fs, configPath, []byte(tt.content), 0644)
			require.NoError(t, err)

			cfg, err := Load(fs, configPath)
			require.NoError(t, err)

			require.Equal(t, tt.wantValue, cfg.Features.AutodocFrequency)
		})
	}
}

// TestLoad_WithOtherConfig_PreservesFeatures verifies features work alongside other config sections
func TestLoad_WithOtherConfig_PreservesFeatures(t *testing.T) {
	content := `
doc = ["docs/api.md", "docs/guide.md"]
no_overwrite = true

[features]
autodoc_session_progress = false
autodoc_session_end = true
autodoc_frequency = 15
`

	fs := afero.NewMemMapFs()
	configPath := "/test/.claudex/config.toml"

	err := afero.WriteFile(fs, configPath, []byte(content), 0644)
	require.NoError(t, err)

	cfg, err := Load(fs, configPath)
	require.NoError(t, err)

	// Assert non-features config
	require.Equal(t, []string{"docs/api.md", "docs/guide.md"}, cfg.Doc)
	require.True(t, cfg.NoOverwrite)

	// Assert features config
	require.False(t, cfg.Features.AutodocSessionProgress)
	require.True(t, cfg.Features.AutodocSessionEnd)
	require.Equal(t, 15, cfg.Features.AutodocFrequency)
}

// TestLoad_MalformedTOML_ReturnsError verifies malformed TOML returns an error
func TestLoad_MalformedTOML_ReturnsError(t *testing.T) {
	content := `[features
autodoc_session_progress = false` // Missing closing bracket

	fs := afero.NewMemMapFs()
	configPath := "/test/.claudex/config.toml"

	err := afero.WriteFile(fs, configPath, []byte(content), 0644)
	require.NoError(t, err)

	_, err = Load(fs, configPath)
	require.Error(t, err, "malformed TOML should return error")
}

// TestLoad_FromTestdataFile verifies loading from actual testdata config file
func TestLoad_FromTestdataFile(t *testing.T) {
	fs := afero.NewOsFs()
	configPath := "../../../testdata/configs/features.toml"

	cfg, err := Load(fs, configPath)
	require.NoError(t, err)

	// This test verifies the actual testdata file content
	// Expected values based on testdata/configs/features.toml
	require.False(t, cfg.Features.AutodocSessionProgress)
	require.True(t, cfg.Features.AutodocSessionEnd)
	require.Equal(t, 10, cfg.Features.AutodocFrequency)
}
