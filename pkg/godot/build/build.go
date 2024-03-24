package build

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
)

const envEncryptionKey = "SCRIPT_AES256_ENCRYPTION_KEY"

/* -------------------------------------------------------------------------- */
/*                                Struct: Build                               */
/* -------------------------------------------------------------------------- */

// Build uniquely specifies a compilation of a Godot export template.
type Build struct {
	// Arch is the CPU architecture of the Godot export template.
	Arch platform.Arch

	// CustomModules is a list of paths to custom modules to include in the
	// template build.
	CustomModules []osutil.Path `hash:"ignore"` // Ignore; paths are separately hashed.

	// CustomPy is a path to a 'custom.py' file which defines export template
	// build options.
	CustomPy osutil.Path `hash:"ignore"` // Ignore; path is separately hashed.

	// DoublePrecision enables double floating-point precision.
	DoublePrecision bool

	// EncryptionKey is an encryption key to embed in the export template.
	//
	// NOTE: While this could just be set in 'Env', exposing it here simplifies
	// setting it externally (i.e. via the 'target' command).
	EncryptionKey string

	// Env is a map of environment variables to set during the build step.
	Env map[string]string

	// Source is the source code specification for the build.
	Source Source

	// Optimize is the level of optimization for the Godot export template.
	Optimize Optimize

	// Platform defines which OS/platform to build for.
	Platform platform.OS

	// Profile is the optimization level of the template.
	Profile Profile

	// SCons contains a specification for how to invoke the compiler.
	SCons SCons
}

/* ------------------------- Function: TemplateName ------------------------- */

// TemplateName returns the base name of the export template defined by the
// specified parameters.
func TemplateName(pl platform.OS, arch platform.Arch, pr Profile) string {
	name := fmt.Sprintf("godot.%s.%s.%s", pl, pr.TargetName(), arch)
	if pl == platform.OSWindows {
		name += ".exe"
	}

	return name
}

/* --------------------- Function: EncryptionKeyFromEnv --------------------- */

// EncryptionKeyFromEnv returns the encryption key set via environment variable.
func EncryptionKeyFromEnv() string {
	return os.Getenv(envEncryptionKey)
}

/* ---------------------------- Method: Filename ---------------------------- */

// Filename returns the base name of the export template generated by this
// 'Binary' specification.
func (b *Build) Filename() string {
	return TemplateName(b.Platform, b.Arch, b.Profile)
}

/* -------------------------- Method: SConsCommand -------------------------- */

// SConsCommand returns the 'SCons' command to build the export template.
func (b *Build) SConsCommand(c *Context) *action.Process { //nolint:cyclop,funlen
	var cmd action.Process

	cmd.Directory = c.PathBuild.String()
	cmd.Verbose = c.Verbose

	scons := b.SCons

	// Define the SCons command, if not yet set.
	if len(scons.Command) == 0 {
		scons.Command = append(scons.Command, "scons")
	}

	// Define the SCons cache path.
	if path := scons.PathCache; path != "" {
		cmd.Environment = append(cmd.Environment, envSConsCache+"="+path.String())
	}

	// Add specified environment variables.
	for k, v := range b.Env {
		cmd.Environment = append(cmd.Environment, k+"="+v)
	}

	// Now pass through all environment variables so that these override
	// previously values.
	cmd.Environment = append(cmd.Environment, os.Environ()...)

	// Set the encryption key on the environment, if one is specified.
	if b.EncryptionKey != "" {
		cmd.Environment = append(cmd.Environment, envEncryptionKey+"="+b.EncryptionKey)
	}

	// Set the SCons cache size limit, if one was set.
	if csl := scons.CacheSizeLimit; csl != nil {
		cmd.Environment = append(
			cmd.Environment,
			fmt.Sprintf("%s=%d", envSConsCacheSizeLimit, *csl),
		)
	}

	// Build the SCons command/argument list.
	var args []string

	// Add multi-core support.
	args = append(args, "-j"+strconv.Itoa(runtime.NumCPU()))

	// Specify the 'platform' argument.
	args = append(args, "platform="+b.Platform.String())

	// Add the achitecture setting (note that this requires the 'platform.Arch'
	// values to match what SCons expects).
	args = append(args, "arch="+b.Arch.String())

	// Specify which target to build.
	args = append(args, "target="+c.Profile.TargetName())

	// Add stricter warning handling.
	args = append(args, "warnings=extra", "werror=yes")

	// Handle a verbose flag.
	if c.Verbose {
		args = append(args, "verbose=yes")
	}

	// Append 'custom_modules' argument.
	if len(b.CustomModules) > 0 {
		modules := make([]string, len(b.CustomModules))
		for i, m := range b.CustomModules {
			modules[i] = m.String()
		}

		args = append(args, fmt.Sprintf(`custom_modules="%s"`, strings.Join(modules, ",")))
	}

	// Append the 'precision' argument.
	if b.DoublePrecision {
		args = append(args, "precision=double")
	}

	// Append profile/optimization-related arguments.
	switch c.Profile {
	case ProfileRelease:
		optimize := OptimizeSpeed
		if b.Optimize != OptimizeUnknown {
			optimize = b.Optimize
		}

		args = append(
			args,
			"production=yes",
			fmt.Sprintf("optimize=%s", optimize),
		)

	case ProfileReleaseDebug:
		optimize := OptimizeSpeedTrace
		if b.Optimize != OptimizeUnknown {
			optimize = b.Optimize
		}

		args = append(
			args,
			"debug_symbols=yes",
			"dev_mode=yes",
			fmt.Sprintf("optimize=%s", optimize),
		)
	default: // ProfileDebug
		optimize := OptimizeDebug
		if b.Optimize != OptimizeUnknown {
			optimize = b.Optimize
		}

		args = append(
			args,
			"debug_symbols=yes",
			"dev_mode=yes",
			fmt.Sprintf("optimize=%s", optimize),
		)
	}

	// Append C/C++ build flags.
	if len(b.SCons.CCFlags) > 0 {
		flags := fmt.Sprintf(`CCFLAGS="%s"`, strings.Join(b.SCons.CCFlags, " "))
		args = append(args, flags)
	}

	if len(b.SCons.CFlags) > 0 {
		flags := fmt.Sprintf(`CFLAGS="%s"`, strings.Join(b.SCons.CFlags, " "))
		args = append(args, flags)
	}

	if len(b.SCons.CXXFlags) > 0 {
		flags := fmt.Sprintf(`CXXFLAGS="%s"`, strings.Join(b.SCons.CXXFlags, " "))
		args = append(args, flags)
	}

	if len(b.SCons.LinkFlags) > 0 {
		flags := fmt.Sprintf(`LINKFLAGS="%s"`, strings.Join(b.SCons.LinkFlags, " "))
		args = append(args, flags)
	}

	// Append extra arguments.
	args = append(args, b.SCons.ExtraArgs...)

	// Attach the command with arguments to the action.
	cmd.Args = append(cmd.Args, scons.Command...)
	cmd.Args = append(cmd.Args, args...)

	return &cmd
}
