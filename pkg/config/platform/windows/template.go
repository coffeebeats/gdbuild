package windows

import (
	"fmt"
	"strings"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/config/platform/common"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                              Struct: Template                              */
/* -------------------------------------------------------------------------- */

type Template struct {
	*common.Template

	// UseMinGW determines whether the MinGW compiler is used.
	UseMinGW *bool `toml:"use_mingw"`

	// PathIcon is a path to a Windows application icon.
	PathIcon osutil.Path `toml:"icon_path"`
}

/* ----------------------------- Impl: Template ----------------------------- */

func (t *Template) Collect(g engine.Source, rc *run.Context) *template.Template {
	out := t.Template.Collect(g, rc)

	out.Builds[0].Platform = platform.OSWindows

	if t.Arch == platform.ArchUnknown {
		out.Builds[0].Arch = platform.ArchAmd64
	}

	scons := &out.Builds[0].SCons
	if rc.Profile.IsRelease() {
		scons.ExtraArgs = append(scons.ExtraArgs, "lto=full")
	}

	if config.Dereference(t.UseMinGW) {
		scons.ExtraArgs = append(scons.ExtraArgs, "use_mingw=yes")
	}

	if t.PathIcon != "" {
		out.RegisterDependencyPath(t.PathIcon)

		// Copy the icon file to the correct location.
		out.Prebuild = action.InOrder(out.Prebuild, NewCopyImageFileAction(t.PathIcon, rc))
	}

	// Register the additional console artifact.
	out.ExtraArtifacts = append(
		out.ExtraArtifacts,
		strings.TrimSuffix(out.Builds[0].Filename(), ".exe")+".console.exe",
	)

	return out
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (t *Template) Configure(rc *run.Context) error {
	if err := t.Template.Configure(rc); err != nil {
		return err
	}

	if err := t.PathIcon.RelTo(rc.PathManifest); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (t *Template) Validate(rc *run.Context) error {
	if err := t.Template.Validate(rc); err != nil {
		return err
	}

	if !t.Arch.IsOneOf(platform.ArchAmd64, platform.ArchI386, platform.ArchUnknown) {
		return fmt.Errorf("%w: unsupport architecture: %s", config.ErrInvalidInput, t.Arch)
	}

	if err := t.PathIcon.CheckIsFileOrEmpty(); err != nil {
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
