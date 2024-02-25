package target

/* -------------------------------------------------------------------------- */
/*                                 Struct: Web                                */
/* -------------------------------------------------------------------------- */

// Web contains 'Web'-specific settings for exporting a Godot project.
type Web struct {
	*Base
}

/* ----------------------------- Impl: Commander ---------------------------- */

func (c *Web) Command() []string {
	return nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Web) Merge(other *Web) error {
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

	return nil
}
