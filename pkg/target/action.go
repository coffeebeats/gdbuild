package target

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/charmbracelet/log"

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
	exportAction, err := xp.Action(rc)
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

	cacheArtifacts, err := NewCacheArtifactsAction(rc, xp, artifacts)
	if err != nil {
		return nil, err
	}

	return action.InOrder(
		xp.RunBefore,
		NewInstallGodotAction(rc, xp.Version, rc.PathWorkspace),
		extractTemplateAction,
		exportAction,
		xp.RunAfter,
		run.NewVerifyArtifactsAction(rc, artifacts),
		cacheArtifacts,
		NewCopyArtifactsAction(rc, artifacts),
	), nil
}

/* -------------------------------------------------------------------------- */
/*                      Function: NewCacheArtifactsAction                     */
/* -------------------------------------------------------------------------- */

// NewCacheArtifactsAction creates an 'action.Action' which caches the generated
// project artifacts in the 'gdbuild' store.
func NewCacheArtifactsAction(
	rc *run.Context,
	xp *export.Export,
	artifacts []string,
) (action.WithDescription[action.Function], error) {
	fn := func(_ context.Context) error {
		pathBin := rc.BinPath()
		if err := pathBin.CheckIsDir(); err != nil {
			return err
		}

		pathStore, err := store.Path()
		if err != nil {
			return err
		}

		var files []string

		for _, a := range artifacts {
			pathArtifact := filepath.Join(pathBin.String(), a)

			log.Debugf("caching artifact in store: %s", a)

			files = append(files, pathArtifact)
		}

		cs, err := export.Checksum(rc, xp)
		if err != nil {
			return err
		}

		pathArchive, err := store.TargetArchive(pathStore, cs)
		if err != nil {
			return err
		}

		return archive.Create(files, pathArchive)
	}

	storePath, err := store.Path()
	if err != nil {
		return action.WithDescription[action.Function]{}, err
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "cache generated artifacts in store: " + storePath,
	}, nil
}

/* -------------------------------------------------------------------------- */
/*                      Function: NewCopyArtifactsAction                      */
/* -------------------------------------------------------------------------- */

// NewCopyArtifactsAction creates an 'action.Action' which moves the generated
// Godot artifacts to the output directory.
func NewCopyArtifactsAction(
	rc *run.Context,
	artifacts []string,
) action.WithDescription[action.Function] {
	fn := func(ctx context.Context) error {
		if rc.PathOut == "" {
			return nil
		}

		pathOut := rc.PathOut.String()
		if err := osutil.EnsureDir(pathOut, osutil.ModeUserRWXGroupRX); err != nil {
			return err
		}

		pathBin := rc.BinPath()
		if err := pathBin.CheckIsDir(); err != nil {
			return err
		}

		for _, a := range artifacts {
			pathArtifact := filepath.Join(pathBin.String(), a)

			log.Debugf("copying artifact %s to directory: %s", a, pathOut)

			if err := osutil.CopyFile(
				ctx,
				pathArtifact,
				filepath.Join(pathOut, a),
			); err != nil {
				return err
			}
		}

		return nil
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "move generated artifacts to output directory: " + rc.PathOut.String(),
	}
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

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "extract cached export template: " + template.Name(rc.Platform, tl.Arch, rc.Profile),
	}, nil
}
