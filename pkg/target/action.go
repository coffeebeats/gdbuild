package target

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/archive"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
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
func Action(rc *run.Context, tl *template.Template, xp *export.Export) (action.Action, error) { //nolint:ireturn
	pathTmp, err := rc.TempDir()
	if err != nil {
		return nil, err
	}

	pathGodot := osutil.Path(filepath.Join(pathTmp, "godot"))

	exportAction, err := xp.Action(rc, pathGodot)
	if err != nil {
		return nil, err
	}

	extractTemplateAction, err := NewExtractTemplateAction(rc, tl)
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
		xp.RunBefore,
		export.NewInstallEditorGodotAction(rc, xp.Version, osutil.Path(pathTmp)),
		extractTemplateAction,
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
	tl *template.Template,
) (action.WithDescription[action.Function], error) {
	storePath, err := store.Path()
	if err != nil {
		return action.WithDescription[action.Function]{}, err
	}

	checksum, err := template.Checksum(tl)
	if err != nil {
		return action.WithDescription[action.Function]{}, err
	}

	pathTmp, err := rc.TempDir()
	if err != nil {
		return action.WithDescription[action.Function]{}, err
	}

	fn := func(ctx context.Context) error {
		pathArchive, err := store.TemplateArchive(storePath, checksum)
		if err != nil {
			return err
		}

		return archive.Extract(ctx, pathArchive, pathTmp)
	}

	pathTemplate := filepath.Join(pathTmp, template.Name(rc.Platform, tl.Arch, rc.Profile))

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "extract cached export template: " + pathTemplate,
	}, nil
}
