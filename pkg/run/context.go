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

	// PathManifest is the path to the GDBuild manifest. This is used to locate
	// relative paths in various other properties.
	PathManifest osutil.Path
	// PathOut is the directory in which to move built artifacts to.
	PathOut osutil.Path
	// PathWorkspace is a working directory in which to run relevant commands.
	// The use case for the directory is command-specific.
	PathWorkspace osutil.Path
}

/* --------------------------- Method: ProjectPath -------------------------- */

// ProjectPath returns the path to the Godot project.
func (c *Context) ProjectPath() osutil.Path {
	return osutil.Path(filepath.Dir(c.PathManifest.String()))
}

/* ------------------------- Method: ProjectManifest ------------------------ */

// ProjectManifest returns the path to the Godot project configuration file.
func (c *Context) ProjectManifest() osutil.Path {
	return osutil.Path(filepath.Join(c.ProjectPath().String(), "project.godot"))
}

/* ----------------------------- Method: BinPath ---------------------------- */

// BinPath returns the path to the Godot template artifacts are compilation.
func (c *Context) BinPath() osutil.Path {
	return osutil.Path(filepath.Join(c.PathWorkspace.String(), "bin"))
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

	// NOTE: PathWorkspace might be generated via hooks, so only verify it's set.
	if c.PathWorkspace == "" {
		return fmt.Errorf("%w: build path", ErrMissingInput)
	}

	return nil
}
