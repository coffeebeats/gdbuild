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
	"github.com/coffeebeats/gdbuild/internal/pathutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/compile"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
	"github.com/coffeebeats/gdbuild/pkg/template"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Windows                              */
/* -------------------------------------------------------------------------- */

type Windows struct {
	*Base

	// UseMinGW determines whether the MinGW compiler is used.
	UseMinGW *bool `toml:"use_mingw"`

	// PathIcon is a path to a Windows application icon.
	PathIcon pathutil.Path `toml:"icon_path"`
}

// Compile-time check that 'Template' is implemented.
var _ Template = (*Windows)(nil)

/* -------------------------- Impl: template.Templater ------------------------- */

func (c *Windows) ToTemplate(g compile.Godot, cc compile.Context) template.Template {
	t := c.Base.ToTemplate(g, cc)

	t.Binaries[0].Platform = platform.OSWindows

	if c.Base.Arch == platform.ArchUnknown {
		t.Binaries[0].Arch = platform.ArchAmd64
	}

	scons := &t.Binaries[0].SCons
	if cc.Profile.IsRelease() {
		scons.ExtraArgs = append(scons.ExtraArgs, "lto=full")
	}

	if config.Dereference(c.UseMinGW) {
		scons.ExtraArgs = append(scons.ExtraArgs, "use_mingw=yes")
	}

	if c.PathIcon != "" {
		t.AddToPaths(c.PathIcon)

		// Copy the icon file to the correct location.
		t.Prebuild = action.InOrder(t.Prebuild, NewCopyImageFileAction(c.PathIcon, &cc.Invoke))
	}

	// Register the additional console artifact.
	t.ExtraArtifacts = append(
		t.ExtraArtifacts,
		strings.TrimSuffix(t.Binaries[0].Filename(), ".exe")+".console.exe",
	)

	return t
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (c *Windows) Configure(cc config.Context) error {
	if err := c.Base.Configure(cc); err != nil {
		return err
	}

	if err := c.PathIcon.RelTo(cc.PathManifest); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *Windows) Validate(cc config.Context) error {
	if err := c.Base.Validate(cc); err != nil {
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
	pathImage pathutil.Path,
	cc *config.Context,
) action.WithDescription[action.Function] {
	pathDst := filepath.Join(cc.PathBuild.String(), "platform/windows/godot.ico")

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
