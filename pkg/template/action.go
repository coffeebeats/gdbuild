package template

import (
	"errors"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
	"github.com/coffeebeats/gdbuild/pkg/run"
	"github.com/coffeebeats/gdbuild/pkg/store"
)

var ErrMissingInput = errors.New("missing input")

/* -------------------------------------------------------------------------- */
/*                              Function: Action                              */
/* -------------------------------------------------------------------------- */

// Action creates a new 'action.Action' which executes the specified processes
// for compiling the export template.
func Action(rc *run.Context, tl *template.Template) (action.Action, error) { //nolint:ireturn,nolintlint
	actions := make([]action.Action, 0)

	actions = append(
		actions,
		tl.Prebuild,
		template.NewVendorGodotAction(&tl.Builds[0].Source, rc),
		tl.Action(rc),
	)

	cs, err := template.Checksum(tl)
	if err != nil {
		return nil, err
	}

	pathBin := rc.BinPath()
	artifacts := tl.Artifacts()

	cacheArtifacts, err := store.NewCacheTemplateAction(rc, pathBin, artifacts, cs)
	if err != nil {
		return nil, err
	}

	actions = append(
		actions,
		tl.Postbuild,
		run.NewVerifyArtifactsAction(rc, pathBin, artifacts),
		cacheArtifacts,
		run.NewCopyArtifactsAction(rc, pathBin, artifacts),
	)

	return action.InOrder(actions...), nil
}
