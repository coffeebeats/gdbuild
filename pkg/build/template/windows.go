package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/command"
	"github.com/coffeebeats/gdbuild/internal/merge"
	"github.com/coffeebeats/gdbuild/pkg/build"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Windows                              */
/* -------------------------------------------------------------------------- */

// Windows contains 'Windows'-specific settings for constructing a custom Godot
// export template.
type Windows struct {
	*Base

	// PathIcon is a path to a Windows application icon.
	PathIcon build.Path `toml:"icon_path"`
}

/* ------------------------- Impl: command.Commander ------------------------ */

func (c *Windows) Command() (*command.Command, error) {
	return nil, ErrUnimplemented
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (c *Windows) Configure(inv *build.Invocation) error {
	if err := c.Base.Configure(inv); err != nil {
		return err
	}

	if err := c.PathIcon.RelTo(inv.PathManifest); err != nil {
		return err
	}

	return nil
}

/* -------------------------- Impl: build.Validate -------------------------- */

func (c *Windows) Validate() error {
	if err := c.Base.Validate(); err != nil {
		return err
	}

	if err := c.PathIcon.CheckIsFileOrEmpty(); err != nil {
		return err
	}

	return nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Windows) Merge(other *Windows) error {
	if c == nil || other == nil {
		return nil
	}

	if err := c.Base.Merge(other.Base); err != nil {
		return err
	}

	if err := merge.Primitive(&c.PathIcon, other.PathIcon); err != nil {
		return fmt.Errorf("%w: icon_path", err)
	}

	return nil
}
