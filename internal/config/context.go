package config

import (
	"fmt"
	"path/filepath"

	"github.com/coffeebeats/gdbuild/internal/osutil"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Context                              */
/* -------------------------------------------------------------------------- */

type Context struct {
	// Verbose determines whether to enable additional logging output.
	Verbose bool

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

/* ----------------------------- Impl: Validator ---------------------------- */

func (c *Context) Validate() error {
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
