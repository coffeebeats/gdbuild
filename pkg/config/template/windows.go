package template

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/build"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Windows                              */
/* -------------------------------------------------------------------------- */

type Windows struct {
	*Base

	// UseMinGW determines whether the MinGW compiler is used.
	UseMinGW bool `toml:"use_mingw"`

	// PathIcon is a path to a Windows application icon.
	PathIcon build.Path `toml:"icon_path"`
}

// Compile-time check that 'Template' is implemented.
var _ Template = (*Windows)(nil)

/* -------------------------- Impl: build.Templater ------------------------- */

func (c *Windows) ToTemplate(g build.Godot, inv build.Invocation) build.Template {
	t := c.Base.ToTemplate(g, inv)

	t.Binaries[0].Platform = build.OSWindows

	if c.Base.Arch == build.ArchUnknown {
		t.Binaries[0].Arch = build.ArchAmd64
	}

	scons := &t.Binaries[0].SCons
	if inv.Profile.IsRelease() {
		scons.ExtraArgs = append(scons.ExtraArgs, "lto=full")
	}

	if c.UseMinGW {
		scons.ExtraArgs = append(scons.ExtraArgs, "use_mingw=yes")
	}

	if c.PathIcon != "" {
		t.AddToPath(c.PathIcon)

		// Copy the icon file to the correct location.
		t.Prebuild = append(t.Prebuild, NewCopyImageFileAction(c.PathIcon, &inv))
	}

	return t
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (c *Windows) Configure(inv build.Invocation) error {
	if err := c.Base.Configure(inv); err != nil {
		return err
	}

	if err := c.PathIcon.RelTo(inv.PathManifest); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *Windows) Validate(inv build.Invocation) error {
	if err := c.Base.Validate(inv); err != nil {
		return err
	}

	if !c.Base.Arch.IsOneOf(build.ArchAmd64, build.ArchI386, build.ArchUnknown) {
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
	pathImage build.Path,
	inv *build.Invocation,
) action.WithDescription[action.Function] {
	pathDst := filepath.Join(inv.PathBuild.String(), "platform/windows/godot.ico")

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
