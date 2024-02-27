package build

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/pkg/build/platform"
)

/* -------------------------------------------------------------------------- */
/*                             Struct: Invocation                             */
/* -------------------------------------------------------------------------- */

// Invocation are build command inputs that are invocation-specific. These need
// to be explicitly set per invocation as they can't be parsed from a GDBuild
// manifest.
type Invocation struct {
	// Verbose determines whether to enable additional logging output.
	Verbose bool

	// Features is the list of feature tags to enable.
	Features []string
	// Platform is the target platform to build for.
	Platform platform.OS
	// Profile is the GDBuild optimization level to build with.
	Profile Profile

	// PathBuild is the directory in which to build the template in. All input
	// artifacts will be copied here and the SCons build command will be
	// executed within this directory. Defaults to a temporary directory.
	PathBuild Path
	// PathManifest is the directory in which the GDBuild manifest is located.
	// This is used to locate relative paths in various other properties.
	PathManifest Path
	// PathOut is the directory in which to move built artifacts to.
	PathOut Path
}

/* ---------------------------- Method: Validate ---------------------------- */

func (c *Invocation) Validate() error {
	if _, err := platform.ParseOS(c.Platform.String()); err != nil {
		return err
	}

	if _, err := ParseProfile(c.Profile.String()); err != nil {
		return err
	}

	if err := c.PathManifest.CheckIsDir(); err != nil {
		return err
	}

	// NOTE: PathBuild might be generated via hooks, so only verify it's set.
	if c.PathBuild == "" {
		return fmt.Errorf("%w: build path", ErrMissingInput)
	}

	// NOTE: PathOut might be generated via hooks, so only verify it's set.
	if c.PathOut == "" {
		return fmt.Errorf("%w: build path", ErrMissingInput)
	}

	return nil
}
