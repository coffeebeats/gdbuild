package template

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/build"
)

var (
	ErrMissingInput  = errors.New("missing input")
	ErrUnimplemented = errors.New("unimplemented")
)

/* -------------------------------------------------------------------------- */
/*                             Interface: Template                            */
/* -------------------------------------------------------------------------- */

// Template is a type which defines settings for compiling a Godot export
// template.
type Template interface {
	action.Actioner
	build.Configurer
	build.Validater
}

/* --------------------- Function: newVendorGodotAction --------------------- */

// newVendorGodotAction creates an 'action.Action' which vendors Godot source
// code into the build directory.
func newVendorGodotAction(g *build.Godot, inv *build.Invocation) action.Action { //nolint:ireturn
	fn := func(ctx context.Context) error {
		if g.IsEmpty() {
			log.Debug("no Godot version set; skipping vendoring of source code")

			return nil
		}

		pathSource, err := filepath.Abs(g.PathSource.String())
		if err != nil {
			return err
		}

		pathBuild, err := filepath.Abs(inv.PathBuild.String())
		if err != nil {
			return err
		}

		if pathSource == pathBuild {
			log.Info("build directory is source directory; skipping vendoring of source code")

			return nil
		}

		if err := osutil.EnsureDir(pathBuild, osutil.ModeUserRWXGroupRX); err != nil {
			return err
		}

		return g.VendorTo(ctx, pathBuild)
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "<go function: vendor godot source code>",
	}
}

/* -------------------- Function: newMoveArtifactsAction -------------------- */

// newMoveArtifactsAction creates an 'action.Action' which moves the generated
// Godot artifacts to the output directory.
func newMoveArtifactsAction(inv *build.Invocation) action.Action { //nolint:ireturn
	fn := func(ctx context.Context) error {
		pathOut := inv.PathOut.String()
		if err := osutil.EnsureDir(pathOut, osutil.ModeUserRWXGroupRX); err != nil {
			return err
		}

		pathBin := inv.BinPath()
		if err := pathBin.CheckIsDir(); err != nil {
			return err
		}

		ff, err := os.ReadDir(pathBin.String())
		if err != nil {
			return err
		}

		for _, f := range ff {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			log.Debugf("moving artifact %s: %s", f.Name(), pathOut)

			if err := os.Rename(
				filepath.Join(pathBin.String(), f.Name()),
				filepath.Join(pathOut, f.Name()),
			); err != nil {
				return err
			}
		}

		return nil
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "<go function: move generated artifacts to output directory>",
	}
}
