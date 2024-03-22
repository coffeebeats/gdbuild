package template

import (
	"fmt"
	"slices"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/exec"
	"github.com/coffeebeats/gdbuild/pkg/build"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: MacOS                               */
/* -------------------------------------------------------------------------- */

type MacOS struct {
	*Base

	// LipoCommand contains arguments used to invoke 'lipo'. Defaults to
	// ["lipo"]. Only used if 'arch' is set to 'build.ArchUniversal'.
	LipoCommand []string `toml:"lipo_command"`

	// Vulkan defines Vulkan-related configuration.
	Vulkan Vulkan `toml:"vulkan"`
}

// Compile-time check that 'Template' is implemented.
var _ Template = (*MacOS)(nil)

/* -------------------------- Impl: build.Templater ------------------------- */

func (c *MacOS) ToTemplate(g build.Godot, inv build.Invocation) build.Template { //nolint:funlen
	switch a := c.Base.Arch; a {
	case build.ArchAmd64, build.ArchArm64:
		t := c.Base.ToTemplate(g, inv)

		t.Binaries[0].Platform = build.OSMacOS

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
	case build.ArchUniversal, build.ArchUnknown:
		// First, create the 'x86_64' binary.
		amd64 := *c
		amd64.Base.Arch = build.ArchAmd64

		templateAmd64 := amd64.ToTemplate(g, inv)

		// Next, create the 'arm64' binary.
		arm64 := *c
		arm64.Base.Arch = build.ArchArm64

		templateArm64 := arm64.ToTemplate(g, inv)

		// Finally, merge the two binaries together.

		lipo := c.LipoCommand
		if len(lipo) == 0 {
			lipo = append(lipo, "lipo")
		}

		targetName := inv.Profile.TargetName()

		cmdLipo := &action.Process{
			Directory:   inv.BinPath().String(),
			Environment: nil,

			Shell:   exec.DefaultShell(),
			Verbose: inv.Verbose,

			Args: append(
				lipo,
				"-create",
				fmt.Sprintf("godot.macos.%s.x86_64", targetName),
				fmt.Sprintf("godot.macos.%s.arm64", targetName),
				"-output",
				fmt.Sprintf("godot.macos.%s.universal", targetName),
			),
		}

		// Construct the output 'Template'. This is because nothing else needs
		// to be copied over from the arch-specific templates and this avoid the
		// need to deduplicate properties.
		t := c.Base.ToTemplate(g, inv)

		t.Binaries = []build.Binary{templateAmd64.Binaries[0], templateArm64.Binaries[0]}
		t.Postbuild = cmdLipo.AndThen(t.Postbuild)

		// Construct a list of paths with duplicates removed. This is preferred
		// over duplicating the code used to decide which paths are dependencies.
		paths := make([]build.Path, 0, len(templateAmd64.Paths)+len(templateArm64.Paths))
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

func (c *MacOS) Configure(inv build.Invocation) error {
	if err := c.Base.Configure(inv); err != nil {
		return err
	}

	if err := c.Vulkan.Configure(inv); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *MacOS) Validate(inv build.Invocation) error {
	if err := c.Base.Validate(inv); err != nil {
		return err
	}

	if !c.Base.Arch.IsOneOf(
		build.ArchAmd64,
		build.ArchArm64,
		build.ArchUniversal,
		build.ArchUnknown,
	) {
		return fmt.Errorf("%w: unsupport architecture: %s", config.ErrInvalidInput, c.Base.Arch)
	}

	// NOTE: Don't check for 'lipo', that should be a runtime check.

	if err := c.Vulkan.Validate(inv); err != nil {
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
	PathSDK build.Path `toml:"sdk_path"`
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (c *Vulkan) Configure(inv build.Invocation) error {
	if err := c.PathSDK.RelTo(inv.PathManifest); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *Vulkan) Validate(_ build.Invocation) error {
	if err := c.PathSDK.CheckIsDir(); err != nil {
		return fmt.Errorf("%w: missing path to Vulkan SDK", err)
	}

	return nil
}
