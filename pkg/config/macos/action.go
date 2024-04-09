package macos

import (
	"archive/zip"
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                        Function: NewAppBundleAction                        */
/* -------------------------------------------------------------------------- */

// NewAppBundleAction creates an 'action.Action' which generates an app bundle
// from the specified export templates.
func NewAppBundleAction(
	rc *run.Context,
	artifacts []string,
) action.WithDescription[action.Function] {
	fn := func(ctx context.Context) error {
		tmp, err := rc.TempDir()
		if err != nil {
			return err
		}

		pathApp := filepath.Join(tmp, "macos_template.app")

		if err := osutil.CopyDir(
			rc.PathWorkspace.Join("misc/dist/macos_template.app").String(),
			pathApp,
		); err != nil {
			return err
		}

		if err := os.MkdirAll(
			filepath.Join(pathApp, "Contents/MacOS"),
			osutil.ModeUserRWX,
		); err != nil {
			return err
		}

		for _, artifact := range artifacts {
			// Expected name for artifact [1]: 'godot_macos_release.universal',
			// or something similar. Double precision builds should have the
			// same naming format.
			//
			// [1] https://docs.godotengine.org/en/stable/contributing/development/compiling/compiling_for_macos.html.
			target := strings.Replace(artifact, ".double", "", 1)
			target = strings.Replace(target, "template_", "", 1)
			target = strings.Replace(target, ".", "_", 2) //nolint:gomnd

			pathDst := filepath.Join(pathApp, "Contents/MacOS", target)

			if err := osutil.CopyFile(
				ctx,
				rc.BinPath().Join(artifact).String(),
				pathDst,
			); err != nil {
				return err
			}
		}

		return zipAppBundle(
			osutil.Path(pathApp),
			rc.BinPath().Join("macos.zip"),
		)
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "create application bundle for artifacts: " + strings.Join(artifacts, ","),
	}
}

/* ------------------------- Function: zipAppBundle ------------------------- */

func zipAppBundle(pathAppBundle, out osutil.Path) error {
	if err := pathAppBundle.CheckIsDir(); err != nil {
		return err
	}

	if err := os.MkdirAll(
		filepath.Dir(out.String()),
		osutil.ModeUserRWX,
	); err != nil {
		return err
	}

	f, err := os.Create(out.String())
	if err != nil {
		return err
	}

	defer f.Close()

	archive := zip.NewWriter(f)

	defer archive.Close()

	root := pathAppBundle.String()

	return fs.WalkDir(
		os.DirFS(root),
		".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				_, err := archive.Create(filepath.Join(filepath.Base(root), path) + "/")

				return err
			}

			w, err := archive.Create(filepath.Join(filepath.Base(root), path))
			if err != nil {
				return err
			}

			f, err := os.Open(filepath.Join(root, path))
			if err != nil {
				return err
			}

			defer f.Close()

			_, err = io.Copy(w, f)

			return err
		},
	)
}
