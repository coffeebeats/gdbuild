package macos

import (
	"errors"
	"fmt"
	"slices"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/exec"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/config/common"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

var ErrInvalidInput = errors.New("invalid input")

/* -------------------------------------------------------------------------- */
/*                              Struct: Template                              */
/* -------------------------------------------------------------------------- */

type Template struct {
	*common.Template

	// LipoCommand contains arguments used to invoke 'lipo'. Defaults to
	// ["lipo"]. Only used if 'arch' is set to 'platform.ArchUniversal'.
	LipoCommand []string `toml:"lipo_command"`

	// Vulkan defines Vulkan-related configuration.
	Vulkan Vulkan `toml:"vulkan"`
}

/* ----------------------------- Impl: Template ----------------------------- */

func (t *Template) Collect(g engine.Source, rc *run.Context) *template.Template { //nolint:funlen
	switch a := t.Arch; a {
	case platform.ArchAmd64, platform.ArchArm64:
		out := t.Template.Collect(g, rc)

		out.Arch = t.Arch
		out.Builds[0].Platform = platform.OSMacOS

		scons := &out.Builds[0].SCons
		if config.Dereference(t.Vulkan.Dynamic) {
			scons.ExtraArgs = append(scons.ExtraArgs, "use_volk=yes")
		} else {
			scons.ExtraArgs = append(scons.ExtraArgs, "use_volk=no")
		}

		if t.Vulkan.PathSDK != "" {
			scons.ExtraArgs = append(scons.ExtraArgs, "vulkan_sdk_path="+t.Vulkan.PathSDK.String())
			out.RegisterDependencyPath(t.Vulkan.PathSDK)
		}

		return out
	case platform.ArchUniversal, platform.ArchUnknown:
		// First, create the 'x86_64' binary.
		amd64 := *t
		amd64.Template.Arch = platform.ArchAmd64

		templateAmd64 := amd64.Collect(g, rc)

		// Next, create the 'arm64' binary.
		arm64 := *t
		arm64.Template.Arch = platform.ArchArm64

		templateArm64 := arm64.Collect(g, rc)

		// Finally, merge the two binaries together.

		lipo := t.LipoCommand
		if len(lipo) == 0 {
			lipo = append(lipo, "lipo")
		}

		templateNameUniversal := template.Name(
			platform.OSMacOS,
			platform.ArchUniversal,
			rc.Profile,
			config.Dereference(t.DoublePrecision),
		)

		cmdLipo := &action.Process{
			Directory:   rc.BinPath().String(),
			Environment: nil,

			Shell:   exec.DefaultShell(),
			Verbose: rc.Verbose,

			Args: append(
				lipo,
				"-create",
				templateAmd64.Builds[0].Filename(),
				templateArm64.Builds[0].Filename(),
				"-output",
				templateNameUniversal,
			),
		}

		// Construct the output 'Template'. This is because nothing else needs
		// to be copied over from the arch-specific templates and this avoid the
		// need to deduplicate properties.
		out := t.Template.Collect(g, rc)

		out.Arch = platform.ArchUniversal
		out.Name = templateNameUniversal

		// Register the additional artifact.
		out.ExtraArtifacts = append(out.ExtraArtifacts, templateNameUniversal)

		out.Builds = []template.Build{templateAmd64.Builds[0], templateArm64.Builds[0]}
		out.Postbuild = cmdLipo.AndThen(out.Postbuild)

		// Construct a list of paths with duplicates removed. This is preferred
		// over duplicating the code used to decide which paths are dependencies.
		paths := make([]osutil.Path, 0, len(templateAmd64.Paths)+len(templateArm64.Paths))
		paths = append(paths, templateAmd64.Paths...)
		paths = append(paths, templateArm64.Paths...)
		slices.Sort(paths)
		paths = slices.Compact(paths)

		out.Paths = paths

		return out

	default:
		panic(fmt.Errorf("%w: unsupported architecture: %s", ErrInvalidInput, a))
	}
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (t *Template) Configure(rc *run.Context) error {
	if err := t.Template.Configure(rc); err != nil {
		return err
	}

	if err := t.Vulkan.Configure(rc); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (t *Template) Validate(rc *run.Context) error {
	if err := t.Template.Validate(rc); err != nil {
		return err
	}

	if !t.Arch.IsOneOf(
		platform.ArchAmd64,
		platform.ArchArm64,
		platform.ArchUniversal,
		platform.ArchUnknown,
	) {
		return fmt.Errorf("%w: unsupport architecture: %s", config.ErrInvalidInput, t.Arch)
	}

	// NOTE: Don't check for 'lipo', that should be a runtime check.

	if err := t.Vulkan.Validate(rc); err != nil {
		return err
	}

	return nil
}

/* --------------------------- Impl: config.Merger -------------------------- */

func (t *Template) MergeInto(other any) error {
	if t == nil || other == nil {
		return nil
	}

	dst, ok := other.(*Template)
	if !ok {
		return fmt.Errorf(
			"%w: expected a '%T' but was '%T'",
			config.ErrInvalidInput,
			new(Template),
			other,
		)
	}

	return config.Merge(dst, *t)
}

/* -------------------------------------------------------------------------- */
/*                               Struct: Vulkan                               */
/* -------------------------------------------------------------------------- */

// Vulkan defines the settings required by the MacOS template for including
// Vulkan support.
type Vulkan struct {
	// Dynamic enables dynamically linking Vulkan to the template.
	Dynamic *bool `toml:"use_volk"`

	// PathSDK is the path to the Vulkan SDK root.
	PathSDK osutil.Path `toml:"sdk_path"`
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (c *Vulkan) Configure(rc *run.Context) error {
	if err := c.PathSDK.RelTo(rc.PathManifest); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *Vulkan) Validate(_ *run.Context) error {
	if err := c.PathSDK.CheckIsDir(); err != nil {
		return fmt.Errorf("%w: missing path to Vulkan SDK", err)
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*                   Struct: TemplateWithFeaturesAndProfile                   */
/* -------------------------------------------------------------------------- */

type TemplateWithFeaturesAndProfile struct {
	*Template

	Feature map[string]TemplateWithProfile `toml:"feature"`
	Profile map[engine.Profile]Template    `toml:"profile"`
}

/* ----------------------- Struct: TemplateWithProfile ---------------------- */

type TemplateWithProfile struct {
	*Template

	Profile map[engine.Profile]Template `toml:"profile"`
}

/* --------------------- Impl: platform.templateBuilder --------------------- */

func (t *TemplateWithFeaturesAndProfile) Build(rc *run.Context, dst *Template) error {
	if t == nil {
		return nil
	}

	// Root-level params
	if err := t.Template.MergeInto(dst); err != nil {
		return err
	}

	// Feature-constrained params
	for _, f := range rc.Features {
		if err := t.Feature[f].Template.MergeInto(dst); err != nil {
			return err
		}
	}

	// Profile-constrained params
	l := t.Profile[rc.Profile]
	if err := l.MergeInto(dst); err != nil {
		return err
	}

	// Feature-and-profile-constrained params
	for _, f := range rc.Features {
		l := t.Feature[f].Profile[rc.Profile]
		if err := l.MergeInto(dst); err != nil {
			return err
		}
	}

	return nil
}
