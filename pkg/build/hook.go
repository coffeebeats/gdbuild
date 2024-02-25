package build

/* -------------------------------------------------------------------------- */
/*                                Struct: Hook                                */
/* -------------------------------------------------------------------------- */

// Hook contains commands to execute before and after a build step.
//
// TODO: Allow per-hook shell settings.
type Hook struct {
	// Pre contains a command to run *before* a build step.
	Pre []string `toml:"prebuild"`
	// Post contains a command to run *after* a build step.
	Post []string `toml:"postbuild"`
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Hook) Merge(other *Hook) error {
	if other == nil {
		return nil
	}

	if c == nil {
		*c = *other

		return nil
	}

	c.Pre = append(c.Pre, other.Pre...)
	c.Post = append(c.Post, other.Post...)

	return nil
}
