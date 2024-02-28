package target

import "github.com/coffeebeats/gdbuild/internal/action"

/* -------------------------------------------------------------------------- */
/*                               Struct: Android                              */
/* -------------------------------------------------------------------------- */

// Android contains 'Android'-specific settings for exporting a Godot project.
type Android struct {
	*Base
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *Android) Action() (action.Action, error) { //nolint:ireturn
	return nil, nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Android) Merge(other *Android) error {
	if c == nil || other == nil {
		return nil
	}

	if err := c.Base.Merge(other.Base); err != nil {
		return err
	}

	return nil
}
