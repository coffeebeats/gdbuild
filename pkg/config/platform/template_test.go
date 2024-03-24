package platform_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/config"
	"github.com/coffeebeats/gdbuild/pkg/config/platform"
	"github.com/coffeebeats/gdbuild/pkg/config/platform/common"
	"github.com/coffeebeats/gdbuild/pkg/config/platform/linux"
	"github.com/coffeebeats/gdbuild/pkg/config/platform/macos"
	"github.com/coffeebeats/gdbuild/pkg/config/platform/windows"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	host "github.com/coffeebeats/gdbuild/pkg/godot/platform"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

func TestTemplateCombine(t *testing.T) {
	tests := []struct {
		name string

		doc string
		rc  run.Context

		want platform.Templater
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

			want: &windows.Template{Template: new(common.Template)},
		},
		{
			name: "base properties are correctly populated",

			rc: run.Context{Platform: host.OSWindows},
			doc: `
			[template]
			arch = "arm64"
			env = { VAR = "123" }
			optimize = "speed_trace"
			custom_py_path = "a/b/custom.py"
			`,

			want: &windows.Template{
				Template: &common.Template{
					Arch:         host.ArchArm64,
					Env:          map[string]string{"VAR": "123"},
					Optimize:     engine.OptimizeSpeedTrace,
					PathCustomPy: osutil.Path("a/b/custom.py"),
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
			[template.profile.release_debug]
			arch = "arm64"

			[template.feature.test]
			env = { VAR = "123" }

			[template.feature.test.profile.release_debug]
			optimize = "size"
			`,

			want: &windows.Template{
				Template: &common.Template{
					Arch:     host.ArchArm64,
					Env:      map[string]string{"VAR": "123"},
					Optimize: engine.OptimizeSize,
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
			[template.platform.windows.profile.release_debug]
			arch = "arm64"

			[template.platform.windows.feature.test]
			env = { VAR = "123" }

			[template.platform.windows.feature.test.profile.release_debug]
			optimize = "size"
			`,

			want: &windows.Template{
				Template: &common.Template{
					Arch:     host.ArchArm64,
					Env:      map[string]string{"VAR": "123"},
					Optimize: engine.OptimizeSize,
				},
			},
		},
		{
			name: "windows-specific properties are correctly populated",

			rc:  run.Context{Platform: host.OSWindows},
			doc: "[template.platform.windows]\nuse_mingw = true",

			want: &windows.Template{
				UseMinGW: pointer(true),
				Template: new(common.Template),
			},
		},
		{
			name: "linux-specific properties with constraints are correctly populated",

			rc: run.Context{
				Platform: host.OSLinux,
				Profile:  engine.ProfileRelease,
			},
			doc: `[template.platform.linux.profile.release]
			use_llvm = true`,

			want: &linux.Template{
				UseLLVM:  pointer(true),
				Template: new(common.Template),
			},
		},
		{
			name: "macos-specific properties with constraints are correctly populated",

			rc: run.Context{
				Features: []string{"test"},
				Platform: host.OSMacOS,
				Profile:  engine.ProfileRelease,
			},
			doc: `
			[template.platform.macos.feature.test]
			lipo_command = ["a"]

			[template.platform.macos.profile.release]
			lipo_command = ["b"]

			[template.platform.macos.feature.test.profile.release]
			lipo_command = ["c"]
			vulkan = { use_volk = true }`,

			want: &macos.Template{
				Template:    &common.Template{},
				LipoCommand: []string{"a", "b", "c"},
				Vulkan:      macos.Vulkan{Dynamic: pointer(true)},
			},
		},
		{
			name: "windows-specific properties with constraints are correctly populated",

			rc: run.Context{
				Features: []string{"test"},
				Platform: host.OSWindows,
				Profile:  engine.ProfileRelease,
			},
			doc: `[template.platform.windows.profile.release]
			use_mingw = true

			[template.platform.windows.feature.test]
			icon_path = "a/b/icon.ico"`,

			want: &windows.Template{
				UseMinGW: pointer(true),
				PathIcon: osutil.Path("a/b/icon.ico"),
				Template: new(common.Template),
			},
		},
	}

	for _, tc := range tests {
		// Given: A 'Manifest' is parsed from the document.
		m, err := config.Parse([]byte(tc.doc))
		require.NoError(t, err)

		// When: The 'Template' type is built from 'Templates'.
		got, err := m.Template.Combine(&tc.rc)

		// Then: The error matches expectations.
		assert.ErrorIs(t, err, tc.err)

		// Then: The returned 'Template' matches expectations.
		assert.Equal(t, tc.want, got)
	}
}

func pointer[T any](value T) *T {
	return &value
}
