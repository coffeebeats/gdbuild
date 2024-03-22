package template

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/mitchellh/hashstructure/v2"
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
	ToTemplate(cg build.Godot, tc build.Context) Template
}

/* -------------------------------------------------------------------------- */
/*                              Struct: Template                              */
/* -------------------------------------------------------------------------- */

// Template defines a Godot export template compilation. Its scope is limited to
// the compilation step.
type Template struct {
	// Binaries is a list of export template compilation definitions that are
	// required by the resulting export template artifact.
	Binaries []Build `hash:"set"`

	// ExtraArtifacts are the base names of export template artifacts which are
	// expected to be found in the 'bin' directory post-compilation. If these
	// are missing, 'gdbuild' will consider the build to have failed. Note that
	// the artifacts pertaining to 'Binaries' do not need to be specified.
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

	for _, b := range t.Binaries {
		artifacts[b.Filename()] = struct{}{}
	}

	for _, a := range t.ExtraArtifacts {
		artifacts[a] = struct{}{}
	}

	return maps.Keys(artifacts)
}

/* --------------------------- Method: AddToPaths --------------------------- */

// AddToPaths is a convenience function for registering a 'Path' dependency, but
// only if it hasn't been added yet.
func (t *Template) AddToPaths(path pathutil.Path) {
	if !slices.Contains(t.Paths, path) {
		t.Paths = append(t.Paths, path)
	}
}

/* ---------------------------- Method: Checksum ---------------------------- */

// Checksum produces a checksum hash of the export template specification. When
// the checksums of two 'Template' definitions matches, the resulting export
// templates will be equivalent.
//
// NOTE: This implementation relies on producers of 'Template' to correctly
// register all file system dependencies within 'Paths'.
func (t *Template) Checksum() (string, error) {
	hash, err := hashstructure.Hash(
		t,
		hashstructure.FormatV2,
		&hashstructure.HashOptions{ //nolint:exhaustruct
			IgnoreZeroValue: true,
			SlicesAsSets:    true,
			ZeroNil:         true,
		},
	)
	if err != nil {
		return "", err
	}

	cs := crc32.New(crc32.IEEETable)

	// Update the 'crc32' hash with the struct hash.
	if _, err := io.Copy(cs, strings.NewReader(strconv.FormatUint(hash, 16))); err != nil {
		return "", err
	}

	for _, p := range t.uniquePaths() {
		root := p.String()

		log.Debugf("hashing files rooted at path: %s", root)

		if err := osutil.HashFiles(cs, root); err != nil {
			return "", err
		}
	}

	return strconv.FormatUint(uint64(cs.Sum32()), 16), nil
}

/* --------------------------- Method: uniquePaths -------------------------- */

// uniquePaths returns the unique list of expanded path dependencies.
func (t *Template) uniquePaths() []pathutil.Path {
	paths := t.Paths

	for _, b := range t.Binaries {
		paths = append(paths, b.CustomModules...)

		if b.CustomPy != "" {
			paths = append(paths, b.CustomPy)
		}

		switch g := b.Godot; {
		case g.PathSource != "":
			paths = append(paths, g.PathSource)
		case g.VersionFile != "":
			paths = append(paths, g.VersionFile)
		}
	}

	slices.Sort(paths)

	return slices.Compact(paths)
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
		2+1+1+len(t.Binaries),
	)

	actions = append(
		actions,
		t.Prebuild,
		build.NewVendorGodotAction(&t.Binaries[0].Godot, &c.context.Invoke),
	)

	for _, b := range t.Binaries {
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
