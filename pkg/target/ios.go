package target

/* -------------------------------------------------------------------------- */
/*                                 Struct: IOS                                */
/* -------------------------------------------------------------------------- */

// IOS contains 'IOS'-specific settings for exporting a Godot project.
type IOS struct {
	*Base
}

/* ----------------------------- Impl: Commander ---------------------------- */

func (c *IOS) Command() []string {
	return nil
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
