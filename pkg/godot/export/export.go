package export

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"golang.org/x/exp/maps"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
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
	// Arch is the architecture of the exported artifacts.
	Arch platform.Arch
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
	// PathTemplateArchive is an optional path to a non-cached export template
	// archive containing the export template to use. If specified, this will
	// take priority over 'Template'.
	PathTemplateArchive osutil.Path `hash:"ignore"`
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
func (x *Export) Action(rc *run.Context, pathGodot osutil.Path) (action.Action, error) { //nolint:ireturn
	presets, err := x.Presets(rc)
	if err != nil {
		return nil, err
	}

	exports := make([]action.Action, 0, 3+len(presets)) //nolint:gomnd

	exports = append(
		exports,
		NewWriteExportPresetsAction(rc, x),
		NewRemoveAllAction(rc.PathWorkspace.Join(".godot").String()),
		NewLoadProjectAction(rc, pathGodot),
	)

	for _, preset := range presets {
		exports = append(exports, NewExportAction(rc, preset, pathGodot))
	}

	return action.InOrder(exports...), nil
}

/* ----------------------------- Method: Presets ---------------------------- */

// Presets constructs the list of 'Preset' types for the specified pack files.
func (x *Export) Presets(rc *run.Context) ([]*Preset, error) {
	presets := make([]*Preset, 0, len(x.PackFiles))

	var embed Preset

	for i, pf := range x.PackFiles {
		preset, err := pf.Preset(rc, x, i)
		if err != nil {
			return nil, err
		}

		if !config.Dereference(pf.Embed) {
			presets = append(presets, &preset)

			continue
		}

		if err := config.Merge(&embed, preset); err != nil {
			return nil, err
		}
	}

	if embed.Platform == platform.OSUnknown {
		return presets, nil
	}

	return append([]*Preset{&embed}, presets...), nil
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
/*                          Function: NewExportAction                         */
/* -------------------------------------------------------------------------- */

// NewExportAction creates a new 'action.Action' which exports the specified
// pack file.
func NewExportAction( //nolint:ireturn
	rc *run.Context,
	preset *Preset,
	pathGodot osutil.Path,
) action.Action {
	var cmd action.Process

	cmd.Verbose = rc.Verbose
	cmd.Directory = rc.PathWorkspace.String()

	if preset.EncryptionKey != "" {
		cmd.Environment = append(
			cmd.Environment,
			"GODOT_SCRIPT_ENCRYPTION_KEY="+preset.EncryptionKey,
		)
	}

	cmd.Args = append(
		cmd.Args,
		pathGodot.String(),
		"--headless",
	)

	if rc.Verbose {
		cmd.Args = append(cmd.Args, "--verbose")
	}

	var command string

	switch {
	case !preset.Embed:
		command = "pack"
	case preset.Embed && rc.Profile.IsRelease():
		command = "release"
	default:
		command = "debug"
	}

	pathArtifact := filepath.Join(rc.PathOut.String(), preset.Name)

	cmd.Args = append(
		cmd.Args,
		"--export-"+command,
		preset.Name,
		pathArtifact,
	)

	return NewMkdirAllAction(filepath.Dir(pathArtifact), osutil.ModeUserRWX).
		AndThen(&cmd)
}

/* -------------------------------------------------------------------------- */
/*                         Function: NewMkdirAllAction                        */
/* -------------------------------------------------------------------------- */

// NewMkdirAllAction creates a new 'action.Action' which creates the specified
// directory if missing.
func NewMkdirAllAction(path string, perm os.FileMode) action.WithDescription[action.Function] {
	fn := func(_ context.Context) error {
		return os.MkdirAll(path, perm)
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "create directory if missing: " + path,
	}
}
