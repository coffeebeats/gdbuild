package target

import "github.com/coffeebeats/gdbuild/internal/action"

/* -------------------------------------------------------------------------- */
/*                                 Struct: Web                                */
/* -------------------------------------------------------------------------- */

// Web contains 'Web'-specific settings for exporting a Godot project.
type Web struct {
	*Base
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *Web) Action() (action.Action, error) { //nolint:ireturn
	return nil, nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Web) Merge(other *Web) error {
	if c == nil || other == nil {
		return nil
	}

	if err := c.Base.Merge(other.Base); err != nil {
		return err
	}

	return nil
}
