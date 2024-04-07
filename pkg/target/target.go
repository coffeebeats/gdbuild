package target

import (
	"context"
	"errors"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/archive"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/run"
	"github.com/coffeebeats/gdbuild/pkg/store"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrMissingInput = errors.New("missing input")
)

/* -------------------------------------------------------------------------- */
/*                              Function: Action                              */
/* -------------------------------------------------------------------------- */

// Action creates a new 'action.Action' which executes the specified processes
// for compiling the export template.
func Action(rc *run.Context, xp *export.Export) (action.Action, error) { //nolint:ireturn
	exportAction, err := xp.Action(rc, rc.GodotPath())
	if err != nil {
		return nil, err
	}

	artifacts, err := xp.Artifacts(rc)
	if err != nil {
		return nil, err
	}

	cs, err := export.Checksum(rc, xp)
	if err != nil {
		return nil, err
	}

	cacheArtifacts, err := store.NewCacheTargetAction(rc, rc.PathOut, artifacts, cs)
	if err != nil {
		return nil, err
	}

	return action.InOrder(
		export.NewInstallEditorGodotAction(rc, xp.Version, rc.GodotPath()),
		xp.RunBefore,
		exportAction,
		xp.RunAfter,
		run.NewVerifyArtifactsAction(rc, rc.PathOut, artifacts),
		cacheArtifacts,
	), nil
}

/* -------------------------------------------------------------------------- */
/*                     Function: NewExtractTemplateAction                     */
/* -------------------------------------------------------------------------- */

// NewExtractTemplateAction creates an 'action.Action' which extracts the cached
// Godot export template into a temporary directory and populates the provided
// string variable 'path' with a path to it.
func NewExtractTemplateAction(
	rc *run.Context,
	pathArchive osutil.Path,
) (action.WithDescription[action.Function], error) {
	pathTmp, err := rc.TempDir()
	if err != nil {
		return action.WithDescription[action.Function]{}, err
	}

	fn := func(ctx context.Context) error {
		return archive.Extract(ctx, pathArchive.String(), pathTmp)
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "extract export template from archive: " + pathArchive.String(),
	}, nil
}
