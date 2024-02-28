package target

import "github.com/coffeebeats/gdbuild/internal/action"

/* -------------------------------------------------------------------------- */
/*                                 Struct: IOS                                */
/* -------------------------------------------------------------------------- */

// IOS contains 'IOS'-specific settings for exporting a Godot project.
type IOS struct {
	*Base
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *IOS) Action() (action.Action, error) { //nolint:ireturn
	return nil, nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *IOS) Merge(other *IOS) error {
	if c == nil || other == nil {
		return nil
	}

	if err := c.Base.Merge(other.Base); err != nil {
		return err
	}

	return nil
}
