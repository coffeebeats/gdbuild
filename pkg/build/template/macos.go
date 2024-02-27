package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/merge"
	"github.com/coffeebeats/gdbuild/pkg/build"
	"github.com/coffeebeats/gdbuild/pkg/build/platform"
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

	if c.Base.Arch == platform.ArchUnknown {
		c.Base.Arch = platform.ArchUniversal
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
	case platform.ArchAmd64, platform.ArchArm64, platform.ArchUniversal:
	case platform.ArchUnknown:
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

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *MacOS) Merge(other *MacOS) error {
	if c == nil || other == nil {
		return nil
	}

	if err := c.Base.Merge(other.Base); err != nil {
		return err
	}

	if err := c.Vulkan.Merge(&other.Vulkan); err != nil {
		return err
	}

	if err := merge.Primitive(&c.PathLipo, other.PathLipo); err != nil {
		return fmt.Errorf("%w: lipo_path", err)
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

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Vulkan) Merge(other *Vulkan) error {
	if c == nil || other == nil {
		return nil
	}

	if err := merge.Pointer(c.Dynamic, other.Dynamic); err != nil {
		return fmt.Errorf("%w: dynamic", err)
	}

	if err := merge.Primitive(&c.PathSDK, other.PathSDK); err != nil {
		return fmt.Errorf("%w: sdk_path", err)
	}

	return nil
}
