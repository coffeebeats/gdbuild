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

	// PathLipo is the path to the 'lipo' executable. Only used if 'arch' is set
	// to 'build.ArchUniversal' and defaults to 'lipo'.
	PathLipo build.Path `toml:"lipo_path"`

	Vulkan Vulkan `toml:"vulkan"`
}

/* ----------------------------- Impl: Template ----------------------------- */

func (c *MacOS) BaseTemplate() *Base {
	return c.Base
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *MacOS) Action() (action.Action, error) { //nolint:ireturn
	return nil, nil
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

	if err := c.PathLipo.RelTo(inv.PathManifest); err != nil {
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

	if err := c.PathLipo.CheckIsFileOrEmpty(); err != nil {
		return err
	}

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
