package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/pkg/build"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Windows                              */
/* -------------------------------------------------------------------------- */

// Windows contains 'Windows'-specific settings for constructing a custom Godot
// export template.
type Windows struct {
	*Base

	// UseMinGW determines whether the MinGW compiler is used.
	UseMinGW bool `toml:"use_mingw"`

	// PathIcon is a path to a Windows application icon.
	PathIcon build.Path `toml:"icon_path"`
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *Windows) Action() (action.Action, error) { //nolint:ireturn
	cmd, err := c.Base.action()
	if err != nil {
		return nil, err
	}

	cmd.Args = append(cmd.Args, "platform="+build.OSWindows.String())

	if c.UseMinGW {
		cmd.Args = append(cmd.Args, "use_mingw=yes")

		if c.Base.Invocation.Profile.IsRelease() {
			cmd.Args = append(cmd.Args, "lto=full")
		}
	}

	return c.wrapBuildCommand(cmd), nil
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (c *Windows) Configure(inv *build.Invocation) error {
	if err := c.Base.Configure(inv); err != nil {
		return err
	}

	if c.Base.Arch == build.ArchUnknown {
		c.Base.Arch = build.ArchAmd64
	}

	if err := c.PathIcon.RelTo(inv.PathManifest); err != nil {
		return err
	}

	return nil
}

/* -------------------------- Impl: build.Validate -------------------------- */

func (c *Windows) Validate() error {
	if err := c.Base.Validate(); err != nil {
		return err
	}

	switch c.Base.Arch {
	case build.ArchAmd64, build.ArchI386:
	case build.ArchUnknown:
	default:
		return fmt.Errorf("%w: unsupport architecture: %s", ErrInvalidInput, c.Base.Arch)
	}

	if err := c.PathIcon.CheckIsFileOrEmpty(); err != nil {
		return err
	}

	return nil
}
