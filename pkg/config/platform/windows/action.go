package windows

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                      Function: NewCopyImageFileAction                      */
/* -------------------------------------------------------------------------- */

// NewCopyImageFileAction creates an 'action.Action' which places the specified
// icon image into the Godot source code.
func NewCopyImageFileAction(
	pathImage osutil.Path,
	rc *run.Context,
) action.WithDescription[action.Function] {
	pathDst := filepath.Join(rc.PathBuild.String(), "platform/windows/godot.ico")

	fn := func(_ context.Context) error {
		dst, err := os.Create(pathDst)
		if err != nil {
			return err
		}

		defer dst.Close()

		src, err := os.Open(pathImage.String())
		if err != nil {
			return err
		}

		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		return nil
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "copy icon into build directory: " + pathDst,
	}
}
