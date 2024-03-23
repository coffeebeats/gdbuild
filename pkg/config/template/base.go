package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/build"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
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
	// Env is a map of environment variables to set during the build step.
	Env map[string]string `toml:"env"`
	// Hook defines commands to be run before or after a build step.
	Hook build.Hook `toml:"hook"`
	// Optimize is the specific optimization level for the template.
	Optimize build.Optimize `toml:"optimize"`
	// PathCustomPy is a path to a 'custom.py' file which defines export
	// template build options.
	PathCustomPy osutil.Path `toml:"custom_py_path"`
	// SCons contains build command-related settings.
	SCons build.SCons `toml:"scons"`
}

// Compile-time check that 'Template' is implemented.
var _ Template = (*Base)(nil)

/* ----------------------------- Impl: Template ----------------------------- */

func (c *Base) Template(src build.Source, bc *build.Context) *build.Template {
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

	return &build.Template{
		Builds: []build.Build{
			{
				Arch:            c.Arch,
				CustomModules:   c.CustomModules,
				CustomPy:        c.PathCustomPy,
				DoublePrecision: config.Dereference(c.DoublePrecision),
				Env:             c.Env,
				Source:          src,
				Optimize:        c.Optimize,
				Platform:        bc.Platform,
				Profile:         bc.Profile,
				SCons:           s,
			},
		},
		ExtraArtifacts: nil,
		Paths:          nil,
		Prebuild:       c.Hook.PreActions(bc),
		Postbuild:      c.Hook.PostActions(bc),
	}
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (c *Base) Configure(bc *build.Context) error {
	if err := c.PathCustomPy.RelTo(bc.PathManifest); err != nil {
		return err
	}

	for i, m := range c.CustomModules {
		if err := m.RelTo(bc.PathManifest); err != nil {
			return err
		}

		c.CustomModules[i] = m
	}

	if err := c.SCons.Configure(bc); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *Base) Validate(bc *build.Context) error {
	for _, m := range c.CustomModules {
		if err := m.CheckIsDirOrEmpty(); err != nil {
			return err
		}
	}

	if err := c.PathCustomPy.CheckIsFileOrEmpty(); err != nil {
		return err
	}

	if err := c.Hook.Validate(bc); err != nil {
		return err
	}

	if err := c.SCons.Validate(bc); err != nil {
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
