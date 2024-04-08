package config_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/exec"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/config"
	"github.com/coffeebeats/gdbuild/pkg/config/common"
	"github.com/coffeebeats/gdbuild/pkg/config/linux"
	"github.com/coffeebeats/gdbuild/pkg/config/macos"
	"github.com/coffeebeats/gdbuild/pkg/config/windows"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

func TestBuildTemplate(t *testing.T) {
	tests := []struct {
		name string

		rc    run.Context
		files map[string]string
		index uint // The root manifest (defaults to '0').

		assert func(t *testing.T, rc *run.Context, tmp string, got *template.Template, err error)
	}{
		{
			name: "empty 'config.extends' returns an error",

			files: map[string]string{
				"gdbuild.toml": `config.extends = ""`,
			},

			assert: func(t *testing.T, rc *run.Context, tmp string, got *template.Template, err error) {
				// Then: There's an error denoting the failure.
				assert.ErrorIs(t, err, config.ErrInvalidInput)

				// Then: The template is empty.
				assert.Equal(t, (*template.Template)(nil), got)
			},
		},
		{
			name: "empty template is correctly converted into default for linux",

			rc: run.Context{
				PathWorkspace: "$TEST_TMPDIR/build",
				PathManifest:  "$TEST_TMPDIR/gdbuild.toml",
				PathOut:       "$TEST_TMPDIR/dist",
				Platform:      platform.OSLinux,
				Profile:       engine.ProfileDebug,
			},
			files: map[string]string{
				"gdbuild.toml": `godot.version = "4.0.0"`,
			},

			assert: func(t *testing.T, rc *run.Context, tmp string, got *template.Template, err error) {
				// Then: There's no error.
				assert.Nil(t, err)

				// Then: The template matches expectations.
				assert.Equal(
					t,
					&template.Template{
						Arch: platform.ArchAmd64,
						Builds: []template.Build{
							{
								Arch:     platform.ArchAmd64,
								Source:   engine.Source{Version: mustParseVersion(t, "4.0.0")},
								Platform: platform.OSLinux,
								Profile:  engine.ProfileDebug,
							},
						},
						Paths:     nil,
						Prebuild:  nil,
						Postbuild: nil,
					},
					got,
				)
			},
		},
		{
			name: "empty template is correctly converted into default for macos",

			rc: run.Context{
				PathManifest:  "$TEST_TMPDIR/gdbuild.toml",
				PathOut:       "$TEST_TMPDIR/dist",
				PathWorkspace: "$TEST_TMPDIR/build",
				Platform:      platform.OSMacOS,
				Profile:       engine.ProfileDebug,
			},
			files: map[string]string{
				"vulkan/": "", // Create an empty directory.

				"gdbuild.toml": `
					godot.version = "4.0.0"

					[template.platform.macos]
					double_precision = true
					vulkan = { sdk_path = "$TEST_TMPDIR/vulkan" }`,
			},

			assert: func(t *testing.T, rc *run.Context, tmp string, got *template.Template, err error) {
				// Then: There's no error.
				assert.Nil(t, err)

				// FIXME: Remove once app bundle action can be asserted on.
				seq, ok := got.Postbuild.(action.Sequence)
				require.True(t, ok)

				seq.Post = nil
				got.Postbuild = seq

				// Then: The template matches expectations.
				assert.Equal(
					t,
					&template.Template{
						Arch: platform.ArchUniversal,
						Builds: []template.Build{
							{
								Arch:            platform.ArchAmd64,
								DoublePrecision: true,
								Source:          engine.Source{Version: mustParseVersion(t, "4.0.0")},
								Platform:        platform.OSMacOS,
								Profile:         engine.ProfileDebug,
								SCons: template.SCons{
									ExtraArgs: []string{
										"use_volk=no",
										"vulkan_sdk_path=" + filepath.Join(tmp, "vulkan"),
									},
								},
							},
							{
								Arch:            platform.ArchArm64,
								DoublePrecision: true,
								Source:          engine.Source{Version: mustParseVersion(t, "4.0.0")},
								Platform:        platform.OSMacOS,
								Profile:         engine.ProfileDebug,
								SCons: template.SCons{
									ExtraArgs: []string{
										"use_volk=no",
										"vulkan_sdk_path=" + filepath.Join(tmp, "vulkan"),
									},
								},
							},
						},
						ExtraArtifacts: []string{
							"godot.macos.template_debug.double.universal",
							"macos.zip",
						},
						NameOverride: "macos.zip",
						Paths:        []osutil.Path{osutil.Path(filepath.Join(tmp, "vulkan"))},
						Prebuild:     nil,
						Postbuild: action.Sequence{
							Action: &action.Process{
								Directory: filepath.Join(tmp, "build/bin"),
								Shell:     exec.DefaultShell(),
								Args: []string{
									"lipo",
									"-create",
									"godot.macos.template_debug.double.x86_64",
									"godot.macos.template_debug.double.arm64",
									"-output",
									"godot.macos.template_debug.double.universal",
								},
							},

							// TODO: Figure out how to test that this works.
							// Post: macos.NewAppBundleAction(rc, []string{"godot.macos.template_debug.double.universal"}),
						},
					},
					got,
				)
			},
		},
		{
			name: "empty template is correctly converted into default for windows",

			rc: run.Context{
				PathWorkspace: "$TEST_TMPDIR/build",
				PathManifest:  "$TEST_TMPDIR/gdbuild.toml",
				PathOut:       "$TEST_TMPDIR/dist",
				Platform:      platform.OSWindows,
				Profile:       engine.ProfileDebug,
			},
			files: map[string]string{
				"gdbuild.toml": `godot.version = "4.0.0"`,
			},

			assert: func(t *testing.T, rc *run.Context, tmp string, got *template.Template, err error) {
				// Then: There's no error.
				assert.Nil(t, err)

				// Then: The template matches expectations.
				assert.Equal(
					t,
					&template.Template{
						Arch: platform.ArchAmd64,
						Builds: []template.Build{
							{
								Arch:     platform.ArchAmd64,
								Source:   engine.Source{Version: mustParseVersion(t, "4.0.0")},
								Platform: platform.OSWindows,
								Profile:  engine.ProfileDebug,
							},
						},
						ExtraArtifacts: []string{"godot.windows.template_debug.x86_64.console.exe"},
					},
					got,
				)
			},
		},
		{
			name: "inherited template is correctly populated",

			rc: run.Context{
				PathManifest:  "$TEST_TMPDIR/gdbuild.toml",
				PathOut:       "$TEST_TMPDIR/dist",
				PathWorkspace: "$TEST_TMPDIR/build",
				Platform:      platform.OSWindows,
				Profile:       engine.ProfileDebug,
			},
			files: map[string]string{
				"parent.toml": `
					godot.version = "4.0.0"

					template.platform.windows.use_mingw = true`,

				"icon.ico": "<image data>", // Create the image file.

				"gdbuild.toml": `
					config.extends = "parent.toml"

					godot.version = "4.2.1"

					[template.platform.windows.profile.debug]
					icon_path = "$TEST_TMPDIR/icon.ico"
					use_mingw = false`,
			},

			assert: func(t *testing.T, rc *run.Context, tmp string, got *template.Template, err error) {
				// Then: There's no error.
				assert.Nil(t, err)

				// Given: The expected icon path.
				image := osutil.Path(filepath.Join(tmp, "icon.ico"))

				// Then: The template matches expectations.

				// NOTE: Function actions can't be checked, so separately test them.
				assert.NotNil(t, got.Prebuild)
				assert.IsType(t, windows.NewCopyImageFileAction(image, rc), got.Prebuild)
				got.Prebuild = nil

				assert.Equal(
					t,
					&template.Template{
						Arch: platform.ArchAmd64,
						Builds: []template.Build{
							{
								Arch:     platform.ArchAmd64,
								Source:   engine.Source{Version: mustParseVersion(t, "4.2.1")},
								Platform: platform.OSWindows,
								Profile:  engine.ProfileDebug,
								SCons:    template.SCons{},
							},
						},
						ExtraArtifacts: []string{"godot.windows.template_debug.x86_64.console.exe"},
						Paths:          []osutil.Path{image},
					},
					got,
				)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Given: A temporary test directory.
			tmp := t.TempDir()

			// Given: The process is updated with the temporary directory variable.
			t.Setenv("TEST_TMPDIR", tmp)

			// Given: The specified configuration files exist.
			for path, contents := range tc.files {
				writeFile(t, tmp, path, contents)
			}

			// Given: The root manifest is parsed.
			doc := tc.files[filepath.Base(tc.rc.PathManifest.String())]
			m, err := config.Parse([]byte(doc))
			require.NoError(t, err)

			// When: The 'Template' is built.
			got, err := config.Template(&tc.rc, m)

			// Then: Results match expectations.
			require.NotNil(t, tc.assert)
			tc.assert(t, &tc.rc, tmp, got, err)
		})
	}
}

