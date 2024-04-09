package macos

import (
	"errors"
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/config/common"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

var ErrMissingInput = errors.New("missing input")

/* -------------------------------------------------------------------------- */
/*                               Struct: Target                               */
/* -------------------------------------------------------------------------- */

type Target struct {
	*common.Target

	// PathIcon is a path to a Windows application icon.
	BundleIdentifier string `toml:"bundle_identifier"`
}

/* ----------------------------- Impl: Exporter ----------------------------- */

func (t *Target) Collect(rc *run.Context, tl *template.Template, ev engine.Version) *export.Export {
	out := t.Target.Collect(rc, tl, ev)

	if out.Options == nil {
		out.Options = map[string]any{}
	}

	out.Options["application/bundle_identifier"] = t.BundleIdentifier

	return out
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (t *Target) Configure(rc *run.Context) error {
	return t.Target.Configure(rc)
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (t *Target) Validate(rc *run.Context) error {
	if err := t.Target.Validate(rc); err != nil {
		return err
	}

	if t.BundleIdentifier == "" {
		return fmt.Errorf("%w: 'bundle_identifier'", ErrMissingInput)
	}

	return nil
}

/* --------------------------- Impl: config.Merger -------------------------- */

func (t *Target) MergeInto(other any) error {
	if t == nil || other == nil {
		return nil
	}

	dst, ok := other.(*Target)
	if !ok {
		return fmt.Errorf(
			"%w: expected a '%T' but was '%T'",
			config.ErrInvalidInput,
			new(Target),
			other,
		)
	}

	return config.Merge(dst, *t)
}

/* -------------------------------------------------------------------------- */
/*                    Struct: TargetWithFeaturesAndProfile                    */
/* -------------------------------------------------------------------------- */

type TargetWithFeaturesAndProfile struct {
	*Target

	Feature map[string]TargetWithProfile `toml:"feature"`
	Profile map[engine.Profile]Target    `toml:"profile"`
}

/* ------------------------ Struct: TargetWithProfile ----------------------- */

type TargetWithProfile struct {
	*Target

	Profile map[engine.Profile]Target `toml:"profile"`
}

/* ---------------------- Impl: platform.targetBuilder ---------------------- */

func (t *TargetWithFeaturesAndProfile) Build(rc *run.Context, dst *Target) error {
	if t == nil {
		return nil
	}

	// Root-level params
	if err := t.Target.MergeInto(dst); err != nil {
		return err
	}

	// Feature-constrained params
	for _, f := range rc.Features {
		if err := t.Feature[f].Target.MergeInto(dst); err != nil {
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
