package target

import "github.com/coffeebeats/gdbuild/internal/action"

/* -------------------------------------------------------------------------- */
/*                               Struct: Windows                              */
/* -------------------------------------------------------------------------- */

// Windows contains 'Windows'-specific settings for exporting a Godot project.
type Windows struct {
	*Base
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *Windows) Action() (action.Action, error) { //nolint:ireturn
	return nil, nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Windows) Merge(other *Windows) error {
	if c == nil || other == nil {
		return nil
	}

	if err := c.Base.Merge(other.Base); err != nil {
		return err
	}

	return nil
}