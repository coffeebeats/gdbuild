package template

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/archive"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/build"
	"github.com/coffeebeats/gdbuild/pkg/store"
)

/* -------------------------------------------------------------------------- */
/*                              Function: Action                              */
/* -------------------------------------------------------------------------- */

// Action creates a new 'action.Action' which executes the specified processes
// for compiling the export template.
func Action(t *build.Template, bc *build.Context) (action.Action, error) { //nolint:ireturn
	actions := make(
		[]action.Action,
		0,
		2+1+1+len(t.Builds),
	)

	actions = append(
		actions,
		t.Prebuild,
		build.NewVendorGodotAction(&t.Builds[0].Source, bc),
	)

	for _, b := range t.Builds {
		actions = append(actions, b.SConsCommand(bc))
	}

	actions = append(
		actions,
		t.Postbuild,
		NewVerifyArtifactsAction(bc, t.Artifacts()),
		NewCacheArtifactsAction(bc, t),
		NewCopyArtifactsAction(bc, t.Artifacts()),
	)

	return action.InOrder(actions...), nil
}

/* -------------------------------------------------------------------------- */
/*                     Function: NewVerifyArtifactsAction                     */
/* -------------------------------------------------------------------------- */

// NewVerifyArtifactsAction creates an 'action.Action' which verifies that all
// required artifacts have been generated.
func NewVerifyArtifactsAction(
	bc *build.Context,
	artifacts []string,
) action.WithDescription[action.Function] {
	fn := func(_ context.Context) error {
		pathBin := bc.BinPath()
		if err := pathBin.CheckIsDir(); err != nil {
			return err
		}

		ff, err := os.ReadDir(pathBin.String())
		if err != nil {
			return err
		}

		found := make(map[string]struct{})

		for _, f := range ff {
			found[f.Name()] = struct{}{}
		}

		for _, a := range artifacts {
			if _, ok := found[a]; !ok {
				return fmt.Errorf(
					"%w: required file not generated: %s",
					ErrMissingInput,
					a,
				)
			}

			log.Debugf(
				"found required artifact: %s",
				filepath.Join(pathBin.String(), a),
			)
		}

		return nil
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "validate generated artifacts: " + strings.Join(artifacts, ", "),
	}
}

/* -------------------------------------------------------------------------- */
/*                      Function: NewCacheArtifactsAction                     */
/* -------------------------------------------------------------------------- */

// NewCacheArtifactsAction creates an 'action.Action' which caches the generated
// Godot artifacts in the 'gdbuild' store.
func NewCacheArtifactsAction(
	bc *build.Context,
	t *build.Template,
) action.WithDescription[action.Function] {
	fn := func(_ context.Context) error {
		pathBin := bc.BinPath()
		if err := pathBin.CheckIsDir(); err != nil {
			return err
		}

		pathStore, err := store.Path()
		if err != nil {
			return err
		}

		var files []string

		for _, a := range t.Artifacts() {
			pathArtifact := filepath.Join(pathBin.String(), a)

			log.Debugf("caching artifact in store: %s", a)

			files = append(files, pathArtifact)
		}

		pathArchive, err := store.TemplateArchive(pathStore, t)
		if err != nil {
			return err
		}

		return archive.Create(files, pathArchive)
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "cache generated artifacts in store",
	}
}

/* -------------------------------------------------------------------------- */
/*                      Function: NewCopyArtifactsAction                      */
/* -------------------------------------------------------------------------- */

// NewCopyArtifactsAction creates an 'action.Action' which moves the generated
// Godot artifacts to the output directory.
func NewCopyArtifactsAction(
	bc *build.Context,
	artifacts []string,
) action.WithDescription[action.Function] {
	fn := func(ctx context.Context) error {
		pathOut := bc.PathOut.String()
		if err := osutil.EnsureDir(pathOut, osutil.ModeUserRWXGroupRX); err != nil {
			return err
		}

		pathBin := bc.BinPath()
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
		Description: "move generated artifacts to output directory: " + bc.PathOut.String(),
	}
}
