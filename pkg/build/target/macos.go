package target

/* -------------------------------------------------------------------------- */
/*                                Struct: MacOS                               */
/* -------------------------------------------------------------------------- */

// MacOS contains 'MacOS'-specific settings for exporting a Godot project.
type MacOS struct {
	*Base
}

/* ----------------------------- Impl: Commander ---------------------------- */

func (c *MacOS) Command() []string {
	return nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *MacOS) Merge(other *MacOS) error {
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
