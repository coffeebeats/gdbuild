package build

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdbuild/internal/osutil"
)

const (
	envSConsCache          = "SCONS_CACHE"
	envSConsCacheSizeLimit = "SCONS_CACHE_LIMIT"
	envSConsFlags          = "SCONSFLAGS"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: SCons                               */
/* -------------------------------------------------------------------------- */

// SCons defines options and settings for use with the Godot build system.
type SCons struct {
	// CCFlags are additional 'CFLAGS' to append to the SCons build command.
	// Note that 'CCFLAGS=...' will be appended *before* 'ExtraArgs'.
	CCFlags []string `hash:"set" toml:"ccflags"`
	// CFlags are additional 'CFLAGS' to append to the SCons build command. Note
	// that 'CFLAGS=...' will be appended *before* 'ExtraArgs'.
	CFlags []string `hash:"set" toml:"cflags"`
	// CXXFlags are additional 'CXXFLAGS' to append to the SCons build command.
	// Note that 'CXXFLAGS=...' will be appended *before* 'ExtraArgs'.
	CXXFlags []string `hash:"set" toml:"cxxflags"`
	// CacheSizeLimit is the limit in MiB.
	CacheSizeLimit *uint32 `hash:"ignore" toml:"cache_size_limit"` // Ignore; doesn't affect binary.
	// Command contains arguments used to invoke SCons. Defaults to ["scons"].
	Command []string `hash:"set" toml:"command"`
	// ExtraArgs are additional arguments to append to the SCons build command.
	ExtraArgs []string `hash:"set" toml:"extra_args"`
	// LinkFlags are additional flags passed to the linker during the SCons
	// build command.
	LinkFlags []string `hash:"set" toml:"link_flags"`
	// PathCache is the path to the SCons cache, relative to the manifest.
	PathCache osutil.Path `hash:"ignore" toml:"cache_path"` // Ignore; doesn't affect binary.
}

/* ---------------------- Method: CacheSizeLimitFromEnv --------------------- */

// CacheSizeLimitFromEnv returns a SCons cache size limit set via environment
// variable.
func (c *SCons) CacheSizeLimitFromEnv() *uint32 {
	cslRaw := os.Getenv(envSConsCacheSizeLimit)
	if cslRaw == "" {
		return nil
	}

	csl, err := strconv.ParseInt(cslRaw, 10, 32)
	if err != nil {
		log.Warnf(
			"found invalid environment variable '%s' (expected u32): %s",
			envSConsCacheSizeLimit,
			cslRaw,
		)

		return nil
	}

	cslU32 := uint32(csl)

	return &cslU32
}

/* ------------------------ Method: ExtraArgsFromEnv ------------------------ */

// ExtraArgsFromEnv returns extra SCons arguments set via environment variable.
func (c *SCons) ExtraArgsFromEnv() []string {
	argsRaw := os.Getenv(envSConsFlags)
	if argsRaw == "" {
		return nil
	}

	return strings.Split(argsRaw, " ")
}

/* ------------------------ Method: PathCacheFromEnv ------------------------ */

// PathCacheFromEnv returns a SCons cache path set via environment variable.
func (c *SCons) PathCacheFromEnv() osutil.Path {
	return osutil.Path(os.Getenv(envSConsCache))
}

/* ---------------------------- config.Configurer --------------------------- */

func (c *SCons) Configure(bc *Context) error {
	if p := os.Getenv(envSConsCache); p != "" {
		c.PathCache = osutil.Path(p)
	}

	if err := c.PathCache.RelTo(bc.PathManifest); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *SCons) Validate(_ *Context) error {
	if err := c.PathCache.CheckIsDirOrEmpty(); err != nil {
		// A missing SCons cache is not a problem.
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	// NOTE: Don't check for 'scons' command, that should be a runtime check.

	return nil
}
