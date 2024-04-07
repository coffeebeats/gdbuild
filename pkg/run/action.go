package run

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/osutil"
)

/* -------------------------------------------------------------------------- */
/*                     Function: NewVerifyArtifactsAction                     */
/* -------------------------------------------------------------------------- */

// NewVerifyArtifactsAction creates an 'action.Action' which verifies that all
// required artifacts have been generated.
func NewVerifyArtifactsAction(
	_ *Context,
	root osutil.Path,
	artifacts []string,
) action.WithDescription[action.Function] {
	fn := func(_ context.Context) error {
		if err := root.CheckIsDir(); err != nil {
			return err
		}

		found := make(map[string]struct{})

		if err := fs.WalkDir(os.DirFS(root.String()), ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			found[path] = struct{}{}

			return nil
		}); err != nil {
			return err
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
				filepath.Join(root.String(), a),
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
func NewCopyArtifactsAction( //nolint:ireturn
	rc *Context,
	root osutil.Path,
	artifacts []string,
) action.Action {
	if rc.PathOut == "" {
		return action.NoOp{}
	}

	fn := func(ctx context.Context) error {
		if rc.PathOut == "" {
			return nil
		}

		pathOut := rc.PathOut.String()
		if err := osutil.EnsureDir(pathOut, osutil.ModeUserRWXGroupRX); err != nil {
			return err
		}

		if err := root.CheckIsDir(); err != nil {
			return err
		}

		for _, a := range artifacts {
			pathArtifact := filepath.Join(root.String(), a)

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
		Description: "move generated artifacts to output directory: " + rc.PathOut.String(),
	}
}

/* -------------------------------------------------------------------------- */
/*                           Function: NewWaitAction                          */
/* -------------------------------------------------------------------------- */

// NewWaitAction creates an 'action.Action' which blocks the thread for the
// specified duration.
func NewWaitAction(d time.Duration) action.WithDescription[action.Function] {
	fn := func(ctx context.Context) error {
		timer := time.NewTimer(d)

		select {
		case <-timer.C:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "wait for the specified duration: " + d.String(),
	}
}
