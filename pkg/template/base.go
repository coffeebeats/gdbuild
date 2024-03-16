package template

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/exec"
	"github.com/coffeebeats/gdbuild/pkg/build"
)

var ErrInvalidInput = errors.New("invalid input")

/* -------------------------------------------------------------------------- */
/*                                Struct: Base                                */
/* -------------------------------------------------------------------------- */

// Base contains platform-agnostic settings for constructing a custom Godot
// export template.
type Base struct {
	// Arch is the CPU architecture of the Godot export template.
	Arch build.Arch `toml:"arch"`
	// CustomModules is a list of paths to custom modules to include in the
	// template build.
	CustomModules []build.Path `toml:"custom_modules"`
	// DoublePrecision enables double floating-point precision.
	DoublePrecision *bool `toml:"double_precision"`
	// Env is a map of environment variables to set during the build step.
	Env map[string]string `toml:"env"`
	// Hook defines commands to be run before or after a build step.
	Hook build.Hook `toml:"hook"`
	// Optimize is the specific optimization level for the template.
	Optimize build.Optimize `toml:"optimize"`
	// PathCustomPy is a path to a 'custom.py' file which defines export
	// template build options.
	PathCustomPy build.Path `toml:"custom_py_path"`
	// SCons contains build command-related settings.
	SCons build.SCons `toml:"scons"`
	// Shell is a shell name to run the commands with. Defaults to 'sh' on Unix
	// hosts and 'PowerShell' on Windows.
	Shell exec.Shell `toml:"shell"`

	// Invocation contains invocation-specific properties. These must be set
	// manually prior to executing a template build.
	build.Invocation
	// Godot contains a specification of which Godot version to use.
	build.Godot
}

/* ----------------------------- Method: action ----------------------------- */