func TestTemplateCombine(t *testing.T) {
	tests := []struct {
		name string

		doc string
		rc  run.Context

		want config.Templater
		err  error
	}{
		{
			name: "invalid platform returns an error",

			rc: run.Context{Platform: platform.OSUnknown},

			err: config.ErrInvalidInput,
		},
		{
			name: "an empty document returns an empty template",

			rc: run.Context{Platform: platform.OSWindows},

			want: &windows.Template{Template: new(common.Template)},
		},
		{
			name: "base properties are correctly populated",

			rc: run.Context{Platform: platform.OSWindows},
			doc: `
			[template]
			arch = "arm64"
			env = { VAR = "123" }
			optimize = "speed_trace"
			custom_py_path = "a/b/custom.py"
			`,

			want: &windows.Template{
				Template: &common.Template{
					Arch:         platform.ArchArm64,
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
				Platform: platform.OSWindows,
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
					Arch:     platform.ArchArm64,
					Env:      map[string]string{"VAR": "123"},
					Optimize: engine.OptimizeSize,
				},
			},
		},
		{
			name: "base properties in platform constraints are correctly populated",

			rc: run.Context{
				Features: []string{"test"},
				Platform: platform.OSWindows,
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
					Arch:     platform.ArchArm64,
					Env:      map[string]string{"VAR": "123"},
					Optimize: engine.OptimizeSize,
				},
			},
		},
		{
			name: "windows-specific properties are correctly populated",

			rc:  run.Context{Platform: platform.OSWindows},
			doc: "[template.platform.windows]\nuse_mingw = true",

			want: &windows.Template{
				UseMinGW: pointer(true),
				Template: new(common.Template),
			},
		},
		{
			name: "linux-specific properties with constraints are correctly populated",

			rc: run.Context{
				Platform: platform.OSLinux,
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
				Platform: platform.OSMacOS,
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
				Platform: platform.OSWindows,
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

func mustParseVersion(t *testing.T, text string) engine.Version {
	var out engine.Version

	err := out.UnmarshalText([]byte(text))
	require.NoError(t, err)

	return out
}

func writeFile(t *testing.T, tmp, path string, doc string) {
	t.Helper()

	isDirectory := strings.HasSuffix(path, "/")
	path = filepath.Join(tmp, path)

	err := os.MkdirAll(filepath.Dir(path), osutil.ModeUserRWX)
	require.NoError(t, err)

	if isDirectory {
		err := os.Mkdir(path, osutil.ModeUserRWX)
		require.NoError(t, err)

		return
	}

	f, err := os.Create(path)
	require.NoError(t, err)

	defer f.Close()

	_, err = io.Copy(f, strings.NewReader(doc))
	require.NoError(t, err)
}
