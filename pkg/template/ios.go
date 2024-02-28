package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/pkg/build"
)

/* -------------------------------------------------------------------------- */
/*                                 Struct: IOS                                */
/* -------------------------------------------------------------------------- */

// IOS contains 'IOS'-specific settings for constructing a custom Godot
// export template.
type IOS struct {
	*Base

	// PathSDK is the path to the IOS SDK root.
	PathSDK build.Path `toml:"sdk_path"`
	// Simulator denotes whether to build for the iOS simulator.
	Simulator *bool `toml:"use_simulator"`
}

/* ----------------------------- Impl: Template ----------------------------- */

func (c *IOS) BaseTemplate() *Base {
	return c.Base
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *IOS) Action() (action.Action, error) { //nolint:ireturn
	return nil, nil
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (c *IOS) Configure(inv *build.Invocation) error {
	if err := c.Base.Configure(inv); err != nil {
		return err
	}

	if c.Base.Arch == build.ArchUnknown {
		c.Base.Arch = build.ArchArm64
	}

	if err := c.PathSDK.RelTo(inv.PathManifest); err != nil {
		return err
	}

	return nil
}

/* -------------------------- Impl: build.Validate -------------------------- */

func (c *IOS) Validate() error {
	if err := c.Base.Validate(); err != nil {
		return err
	}

	switch c.Base.Arch {
	case build.ArchAmd64, build.ArchArm64:
	case build.ArchUnknown:
	default:
		return fmt.Errorf("%w: unsupport architecture: %s", ErrInvalidInput, c.Base.Arch)
	}

	if err := c.PathSDK.CheckIsDirOrEmpty(); err != nil {
		return err
	}

	return nil
}