// action builds an 'action.Action' which executes the entire workflow for
// compiling the Godot export template and moving artifacts to the output path.
// A reference to the 'action.Process' containing the SCons command is provided
// so that platform-specific 'Template' types can modify it prior to execution.
func (c *Base) action() (*action.Process, error) { //nolint:cyclop,funlen
	var cmd action.Process

	cmd.Directory = c.Invocation.PathBuild.String()

	cmd.Verbose = c.Invocation.Verbose

	// Define the SCons cache path.
	if path := c.SCons.PathCache; path != "" {
		cmd.Environment = append(cmd.Environment, build.EnvSConsCache+"="+path.String())
	}

	// Add specified environment variables.
	for k, v := range c.Env {
		cmd.Environment = append(cmd.Environment, k+"="+v)
	}

	// Now pass through all environment variables so that these override
	// previously values.
	cmd.Environment = append(cmd.Environment, os.Environ()...)

	// Build the SCons command/argument list.
	cmd.Args = append(cmd.Args, c.SCons.Command...)
	cmd.Args = append(cmd.Args, "-j"+strconv.Itoa(runtime.NumCPU()))

	// Add stricter warning handling.
	cmd.Args = append(cmd.Args, "warnings=extra", "werror=yes")

	// Handle a verbose flag.
	if c.Invocation.Verbose {
		cmd.Args = append(cmd.Args, "verbose=yes")
	}

	// Add the achitecture setting (note that this requires the 'build.Arch'
	// values to match what SCons expects).
	cmd.Args = append(cmd.Args, "arch="+c.Arch.String())

	// Append 'custom_modules' argument.
	if len(c.CustomModules) > 0 {
		modules := make([]string, len(c.CustomModules))
		for i, m := range c.CustomModules {
			modules[i] = m.String()
		}

		cmd.Args = append(cmd.Args, fmt.Sprintf(`custom_modules="%s"`, strings.Join(modules, ",")))
	}

	// Append the 'precision' argument.
	if c.DoublePrecision != nil && *c.DoublePrecision {
		cmd.Args = append(cmd.Args, "precision=double")
	}

	// Append the 'custom.py' argument.
	if path := c.PathCustomPy; path != "" {
		cmd.Args = append(cmd.Args, "custom="+path.String())
	}

	cmd.Args = append(cmd.Args, "target="+c.targetName())

	// Append profile/optimization-related arguments.
	switch c.Invocation.Profile {
	case build.ProfileRelease:
		optimize := build.OptimizeSpeed
		if c.Optimize != build.OptimizeUnknown {
			optimize = c.Optimize
		}

		cmd.Args = append(
			cmd.Args,
			"production=yes",
			fmt.Sprintf("optimize=%s", optimize),
		)

	case build.ProfileReleaseDebug:
		optimize := build.OptimizeSpeedTrace
		if c.Optimize != build.OptimizeUnknown {
			optimize = c.Optimize
		}

		cmd.Args = append(
			cmd.Args,
			"debug_symbols=yes",
			"dev_mode=yes",
			fmt.Sprintf("optimize=%s", optimize),
		)
	case build.ProfileDebug:
		optimize := build.OptimizeDebug
		if c.Optimize != build.OptimizeUnknown {
			optimize = c.Optimize
		}

		cmd.Args = append(
			cmd.Args,
			"debug_symbols=yes",
			"dev_mode=yes",
			fmt.Sprintf("optimize=%s", optimize),
		)
	default:
		return nil, fmt.Errorf("%w: profile: %s", ErrInvalidInput, c.Invocation.Profile)
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

	return &cmd, nil
}

/* --------------------------- Method: targetName --------------------------- */

// targetName returns the name of the SCons build target based on the profile.
func (c *Base) targetName() string {
	if c.Invocation.Profile == build.ProfileRelease {
		return "template_release"
	}

	return "template_debug"
}

/* ------------------------ Method: wrapBuildCommand ------------------------ */

// wrapBuildCommand wraps the provided build command in the 'Base' template's
// hook, vendor, and artifact move actions.
func (c *Base) wrapBuildCommand(cmd action.Action) action.Sequence {
	return action.Sequence{
		Pre: action.Commands{
			Commands: c.Hook.Pre,
			Shell:    c.Hook.Shell,
			Verbose:  c.Invocation.Verbose,
		},
		Action: cmd.
			// Vendor the Godot source code prior to executing the build action.
			After(newVendorGodotAction(&c.Godot, &c.Invocation)).
			// Move the generated Godot export template artifacts after  executing the build action.
			AndThen(newMoveArtifactsAction(&c.Invocation)),
		Post: action.Commands{
			Commands: c.Hook.Post,
			Shell:    c.Hook.Shell,
			Verbose:  c.Invocation.Verbose,
		},
	}
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (c *Base) Configure(inv *build.Invocation) error {
	if err := c.PathCustomPy.RelTo(inv.PathManifest); err != nil {
		return err
	}

	for i, m := range c.CustomModules {
		if err := m.RelTo(inv.PathManifest); err != nil {
			return err
		}

		c.CustomModules[i] = m
	}

	if err := c.SCons.Configure(inv); err != nil {
		return err
	}

	if err := c.Godot.Configure(inv); err != nil {
		return err
	}

	return nil
}

/* -------------------------- Impl: build.Validate -------------------------- */

func (c *Base) Validate() error {
	if err := c.Invocation.Validate(); err != nil {
		return err
	}

	if c.Shell != exec.ShellUnknown {
		if _, err := exec.ParseShell(c.Shell.String()); err != nil {
			return fmt.Errorf("%w: unsupported shell: %s", ErrInvalidInput, c.Shell)
		}
	}

	for _, m := range c.CustomModules {
		if err := m.CheckIsDirOrEmpty(); err != nil {
			return err
		}
	}

	if err := c.PathCustomPy.CheckIsFileOrEmpty(); err != nil {
		return err
	}

	if err := c.Godot.Validate(); err != nil {
		return err
	}

	if err := c.SCons.Validate(); err != nil {
		return err
	}

	return nil
}
