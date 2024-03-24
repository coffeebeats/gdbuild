package engine

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdenv/pkg/godot/version"
	"github.com/coffeebeats/gdenv/pkg/install"
	"github.com/coffeebeats/gdenv/pkg/store"

	"github.com/coffeebeats/gdbuild/internal/osutil"
)

var (
	ErrConflictingValue = errors.New("conflicting setting")
	ErrInvalidInput     = errors.New("invalid input")
	ErrMissingInput     = errors.New("missing input")
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Source                               */
/* -------------------------------------------------------------------------- */

// Source defines options and settings for which Godot version to use. Note that
// only one of these options can be used at a time, but one *must* be specified.
type Source struct {
	// PathSource is a path to a directory containing the Godot source code.
	PathSource osutil.Path `hash:"ignore" toml:"src_path"`
	// Version is a specific version label to download.
	Version Version `toml:"version"`
	// VersionFile is a file containing just the a version label to download.
	VersionFile osutil.Path `hash:"ignore" toml:"version_file"`
}

/* ----------------------------- Method: IsEmpty ---------------------------- */

// IsEmpty returns whether all properties are unset, implying there is no need
// to vendor Godot source code.
func (c *Source) IsEmpty() bool {
	return c.PathSource == "" && c.Version.IsZero() && c.VersionFile == ""
}

/* -------------------------- Method: ParseVersion -------------------------- */

// ParseVersion determines the 'Godot' version from the underlying 'Source'
// configuration. Returns an 'ErrConflictingValue' when a source path is set
// on the struct instead of a version specification.
func (c *Source) ParseVersion() (Version, error) {
	if c.PathSource != "" {
		return Version{}, ErrConflictingValue
	}

	if !c.Version.IsZero() {
		if c.VersionFile != "" {
			return Version{}, fmt.Errorf(
				"%w: cannot specify 'src_path' or 'version_file' with 'version'",
				ErrConflictingValue,
			)
		}

		return c.Version, nil
	}

	bb, err := os.ReadFile(c.VersionFile.String())
	if err != nil {
		return Version{}, err
	}

	v, err := version.Parse(string(bb))
	if err != nil {
		return Version{}, fmt.Errorf("%w: %w: %s", ErrInvalidInput, err, bb)
	}

	return Version(v), nil
}

/* ---------------------------- Method: VendorTo ---------------------------- */

// VendorTo vendors the Godot source code to the specified directory.
func (c *Source) VendorTo(ctx context.Context, out string) error {
	if c.IsEmpty() {
		return fmt.Errorf("%w: no Godot version or source path set", ErrMissingInput)
	}

	if c.PathSource != "" {
		return osutil.CopyDir(c.PathSource.String(), out)
	}

	v, err := c.ParseVersion()
	if err != nil {
		return err
	}

	storePath, err := store.Path() // Use the 'gdenv' store.
	if err != nil {
		tmp, err := os.MkdirTemp("", "gdbuild-*")
		if err != nil {
			return err
		}

		log.Debugf("no 'gdenv' store found; using temporary directory: %s", tmp)

		storePath = tmp
	}

	return install.Vendor(
		ctx,
		storePath,
		version.Version(v),
		out,
		/* force= */ false,
	)
}
