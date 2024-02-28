package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/action"
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

var _ Template = (*Android)(nil)

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *Android) Action() (action.Action, error) { //nolint:ireturn
	return nil, ErrUnimplemented
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (c *Android) Configure(inv *build.Invocation) error {
	if err := c.Base.Configure(inv); err != nil {
		return err
	}

	if c.Base.Arch == build.ArchUnknown {
		c.Base.Arch = build.ArchUniversal
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
	case build.ArchArm32, build.ArchArm64:
	case build.ArchI386, build.ArchAmd64:
	case build.ArchUnknown:
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
