package config

import (
	"fmt"
	"path/filepath"

	"github.com/coffeebeats/gdbuild/internal/pathutil"
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
	PathBuild pathutil.Path
	// PathManifest is the path to the GDBuild manifest. This is used to locate
	// relative paths in various other properties.
	PathManifest pathutil.Path
	// PathOut is the directory in which to move built artifacts to.
	PathOut pathutil.Path
}

/* ----------------------------- Method: BinPath ---------------------------- */

// BinPath returns the path to the Godot template artifacts are compilation.
func (c *Context) BinPath() pathutil.Path {
	return pathutil.Path(filepath.Join(c.PathBuild.String(), "bin"))
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
