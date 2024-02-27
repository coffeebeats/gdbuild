package target

/* -------------------------------------------------------------------------- */
/*                               Struct: Windows                              */
/* -------------------------------------------------------------------------- */

// Windows contains 'Windows'-specific settings for exporting a Godot project.
type Windows struct {
	*Base
}

/* ----------------------------- Impl: Commander ---------------------------- */

func (c *Windows) Command() []string {
	return nil
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
