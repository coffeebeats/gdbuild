package template

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/merge"
	"github.com/coffeebeats/gdbuild/pkg/build"
)

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnimplemented = errors.New("unimplemented")
)

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
	// Shell is a shell name to run the commands with. Defaults to 'sh'.
	Shell build.Shell `toml:"shell"`

	// Invocation contains invocation-specific properties. These must be set
	// manually prior to executing a template build.
	build.Invocation
	// Godot contains a specification of which Godot version to use.
	build.Godot
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *Base) Action() (action.Action, error) { //nolint:cyclop,funlen,ireturn
	var cmd action.Process

	cmd.Directory = string(c.Invocation.PathBuild)

	cmd.Environment = c.Env
	if cmd.Environment == nil {
		cmd.Environment = map[string]string{}
	}

	// Define the SCons cache path.
	if path := c.SCons.PathCache; path != "" {
		cmd.Environment["SCONS_CACHE"] = path.String()
	}

	cmd.Args = append(cmd.Args, "scons")

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

	// Append profile/optimization-related arguments.
	switch c.Invocation.Profile {
	case build.ProfileRelease:
		optimize := build.OptimizeSpeed
		if c.Optimize != build.OptimizeUnknown {
			optimize = c.Optimize
		}

		cmd.Args = append(
			cmd.Args,
			"target=template_release",
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
			"target=template_debug",
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
			"target=template_debug",
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

/* ------------------------- Impl: build.Configurer ------------------------- */

func (c *Base) Configure(inv *build.Invocation) error {
	if err := c.SCons.PathCache.RelTo(inv.PathManifest); err != nil {
		return err
	}

	if err := c.PathCustomPy.RelTo(inv.PathManifest); err != nil {
		return err
	}

	for i, m := range c.CustomModules {
		if err := m.RelTo(inv.PathManifest); err != nil {
			return err
		}

		c.CustomModules[i] = m
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

	if c.Shell != build.ShellUnknown {
		return fmt.Errorf("%w: unsupported shell: %s", ErrInvalidInput, c.Shell)
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

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Base) Merge(other *Base) error {
	if c == nil || other == nil {
		return nil
	}

	c.CustomModules = append(c.CustomModules, other.CustomModules...)

	if err := merge.Pointer(c.DoublePrecision, other.DoublePrecision); err != nil {
		return fmt.Errorf("%w: double_precision", err)
	}

	if err := merge.Map(&c.Env, other.Env); err != nil {
		return fmt.Errorf("%w: env", err)
	}

	if err := merge.Primitive(&c.Optimize, other.Optimize); err != nil {
		return fmt.Errorf("%w: optimize", err)
	}

	if err := merge.Primitive(&c.PathCustomPy, other.PathCustomPy); err != nil {
		return fmt.Errorf("%w: custom_py_path", err)
	}

	if err := c.Hook.Merge(&other.Hook); err != nil {
		return err
	}

	if err := c.SCons.Merge(&other.SCons); err != nil {
		return err
	}

	return nil
}
