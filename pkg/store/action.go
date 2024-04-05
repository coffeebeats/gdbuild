package store

import (
	"context"
	"os"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/archive"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                       Function: NewCacheTargetAction                       */
/* -------------------------------------------------------------------------- */

// NewCacheTargetAction creates an 'action.Action' which caches the generated
// project artifacts in the 'gdbuild' store.
func NewCacheTargetAction(
	_ *run.Context,
	root osutil.Path,
	artifacts []string,
	checksum string,
) (action.WithDescription[action.Function], error) {
	storePath, err := Path()
	if err != nil {
		return action.WithDescription[action.Function]{}, err
	}

	fn := func(ctx context.Context) error {
		pathArchive, err := TargetArchive(storePath, checksum)
		if err != nil {
			return err
		}

		return archiveArtifacts(ctx, root, artifacts, pathArchive)
	}

	return describeCacheFn(fn)
}

/* -------------------------------------------------------------------------- */
/*                      Function: NewCacheTemplateAction                      */
/* -------------------------------------------------------------------------- */

// NewCacheTemplateAction creates an 'action.Action' which caches the generated
// export template in the 'gdbuild' store.
func NewCacheTemplateAction(
	_ *run.Context,
	root osutil.Path,
	artifacts []string,
	checksum string,
) (action.WithDescription[action.Function], error) {
	storePath, err := Path()
	if err != nil {
		return action.WithDescription[action.Function]{}, err
	}

	fn := func(ctx context.Context) error {
		pathArchive, err := TemplateArchive(storePath, checksum)
		if err != nil {
			return err
		}

		return archiveArtifacts(ctx, root, artifacts, pathArchive)
	}

	return describeCacheFn(fn)
}

/* ------------------------ Function: describeCacheFn ----------------------- */

func describeCacheFn(fn action.Function) (action.WithDescription[action.Function], error) {
	storePath, err := Path()
	if err != nil {
		return action.WithDescription[action.Function]{}, err
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "cache generated artifacts in store: " + storePath,
	}, nil
}

/* ----------------------- Function: archiveArtifacts ----------------------- */

func archiveArtifacts(_ context.Context, root osutil.Path, artifacts []string, out string) error {
	if err := root.CheckIsDir(); err != nil {
		return err
	}

	files := make([]string, 0, len(artifacts))

	for _, a := range artifacts {
		log.Debugf("archiving artifact in store: %s", a)

		files = append(files, a)
	}

	return archive.Create(os.DirFS(root.String()), files, out)
}
