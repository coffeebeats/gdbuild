package build

import (
	"context"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/mitchellh/hashstructure/v2"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/osutil"
)

/* -------------------------------------------------------------------------- */
/*                            Interface: Templater                            */
/* -------------------------------------------------------------------------- */

type Templater interface {
	ToTemplate(godot Godot, inv Invocation) Template
}

/* -------------------------------------------------------------------------- */
/*                              Struct: Template                              */
/* -------------------------------------------------------------------------- */

// Template defines a Godot export template compilation. Its scope is limited to
// the compilation step.
type Template struct {
	// Binaries is a list of export template compilation definitions that are
	// required by the resulting export template artifact.
	Binaries []Binary `hash:"set"`

	// Paths is a list of additional files and folders which this template
	// depends on. Useful for recording dependencies which are defined in
	// otherwise opaque properties like 'Hook'.
	Paths []Path `hash:"set"`

	// Prebuild contains an ordered list of actions to execute prior to
	// compilation of the export templates.
	Prebuild []action.Action `hash:"set,string"`

	// Postbuild contains an ordered list of actions to execute after
	// compilation of the export templates.
	Postbuild []action.Action `hash:"set,string"`
}

/* --------------------------- Method: AddToPaths --------------------------- */

// AddToPaths is a convenience function for registering a 'Path' dependency, but
// only if it hasn't been added yet.
func (t *Template) AddToPaths(path Path) {
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
func (t *Template) Checksum(inv *Invocation) (string, error) {
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

	for _, p := range t.uniquePaths(inv) {
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
func (t *Template) uniquePaths(_ *Invocation) []Path {
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
/*                               Struct: Binary                               */
/* -------------------------------------------------------------------------- */

// Binary uniquely specifies a compilation of a Godot export template.
type Binary struct {
	// Arch is the CPU architecture of the Godot export template.
	Arch Arch

	// CustomModules is a list of paths to custom modules to include in the
	// template build.
	CustomModules []Path `hash:"ignore"`

	// CustomPy is a path to a 'custom.py' file which defines export template
	// build options.
	CustomPy Path `hash:"ignore"`

	// DoublePrecision enables double floating-point precision.
	DoublePrecision bool

	// Env is a map of environment variables to set during the build step.
	Env map[string]string

	// Godot is the source code specification for the build.
	Godot Godot

	// Optimize is the level of optimization for the Godot export template.
	Optimize Optimize

	// Platform defines which OS/platform to build for.
	Platform OS

	// Profile is the optimization level of the template.
	Profile Profile

	// SCons contains a specification for how to invoke the compiler.
	SCons SCons
}

/* -------------------------- Method: SConsCommand -------------------------- */

// SConsCommand returns the 'SCons' command to build the export template.
func (b *Binary) SConsCommand(inv *Invocation) *action.Process { //nolint:cyclop,funlen
	var cmd action.Process

	cmd.Directory = inv.PathBuild.String()
	cmd.Verbose = inv.Verbose

	scons := b.SCons

	// Define the SCons command, if not yet set.
	if len(scons.Command) == 0 {
		scons.Command = append(scons.Command, "scons")
	}

	// Define the SCons cache path.
	if path := scons.PathCache; path != "" {
		cmd.Environment = append(cmd.Environment, EnvSConsCache+"="+path.String())
	}

	// Add specified environment variables.
	for k, v := range b.Env {
		cmd.Environment = append(cmd.Environment, k+"="+v)
	}

	// Now pass through all environment variables so that these override
	// previously values.
	cmd.Environment = append(cmd.Environment, os.Environ()...)

	// Set the SCons cache size limit, if one was set.
	if l := scons.CacheSizeLimit; l != nil {
		cmd.Environment = append(
			cmd.Environment,
			fmt.Sprintf("SCONS_CACHE_LIMIT=%d", *scons.CacheSizeLimit),
		)
	}

	// Build the SCons command/argument list.
	cmd.Args = append(cmd.Args, scons.Command...)
	cmd.Args = append(cmd.Args, "-j"+strconv.Itoa(runtime.NumCPU()))

	// Specify the 'platform' argument.
	cmd.Args = append(cmd.Args, b.Platform.String())

	// Add the achitecture setting (note that this requires the 'build.Arch'
	// values to match what SCons expects).
	cmd.Args = append(cmd.Args, "arch="+b.Arch.String())

	// Specify which target to build.
	cmd.Args = append(cmd.Args, "target="+inv.Profile.TargetName())

	// Add stricter warning handling.
	cmd.Args = append(cmd.Args, "warnings=extra", "werror=yes")

	// Handle a verbose flag.
	if inv.Verbose {
		cmd.Args = append(cmd.Args, "verbose=yes")
	}

	// Append 'custom_modules' argument.
	if len(b.CustomModules) > 0 {
		modules := make([]string, len(b.CustomModules))
		for i, m := range b.CustomModules {
			modules[i] = m.String()
		}

		cmd.Args = append(cmd.Args, fmt.Sprintf(`custom_modules="%s"`, strings.Join(modules, ",")))
	}

	// Append the 'precision' argument.
	if b.DoublePrecision {
		cmd.Args = append(cmd.Args, "precision=double")
	}

	// Append profile/optimization-related arguments.
	switch inv.Profile {
	case ProfileRelease:
		optimize := OptimizeSpeed
		if b.Optimize != OptimizeUnknown {
			optimize = b.Optimize
		}

		cmd.Args = append(
			cmd.Args,
			"production=yes",
			fmt.Sprintf("optimize=%s", optimize),
		)

	case ProfileReleaseDebug:
		optimize := OptimizeSpeedTrace
		if b.Optimize != OptimizeUnknown {
			optimize = b.Optimize
		}

		cmd.Args = append(
			cmd.Args,
			"debug_symbols=yes",
			"dev_mode=yes",
			fmt.Sprintf("optimize=%s", optimize),
		)
	default: // ProfileDebug
		optimize := OptimizeDebug
		if b.Optimize != OptimizeUnknown {
			optimize = b.Optimize
		}

		cmd.Args = append(
			cmd.Args,
			"debug_symbols=yes",
			"dev_mode=yes",
			fmt.Sprintf("optimize=%s", optimize),
		)
	}

	// Append C/C++ build flags.
	if len(b.SCons.CCFlags) > 0 {
		flags := fmt.Sprintf(`CCFLAGS="%s"`, strings.Join(b.SCons.CCFlags, " "))
		cmd.Args = append(cmd.Args, flags)
	}

	if len(b.SCons.CFlags) > 0 {
		flags := fmt.Sprintf(`CFLAGS="%s"`, strings.Join(b.SCons.CFlags, " "))
		cmd.Args = append(cmd.Args, flags)
	}

	if len(b.SCons.CXXFlags) > 0 {
		flags := fmt.Sprintf(`CXXFLAGS="%s"`, strings.Join(b.SCons.CXXFlags, " "))
		cmd.Args = append(cmd.Args, flags)
	}

	if len(b.SCons.LinkFlags) > 0 {
		flags := fmt.Sprintf(`LINKFLAGS="%s"`, strings.Join(b.SCons.LinkFlags, " "))
		cmd.Args = append(cmd.Args, flags)
	}

	// Append extra arguments.
	cmd.Args = append(cmd.Args, b.SCons.ExtraArgs...)

	return &cmd
}

/* -------------------------------------------------------------------------- */
/*                      Function: NewMoveArtifactsAction                      */
/* -------------------------------------------------------------------------- */

// NewMoveArtifactsAction creates an 'action.Action' which moves the generated
// Godot artifacts to the output directory.
func NewMoveArtifactsAction(inv *Invocation) action.Action { //nolint:ireturn
	fn := func(ctx context.Context) error {
		pathOut := inv.PathOut.String()
		if err := osutil.EnsureDir(pathOut, osutil.ModeUserRWXGroupRX); err != nil {
			return err
		}

		pathBin := inv.BinPath()
		if err := pathBin.CheckIsDir(); err != nil {
			return err
		}

		ff, err := os.ReadDir(pathBin.String())
		if err != nil {
			return err
		}

		for _, f := range ff {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			log.Debugf("moving artifact %s: %s", f.Name(), pathOut)

			if err := os.Rename(
				filepath.Join(pathBin.String(), f.Name()),
				filepath.Join(pathOut, f.Name()),
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
func Compile(t *Template, inv *Invocation) (action.Action, error) { //nolint:ireturn
	return compilation{inv, t}.Action()
}

/* --------------------------- Struct: compilation -------------------------- */

type compilation struct {
	invocation *Invocation
	template   *Template
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c compilation) Action() (action.Action, error) { //nolint:ireturn
	t := c.template
	inv := c.invocation

	actions := make(
		[]action.Action,
		0,
		2+len(t.Prebuild)+len(t.Postbuild)+len(t.Binaries),
	)

	actions = append(actions, t.Prebuild...)
	actions = append(actions, NewVendorGodotAction(&t.Binaries[0].Godot, inv))

	for _, b := range t.Binaries {
		actions = append(actions, b.SConsCommand(inv))
	}

	actions = append(actions, t.Postbuild...)
	actions = append(actions, NewMoveArtifactsAction(inv))

	return action.InOrder(actions...), nil
}
