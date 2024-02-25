package target

/* -------------------------------------------------------------------------- */
/*                               Struct: Android                              */
/* -------------------------------------------------------------------------- */

// Android contains 'Android'-specific settings for exporting a Godot project.
type Android struct {
	*Base
}

/* ----------------------------- Impl: Commander ---------------------------- */

func (c *Android) Command() []string {
	return nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Android) Merge(other *Android) error {
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
