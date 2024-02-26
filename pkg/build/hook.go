package build

import "github.com/coffeebeats/gdbuild/internal/action"

/* -------------------------------------------------------------------------- */
/*                                Struct: Hook                                */
/* -------------------------------------------------------------------------- */

// Hook contains commands to execute before and after a build step.
//
// TODO: Allow per-hook shell settings.
type Hook struct {
	// Pre contains a command to run *before* a build step.
	Pre []action.Command `toml:"prebuild"`
	// Post contains a command to run *after* a build step.
	Post []action.Command `toml:"postbuild"`
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Hook) Merge(other *Hook) error {
	if c == nil || other == nil {
		return nil
	}

	c.Pre = append(c.Pre, other.Pre...)
	c.Post = append(c.Post, other.Post...)

	return nil
}
