package build

import (
	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Context                              */
/* -------------------------------------------------------------------------- */

// config.Context are build command inputs that are invocation-specific. These need
// to be explicitly set per invocation as they can't be parsed from a GDBuild
// manifest.
type Context struct {
	// Invoke contains the application-wide context pertaining to the
	// 'gdbuild' invocation.
	Invoke config.Context

	// Features is the list of feature tags to enable.
	Features []string
	// Platform is the target platform to build for.
	Platform platform.OS
	// Profile is the GDBuild optimization level to build with.
	Profile Profile
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *Context) Validate() error {
	if err := c.Invoke.Validate(); err != nil {
		return err
	}

	if _, err := platform.ParseOS(c.Platform.String()); err != nil {
		return err
	}

	if _, err := ParseProfile(c.Profile.String()); err != nil {
		return err
	}

	return nil
}
