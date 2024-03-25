package export

import (
	"errors"
	"hash/crc64"
	"io"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/mitchellh/hashstructure/v2"
	"golang.org/x/exp/maps"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
	"github.com/coffeebeats/gdbuild/pkg/run"
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

/* ----------------------------- Method: Actions ---------------------------- */

// Action creates an 'action.Action' for running the export action.
func (x *Export) Action() action.Action { //nolint:ireturn
	return action.NoOp{}
}

/* ---------------------------- Method: Artifacts --------------------------- */

// Artifacts returns the set of exported project artifacts required by the
// underlying target definition.
func (x *Export) Artifacts() []string {
	artifacts := make(map[string]struct{})

	// TODO: Implement this.

	return maps.Keys(artifacts)
}

/* -------------------------------------------------------------------------- */
/*                             Function: Checksum                             */
/* -------------------------------------------------------------------------- */

// Checksum produces a checksum hash of the export specification. When the
// checksums of two 'Export' definitions matches, the resulting exported
// artifacts will be equivalent.
func (x *Export) Checksum(rc *run.Context) (string, error) {
	hash, err := hashstructure.Hash(
		x,
		hashstructure.FormatV2,
		&hashstructure.HashOptions{ //nolint:exhaustruct
			IgnoreZeroValue: true,
			SlicesAsSets:    true,
			ZeroNil:         true,
		},
	)
	if err != nil {
		return "", err
	}

	cs := crc64.New(crc64.MakeTable(crc64.ECMA))

	// Update the 'crc64' hash with the struct hash.
	if _, err := io.Copy(cs, strings.NewReader(strconv.FormatUint(hash, 16))); err != nil {
		return "", err
	}

	files := make([]osutil.Path, 0)
	pathRoot := osutil.Path(filepath.Dir(rc.PathManifest.String()))

	for _, pck := range x.PackFiles {
		ff, err := pck.Files(pathRoot)
		if err != nil {
			return "", err
		}

		files = append(files, ff...)
	}

	// Make the path list unique and sorted.
	slices.Sort(files)
	files = slices.Compact(files)

	for _, path := range files {
		if err := osutil.HashFiles(cs, path.String()); err != nil {
			return "", err
		}
	}

	return strconv.FormatUint(cs.Sum64(), 16), nil
}
