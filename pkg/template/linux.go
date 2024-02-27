package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/pkg/build"
	"github.com/coffeebeats/gdbuild/pkg/build/platform"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: Linux                               */
/* -------------------------------------------------------------------------- */

// Linux contains 'linux'-specific settings for constructing a custom Godot
// export template.
type Linux struct {
	*Base
}

/* ----------------------------- Impl: Template ----------------------------- */

func (c *Linux) BaseTemplate() *Base {
	return c.Base
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *Linux) Action() (action.Action, error) { //nolint:ireturn
	cmd, err := c.Base.action()
	if err != nil {
		return nil, err
	}

	cmd.process.Args = append(cmd.process.Args, "platform="+platform.OSLinux.String())

	return cmd.action, nil
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (c *Linux) Configure(inv *build.Invocation) error {
	if err := c.Base.Configure(inv); err != nil {
		return err
	}

	if c.Base.Arch == platform.ArchUnknown {
		c.Base.Arch = platform.ArchAmd64
	}

	return nil
}

/* -------------------------- Impl: build.Validate -------------------------- */

func (c *Linux) Validate() error {
	if err := c.Base.Validate(); err != nil {
		return err
	}

	switch c.Base.Arch {
	case platform.ArchI386, platform.ArchAmd64:
	case platform.ArchUnknown:
	default:
		return fmt.Errorf("%w: unsupport architecture: %s", ErrInvalidInput, c.Base.Arch)
	}

	return nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Linux) Merge(other *Linux) error {
	if c == nil || other == nil {
		return nil
	}

	if err := c.Base.Merge(other.Base); err != nil {
		return err
	}

	return nil
}
