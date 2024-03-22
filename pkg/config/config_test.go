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
	"github.com/coffeebeats/gdbuild/pkg/build"
	"github.com/coffeebeats/gdbuild/pkg/config"
	"github.com/coffeebeats/gdbuild/pkg/config/template"
)

func TestBuildTemplate(t *testing.T) {
	tests := []struct {
		name string

		inv   build.Invocation
		files map[string]string
		index uint // The root manifest (defaults to '0').

		assert func(t *testing.T, inv *build.Invocation, tmp string, got build.Template, err error)
	}{
		{
			name: "empty 'config.extends' returns an error",

			files: map[string]string{
				"gdbuild.toml": `config.extends = ""`,
			},

			assert: func(t *testing.T, inv *build.Invocation, tmp string, got build.Template, err error) {
				// Then: There's an error denoting the failure.
				assert.ErrorIs(t, err, config.ErrInvalidInput)

				// Then: The template is empty.
				assert.Equal(t, build.Template{}, got)
			},
		},
		{
			name: "empty template is correctly converted into default for linux",

			inv: build.Invocation{
				PathBuild:    "$TEST_TMPDIR/build",
				PathManifest: "$TEST_TMPDIR/gdbuild.toml",
				PathOut:      "$TEST_TMPDIR/dist",
				Platform:     platform.OSLinux,
				Profile:      build.ProfileDebug,
			},
			files: map[string]string{
				"gdbuild.toml": `godot.version = "4.0.0"`,
			},

			assert: func(t *testing.T, inv *build.Invocation, tmp string, got build.Template, err error) {
				// Then: There's no error.
				assert.Nil(t, err)

				// Then: The template matches expectations.
				assert.Equal(
					t,
					build.Template{
						Binaries: []build.Binary{
							{
								Arch:     platform.ArchAmd64,
								Godot:    build.Godot{Version: "4.0.0"},
								Platform: platform.OSLinux,
								Profile:  build.ProfileDebug,
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

			inv: build.Invocation{
				PathBuild:    "$TEST_TMPDIR/build",
				PathManifest: "$TEST_TMPDIR/gdbuild.toml",
				PathOut:      "$TEST_TMPDIR/dist",
				Platform:     platform.OSMacOS,
				Profile:      build.ProfileDebug,
			},
			files: map[string]string{
				"vulkan/": "", // Create an empty directory.

				"gdbuild.toml": `
					godot.version = "4.0.0"

					[template.platform.macos]
					vulkan = { sdk_path = "$TEST_TMPDIR/vulkan" }`,
			},

			assert: func(t *testing.T, inv *build.Invocation, tmp string, got build.Template, err error) {
				// Then: There's no error.
				assert.Nil(t, err)

				// Then: The template matches expectations.
				assert.Equal(
					t,
					build.Template{
						Binaries: []build.Binary{
							{
								Arch:     platform.ArchAmd64,
								Godot:    build.Godot{Version: "4.0.0"},
								Platform: platform.OSMacOS,
								Profile:  build.ProfileDebug,
								SCons: build.SCons{
									ExtraArgs: []string{
										"use_volk=no",
										"vulkan_sdk_path=" + filepath.Join(tmp, "vulkan"),
									},
								},
							},
							{
								Arch:     platform.ArchArm64,
								Godot:    build.Godot{Version: "4.0.0"},
								Platform: platform.OSMacOS,
								Profile:  build.ProfileDebug,
								SCons: build.SCons{
									ExtraArgs: []string{
										"use_volk=no",
										"vulkan_sdk_path=" + filepath.Join(tmp, "vulkan"),
									},
								},
							},
						},
						ExtraArtifacts: []string{"godot.macos.template_debug.universal"},
						Paths:          []build.Path{build.Path(filepath.Join(tmp, "vulkan"))},
						Prebuild:       nil,
						Postbuild: &action.Process{
							Directory: filepath.Join(tmp, "build/bin"),
							Shell:     exec.DefaultShell(),
							Args: []string{
								"lipo",
								"-create",
								"godot.macos.template_debug.x86_64",
								"godot.macos.template_debug.arm64",
								"-output",
								"godot.macos.template_debug.universal",
							},
						},
					},
					got,
				)
			},
		},
		{
			name: "empty template is correctly converted into default for windows",

			inv: build.Invocation{
				PathBuild:    "$TEST_TMPDIR/build",
				PathManifest: "$TEST_TMPDIR/gdbuild.toml",
				PathOut:      "$TEST_TMPDIR/dist",
				Platform:     platform.OSWindows,
				Profile:      build.ProfileDebug,
			},
			files: map[string]string{
				"gdbuild.toml": `godot.version = "4.0.0"`,
			},

			assert: func(t *testing.T, inv *build.Invocation, tmp string, got build.Template, err error) {
				// Then: There's no error.
				assert.Nil(t, err)

				// Then: The template matches expectations.
				assert.Equal(
					t,
					build.Template{
						Binaries: []build.Binary{
							{
								Arch:     platform.ArchAmd64,
								Godot:    build.Godot{Version: "4.0.0"},
								Platform: platform.OSWindows,
								Profile:  build.ProfileDebug,
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

			inv: build.Invocation{
				PathBuild:    "$TEST_TMPDIR/build",
				PathManifest: "$TEST_TMPDIR/gdbuild.toml",
				PathOut:      "$TEST_TMPDIR/dist",
				Platform:     platform.OSWindows,
				Profile:      build.ProfileDebug,
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

			assert: func(t *testing.T, inv *build.Invocation, tmp string, got build.Template, err error) {
				// Then: There's no error.
				assert.Nil(t, err)

				// Given: The expected icon path.
				image := build.Path(filepath.Join(tmp, "icon.ico"))

				// Then: The template matches expectations.

				// NOTE: Function actions can't be checked, so separately test them.
				assert.NotNil(t, got.Prebuild)
				assert.IsType(t, template.NewCopyImageFileAction(image, inv), got.Prebuild)
				got.Prebuild = nil

				assert.Equal(
					t,
					build.Template{
						Binaries: []build.Binary{
							{
								Arch:     platform.ArchAmd64,
								Godot:    build.Godot{Version: "4.2.1"},
								Platform: platform.OSWindows,
								Profile:  build.ProfileDebug,
								SCons:    build.SCons{},
							},
						},
						ExtraArtifacts: []string{"godot.windows.template_debug.x86_64.console.exe"},
						Paths:          []build.Path{image},
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
			doc := tc.files[filepath.Base(tc.inv.PathManifest.String())]
			m, err := config.Parse([]byte(doc))
			require.NoError(t, err)

			// When: The 'Template' is built.
			got, err := m.BuildTemplate(tc.inv)

			// Then: Results match expectations.
			require.NotNil(t, tc.assert)
			tc.assert(t, &tc.inv, tmp, got, err)
		})
	}
}

type configuration struct {
	filename string
	manifest string
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
