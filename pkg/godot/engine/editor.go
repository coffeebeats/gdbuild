package engine

import (
	"context"
	"os"
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/pkg/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
	"github.com/coffeebeats/gdenv/pkg/install"
	"github.com/coffeebeats/gdenv/pkg/store"

	"github.com/coffeebeats/gdbuild/internal/osutil"
)

/* -------------------------------------------------------------------------- */
/*                            Function: EditorName                            */
/* -------------------------------------------------------------------------- */

// EditorName returns the expected name of the Godot editor executable based on
// the host platform.
func EditorName() string {
	name := "godot"
	if runtime.GOOS == "windows" {
		name = strings.ToTitle(name)
		name += ".exe"
	}

	return name
}

/* -------------------------------------------------------------------------- */
/*                           Function: InstallEditor                          */
/* -------------------------------------------------------------------------- */

// InstallEditor installs the Godot editor to the specified directory.
func InstallEditor(ctx context.Context, v Version, out string) error {
	storePath, err := store.Path() // Use the 'gdenv' store.
	if err != nil {
		tmp, err := os.MkdirTemp("", "gdbuild-*")
		if err != nil {
			return err
		}

		log.Debugf("no 'gdenv' store found; using temporary directory: %s", tmp)

		storePath = tmp
	}

	pl, err := platform.Detect()
	if err != nil {
		return err
	}

	ex := executable.New(version.Version(v), pl)

	if err := install.Executable(
		ctx,
		storePath,
		ex,
		/* force= */ false,
	); err != nil {
		return err
	}

	godot, err := store.Executable(storePath, ex)
	if err != nil {
		return err
	}

	return osutil.CopyFile(ctx, godot, out)
}
