package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/command"
	"github.com/coffeebeats/gdbuild/internal/merge"
)

/* -------------------------------------------------------------------------- */
/*                                 Struct: Web                                */
/* -------------------------------------------------------------------------- */

// Web contains 'Web'-specific settings for constructing a custom Godot
// export template.
type Web struct {
	*Base

	// EnableEval defines whether to enable Javascript "eval()" calls.
	EnableEval *bool `toml:"enable_eval"`
}

/* ----------------------------- Impl: Commander ---------------------------- */

func (c *Web) Command() (*command.Command, error) {
	return nil, ErrUnimplemented
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

	if err := merge.Pointer(c.EnableEval, other.EnableEval); err != nil {
		return fmt.Errorf("%w: enable_eval", err)
	}

	return nil
}
