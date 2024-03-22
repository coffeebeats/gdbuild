package template

import (
	"fmt"
	"slices"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/exec"
	"github.com/coffeebeats/gdbuild/internal/pathutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/compile"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
	"github.com/coffeebeats/gdbuild/pkg/godot/scons"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: MacOS                               */
/* -------------------------------------------------------------------------- */

type MacOS struct {
	*Base

	// LipoCommand contains arguments used to invoke 'lipo'. Defaults to
	// ["lipo"]. Only used if 'arch' is set to 'platform.ArchUniversal'.
	LipoCommand []string `toml:"lipo_command"`

	// Vulkan defines Vulkan-related configuration.
	Vulkan Vulkan `toml:"vulkan"`
}

// Compile-time check that 'Template' is implemented.
var _ Template = (*MacOS)(nil)

/* -------------------------- Impl: template.Templater ------------------------- */

func (c *MacOS) ToTemplate(g compile.Godot, cc compile.Context) template.Template { //nolint:funlen
	switch a := c.Base.Arch; a {
	case platform.ArchAmd64, platform.ArchArm64:
		t := c.Base.ToTemplate(g, cc)

		t.Binaries[0].Platform = platform.OSMacOS

		scons := &t.Binaries[0].SCons
		if config.Dereference(c.Vulkan.Dynamic) {
			scons.ExtraArgs = append(scons.ExtraArgs, "use_volk=yes")
		} else {
			scons.ExtraArgs = append(scons.ExtraArgs, "use_volk=no")
		}

		if c.Vulkan.PathSDK != "" {
			scons.ExtraArgs = append(scons.ExtraArgs, "vulkan_sdk_path="+c.Vulkan.PathSDK.String())
			t.AddToPaths(c.Vulkan.PathSDK)
		}

		return t
	case platform.ArchUniversal, platform.ArchUnknown:
		// First, create the 'x86_64' binary.
		amd64 := *c
		amd64.Base.Arch = platform.ArchAmd64

		templateAmd64 := amd64.ToTemplate(g, cc)

		// Next, create the 'arm64' binary.
		arm64 := *c
		arm64.Base.Arch = platform.ArchArm64

		templateArm64 := arm64.ToTemplate(g, cc)

		// Finally, merge the two binaries together.

		lipo := c.LipoCommand
		if len(lipo) == 0 {
			lipo = append(lipo, "lipo")
		}

		templateNameUniversal := scons.TemplateName(
			platform.OSMacOS,
			platform.ArchUniversal,
			cc.Profile,
		)

		cmdLipo := &action.Process{
			Directory:   cc.Invoke.BinPath().String(),
			Environment: nil,

			Shell:   exec.DefaultShell(),
			Verbose: cc.Invoke.Verbose,

			Args: append(
				lipo,
				"-create",
				templateAmd64.Binaries[0].Filename(),
				templateArm64.Binaries[0].Filename(),
				"-output",
				templateNameUniversal,
			),
		}

		// Construct the output 'Template'. This is because nothing else needs
		// to be copied over from the arch-specific templates and this avoid the
		// need to deduplicate properties.
		t := c.Base.ToTemplate(g, cc)

		// Register the additional artifact.
		t.ExtraArtifacts = append(t.ExtraArtifacts, templateNameUniversal)

		t.Binaries = []scons.Build{templateAmd64.Binaries[0], templateArm64.Binaries[0]}
		t.Postbuild = cmdLipo.AndThen(t.Postbuild)

		// Construct a list of paths with duplicates removed. This is preferred
		// over duplicating the code used to decide which paths are dependencies.
		paths := make([]pathutil.Path, 0, len(templateAmd64.Paths)+len(templateArm64.Paths))
		paths = append(paths, templateAmd64.Paths...)
		paths = append(paths, templateArm64.Paths...)
		slices.Sort(paths)
		paths = slices.Compact(paths)

		t.Paths = paths

		return t

	default:
		panic(fmt.Errorf("%w: unsupported architecture: %s", ErrInvalidInput, a))
	}
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (c *MacOS) Configure(cc config.Context) error {
	if err := c.Base.Configure(cc); err != nil {
		return err
	}

	if err := c.Vulkan.Configure(cc); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *MacOS) Validate(cc config.Context) error {
	if err := c.Base.Validate(cc); err != nil {
		return err
	}

	if !c.Base.Arch.IsOneOf(
		platform.ArchAmd64,
		platform.ArchArm64,
		platform.ArchUniversal,
		platform.ArchUnknown,
	) {
		return fmt.Errorf("%w: unsupport architecture: %s", config.ErrInvalidInput, c.Base.Arch)
	}

	// NOTE: Don't check for 'lipo', that should be a runtime check.

	if err := c.Vulkan.Validate(cc); err != nil {
		return err
	}

	return nil
}

/* --------------------------- Impl: config.Merger -------------------------- */

func (c *MacOS) MergeInto(other any) error {
	if c == nil || other == nil {
		return nil
	}

	dst, ok := other.(*MacOS)
	if !ok {
		return fmt.Errorf(
			"%w: expected a '%T' but was '%T'",
			config.ErrInvalidInput,
			new(MacOS),
			other,
		)
	}

	return config.Merge(dst, *c)
}

/* -------------------------------------------------------------------------- */
/*                               Struct: Vulkan                               */
/* -------------------------------------------------------------------------- */

// Vulkan defines the settings required by the MacOS template for including
// Vulkan support.
type Vulkan struct {
	// Dynamic enables dynamically linking Vulkan to the template.
	Dynamic *bool `toml:"use_volk"`

	// PathSDK is the path to the Vulkan SDK root.
	PathSDK pathutil.Path `toml:"sdk_path"`
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (c *Vulkan) Configure(cc config.Context) error {
	if err := c.PathSDK.RelTo(cc.PathManifest); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *Vulkan) Validate(_ config.Context) error {
	if err := c.PathSDK.CheckIsDir(); err != nil {
		return fmt.Errorf("%w: missing path to Vulkan SDK", err)
	}

	return nil
}
