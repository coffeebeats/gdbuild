package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/command"
	"github.com/coffeebeats/gdbuild/pkg/build"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: Linux                               */
/* -------------------------------------------------------------------------- */

// Linux contains 'linux'-specific settings for constructing a custom Godot
// export template.
type Linux struct {
	*Base
}

/* ----------------------------- Impl: Commander ---------------------------- */

func (c *Linux) Command() (*command.Command, error) {
	base, err := c.Base.Command()
	if err != nil {
		return nil, err
	}

	arch := build.ArchAmd64
	if c.Base.Arch != build.ArchUnknown {
		arch = c.Base.Arch
	}

	switch arch {
	case build.ArchAmd64:
		base.Args = append(base.Args, "arch="+arch.String())
	default:
		return nil, fmt.Errorf("%w: unsupport architecture: %s", ErrInvalidInput, arch)
	}

	return base, nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Linux) Merge(other *Linux) error {
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
