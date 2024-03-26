package common

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
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

func (t *Target) Collect(rc *run.Context, tl *template.Template, ev engine.Version) *export.Export {
	// Set the encryption key environment variable; see
	// https://docs.godotengine.org/en/stable/contributing/development/compiling/compiling_with_script_encryption_key.html.
	var encryptionKey string
	if ek := template.EncryptionKeyFromEnv(); ek != "" {
		encryptionKey = ek
	} else if t.EncryptionKey != "" {
		ek := os.ExpandEnv(t.EncryptionKey)
		if ek != "" {
			encryptionKey = ek
		} else {
			log.Warnf(
				"encryption key set in manifest, but value was empty: %s",
				t.EncryptionKey,
			)
		}
	}

	ff := make([]string, 0, len(t.DefaultFeatures)+len(rc.Features))
	ff = append(ff, t.DefaultFeatures...)
	ff = append(ff, rc.Features...)

	return &export.Export{
		EncryptionKey: encryptionKey,
		Features:      ff,
		Options:       t.Options,
		PackFiles:     t.PackFiles,
		RunBefore:     t.Hook.PreActions(rc),
		RunAfter:      t.Hook.PostActions(rc),
		Runnable:      config.Dereference(t.Runnable),
		Server:        config.Dereference(t.Server),
		Template:      tl,
		Version:       ev,
	}
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (t *Target) Configure(rc *run.Context) error {
	for _, pf := range t.PackFiles {
		if err := pf.Configure(rc); err != nil {
			return err
		}
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (t *Target) Validate(rc *run.Context) error {
	if err := t.Hook.Validate(rc); err != nil {
		return err
	}

	for _, pf := range t.PackFiles {
		if err := pf.Validate(rc); err != nil {
			return err
		}
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
