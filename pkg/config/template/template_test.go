package template_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/config"
	"github.com/coffeebeats/gdbuild/pkg/config/template"
	"github.com/coffeebeats/gdbuild/pkg/godot/build"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
)

func TestTemplateBuild(t *testing.T) {
	tests := []struct {
		name string

		doc string
		bc  build.Context

		want template.Template
		err  error
	}{
		{
			name: "invalid platform returns an error",

			bc: build.Context{Platform: platform.OSUnknown},

			err: template.ErrInvalidInput,
		},
		{
			name: "an empty document returns an empty template",

			bc: build.Context{Platform: platform.OSWindows},

			want: &template.Windows{Base: &template.Base{}},
		},
		{
			name: "base properties are correctly populated",

			bc: build.Context{Platform: platform.OSWindows},
			doc: `
			[template]
			arch = "arm64"
			env = { VAR = "123" }
			optimize = "speed_trace"
			custom_py_path = "a/b/custom.py"
			`,

			want: &template.Windows{
				Base: &template.Base{
					Arch:         platform.ArchArm64,
					Env:          map[string]string{"VAR": "123"},
					Optimize:     build.OptimizeSpeedTrace,
					PathCustomPy: osutil.Path("a/b/custom.py"),
				},
			},
		},
		{
			name: "base properties with constraints are correctly populated",

			bc: build.Context{
				Features: []string{"test"},
				Platform: platform.OSWindows,
				Profile:  build.ProfileReleaseDebug,
			},
			doc: `
			[template.profile.release_debug]
			arch = "arm64"

			[template.feature.test]
			env = { VAR = "123" }

			# [template.feature.test.profile.release_debug]
			optimize = "speed_trace"
			`,

			want: &template.Windows{
				Base: &template.Base{
					Arch:     platform.ArchArm64,
					Env:      map[string]string{"VAR": "123"},
					Optimize: build.OptimizeSpeedTrace,
				},
			},
		},
		{
			name: "base properties in platform constraints are correctly populated",

			bc: build.Context{
				Features: []string{"test"},
				Platform: platform.OSWindows,
				Profile:  build.ProfileReleaseDebug,
			},
			doc: `
			[template.platform.windows.profile.release_debug]
			arch = "arm64"

			[template.platform.windows.feature.test]
			env = { VAR = "123" }

			# [template.platform.windows.feature.test.profile.release_debug]
			optimize = "speed_trace"
			`,

			want: &template.Windows{
				Base: &template.Base{
					Arch:     platform.ArchArm64,
					Env:      map[string]string{"VAR": "123"},
					Optimize: build.OptimizeSpeedTrace,
				},
			},
		},
		{
			name: "windows-specific properties are correctly populated",

			bc:  build.Context{Platform: platform.OSWindows},
			doc: "[template.platform.windows]\nuse_mingw = true",

			want: &template.Windows{
				UseMinGW: pointer(true),
				Base:     &template.Base{},
			},
		},
		{
			name: "linux-specific properties with constraints are correctly populated",

			bc: build.Context{
				Platform: platform.OSLinux,
				Profile:  build.ProfileRelease,
			},
			doc: `[template.platform.linux.profile.release]
			use_llvm = true`,

			want: &template.Linux{
				UseLLVM: pointer(true),
				Base:    &template.Base{},
			},
		},
		{
			name: "macos-specific properties with constraints are correctly populated",

			bc: build.Context{
				Features: []string{"test"},
				Platform: platform.OSMacOS,
				Profile:  build.ProfileRelease,
			},
			doc: `
			[template.platform.macos.feature.test]
			lipo_command = ["a"]

			[template.platform.macos.profile.release]
			lipo_command = ["b"]

			[template.platform.macos.feature.test.profile.release]
			lipo_command = ["c"]
			vulkan = { use_volk = true }`,

			want: &template.MacOS{
				Base:        &template.Base{},
				LipoCommand: []string{"a", "b", "c"},
				Vulkan:      template.Vulkan{Dynamic: pointer(true)},
			},
		},
		{
			name: "windows-specific properties with constraints are correctly populated",

			bc: build.Context{
				Features: []string{"test"},
				Platform: platform.OSWindows,
				Profile:  build.ProfileRelease,
			},
			doc: `[template.platform.windows.profile.release]
			use_mingw = true

			[template.platform.windows.feature.test]
			icon_path = "a/b/icon.ico"`,

			want: &template.Windows{
				UseMinGW: pointer(true),
				PathIcon: osutil.Path("a/b/icon.ico"),
				Base:     &template.Base{},
			},
		},
	}

	for _, tc := range tests {
		// Given: A 'Manifest' is parsed from the document.
		m, err := config.Parse([]byte(tc.doc))
		require.NoError(t, err)

		// When: The 'Template' type is built from 'Templates'.
		got, err := m.Template.Build(&tc.bc)

		// Then: The error matches expectations.
		assert.ErrorIs(t, err, tc.err)

		// Then: The returned 'Template' matches expectations.
		assert.Equal(t, tc.want, got)
	}
}

func pointer[T any](value T) *T {
	return &value
}
