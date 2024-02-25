package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/command"
	"github.com/coffeebeats/gdbuild/internal/merge"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Android                              */
/* -------------------------------------------------------------------------- */

// Android contains 'Android'-specific settings for constructing a custom Godot
// export template.
type Android struct {
	*Base

	// PathGradle is the path to a Gradle wrapper executable.
	PathGradlew string `toml:"gradlew_path"`
	// PathSDK is the path to the Android SDK root.
	PathSDK string `toml:"sdk_path"`
}

/* ----------------------------- Impl: Commander ---------------------------- */

func (c *Android) Command() (*command.Command, error) {
	return nil, ErrUnimplemented
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Android) Merge(other *Android) error {
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

	if err := merge.Primitive(&c.PathGradlew, other.PathGradlew); err != nil {
		return fmt.Errorf("%w: gradlew_path", err)
	}

	if err := merge.Primitive(&c.PathSDK, other.PathSDK); err != nil {
		return fmt.Errorf("%w: sdk_path", err)
	}

	return nil
}
