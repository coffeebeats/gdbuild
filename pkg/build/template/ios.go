package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/command"
	"github.com/coffeebeats/gdbuild/internal/merge"
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

/* ------------------------- Impl: command.Commander ------------------------ */

func (c *IOS) Command() (*command.Command, error) {
	return nil, ErrUnimplemented
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (c *IOS) Configure(inv *build.Invocation) error {
	if err := c.Base.Configure(inv); err != nil {
		return err
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

	if err := c.PathSDK.CheckIsDirOrEmpty(); err != nil {
		return err
	}

	return nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *IOS) Merge(other *IOS) error {
	if c == nil || other == nil {
		return nil
	}

	if err := c.Base.Merge(other.Base); err != nil {
		return err
	}

	if err := merge.Pointer(c.Simulator, other.Simulator); err != nil {
		return fmt.Errorf("%w: use_simulator", err)
	}

	if err := merge.Primitive(&c.PathSDK, other.PathSDK); err != nil {
		return fmt.Errorf("%w: sdk_path", err)
	}

	return nil
}
