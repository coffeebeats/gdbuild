package template

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/build"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Windows                              */
/* -------------------------------------------------------------------------- */

type Windows struct {
	*Base

	// UseMinGW determines whether the MinGW compiler is used.
	UseMinGW *bool `toml:"use_mingw"`

	// PathIcon is a path to a Windows application icon.
	PathIcon osutil.Path `toml:"icon_path"`
}

// Compile-time check that 'Template' is implemented.
var _ Template = (*Windows)(nil)

/* ----------------------------- Impl: Template ----------------------------- */

func (c *Windows) Template(g build.Source, rc *run.Context) *build.Template {
	t := c.Base.Template(g, rc)

	t.Builds[0].Platform = platform.OSWindows

	if c.Base.Arch == platform.ArchUnknown {
		t.Builds[0].Arch = platform.ArchAmd64
	}

	scons := &t.Builds[0].SCons
	if rc.Profile.IsRelease() {
		scons.ExtraArgs = append(scons.ExtraArgs, "lto=full")
	}

	if config.Dereference(c.UseMinGW) {
		scons.ExtraArgs = append(scons.ExtraArgs, "use_mingw=yes")
	}

	if c.PathIcon != "" {
		t.RegisterDependencyPath(c.PathIcon)

		// Copy the icon file to the correct location.
		t.Prebuild = action.InOrder(t.Prebuild, NewCopyImageFileAction(c.PathIcon, rc))
	}

	// Register the additional console artifact.
	t.ExtraArtifacts = append(
		t.ExtraArtifacts,
		strings.TrimSuffix(t.Builds[0].Filename(), ".exe")+".console.exe",
	)

	return t
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (c *Windows) Configure(rc *run.Context) error {
	if err := c.Base.Configure(rc); err != nil {
		return err
	}

	if err := c.PathIcon.RelTo(rc.PathManifest); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *Windows) Validate(rc *run.Context) error {
	if err := c.Base.Validate(rc); err != nil {
		return err
	}

	if !c.Base.Arch.IsOneOf(platform.ArchAmd64, platform.ArchI386, platform.ArchUnknown) {
		return fmt.Errorf("%w: unsupport architecture: %s", config.ErrInvalidInput, c.Base.Arch)
	}

	if err := c.PathIcon.CheckIsFileOrEmpty(); err != nil {
		return err
	}

	return nil
}

/* --------------------------- Impl: config.Merger -------------------------- */

func (c *Windows) MergeInto(other any) error {
	if c == nil || other == nil {
		return nil
	}

	dst, ok := other.(*Windows)
	if !ok {
		return fmt.Errorf(
			"%w: expected a '%T' but was '%T'",
			config.ErrInvalidInput,
			new(Windows),
			other,
		)
	}

	return config.Merge(dst, *c)
}

/* -------------------------------------------------------------------------- */
/*                      Function: NewCopyImageFileAction                      */
/* -------------------------------------------------------------------------- */

// NewCopyImageFileAction creates an 'action.Action' which places the specified
// icon image into the Godot source code.
func NewCopyImageFileAction(
	pathImage osutil.Path,
	rc *run.Context,
) action.WithDescription[action.Function] {
	pathDst := filepath.Join(rc.PathBuild.String(), "platform/windows/godot.ico")

	fn := func(_ context.Context) error {
		dst, err := os.Create(pathDst)
		if err != nil {
			return err
		}

		defer dst.Close()

		src, err := os.Open(pathImage.String())
		if err != nil {
			return err
		}

		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		return nil
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "copy icon into build directory: " + pathDst,
	}
}
