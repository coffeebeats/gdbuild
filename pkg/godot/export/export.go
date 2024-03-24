package export

import (
	"errors"

	"github.com/coffeebeats/gdbuild/pkg/godot/build"
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
	// Hook defines commands to be run before or after the target artifact is
	// generated.
	Hook build.Hook
	// Options are 'export_presets.cfg' overrides, specifically the preset
	// 'options' table, for the exported artifact.
	Options map[string]any
	// PackFiles defines the game files exported as part of this artifact.
	PackFiles []PackFile
	// Runnable is whether the export artifact should be executable. This should
	// be true for client and server targets and false for artifacts like DLC.
	Runnable bool
	// Server configures the target as a server-only executable, enabling some
	// optimizations like disabling graphics.
	Server bool
}
