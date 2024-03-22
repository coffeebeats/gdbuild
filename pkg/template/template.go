package template

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"golang.org/x/exp/maps"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/internal/pathutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/build"
)

var ErrMissingInput = errors.New("missing input")

/* -------------------------------------------------------------------------- */
/*                            Interface: Templater                            */
/* -------------------------------------------------------------------------- */

type Templater interface {
	ToTemplate(cg build.Source, tc build.Context) *Template
}

/* -------------------------------------------------------------------------- */
/*                              Struct: Template                              */
/* -------------------------------------------------------------------------- */

// Template defines a Godot export template compilation. Its scope is limited to
// the compilation step.
type Template struct {
	// Builds is a list of export template compilation definitions that are
	// required by the resulting export template artifact.
	Builds []build.Build `hash:"set"`

	// ExtraArtifacts are the base names of export template artifacts which are
	// expected to be found in the 'bin' directory post-compilation. If these
	// are missing, 'gdbuild' will consider the build to have failed. Note that
	// the artifacts pertaining to 'Builds' do not need to be specified.
	ExtraArtifacts []string `hash:"ignore"`

	// Paths is a list of additional files and folders which this template
	// depends on. Useful for recording dependencies which are defined in
	// otherwise opaque properties like 'Hook'.
	Paths []pathutil.Path `hash:"set"`

	// Prebuild contains an ordered list of actions to execute prior to
	// compilation of the export templates.
	Prebuild action.Action `hash:"string"`

	// Postbuild contains an ordered list of actions to execute after
	// compilation of the export templates.
	Postbuild action.Action `hash:"string"`
}

/* ---------------------------- Method: Artifacts --------------------------- */

// Artifacts returns the set of export template artifacts required by the
// underlying template build definition. This will join the files generated by
// the included 'Binary' definitions with those added in 'ExtraArtifacts'.
func (t *Template) Artifacts() []string {
	artifacts := make(map[string]struct{})

	for _, b := range t.Builds {
		artifacts[b.Filename()] = struct{}{}
	}

	for _, a := range t.ExtraArtifacts {
		artifacts[a] = struct{}{}
	}

	return maps.Keys(artifacts)
}

/* --------------------- Method: RegisterDependencyPath --------------------- */

// RegisterDependencyPath is a convenience function for registering a 'Path'
// dependency, but only if it hasn't been added yet.
func (t *Template) RegisterDependencyPath(path pathutil.Path) {
	if !slices.Contains(t.Paths, path) {
		t.Paths = append(t.Paths, path)
	}
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
	inv *config.Context,
	artifacts []string,
) action.WithDescription[action.Function] {
	fn := func(ctx context.Context) error {
		pathOut := inv.PathOut.String()
		if err := osutil.EnsureDir(pathOut, osutil.ModeUserRWXGroupRX); err != nil {
			return err
		}

		pathBin := inv.BinPath()
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
		Description: "move generated artifacts to output directory: " + inv.PathOut.String(),
	}
}

/* -------------------------------------------------------------------------- */
/*                              Function: Compile                             */
/* -------------------------------------------------------------------------- */

// Compile creates a new 'action.Action' which executes the specified processes
// for compiling the export template.
func Compile(t *Template, bc *build.Context) (action.Action, error) { //nolint:ireturn
	return compilation{context: bc, template: t}.Action()
}

/* --------------------------- Struct: compilation -------------------------- */

type compilation struct {
	context  *build.Context
	template *Template
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c compilation) Action() (action.Action, error) { //nolint:ireturn
	t := c.template

	actions := make(
		[]action.Action,
		0,
		2+1+1+len(t.Builds),
	)

	actions = append(
		actions,
		t.Prebuild,
		build.NewVendorGodotAction(&t.Builds[0].Source, &c.context.Invoke),
	)

	for _, b := range t.Builds {
		actions = append(actions, b.SConsCommand(c.context))
	}

	actions = append(
		actions,
		t.Postbuild,
		NewVerifyArtifactsAction(&c.context.Invoke, t.Artifacts()),
		NewCopyArtifactsAction(&c.context.Invoke, t.Artifacts()),
	)

	return action.InOrder(actions...), nil
}
