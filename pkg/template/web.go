package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/action"
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

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *Web) Action() (action.Action, error) { //nolint:ireturn
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

	switch c.Base.Arch {
	case build.ArchUnknown:
	default:
		return fmt.Errorf("%w: unsupport architecture: %s", ErrInvalidInput, c.Base.Arch)
	}

	return nil
}
