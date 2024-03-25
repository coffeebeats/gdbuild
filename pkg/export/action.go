package export

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/archive"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/run"
	"github.com/coffeebeats/gdbuild/pkg/store"
)

var ErrMissingInput = errors.New("missing input")

/* -------------------------------------------------------------------------- */
/*                              Function: Action                              */
/* -------------------------------------------------------------------------- */

// Action creates a new 'action.Action' which executes the specified processes
// for compiling the export template.
func Action(rc *run.Context, x *export.Export) (action.Action, error) { //nolint:ireturn
	actions := make([]action.Action, 0)

	actions = append(
		actions,
		x.RunBefore,
		NewInstallGodotAction(rc, x.Version, rc.PathBuild),
		x.Action(),
	)

	cacheArtifacts, err := NewCacheArtifactsAction(rc, x)
	if err != nil {
		return nil, err
	}

	actions = append(
		actions,
		x.RunAfter,
		run.NewVerifyArtifactsAction(rc, x.Artifacts()),
		cacheArtifacts,
		NewCopyArtifactsAction(rc, x.Artifacts()),
	)

	return action.InOrder(actions...), nil
}

/* -------------------------------------------------------------------------- */
/*                      Function: NewCacheArtifactsAction                     */
/* -------------------------------------------------------------------------- */

// NewCacheArtifactsAction creates an 'action.Action' which caches the generated
// project artifacts in the 'gdbuild' store.
func NewCacheArtifactsAction(
	rc *run.Context,
	x *export.Export,
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

		for _, a := range x.Artifacts() {
			pathArtifact := filepath.Join(pathBin.String(), a)

			log.Debugf("caching artifact in store: %s", a)

			files = append(files, pathArtifact)
		}

		pathArchive, err := store.TargetArchive(pathStore, rc, x)
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
