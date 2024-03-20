package build

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
	"github.com/coffeebeats/gdenv/pkg/install"
	"github.com/coffeebeats/gdenv/pkg/store"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/osutil"
)

// ErrConflictingValue is returned when two settings conflict with eachother.
var ErrConflictingValue = errors.New("conflicting setting")

/* -------------------------------------------------------------------------- */
/*                                Struct: Godot                               */
/* -------------------------------------------------------------------------- */

// Godot defines options and settings for which Godot version to use. Note that
// only one of these options can be used at a time, but one *must* be specified.
type Godot struct {
	// PathSource is a path to a directory containing the Godot source code.
	PathSource Path `toml:"src_path"`
	// Version is a specific version label to download.
	Version string `toml:"version"`
	// VersionFile is a file containing just the a version label to download.
	VersionFile Path `toml:"version_file"`
}

/* ----------------------------- Method: IsEmpty ---------------------------- */

// IsEmpty returns whether all properties are unset, implying there is no need
// to vendor Godot source code.
func (c *Godot) IsEmpty() bool {
	return c.PathSource == "" && c.Version == "" && c.VersionFile == ""
}

/* ---------------------------- Method: VendorTo ---------------------------- */

// VendorTo vendors the Godot source code to the specified directory.
func (c *Godot) VendorTo(ctx context.Context, out string) error {
	if c.IsEmpty() {
		return fmt.Errorf("%w: no Godot version or source path set", ErrMissingInput)
	}

	if c.PathSource != "" {
		return osutil.CopyDir(c.PathSource.String(), out)
	}

	var input string

	switch {
	case c.Version != "":
		input = c.Version
	case c.VersionFile != "":
		contents, err := os.ReadFile(c.VersionFile.String())
		if err != nil {
			return err
		}

		input = string(contents)
	}

	v, err := version.Parse(input)
	if err != nil {
		return fmt.Errorf("%w: %w: %s", ErrInvalidInput, err, c.Version)
	}

	storePath, err := store.Path()
	if err != nil {
		tmp, err := os.MkdirTemp("", "gdbuild-*")
		if err != nil {
			return err
		}

		log.Debugf("no 'gdenv' store found; using temporary directory: %s", tmp)

		storePath = tmp
	}

	return install.Vendor(ctx, storePath, v, out /* force= */, false)
}

/* ---------------------------- config.Configurer --------------------------- */

func (c *Godot) Configure(inv Invocation) error {
	if err := c.PathSource.RelTo(inv.PathManifest); err != nil {
		return err
	}

	if err := c.VersionFile.RelTo(inv.PathManifest); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *Godot) Validate(_ Invocation) error { //nolint:cyclop,funlen
	if c.IsEmpty() {
		return fmt.Errorf("%w: no Godot version specified in manifest", ErrMissingInput)
	}

	if c.PathSource != "" {
		if c.Version != "" || c.VersionFile != "" {
			return fmt.Errorf(
				"%w: cannot specify 'version' or 'version_file' with 'src_path'",
				ErrConflictingValue,
			)
		}

		// NOTE: Can't fully validate that directory exists in case they're
		// generated by a hook.
		if err := c.PathSource.CheckIsDirOrEmpty(); err != nil {
			return err
		}

		return nil
	}

	var input string

	if c.Version != "" {
		if c.PathSource != "" || c.VersionFile != "" {
			return fmt.Errorf(
				"%w: cannot specify 'src_path' or 'version_file' with 'version'",
				ErrConflictingValue,
			)
		}

		input = c.Version
	}

	if c.VersionFile != "" {
		if c.PathSource != "" || c.Version != "" {
			return fmt.Errorf(
				"%w: cannot specify 'src_path' or 'version' with 'version_file'",
				ErrConflictingValue,
			)
		}

		// NOTE: Can't fully validate that the file exists in case it's
		// generated by a hook.
		if err := c.VersionFile.CheckIsFileOrEmpty(); err != nil {
			return err
		}

		contents, err := os.ReadFile(c.VersionFile.String())
		if err != nil {
			return err
		}

		input = string(contents)
	}

	if _, err := version.Parse(input); err != nil {
		return fmt.Errorf("%w: %w: %s", ErrInvalidInput, err, c.Version)
	}

	return nil
}

/* --------------------- Function: NewVendorGodotAction --------------------- */

// NewVendorGodotAction creates an 'action.Action' which vendors Godot source
// code into the build directory.
func NewVendorGodotAction(g *Godot, inv *Invocation) action.WithDescription[action.Function] {
	fn := func(ctx context.Context) error {
		if g.IsEmpty() {
			log.Debug("no Godot version set; skipping vendoring of source code")

			return nil
		}

		pathSource, err := filepath.Abs(g.PathSource.String())
		if err != nil {
			return err
		}

		pathBuild, err := filepath.Abs(inv.PathBuild.String())
		if err != nil {
			return err
		}

		if pathSource == pathBuild {
			log.Info("build directory is source directory; skipping vendoring of source code")

			return nil
		}

		if err := osutil.EnsureDir(pathBuild, osutil.ModeUserRWXGroupRX); err != nil {
			return err
		}

		return g.VendorTo(ctx, pathBuild)
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "<go function: vendor godot source code>",
	}
}
