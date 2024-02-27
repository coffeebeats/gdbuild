package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/merge"
	"github.com/coffeebeats/gdbuild/pkg/build"
	"github.com/coffeebeats/gdbuild/pkg/build/platform"
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

var _ Template = (*Android)(nil)

/* ----------------------------- Impl: Template ----------------------------- */

func (c *Android) BaseTemplate() *Base {
	return c.Base
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *Android) Action() (action.Action, error) { //nolint:ireturn
	return nil, nil
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (c *Android) Configure(inv *build.Invocation) error {
	if err := c.Base.Configure(inv); err != nil {
		return err
	}

	if c.Base.Arch == platform.ArchUnknown {
		c.Base.Arch = platform.ArchUniversal
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

	switch c.Base.Arch {
	case platform.ArchArm32, platform.ArchArm64:
	case platform.ArchI386, platform.ArchAmd64:
	case platform.ArchUnknown:
	default:
		return fmt.Errorf("%w: unsupport architecture: %s", ErrInvalidInput, c.Base.Arch)
	}

	if err := c.PathGradlew.CheckIsFileOrEmpty(); err != nil {
		return err
	}

	if err := c.PathSDK.CheckIsDirOrEmpty(); err != nil {
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
