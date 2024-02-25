package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/command"
	"github.com/coffeebeats/gdbuild/internal/merge"
	"github.com/coffeebeats/gdbuild/pkg/build"
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

/* ------------------------- Impl: command.Commander ------------------------ */

func (c *Web) Command() (*command.Command, error) {
	return nil, ErrUnimplemented
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (c *Web) Configure(inv *build.Invocation) error {
	if err := c.Base.Configure(inv); err != nil {
		return err
	}

	return nil
}

/* -------------------------- Impl: build.Validate -------------------------- */

func (c *Web) Validate() error {
	if err := c.Base.Validate(); err != nil {
		return err
	}

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

	if err := merge.Pointer(c.EnableEval, other.EnableEval); err != nil {
		return fmt.Errorf("%w: enable_eval", err)
	}

	return nil
}
