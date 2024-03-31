package export

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/config"
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
	// PathTemplate is a path to the export template to use during exporting.
	PathTemplate osutil.Path `hash:"ignore"`
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
func (x *Export) Action(rc *run.Context) (action.Action, error) { //nolint:ireturn
	var out action.Action = action.NoOp{}

	out = out.AndThen(NewWriteExportPresetsAction(rc, x))

	presets, err := x.Presets(rc)
	if err != nil {
		return nil, err
	}

	for _, preset := range presets {
		preset := preset
		out = out.AndThen(NewExportAction(rc, &preset))
	}

	return out, nil
}

/* ----------------------------- Method: Presets ---------------------------- */

// Presets constructs the list of 'Preset' types for the specified pack files.
func (x *Export) Presets(rc *run.Context) ([]Preset, error) {
	presets := make([]Preset, 0, len(x.PackFiles))

	var embed Preset

	for i, pf := range x.PackFiles {
		preset, err := pf.Preset(rc, x, i)
		if err != nil {
			return nil, err
		}

		if !config.Dereference(pf.Embed) {
			presets = append(presets, preset)

			continue
		}

		if err := config.Merge(&embed, preset); err != nil {
			return nil, err
		}
	}

	return append([]Preset{embed}, presets...), nil
}

/* ---------------------------- Method: Artifacts --------------------------- */

// Artifacts returns the set of exported project artifacts required by the
// underlying target definition.
func (x *Export) Artifacts(rc *run.Context) ([]string, error) {
	artifacts := make(map[string]struct{})

	presets, err := x.Presets(rc)
	if err != nil {
		return nil, err
	}

	for _, preset := range presets {
		artifacts[preset.Name] = struct{}{}
	}

	return maps.Keys(artifacts), nil
}

/* -------------------------------------------------------------------------- */
/*                    Function: NewWriteExportPresetsAction                   */
/* -------------------------------------------------------------------------- */

// NewWriteExportPresetsAction creates a new 'action.Action' which constructs an
// 'export_presets.cfg' file based on the target. It will be written to the
// workspace directory and overwrite any existing files.
func NewWriteExportPresetsAction(
	rc *run.Context,
	x *Export,
) action.WithDescription[action.Function] {
	path := filepath.Join(rc.PathWorkspace.String(), "export_presets.cfg")

	fn := func(_ context.Context) error {
		presets, err := x.Presets(rc)
		if err != nil {
			return err
		}

		var cfg strings.Builder

		for i, preset := range presets {
			if err := preset.Marshal(&cfg, i); err != nil {
				return err
			}
		}

		f, err := os.Create(filepath.Join(rc.PathWorkspace.String(), "export_presets.cfg"))
		if err != nil {
			return err
		}

		defer f.Close()

		if _, err := io.Copy(f, strings.NewReader(cfg.String()+"\n")); err != nil {
			return err
		}

		return nil
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "generate export presets file: " + path,
	}
}

/* -------------------------------------------------------------------------- */
/*                          Function: NewExportAction                         */
/* -------------------------------------------------------------------------- */

// NewExportAction creates a new 'action.Action' which exports the specified
// pack file.
func NewExportAction(
	rc *run.Context,
	preset *Preset,
) *action.Process {
	var cmd action.Process

	cmd.Verbose = rc.Verbose

	cmd.Directory = rc.PathWorkspace.String()
	cmd.Environment = os.Environ()

	cmd.Args = append(
		cmd.Args,
		"godot",
		"--headless",
	)

	if rc.Verbose {
		cmd.Args = append(cmd.Args, "--verbose")
	}

	profile := "release"
	if rc.Profile == engine.ProfileDebug {
		profile = "debug"
	}

	cmd.Args = append(
		cmd.Args,
		"--export-"+profile,
		preset.Name,
		filepath.Join(rc.PathOut.String(), preset.PathExport),
	)

	return &cmd
}
