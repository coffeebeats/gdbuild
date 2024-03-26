package export

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                       Function: NewInstallGodotAction                      */
/* -------------------------------------------------------------------------- */

// NewInstallGodotAction creates an 'action.Action' which installs the Godot
// editor into the build directory.
func NewInstallGodotAction(
	rc *run.Context,
	ev engine.Version,
	path osutil.Path,
) action.WithDescription[action.Function] {
	fn := func(ctx context.Context) error {
		pathBuild, err := filepath.Abs(rc.PathWorkspace.String())
		if err != nil {
			return err
		}

		if err := osutil.EnsureDir(pathBuild, osutil.ModeUserRWXGroupRX); err != nil {
			return err
		}

		pathGodot := path.String()

		info, err := os.Stat(pathGodot)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		}

		if info != nil && info.IsDir() {
			pathGodot = filepath.Join(pathGodot, engine.EditorName())
		}

		return engine.InstallEditor(ctx, ev, pathGodot)
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "install godot editor: " + path.String(),
	}
}
