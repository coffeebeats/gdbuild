package template

import (
	"context"
	"path/filepath"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                       Function: NewVendorGodotAction                       */
/* -------------------------------------------------------------------------- */

// NewVendorGodotAction creates an 'action.Action' which vendors Godot source
// code into the build directory.
func NewVendorGodotAction(src *engine.Source, rc *run.Context) action.WithDescription[action.Function] {
	fn := func(ctx context.Context) error {
		if src.IsEmpty() {
			log.Debug("no Godot version set; skipping vendoring of source code")

			return nil
		}

		pathSource, err := filepath.Abs(src.PathSource.String())
		if err != nil {
			return err
		}

		pathBuild, err := filepath.Abs(rc.PathWorkspace.String())
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

		return src.VendorTo(ctx, pathBuild)
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "vendor godot source code: " + rc.PathWorkspace.String(),
	}
}
