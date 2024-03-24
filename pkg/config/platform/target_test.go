package platform_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/pkg/config"
	"github.com/coffeebeats/gdbuild/pkg/config/platform"
	"github.com/coffeebeats/gdbuild/pkg/config/platform/common"
	"github.com/coffeebeats/gdbuild/pkg/config/platform/windows"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	host "github.com/coffeebeats/gdbuild/pkg/godot/platform"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

func TestTargetCombine(t *testing.T) {
	tests := []struct {
		name string

		doc string
		rc  run.Context

		want platform.Exporter
		err  error
	}{
		{
			name: "invalid platform returns an error",

			rc: run.Context{Platform: host.OSUnknown},

			err: config.ErrInvalidInput,
		},
		{
			name: "an empty document returns an empty template",

			rc: run.Context{Platform: host.OSWindows},

			want: &windows.Target{Target: &common.Target{}},
		},
		{
			name: "base properties are correctly populated",

			rc: run.Context{Platform: host.OSWindows},
			doc: `
			[target.target]
			default_features = ["feature1", "feature2"]
			encryption_key = "abcdefg"
			hook = { run_before = ["echo before"] }
			options = {option-name = "option-value"}
			pack_files = [{glob = ["*"], encrypt = true}]
			runnable = true
			server = false
			`,

			want: &windows.Target{
				Target: &common.Target{
					DefaultFeatures: []string{"feature1", "feature2"},
					EncryptionKey:   "abcdefg",
					Hook:            run.Hook{Pre: []action.Command{"echo before"}},
					Options:         map[string]any{"option-name": "option-value"},
					PackFiles: []export.PackFile{
						{
							Glob:    []string{"*"},
							Encrypt: pointer(true),
						},
					},
					Runnable: pointer(true),
					Server:   pointer(false),
				},
			},
		},
		{
			name: "base properties with constraints are correctly populated",

			rc: run.Context{
				Features: []string{"test"},
				Platform: host.OSWindows,
				Profile:  engine.ProfileReleaseDebug,
			},
			doc: `
			[target.target.profile.release_debug]
			default_features = ["feature1"]

			[target.target.feature.test]
			hook = { run_before = ["echo before"] }

			[target.target.feature.test.profile.release_debug]
			runnable = true
			`,

			want: &windows.Target{
				Target: &common.Target{
					DefaultFeatures: []string{"feature1"},
					Hook:            run.Hook{Pre: []action.Command{"echo before"}},
					Runnable:        pointer(true),
				},
			},
		},
		{
			name: "base properties in platform constraints are correctly populated",

			rc: run.Context{
				Features: []string{"test"},
				Platform: host.OSWindows,
				Profile:  engine.ProfileReleaseDebug,
			},
			doc: `
			[target.target.platform.windows.profile.release_debug]
			default_features = ["feature1"]

			[target.target.platform.windows.feature.test]
			hook = { run_before = ["echo before"] }

			[target.target.platform.windows.feature.test.profile.release_debug]
			runnable = true
			`,

			want: &windows.Target{
				Target: &common.Target{
					DefaultFeatures: []string{"feature1"},
					Hook:            run.Hook{Pre: []action.Command{"echo before"}},
					Runnable:        pointer(true),
				},
			},
		},
	}

	for _, tc := range tests {
		// Given: A 'Manifest' is parsed from the document.
		m, err := config.Parse([]byte(tc.doc))
		require.NoError(t, err)

		// When: The 'Target' type is built from 'Targets'.
		got, err := m.Target["target"].Combine(&tc.rc)

		// Then: The error matches expectations.
		assert.ErrorIs(t, err, tc.err)

		// Then: The returned 'Target' matches expectations.
		assert.Equal(t, tc.want, got)
	}
}
