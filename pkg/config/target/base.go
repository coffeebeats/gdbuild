package target

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: Base                                */
/* -------------------------------------------------------------------------- */

// Base specifies a single, platform-agnostic exportable artifact within the
// Godot project.
type Base struct {
	// Name is the display name of the target. Not used by Godot.
	Name string

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

func (b *Base) Export(_ engine.Source, _ *run.Context) *export.Export {
	return nil
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (b *Base) Configure(_ *run.Context) error {
	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (b *Base) Validate(_ *run.Context) error {
	return nil
}

/* --------------------------- Impl: config.Merger -------------------------- */

func (b *Base) MergeInto(other any) error {
	if b == nil || other == nil {
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

	return config.Merge(dst, *b)
}
