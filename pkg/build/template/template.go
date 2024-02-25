package template

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/coffeebeats/gdbuild/internal/command"
	"github.com/coffeebeats/gdbuild/internal/merge"
	"github.com/coffeebeats/gdbuild/pkg/build"
)

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnimplemented = errors.New("unimplemented")
)

/* -------------------------------------------------------------------------- */
/*                            Interface: Commander                            */
/* -------------------------------------------------------------------------- */

type Commander interface {
	Command(pl build.OS, pr build.Profile, ff ...string) ([]string, error)
}

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
	CustomModules []string `toml:"custom_modules"`
	// DoublePrecision enables double floating-point precision.
	DoublePrecision *bool `toml:"double_precision"`
	// Env is a map of environment variables to set during the build step.
	Env map[string]string `toml:"env"`
	// Hook defines commands to be run before or after a build step.
	Hook build.Hook `json:"hook" toml:"hook"`
	// Optimize is the specific optimization level for the template.
	Optimize build.Optimize `toml:"optimize"`

	// PathCustomPy is a path to a 'custom.py' file which defines export
	// template build options. Defaults to 'custom.py', but will be ignored if
	// one isn't found.
	PathCustomPy string `toml:"custom_py_path"`
	// PathGodotSource is the path to the Godot source code.
	PathGodotSource string `toml:"godot_src_path"`
	// SCons contains build command-related settings.
	SCons SCons `json:"scons" toml:"scons"`

	// Execution contains invocation-specific properties. These must be set
	// manually prior to executing a template build.
	Execution Execution
}

/* ------------------------- Impl: command.Commander ------------------------ */

func (c *Base) Command() (*command.Command, error) { //nolint:cyclop,funlen
	var cmd command.Command

	cmd.Directory = c.Execution.PathBuild
	cmd.Environment = c.Env

	if c.SCons.PathCache != "" {
		path, err := c.Execution.makeAbsolute(c.SCons.PathCache)
		if err != nil {
			return nil, err
		}

		cmd.Environment["SCONS_CACHE"] = path
	}

	cmd.Shell = command.ShellSh
	if c.Execution.Shell != command.ShellUnknown {
		cmd.Shell = c.Execution.Shell
	}

	cmd.Args = append(cmd.Args, "scons")

	modules := make([]string, 0, len(c.CustomModules))

	for _, m := range c.CustomModules {
		path, err := c.Execution.makeAbsolute(m)
		if err != nil {
			return nil, err
		}

		modules = append(modules, path)
	}

	if len(modules) > 0 {
		cmd.Args = append(cmd.Args, fmt.Sprintf(`custom_modules="%s"`, strings.Join(modules, ",")))
	}

	if c.DoublePrecision != nil && *c.DoublePrecision {
		cmd.Args = append(cmd.Args, "precision=double")
	}

	if path := c.PathCustomPy; path != "" {
		path, err := c.Execution.makeAbsolute(path)
		if err != nil {
			return nil, err
		}

		cmd.Args = append(cmd.Args, "custom="+path)
	}

	switch c.Execution.Profile {
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
		return nil, fmt.Errorf("%w: profile: %s", ErrInvalidInput, c.Execution.Profile)
	}

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

	cmd.Args = append(cmd.Args, c.SCons.ExtraArgs...)

	return &cmd, nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Base) Merge(other *Base) error {
	if other == nil {
		return nil
	}

	if c == nil {
		*c = *other

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

	if err := merge.Primitive(&c.PathGodotSource, other.PathGodotSource); err != nil {
		return fmt.Errorf("%w: godot_src_path", err)
	}

	if err := c.Hook.Merge(&other.Hook); err != nil {
		return err
	}

	if err := c.SCons.Merge(&other.SCons); err != nil {
		return err
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*                              Struct: Execution                             */
/* -------------------------------------------------------------------------- */

// Execution are the template build inputs that are execution-specific. These
// need to be explicitly set per invocation as they can't be parsed from a
// GDBuild manifest.
type Execution struct {
	// Features is the list of feature tags to enable.
	Features []string
	// Platform is the target platform to build for.
	Platform build.OS
	// Profile is the GDBuild optimization level to build with.
	Profile build.Profile

	// PathManifest is the directory in which the GDBuild manifest is located.
	// This is used to locate relative paths in various other properties.
	PathManifest string
	// PathBuild is the directory in which to build the template in. All input
	// artifacts will be copied here and the SCons build command will be
	// executed within this directory. Defaults to a temporary directory.
	PathBuild string

	// Shell is the name of the shell to build the template with. Defaults to
	// 'command.ShellSh'.
	Shell command.Shell
}

/* -------------------------- Method: makeAbsolute -------------------------- */

// makeAbsolute converts the provided path into an absolute path, resolving any
// relative paths against the configured manifest path.
func (c *Execution) makeAbsolute(path string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}

	if c.PathManifest == "" {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		c.PathManifest = wd
	}

	path, err := filepath.Rel(c.PathManifest, path)
	if err != nil {
		return "", err
	}

	return path, nil
}

/* -------------------------------------------------------------------------- */
/*                                Struct: SCons                               */
/* -------------------------------------------------------------------------- */

// SCons defines options and settings for use with the _Godot_ build system.
type SCons struct {
	// CCFlags are additional 'CFLAGS' to append to the SCons build command.
	// Note that 'CCFLAGS=...' will be appended *before* 'ExtraArgs'.
	CCFlags []string `toml:"ccflags"`
	// CFlags are additional 'CFLAGS' to append to the SCons build command. Note
	// that 'CFLAGS=...' will be appended *before* 'ExtraArgs'.
	CFlags []string `toml:"cflags"`
	// CXXFlags are additional 'CXXFLAGS' to append to the SCons build command.
	// Note that 'CXXFLAGS=...' will be appended *before* 'ExtraArgs'.
	CXXFlags []string `toml:"cxxflags"`
	// ExtraArgs are additional arguments to append to the SCons build command.
	ExtraArgs []string `toml:"extra_args"`
	// PathCache is the path to the SCons cache, relative to the manifest.
	PathCache string `toml:"cache_path"`
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *SCons) Merge(other *SCons) error {
	if other == nil {
		return nil
	}

	if c == nil {
		*c = *other

		return nil
	}

	c.CCFlags = append(c.CCFlags, other.CCFlags...)
	c.CFlags = append(c.CFlags, other.CFlags...)
	c.CXXFlags = append(c.CXXFlags, other.CXXFlags...)
	c.ExtraArgs = append(c.ExtraArgs, other.ExtraArgs...)

	if err := merge.Primitive(&c.PathCache, other.PathCache); err != nil {
		return fmt.Errorf("%w: cache_path", err)
	}

	return nil
}
