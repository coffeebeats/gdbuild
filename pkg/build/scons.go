package build

import (
	"errors"
	"os"
)

const (
	EnvSConsCache = "SCONS_CACHE"
	EnvSConsFlags = "SCONSFLAGS"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: SCons                               */
/* -------------------------------------------------------------------------- */

// SCons defines options and settings for use with the Godot build system.
type SCons struct {
	// CCFlags are additional 'CFLAGS' to append to the SCons build command.
	// Note that 'CCFLAGS=...' will be appended *before* 'ExtraArgs'.
	CCFlags []string `toml:"ccflags"`
	// CFlags are additional 'CFLAGS' to append to the SCons build command. Note
	// that 'CFLAGS=...' will be appended *before* 'ExtraArgs'.
	CFlags []string `toml:"cflags"`
	// CXXFlags are additional 'CXXFLAGS' to append to the SCons build command.
	// Note that 'CXXFLAGS=...' will be appended *before* 'ExtraArgs'.
	CXXFlags []string `toml:"cxxflags"`
	// Command contains arguments used to invoke SCons. Defaults to ["scons"].
	Command []string `toml:"command"`
	// ExtraArgs are additional arguments to append to the SCons build command.
	ExtraArgs []string `toml:"extra_args"`
	// LinkFlags are additional flags passed to the linker during the SCons
	// build command.
	LinkFlags []string `toml:"link_flags"`
	// PathCache is the path to the SCons cache, relative to the manifest.
	PathCache Path `toml:"cache_path"`
	// CacheSizeLimit is the limit in MiB.
	CacheSizeLimit *uint `toml:"cache_size_limit"`
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (c *SCons) Configure(inv *Invocation) error {
	if len(c.Command) == 0 {
		c.Command = append(c.Command, "scons")
	}

	if p := os.Getenv(EnvSConsCache); p != "" {
		c.PathCache = Path(p)
	}

	if err := c.PathCache.RelTo(inv.PathManifest); err != nil {
		return err
	}

	return nil
}

/* -------------------------- Impl: build.Validater ------------------------- */

func (c *SCons) Validate() error {
	if err := c.PathCache.CheckIsDirOrEmpty(); err != nil {
		// A missing SCons cache is not a problem.
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	// NOTE: Don't check for 'scons' command, that should be a runtime check.

	return nil
}
