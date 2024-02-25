package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/command"
	"github.com/coffeebeats/gdbuild/internal/merge"
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
	PathLipo string `toml:"lipo_path"`

	Vulkan Vulkan `toml:"vulkan"`
}

/* ----------------------------- Impl: Commander ---------------------------- */

func (c *MacOS) Command() (*command.Command, error) {
	return nil, ErrUnimplemented
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *MacOS) Merge(other *MacOS) error {
	if other == nil {
		return nil
	}

	if c == nil {
		*c = *other

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
	PathSDK string `toml:"sdk_path"`
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Vulkan) Merge(other *Vulkan) error {
	if other == nil {
		return nil
	}

	if c == nil {
		*c = *other

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
