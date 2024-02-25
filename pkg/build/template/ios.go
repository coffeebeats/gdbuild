package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/command"
	"github.com/coffeebeats/gdbuild/internal/merge"
)

/* -------------------------------------------------------------------------- */
/*                                 Struct: IOS                                */
/* -------------------------------------------------------------------------- */

// IOS contains 'IOS'-specific settings for constructing a custom Godot
// export template.
type IOS struct {
	*Base

	// PathSDK is the path to the IOS SDK root.
	PathSDK string `toml:"sdk_path"`
	// Simulator denotes whether to build for the iOS simulator.
	Simulator *bool `toml:"use_simulator"`
}

/* ----------------------------- Impl: Commander ---------------------------- */

func (c *IOS) Command() (*command.Command, error) {
	return nil, ErrUnimplemented
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *IOS) Merge(other *IOS) error {
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

	if err := merge.Pointer(c.Simulator, other.Simulator); err != nil {
		return fmt.Errorf("%w: use_simulator", err)
	}

	if err := merge.Primitive(&c.PathSDK, other.PathSDK); err != nil {
		return fmt.Errorf("%w: sdk_path", err)
	}

	return nil
}
