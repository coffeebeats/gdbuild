package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"

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
	Binaries []Compilation

	// Paths is a list of additional files and folders which this template
	// depends on. Useful for recording dependencies which are defined in
	// otherwise opaque properties like 'Hook'.
	Paths []Path

	// Prebuild contains an ordered list of actions to execute prior to
	// compilation of the export templates.
	Prebuild []action.Action

	// Postbuild contains an ordered list of actions to execute after
	// compilation of the export templates.
	Postbuild []action.Action
}

/* ---------------------------- Method: AddToPath --------------------------- */

// AddToPath is a convenience function for registering a 'Path' dependency, but
// only if it hasn't been added yet.
func (t *Template) AddToPath(path Path) {
	if !slices.Contains(t.Paths, path) {
		t.Paths = append(t.Paths, path)
	}
}

/* -------------------------------------------------------------------------- */
/*                             Struct: Compilation                            */
/* -------------------------------------------------------------------------- */

// Compilation uniquely specifies a compilation of a Godot export template.
type Compilation struct {
	// Arch is the CPU architecture of the Godot export template.
	Arch Arch

	// CustomModules is a list of paths to custom modules to include in the
	// template build.
	CustomModules []Path

	// CustomPy is a path to a 'custom.py' file which defines export template
	// build options.
	CustomPy Path

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
func (c *Compilation) SConsCommand(inv *Invocation) *action.Process { //nolint:cyclop,funlen
	var cmd action.Process

	cmd.Directory = inv.PathBuild.String()
	cmd.Verbose = inv.Verbose

	scons := c.SCons

	// Define the SCons command, if not yet set.
	if len(scons.Command) == 0 {
		scons.Command = append(scons.Command, "scons")
	}

	// Define the SCons cache path.
	if path := scons.PathCache; path != "" {
		cmd.Environment = append(cmd.Environment, EnvSConsCache+"="+path.String())
	}

	// Add specified environment variables.
	for k, v := range c.Env {
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
	cmd.Args = append(cmd.Args, c.Platform.String())

	// Add the achitecture setting (note that this requires the 'build.Arch'
	// values to match what SCons expects).
	cmd.Args = append(cmd.Args, "arch="+c.Arch.String())

	// Specify which target to build.
	cmd.Args = append(cmd.Args, "target="+inv.Profile.TargetName())

	// Add stricter warning handling.
	cmd.Args = append(cmd.Args, "warnings=extra", "werror=yes")

	// Handle a verbose flag.
	if inv.Verbose {
		cmd.Args = append(cmd.Args, "verbose=yes")
	}

	// Append 'custom_modules' argument.
	if len(c.CustomModules) > 0 {
		modules := make([]string, len(c.CustomModules))
		for i, m := range c.CustomModules {
			modules[i] = m.String()
		}

		cmd.Args = append(cmd.Args, fmt.Sprintf(`custom_modules="%s"`, strings.Join(modules, ",")))
	}

	// Append the 'precision' argument.
	if c.DoublePrecision {
		cmd.Args = append(cmd.Args, "precision=double")
	}

	// Append profile/optimization-related arguments.
	switch inv.Profile {
	case ProfileRelease:
		optimize := OptimizeSpeed
		if c.Optimize != OptimizeUnknown {
			optimize = c.Optimize
		}

		cmd.Args = append(
			cmd.Args,
			"production=yes",
			fmt.Sprintf("optimize=%s", optimize),
		)

	case ProfileReleaseDebug:
		optimize := OptimizeSpeedTrace
		if c.Optimize != OptimizeUnknown {
			optimize = c.Optimize
		}

		cmd.Args = append(
			cmd.Args,
			"debug_symbols=yes",
			"dev_mode=yes",
			fmt.Sprintf("optimize=%s", optimize),
		)
	default: // ProfileDebug
		optimize := OptimizeDebug
		if c.Optimize != OptimizeUnknown {
			optimize = c.Optimize
		}

		cmd.Args = append(
			cmd.Args,
			"debug_symbols=yes",
			"dev_mode=yes",
			fmt.Sprintf("optimize=%s", optimize),
		)
	}

	// Append C/C++ build flags.
	if len(c.SCons.CCFlags) > 0 {
		flags := fmt.Sprintf(`CCFLAGS="%s"`, strings.Join(c.SCons.CCFlags, " "))
		cmd.Args = append(cmd.Args, flags)
	}

	if len(c.SCons.CFlags) > 0 {
		flags := fmt.Sprintf(`CFLAGS="%s"`, strings.Join(c.SCons.CFlags, " "))
		cmd.Args = append(cmd.Args, flags)
	}

	if len(c.SCons.CXXFlags) > 0 {
		flags := fmt.Sprintf(`CXXFLAGS="%s"`, strings.Join(c.SCons.CXXFlags, " "))
		cmd.Args = append(cmd.Args, flags)
	}

	if len(c.SCons.LinkFlags) > 0 {
		flags := fmt.Sprintf(`LINKFLAGS="%s"`, strings.Join(c.SCons.LinkFlags, " "))
		cmd.Args = append(cmd.Args, flags)
	}

	// Append extra arguments.
	cmd.Args = append(cmd.Args, c.SCons.ExtraArgs...)

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
