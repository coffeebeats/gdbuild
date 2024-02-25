package build

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/merge"
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
	// ExtraArgs are additional arguments to append to the SCons build command.
	ExtraArgs []string `toml:"extra_args"`
	// LinkFlags are additional flags passed to the linker during the SCons
	// build command.
	LinkFlags []string `toml:"link_flags"`
	// PathCache is the path to the SCons cache, relative to the manifest.
	PathCache Path `toml:"cache_path"`
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *SCons) Merge(other *SCons) error {
	if c == nil || other == nil {
		return nil
	}

	c.CCFlags = append(c.CCFlags, other.CCFlags...)
	c.CFlags = append(c.CFlags, other.CFlags...)
	c.CXXFlags = append(c.CXXFlags, other.CXXFlags...)
	c.ExtraArgs = append(c.ExtraArgs, other.ExtraArgs...)

	if err := merge.Primitive(&c.PathCache, other.PathCache); err != nil {
		return fmt.Errorf("%w: cache_path", err)
	}

	return nil
}
