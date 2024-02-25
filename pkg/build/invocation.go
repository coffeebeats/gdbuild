package build

import (
	"fmt"
	"os"
	"path/filepath"
)

/* -------------------------------------------------------------------------- */
/*                             Struct: Invocation                             */
/* -------------------------------------------------------------------------- */

// Invocation are build command inputs that are invocation-specific. These need
// to be explicitly set per invocation as they can't be parsed from a GDBuild
// manifest.
type Invocation struct {
	// Features is the list of feature tags to enable.
	Features []string
	// Platform is the target platform to build for.
	Platform OS
	// Profile is the GDBuild optimization level to build with.
	Profile Profile

	// PathManifest is the directory in which the GDBuild manifest is located.
	// This is used to locate relative paths in various other properties.
	PathManifest Path
	// PathBuild is the directory in which to build the template in. All input
	// artifacts will be copied here and the SCons build command will be
	// executed within this directory. Defaults to a temporary directory.
	PathBuild Path
}

/* ---------------------------- Method: Validate ---------------------------- */

func (c *Invocation) Validate() error {
	if _, err := ParseOS(c.Platform.String()); err != nil {
		return err
	}

	if _, err := ParseProfile(c.Profile.String()); err != nil {
		return err
	}

	info, err := os.Stat(string(c.PathManifest))
	if err != nil {
		return fmt.Errorf("%w: manifest directory: %s", err, c.PathManifest)
	}

	if !info.IsDir() {
		return fmt.Errorf(
			"%w: manifest directory: not a directory: %s",
			ErrInvalidInput,
			c.PathManifest,
		)
	}

	c.PathManifest = Path(filepath.Clean(string(c.PathManifest)))

	info, err = os.Stat(string(c.PathBuild))
	if err != nil {
		return fmt.Errorf("%w: manifest directory: %s", err, c.PathManifest)
	}

	if !info.IsDir() {
		return fmt.Errorf(
			"%w: manifest directory: not a directory: %s",
			ErrInvalidInput,
			c.PathManifest,
		)
	}

	c.PathBuild = Path(filepath.Clean(string(c.PathBuild)))

	return nil
}
