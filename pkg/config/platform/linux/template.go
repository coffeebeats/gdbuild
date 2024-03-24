package linux

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/config/platform/common"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: Template                               */
/* -------------------------------------------------------------------------- */

type Template struct {
	*common.Template

	// UseLLVM determines whether the LLVM compiler is used.
	UseLLVM *bool `toml:"use_llvm"`
}

/* ------------------------- Impl: config.Templater ------------------------- */

func (t *Template) Collect(g engine.Source, rc *run.Context) *template.Template {
	out := t.Template.Collect(g, rc)

	out.Builds[0].Platform = platform.OSLinux

	if t.Arch == platform.ArchUnknown {
		out.Builds[0].Arch = platform.ArchAmd64
	}

	scons := &out.Builds[0].SCons
	if config.Dereference(t.UseLLVM) {
		scons.ExtraArgs = append(scons.ExtraArgs, "use_llvm=yes")
	} else if rc.Profile.IsRelease() { // Only valid with GCC.
		scons.ExtraArgs = append(scons.ExtraArgs, "lto=full")
	}

	return out
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (t *Template) Configure(rc *run.Context) error {
	if err := t.Template.Configure(rc); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (t *Template) Validate(rc *run.Context) error {
	if err := t.Template.Validate(rc); err != nil {
		return err
	}

	if !t.Arch.IsOneOf(platform.ArchI386, platform.ArchAmd64, platform.ArchUnknown) {
		return fmt.Errorf("%w: unsupport architecture: %s", config.ErrInvalidInput, t.Arch)
	}

	switch t.Arch {
	case platform.ArchI386, platform.ArchAmd64:
	case platform.ArchUnknown:
	default:
		return fmt.Errorf("%w: unsupport architecture: %s", config.ErrInvalidInput, t.Arch)
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
