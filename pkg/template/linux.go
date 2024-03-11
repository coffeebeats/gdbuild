package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/pkg/build"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: Linux                               */
/* -------------------------------------------------------------------------- */

// Linux contains 'linux'-specific settings for constructing a custom Godot
// export template.
type Linux struct {
	*Base

	// UseLLVM determines whether the LLVM compiler is used.
	UseLLVM bool `toml:"use_llvm"`
}

/* -------------------------- Impl: action.Actioner ------------------------- */

func (c *Linux) Action() (action.Action, error) { //nolint:ireturn
	cmd, err := c.Base.action()
	if err != nil {
		return nil, err
	}

	cmd.Args = append(cmd.Args, "platform="+build.OSLinux.String()) //nolint:goconst

	if c.UseLLVM {
		cmd.Args = append(cmd.Args, "use_llvm=yes")
	} else if c.Base.Invocation.Profile.IsRelease() { // Only valid with GCC.
		cmd.Args = append(cmd.Args, "lto=full")
	}

	return c.wrapBuildCommand(cmd), nil
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (c *Linux) Configure(inv *build.Invocation) error {
	if err := c.Base.Configure(inv); err != nil {
		return err
	}

	if c.Base.Arch == build.ArchUnknown {
		c.Base.Arch = build.ArchAmd64
	}

	return nil
}

/* -------------------------- Impl: build.Validate -------------------------- */

func (c *Linux) Validate() error {
	if err := c.Base.Validate(); err != nil {
		return err
	}

	switch c.Base.Arch {
	case build.ArchI386, build.ArchAmd64:
	case build.ArchUnknown:
	default:
		return fmt.Errorf("%w: unsupport architecture: %s", ErrInvalidInput, c.Base.Arch)
	}

	return nil
}
