package target

import "github.com/coffeebeats/gdbuild/internal/action"

/* -------------------------------------------------------------------------- */
/*                                Struct: MacOS                               */
/* -------------------------------------------------------------------------- */

// MacOS contains 'MacOS'-specific settings for exporting a Godot project.
type MacOS struct {
	*Base
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *MacOS) Action() (action.Action, error) { //nolint:ireturn
	return nil, nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *MacOS) Merge(other *MacOS) error {
	if c == nil || other == nil {
		return nil
	}

	if err := c.Base.Merge(other.Base); err != nil {
		return err
	}

	return nil
}
