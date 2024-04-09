package run

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"

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
	// DryRun determines whether the specified command will be executed.
	DryRun bool
	// Target is the name of the target specified during an export command.
	//
	// HACK: Context should realistically be divided into per-command contexts
	// so that properties can be specialized. For now, deal with this. It's only
	// set for 'gdbuild target'.
	Target string

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

	tmp string
}

/* ----------------------------- Method: BinPath ---------------------------- */

// BinPath returns the path to the Godot template artifacts are compilation.
func (c *Context) BinPath() osutil.Path {
	return osutil.Path(filepath.Join(c.PathWorkspace.String(), "bin"))
}

/* ---------------------------- Method: GodotPath --------------------------- */

// GodotPath returns the path to the Godot editor executable.
func (c *Context) GodotPath() osutil.Path {
	pathTmp, err := c.TempDir()
	if err != nil {
		return ""
	}

	return osutil.Path(filepath.Join(pathTmp, engine.EditorName()))
}

/* --------------------------- Method: HasTempDir --------------------------- */

// HasTempDir returns whether a run-specific temporary directory has been
// created.
func (c *Context) HasTempDir() bool {
	return c.tmp != ""
}

/* -------------------- Method: GodotProjectManifestPath -------------------- */

// GodotProjectManifestPath returns the path to the Godot project manifest file.
func (c *Context) GodotProjectManifestPath() osutil.Path {
	return osutil.Path(filepath.Join(c.PathWorkspace.String(), "project.godot"))
}

/* ----------------------------- Method: TempDir ---------------------------- */

// TempDir constructs a temporary directory for the duration of the context.
// Only one directory will ever be created.
func (c *Context) TempDir() (string, error) {
	if c.DryRun {
		return filepath.Join(os.TempDir(), "gdbuild-*"), nil
	}

	if c.tmp != "" {
		return c.tmp, nil
	}

	tmp, err := os.MkdirTemp("", "gdbuild-*")
	if err != nil {
		return "", err
	}

	log.Debugf("created temporary directory: %s", tmp)

	c.tmp = tmp

	return tmp, nil
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
