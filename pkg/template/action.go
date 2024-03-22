package template

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/build"
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
		build.NewVendorGodotAction(&t.Builds[0].Source, &bc.Invoke),
	)

	for _, b := range t.Builds {
		actions = append(actions, b.SConsCommand(bc))
	}

	actions = append(
		actions,
		t.Postbuild,
		NewVerifyArtifactsAction(&bc.Invoke, t.Artifacts()),
		NewCopyArtifactsAction(&bc.Invoke, t.Artifacts()),
	)

	return action.InOrder(actions...), nil
}

/* -------------------------------------------------------------------------- */
/*                     Function: NewVerifyArtifactsAction                     */
/* -------------------------------------------------------------------------- */

// NewVerifyArtifactsAction creates an 'action.Action' which verifies that all
// required artifacts have been generated.
func NewVerifyArtifactsAction(
	bc *config.Context,
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
/*                      Function: NewCopyArtifactsAction                      */
/* -------------------------------------------------------------------------- */

// NewCopyArtifactsAction creates an 'action.Action' which moves the generated
// Godot artifacts to the output directory.
func NewCopyArtifactsAction(
	cc *config.Context,
	artifacts []string,
) action.WithDescription[action.Function] {
	fn := func(ctx context.Context) error {
		pathOut := cc.PathOut.String()
		if err := osutil.EnsureDir(pathOut, osutil.ModeUserRWXGroupRX); err != nil {
			return err
		}

		pathBin := cc.BinPath()
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
		Description: "move generated artifacts to output directory: " + cc.PathOut.String(),
	}
}
