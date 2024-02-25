package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/command"
	"github.com/coffeebeats/gdbuild/internal/merge"
	"github.com/coffeebeats/gdbuild/pkg/build"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Android                              */
/* -------------------------------------------------------------------------- */

// Android contains 'Android'-specific settings for constructing a custom Godot
// export template.
type Android struct {
	*Base

	// PathGradle is the path to a Gradle wrapper executable.
	PathGradlew build.Path `toml:"gradlew_path"`
	// PathSDK is the path to the Android SDK root.
	PathSDK build.Path `toml:"sdk_path"`
}

/* ------------------------- Impl: command.Commander ------------------------ */

func (c *Android) Command() (*command.Command, error) {
	return nil, ErrUnimplemented
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (c *Android) Configure(inv *build.Invocation) error {
	if err := c.Base.Configure(inv); err != nil {
		return err
	}

	if err := c.PathGradlew.RelTo(inv.PathManifest); err != nil {
		return err
	}

	if err := c.PathSDK.RelTo(inv.PathManifest); err != nil {
		return err
	}

	return nil
}

/* -------------------------- Impl: build.Validate -------------------------- */

func (c *Android) Validate() error {
	if err := c.Base.Validate(); err != nil {
		return err
	}

	return nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Android) Merge(other *Android) error {
	if c == nil || other == nil {
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
