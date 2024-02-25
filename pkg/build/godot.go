package build

import (
	"fmt"
)

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

/* ---------------------------- Impl: Configurer ---------------------------- */

func (c *Godot) Configure(inv *Invocation) error {
	if err := c.PathSource.RelTo(inv.PathManifest); err != nil {
		return err
	}

	if err := c.VersionFile.RelTo(inv.PathManifest); err != nil {
		return err
	}

	return nil
}

/* ---------------------------- Method: Validate ---------------------------- */

func (c *Godot) Validate() error {
	if c.PathSource != "" {
		if c.Version != "" || c.VersionFile != "" {
			return fmt.Errorf(
				"%w: cannot specify 'version' or 'version_file' with 'src_path'",
				ErrConflictingValue,
			)
		}
	}

	if c.Version != "" {
		if c.PathSource != "" || c.VersionFile != "" {
			return fmt.Errorf(
				"%w: cannot specify 'src_path' or 'version_file' with 'version'",
				ErrConflictingValue,
			)
		}
	}

	if c.VersionFile != "" {
		if c.PathSource != "" || c.Version != "" {
			return fmt.Errorf(
				"%w: cannot specify 'src_path' or 'version' with 'version_file'",
				ErrConflictingValue,
			)
		}
	}

	return nil
}
