package common

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                              Struct: Template                              */
/* -------------------------------------------------------------------------- */

// Template contains platform-agnostic settings for constructing a custom Godot
// export template.
type Template struct {
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
	Hook run.Hook `toml:"hook"`
	// Optimize is the specific optimization level for the template.
	Optimize engine.Optimize `toml:"optimize"`
	// PathCustomPy is a path to a 'custom.py' file which defines export
	// template build options.
	PathCustomPy osutil.Path `toml:"custom_py_path"`
	// SCons contains build command-related settings.
	SCons template.SCons `toml:"scons"`
}

/* ----------------------------- Impl: Template ----------------------------- */

func (t Template) Collect(src engine.Source, rc *run.Context) *template.Template {
	s := t.SCons

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

	return &template.Template{
		Arch: t.Arch,
		Builds: []template.Build{{
			Arch:            t.Arch,
			CustomModules:   t.CustomModules,
			CustomPy:        t.PathCustomPy,
			DoublePrecision: config.Dereference(t.DoublePrecision),
			EncryptionKey:   template.EncryptionKeyFromEnv(),
			Env:             t.Env,
			Source:          src,
			Optimize:        t.Optimize,
			Platform:        rc.Platform,
			Profile:         rc.Profile,
			SCons:           s,
		}},
		ExtraArtifacts: nil,
		NameOverride:   "",
		Paths:          nil,
		Prebuild:       t.Hook.PreActions(rc),
		Postbuild:      t.Hook.PostActions(rc),
	}
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (t *Template) Configure(rc *run.Context) error {
	if err := t.PathCustomPy.RelTo(rc.PathManifest); err != nil {
		return err
	}

	for i, m := range t.CustomModules {
		if err := m.RelTo(rc.PathManifest); err != nil {
			return err
		}

		t.CustomModules[i] = m
	}

	if err := t.SCons.Configure(rc); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (t *Template) Validate(rc *run.Context) error {
	for _, m := range t.CustomModules {
		if err := m.CheckIsDirOrEmpty(); err != nil {
			return err
		}
	}

	if err := t.PathCustomPy.CheckIsFileOrEmpty(); err != nil {
		return err
	}

	if err := t.Hook.Validate(rc); err != nil {
		return err
	}

	if err := t.SCons.Validate(rc); err != nil {
		return err
	}

	return nil
}

/* --------------------------- Impl: config.Merger -------------------------- */

func (t *Template) MergeInto(other any) error {
	if t == nil || other == nil {
		return nil
	}

	dst, ok := other.(*Template)
	if !ok {
		return fmt.Errorf(
			"%w: expected a '%T' but was '%T'",
			config.ErrInvalidInput,
			new(Template),
			other,
		)
	}

	return config.Merge(dst, *t)
}

/* -------------------------------------------------------------------------- */
/*                   Struct: TemplateWithFeaturesAndProfile                   */
/* -------------------------------------------------------------------------- */

type TemplateWithFeaturesAndProfile struct {
	*Template

	Feature map[string]TemplateWithProfile `toml:"feature"`
	Profile map[engine.Profile]Template    `toml:"profile"`
}

/* ----------------------- Struct: TemplateWithProfile ---------------------- */

type TemplateWithProfile struct {
	*Template

	Profile map[engine.Profile]Template `toml:"profile"`
}

/* --------------------- Impl: platform.templateBuilder --------------------- */

func (t *TemplateWithFeaturesAndProfile) Build(rc *run.Context, dst *Template) error {
	if t == nil {
		return nil
	}

	// Root-level params
	if err := t.Template.MergeInto(dst); err != nil {
		return err
	}

	// Feature-constrained params
	for _, f := range rc.Features {
		if err := t.Feature[f].Template.MergeInto(dst); err != nil {
			return err
		}
	}

	// Profile-constrained params
	l := t.Profile[rc.Profile]
	if err := l.MergeInto(dst); err != nil {
		return err
	}

	// Feature-and-profile-constrained params
	for _, f := range rc.Features {
		l := t.Feature[f].Profile[rc.Profile]
		if err := l.MergeInto(dst); err != nil {
			return err
		}
	}

	return nil
}
