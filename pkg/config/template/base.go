package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/build"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: Base                                */
/* -------------------------------------------------------------------------- */

// Base contains platform-agnostic settings for constructing a custom Godot
// export template.
type Base struct {
	// Arch is the CPU architecture of the Godot export template.
	Arch build.Arch `toml:"arch"`
	// CustomModules is a list of paths to custom modules to include in the
	// template build.
	CustomModules []build.Path `toml:"custom_modules"`
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
	PathCustomPy build.Path `toml:"custom_py_path"`
	// SCons contains build command-related settings.
	SCons build.SCons `toml:"scons"`
}

// Compile-time check that 'Template' is implemented.
var _ Template = (*Base)(nil)

/* -------------------------- Impl: build.Templater ------------------------- */

func (c *Base) ToTemplate(g build.Godot, inv build.Invocation) build.Template {
	scons := c.SCons

	// Append environment-specified arguments.
	scons.ExtraArgs = append(scons.ExtraArgs, scons.ExtraArgsFromEnv()...)

	// Override the cache path using an environment-specified path.
	if pc := scons.PathCacheFromEnv(); pc != "" {
		scons.PathCache = pc
	}

	// Override the cache size limit using an environment-specified path.
	if csl := scons.CacheSizeLimitFromEnv(); csl != nil {
		scons.CacheSizeLimit = csl
	}

	return build.Template{
		Binaries: []build.Binary{
			{
				Arch:            c.Arch,
				CustomModules:   c.CustomModules,
				CustomPy:        c.PathCustomPy,
				DoublePrecision: config.Dereference(c.DoublePrecision),
				Env:             c.Env,
				Godot:           g,
				Optimize:        c.Optimize,
				Platform:        inv.Platform,
				Profile:         inv.Profile,
				SCons:           scons,
			},
		},
		ExtraArtifacts: nil,
		Paths:          nil,
		Prebuild:       c.Hook.PreActions(inv),
		Postbuild:      c.Hook.PostActions(inv),
	}
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (c *Base) Configure(inv build.Invocation) error {
	if err := c.PathCustomPy.RelTo(inv.PathManifest); err != nil {
		return err
	}

	for i, m := range c.CustomModules {
		if err := m.RelTo(inv.PathManifest); err != nil {
			return err
		}

		c.CustomModules[i] = m
	}

	if err := c.SCons.Configure(inv); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *Base) Validate(inv build.Invocation) error {
	if err := inv.Validate(); err != nil {
		return err
	}

	for _, m := range c.CustomModules {
		if err := m.CheckIsDirOrEmpty(); err != nil {
			return err
		}
	}

	if err := c.PathCustomPy.CheckIsFileOrEmpty(); err != nil {
		return err
	}

	if err := c.Hook.Validate(inv); err != nil {
		return err
	}

	if err := c.SCons.Validate(inv); err != nil {
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
