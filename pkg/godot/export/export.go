package export

import (
	"errors"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
)

var (
	ErrConflictingValue = errors.New("conflicting setting")
	ErrInvalidInput     = errors.New("invalid input")
	ErrMissingInput     = errors.New("missing input")
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Export                               */
/* -------------------------------------------------------------------------- */

// Export specifies a single, platform-agnostic exportable artifact within the
// Godot project.
type Export struct {
	// EncryptionKey is an encryption key used to encrypt game assets with.
	EncryptionKey string
	// Features contains the slice of Godot project feature tags to build with.
	Features []string
	// Options are 'export_presets.cfg' overrides, specifically the preset
	// 'options' table, for the exported artifact.
	Options map[string]any
	// PackFiles defines the game files exported as part of this artifact.
	PackFiles []PackFile
	// Template specifies the export template to use.
	Template *template.Template `hash:"string"`
	// RunBefore contains an ordered list of actions to execute prior to
	// exporting the target.
	RunBefore action.Action `hash:"string"`
	// RunAfter contains an ordered list of actions to execute after exporting
	// the target.
	RunAfter action.Action `hash:"string"`
	// Runnable is whether the export artifact should be executable. This should
	// be true for client and server targets and false for artifacts like DLC.
	Runnable bool
	// Server configures the target as a server-only executable, enabling some
	// optimizations like disabling graphics.
	Server bool
	// Version is the editor version to use for exporting.
	Version engine.Version
}
