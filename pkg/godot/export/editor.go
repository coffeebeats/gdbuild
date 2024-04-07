package export

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                    Function: NewInstallEditorGodotAction                   */
/* -------------------------------------------------------------------------- */

// NewInstallEditorGodotAction creates an 'action.Action' which installs the
// Godot editor into the build directory.
func NewInstallEditorGodotAction(
	_ *run.Context,
	ev engine.Version,
	pathGodot osutil.Path,
) action.WithDescription[action.Function] {
	fn := func(ctx context.Context) error {
		info, err := os.Stat(pathGodot.String())
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		}

		if info != nil && info.IsDir() {
			return fmt.Errorf(
				"%w: expected a path to a Godot executable file: %s",
				ErrInvalidInput,
				pathGodot,
			)
		}

		if err := engine.InstallEditor(ctx, ev, pathGodot.String()); err != nil {
			return err
		}

		// NOTE: Make the editor run in self-contained mode.
		f, err := os.Create(filepath.Join(filepath.Dir(pathGodot.String()), "._sc_"))
		if err != nil {
			return err
		}

		defer f.Close()

		return nil
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: fmt.Sprintf("install godot '%s' editor: %s", ev.String(), pathGodot),
	}
}

/* -------------------------------------------------------------------------- */
/*                       Function: NewLoadProjectAction                       */
/* -------------------------------------------------------------------------- */

// NewLoadProjectAction creates an 'action.Action' which opens the Godot project
// is the editor for the purpose of generating import files.
func NewLoadProjectAction(
	rc *run.Context,
	pathGodotEditor osutil.Path,
) *action.Process {
	var cmd action.Process

	cmd.Directory = rc.PathWorkspace.String()
	cmd.Verbose = rc.Verbose

	cmd.Args = []string{
		pathGodotEditor.String(),
		"--editor",
		"--headless",
		"--quit-after",
		"2",
	}

	return &cmd
}

/* -------------------------------------------------------------------------- */
/*                        Function: NewRemoveAllAction                        */
/* -------------------------------------------------------------------------- */

// NewRemoveAllAction creates a new 'action.Action' which removes the specified
// directory if preset.
func NewRemoveAllAction(path string) action.WithDescription[action.Function] {
	fn := func(_ context.Context) error {
		return os.RemoveAll(path)
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "remove directory if preset: " + path,
	}
}
