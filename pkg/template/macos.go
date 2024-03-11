package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/pkg/build"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: MacOS                               */
/* -------------------------------------------------------------------------- */

// MacOS contains 'macos'-specific settings for constructing a custom Godot
// export template.
type MacOS struct {
	*Base

	// LipoCommand contains arguments used to invoke 'lipo'. Defaults to
	// ["lipo"]. Only used if 'arch' is set to 'build.ArchUniversal'.
	LipoCommand []string `toml:"lipo_command"`

	// Vulkan defines Vulkan-related configuration.
	Vulkan Vulkan `toml:"vulkan"`
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *MacOS) Action() (action.Action, error) { //nolint:cyclop,funlen,ireturn
	switch a := c.Base.Arch; a {
	case build.ArchAmd64, build.ArchArm64:
		cmd, err := c.Base.action()
		if err != nil {
			return nil, err
		}

		cmd.Args = append(cmd.Args, "platform="+build.OSMacOS.String())

		if volk := c.Vulkan.Dynamic; volk != nil && *volk {
			cmd.Args = append(cmd.Args, "use_volk=yes")
		}

		if vulkan := c.Vulkan.PathSDK; vulkan != "" {
			cmd.Args = append(cmd.Args, "vulkan_sdk_path="+string(vulkan))
		}

		return c.wrapMacOSBuildCommand(cmd), nil
	case build.ArchUniversal:
		// First, create the 'x86_64' binary.
		templateAmd64 := *c
		templateAmd64.Base.Arch = build.ArchAmd64

		buildAmd64, err := templateAmd64.Action()
		if err != nil {
			return nil, err
		}

		cmdAmd64, ok := buildAmd64.(action.Sequence).Unwrap().(*action.Process)
		if !ok {
			return nil, fmt.Errorf("%w: failed to unwrap action", ErrInvalidInput)
		}

		// Next, create the 'arm64' binary.
		templateArm64 := *c
		templateArm64.Base.Arch = build.ArchArm64

		buildArm64, err := templateArm64.Action()
		if err != nil {
			return nil, err
		}

		cmdArm64, ok := buildArm64.(action.Sequence).Unwrap().(*action.Process)
		if !ok {
			return nil, fmt.Errorf("%w: failed to unwrap action", ErrInvalidInput)
		}

		// Finally, merge the two binaries together.
		lipo := c.LipoCommand
		if len(lipo) == 0 {
			lipo = append(lipo, "lipo")
		}

		targetName := c.Base.targetName()

		cmdLipo := &action.Process{
			Directory:   string(c.Invocation.BinPath()),
			Environment: nil,

			Shell: cmdArm64.Shell,

			Verbose: c.Invocation.Verbose,

			Args: append(
				lipo,
				"-create",
				fmt.Sprintf("godot.macos.%s.x86_64", targetName),
				fmt.Sprintf("godot.macos.%s.arm64", targetName),
				"-output",
				fmt.Sprintf("godot.macos.%s.universal", targetName),
			),
		}

		return c.wrapMacOSBuildCommand(
			cmdAmd64.
				AndThen(cmdArm64).
				AndThen(cmdLipo),
		), nil

	default:
		return nil, fmt.Errorf("%w: unsupported architecture: %s", ErrInvalidInput, a)
	}
}

/* ---------------------- Method: wrapMacOSBuildCommand --------------------- */

func (c *MacOS) wrapMacOSBuildCommand(cmd action.Action) action.Sequence {
	return c.wrapBuildCommand(cmd)
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (c *MacOS) Configure(inv *build.Invocation) error {
	if err := c.Base.Configure(inv); err != nil {
		return err
	}

	if c.Base.Arch == build.ArchUnknown {
		c.Base.Arch = build.ArchUniversal
	}

	if err := c.Vulkan.Configure(inv); err != nil {
		return err
	}

	return nil
}

/* -------------------------- Impl: build.Validate -------------------------- */

func (c *MacOS) Validate() error {
	if err := c.Base.Validate(); err != nil {
		return err
	}

	switch c.Base.Arch {
	case build.ArchAmd64, build.ArchArm64, build.ArchUniversal:
	case build.ArchUnknown:
	default:
		return fmt.Errorf("%w: unsupport architecture: %s", ErrInvalidInput, c.Base.Arch)
	}

	// NOTE: Don't check for 'lipo', that should be a runtime check.

	if err := c.Vulkan.Validate(); err != nil {
		return err
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*                               Struct: Vulkan                               */
/* -------------------------------------------------------------------------- */

// Vulkan defines the settings required by the MacOS template for including
// Vulkan support.
type Vulkan struct {
	// Dynamic enables dynamically linking Vulkan to the template.
	Dynamic *bool `toml:"dynamic"`
	// PathSDK is the path to the Vulkan SDK root.
	PathSDK build.Path `toml:"sdk_path"`
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (c *Vulkan) Configure(inv *build.Invocation) error {
	if err := c.PathSDK.RelTo(inv.PathManifest); err != nil {
		return err
	}

	return nil
}

/* -------------------------- Impl: build.Validate -------------------------- */

func (c *Vulkan) Validate() error {
	if err := c.PathSDK.CheckIsDirOrEmpty(); err != nil {
		return err
	}

	return nil
}
