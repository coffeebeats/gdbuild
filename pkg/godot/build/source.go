package build

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdenv/pkg/godot/version"
	"github.com/coffeebeats/gdenv/pkg/install"
	"github.com/coffeebeats/gdenv/pkg/store"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/osutil"
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

/* -------------------------- Method: parseVersion -------------------------- */

// parseVersion determines the 'Godot' version from the underlying 'Source'
// configuration. Returns an 'ErrConflictingValue' when a source path is set
// on the struct instead of a version specification.
func (c *Source) parseVersion() (Version, error) {
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

	v, err := c.parseVersion()
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

/* ---------------------------- config.Configurer --------------------------- */

func (c *Source) Configure(bc *Context) error {
	if err := c.PathSource.RelTo(bc.PathManifest); err != nil {
		return err
	}

	if err := c.VersionFile.RelTo(bc.PathManifest); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *Source) Validate(_ *Context) error {
	if c.IsEmpty() {
		return fmt.Errorf("%w: no Godot version specified in manifest", ErrMissingInput)
	}

	if c.PathSource != "" {
		if !c.Version.IsZero() || c.VersionFile != "" {
			return fmt.Errorf(
				"%w: cannot specify 'version' or 'version_file' with 'src_path'",
				ErrConflictingValue,
			)
		}

		if err := c.PathSource.CheckIsDir(); err != nil {
			// NOTE: A hook might generate this file, so don't return an error.
			if !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		}

		return nil
	}

	if _, err := c.parseVersion(); err != nil {
		// NOTE: A hook might generate this file, so don't return an error.
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		return err
	}

	return nil
}

/* --------------------------- Impl: config.Merger -------------------------- */

func (c *Source) MergeInto(other any) error {
	if c == nil || other == nil {
		return nil
	}

	dst, ok := other.(*Source)
	if !ok {
		return fmt.Errorf(
			"%w: expected a '%T' but was '%T'",
			config.ErrInvalidInput,
			new(Source),
			other,
		)
	}

	return config.Merge(dst, *c)
}

/* -------------------------------------------------------------------------- */
/*                               Struct: Version                              */
/* -------------------------------------------------------------------------- */

// Version is a type wrapper around a 'version.Version' struct from 'gdenv'.
type Version version.Version

/* ----------------------------- Method: IsZero ----------------------------- */

// IsZero returns whether the 'Version' struct is equal to the zero value.
func (v Version) IsZero() bool {
	return v == Version{} //nolint:exhaustruct
}

/* --------------------------- Impl: fmt.Stringer --------------------------- */

func (v Version) String() string {
	return version.Version(v).String()
}

/* ---------------------- Impl: encoding.UnmarshalText ---------------------- */

func (v *Version) UnmarshalText(text []byte) error {
	parsed, err := version.Parse(string(text))
	if err != nil {
		if !errors.Is(err, version.ErrMissing) {
			return err
		}

		*v = Version{} //nolint:exhaustruct

		return nil
	}

	*v = Version(parsed)

	return nil
}

/* -------------------------------------------------------------------------- */
/*                       Function: NewVendorGodotAction                       */
/* -------------------------------------------------------------------------- */

// NewVendorGodotAction creates an 'action.Action' which vendors Godot source
// code into the build directory.
func NewVendorGodotAction(src *Source, cc *Context) action.WithDescription[action.Function] {
	fn := func(ctx context.Context) error {
		if src.IsEmpty() {
			log.Debug("no Godot version set; skipping vendoring of source code")

			return nil
		}

		pathSource, err := filepath.Abs(src.PathSource.String())
		if err != nil {
			return err
		}

		pathBuild, err := filepath.Abs(cc.PathBuild.String())
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

		return src.VendorTo(ctx, pathBuild)
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "vendor godot source code: " + cc.PathBuild.String(),
	}
}
