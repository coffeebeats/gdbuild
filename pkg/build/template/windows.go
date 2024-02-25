package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/command"
	"github.com/coffeebeats/gdbuild/internal/merge"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Windows                              */
/* -------------------------------------------------------------------------- */

// Windows contains 'Windows'-specific settings for constructing a custom Godot
// export template.
type Windows struct {
	*Base

	// PathIcon is a path to a Windows application icon.
	PathIcon string `toml:"icon_path"`
}

/* ----------------------------- Impl: Commander ---------------------------- */

func (c *Windows) Command() (*command.Command, error) {
	return nil, ErrUnimplemented
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Windows) Merge(other *Windows) error {
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

	if err := merge.Primitive(&c.PathIcon, other.PathIcon); err != nil {
		return fmt.Errorf("%w: icon_path", err)
	}

	return nil
}
