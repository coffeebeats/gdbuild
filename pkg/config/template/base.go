package template

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/build"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: Base                                */
/* -------------------------------------------------------------------------- */

// Base contains platform-agnostic settings for constructing a custom Godot
// export template.
type Base struct {
	// Arch is the CPU architecture of the Godot export template.
	Arch platform.Arch `toml:"arch"`
	// CustomModules is a list of paths to custom modules to include in the
	// template build.
	CustomModules []osutil.Path `toml:"custom_modules"`
	// DoublePrecision enables double floating-point precision.
	DoublePrecision *bool `toml:"double_precision"`
	// EncryptionKey is the encryption key to embed in the export template.
	EncryptionKey string `toml:"encryption_key"`
	// Env is a map of environment variables to set during the build step.
	Env map[string]string `toml:"env"`
	// Hook defines commands to be run before or after a build step.
	Hook run.Hook `toml:"hook"`
	// Optimize is the specific optimization level for the template.
	Optimize engine.Optimize `toml:"optimize"`
	// PathCustomPy is a path to a 'custom.py' file which defines export
	// template build options.
	PathCustomPy osutil.Path `toml:"custom_py_path"`
	// SCons contains build command-related settings.
	SCons build.SCons `toml:"scons"`
}

// Compile-time check that 'Template' is implemented.
var _ Template = (*Base)(nil)

/* ----------------------------- Impl: Template ----------------------------- */

func (c *Base) Template(src engine.Source, rc *run.Context) *template.Template {
	s := c.SCons

	// Append environment-specified arguments.
	s.ExtraArgs = append(s.ExtraArgs, s.ExtraArgsFromEnv()...)

	// Override the cache path using an environment-specified path.
	if pc := s.PathCacheFromEnv(); pc != "" {
		s.PathCache = pc
	}

	// Override the cache size limit using an environment-specified path.
	if csl := s.CacheSizeLimitFromEnv(); csl != nil {
		s.CacheSizeLimit = csl
	}

	// Set the encryption key environment variable; see
	// https://docs.godotengine.org/en/stable/contributing/development/compiling/compiling_with_script_encryption_key.html.
	var encryptionKey string
	if ek := build.EncryptionKeyFromEnv(); ek != "" {
		encryptionKey = ek
	} else if c.EncryptionKey != "" {
		ek := os.ExpandEnv(c.EncryptionKey)
		if ek != "" {
			encryptionKey = ek
		} else {
			log.Warnf("encryption key set in manifest, but value was empty: %s", c.EncryptionKey)
		}
	}

	return &template.Template{
		Builds: []build.Build{
			{
				Arch:            c.Arch,
				CustomModules:   c.CustomModules,
				CustomPy:        c.PathCustomPy,
				DoublePrecision: config.Dereference(c.DoublePrecision),
				EncryptionKey:   encryptionKey,
				Env:             c.Env,
				Source:          src,
				Optimize:        c.Optimize,
				Platform:        rc.Platform,
				Profile:         rc.Profile,
				SCons:           s,
			},
		},
		ExtraArtifacts: nil,
		Paths:          nil,
		Prebuild:       c.Hook.PreActions(rc),
		Postbuild:      c.Hook.PostActions(rc),
	}
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (c *Base) Configure(rc *run.Context) error {
	if err := c.PathCustomPy.RelTo(rc.PathManifest); err != nil {
		return err
	}

	for i, m := range c.CustomModules {
		if err := m.RelTo(rc.PathManifest); err != nil {
			return err
		}

		c.CustomModules[i] = m
	}

	if err := c.SCons.Configure(rc); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *Base) Validate(rc *run.Context) error {
	for _, m := range c.CustomModules {
		if err := m.CheckIsDirOrEmpty(); err != nil {
			return err
		}
	}

	if err := c.PathCustomPy.CheckIsFileOrEmpty(); err != nil {
		return err
	}

	if err := c.Hook.Validate(rc); err != nil {
		return err
	}

	if err := c.SCons.Validate(rc); err != nil {
		return err
	}

	return nil
}

/* --------------------------- Impl: config.Merger -------------------------- */

func (c *Base) MergeInto(other any) error {
	if c == nil || other == nil {
		return nil
	}

	dst, ok := other.(*Base)
	if !ok {
		return fmt.Errorf(
			"%w: expected a '%T' but was '%T'",
			config.ErrInvalidInput,
			new(Base),
			other,
		)
	}

	return config.Merge(dst, *c)
}
