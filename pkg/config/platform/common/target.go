package common

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Target                               */
/* -------------------------------------------------------------------------- */

// Target specifies a single, platform-agnostic exportable artifact within the
// Godot project.
type Target struct {
	// DefaultFeatures contains the slice of Godot project feature tags to build
	// with.
	DefaultFeatures []string `toml:"default_features"`
	// EncryptionKey is the encryption key to encrypt game assets with.
	EncryptionKey string `toml:"encryption_key"`
	// Hook defines commands to be run before or after the target artifact is
	// generated.
	Hook run.Hook `toml:"hook"`
	// Options are 'export_presets.cfg' overrides, specifically the preset
	// 'options' table, for the exported artifact.
	Options map[string]any `toml:"options"`
	// PackFiles defines the game files exported as part of this artifact.
	PackFiles []export.PackFile `toml:"pack_files"`
	// Runnable is whether the export artifact should be executable. This should
	// be true for client and server targets and false for artifacts like DLC.
	Runnable *bool `toml:"runnable"`
	// Server configures the target as a server-only executable, enabling some
	// optimizations like disabling graphics.
	Server *bool `toml:"server"`
}

/* ----------------------------- Impl: Exporter ----------------------------- */

func (t *Target) Collect(_ engine.Source, _ *run.Context) *export.Export {
	return nil
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (t *Target) Configure(_ *run.Context) error {
	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (t *Target) Validate(_ *run.Context) error {
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
