package run

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
)

var ErrMissingInput = errors.New("missing input")

/* -------------------------------------------------------------------------- */
/*                               Struct: Context                              */
/* -------------------------------------------------------------------------- */

// Context contains build command inputs that are invocation-specific. These
// need to be explicitly set per invocation as they can't be parsed from a
// GDBuild manifest.
type Context struct {
	// Verbose determines whether to enable additional logging output.
	Verbose bool

	// Features is the list of feature tags to enable.
	Features []string
	// Platform is the target platform to build for.
	Platform platform.OS
	// Profile is the GDBuild optimization level to build with.
	Profile engine.Profile

	// PathBuild is the directory in which to build the template in. All input
	// artifacts will be copied here and the SCons build command will be
	// executed within this directory. Defaults to a temporary directory.
	PathBuild osutil.Path
	// PathManifest is the path to the GDBuild manifest. This is used to locate
	// relative paths in various other properties.
	PathManifest osutil.Path
	// PathOut is the directory in which to move built artifacts to.
	PathOut osutil.Path
}

/* ----------------------------- Method: BinPath ---------------------------- */

// BinPath returns the path to the Godot template artifacts are compilation.
func (c *Context) BinPath() osutil.Path {
	return osutil.Path(filepath.Join(c.PathBuild.String(), "bin"))
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *Context) Validate() error {
	if _, err := platform.ParseOS(c.Platform.String()); err != nil {
		return err
	}

	if _, err := engine.ParseProfile(c.Profile.String()); err != nil {
		return err
	}

	if err := c.PathManifest.CheckIsFile(); err != nil {
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
