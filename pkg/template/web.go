package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/merge"
	"github.com/coffeebeats/gdbuild/pkg/build"
	"github.com/coffeebeats/gdbuild/pkg/build/platform"
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

/* ----------------------------- Impl: Template ----------------------------- */

func (c *Web) BaseTemplate() *Base {
	return c.Base
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *Web) Action() (action.Action, error) { //nolint:ireturn
	return nil, nil
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
	case platform.ArchUnknown:
	default:
		return fmt.Errorf("%w: unsupport architecture: %s", ErrInvalidInput, c.Base.Arch)
	}

	return nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Web) Merge(other *Web) error {
	if c == nil || other == nil {
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
