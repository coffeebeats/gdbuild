package target

import "github.com/coffeebeats/gdbuild/internal/action"

/* -------------------------------------------------------------------------- */
/*                                Struct: Linux                               */
/* -------------------------------------------------------------------------- */

// Linux contains 'Linux'-specific settings for exporting a Godot project.
type Linux struct {
	*Base
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *Linux) Action() (action.Action, error) { //nolint:ireturn
	return nil, nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Linux) Merge(other *Linux) error {
	if c == nil || other == nil {
		return nil
	}

	if err := c.Base.Merge(other.Base); err != nil {
		return err
	}

	return nil
}
