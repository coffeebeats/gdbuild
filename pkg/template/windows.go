package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/pkg/build"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Windows                              */
/* -------------------------------------------------------------------------- */

// Windows contains 'Windows'-specific settings for constructing a custom Godot
// export template.
type Windows struct {
	*Base

	// PathIcon is a path to a Windows application icon.
	PathIcon build.Path `toml:"icon_path"`
}

/* ----------------------------- Impl: Template ----------------------------- */

func (c *Windows) BaseTemplate() *Base {
	return c.Base
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *Windows) Action() (action.Action, error) { //nolint:ireturn
	return nil, nil
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (c *Windows) Configure(inv *build.Invocation) error {
	if err := c.Base.Configure(inv); err != nil {
		return err
	}

	if c.Base.Arch == build.ArchUnknown {
		c.Base.Arch = build.ArchAmd64
	}

	if err := c.PathIcon.RelTo(inv.PathManifest); err != nil {
		return err
	}

	return nil
}

/* -------------------------- Impl: build.Validate -------------------------- */

func (c *Windows) Validate() error {
	if err := c.Base.Validate(); err != nil {
		return err
	}

	switch c.Base.Arch {
	case build.ArchAmd64, build.ArchI386:
	case build.ArchUnknown:
	default:
		return fmt.Errorf("%w: unsupport architecture: %s", ErrInvalidInput, c.Base.Arch)
	}

	if err := c.PathIcon.CheckIsFileOrEmpty(); err != nil {
		return err
	}

	return nil
}
