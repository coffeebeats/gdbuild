package template

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/build"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: Linux                               */
/* -------------------------------------------------------------------------- */

type Linux struct {
	*Base

	// UseLLVM determines whether the LLVM compiler is used.
	UseLLVM *bool `toml:"use_llvm"`
}

// Compile-time check that 'Template' is implemented.
var _ Template = (*Linux)(nil)

/* -------------------------- Impl: build.Templater ------------------------- */

func (c *Linux) ToTemplate(g build.Godot, inv build.Invocation) build.Template {
	t := c.Base.ToTemplate(g, inv)

	t.Binaries[0].Platform = platform.OSLinux

	if c.Base.Arch == platform.ArchUnknown {
		t.Binaries[0].Arch = platform.ArchAmd64
	}

	scons := &t.Binaries[0].SCons
	if config.Dereference(c.UseLLVM) {
		scons.ExtraArgs = append(scons.ExtraArgs, "use_llvm=yes")
	} else if inv.Profile.IsRelease() { // Only valid with GCC.
		scons.ExtraArgs = append(scons.ExtraArgs, "lto=full")
	}

	return t
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (c *Linux) Configure(inv build.Invocation) error {
	if err := c.Base.Configure(inv); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *Linux) Validate(inv build.Invocation) error {
	if err := c.Base.Validate(inv); err != nil {
		return err
	}

	if !c.Base.Arch.IsOneOf(platform.ArchI386, platform.ArchAmd64, platform.ArchUnknown) {
		return fmt.Errorf("%w: unsupport architecture: %s", config.ErrInvalidInput, c.Base.Arch)
	}

	switch c.Base.Arch {
	case platform.ArchI386, platform.ArchAmd64:
	case platform.ArchUnknown:
	default:
		return fmt.Errorf("%w: unsupport architecture: %s", config.ErrInvalidInput, c.Base.Arch)
	}

	return nil
}

/* --------------------------- Impl: config.Merger -------------------------- */

func (c *Linux) MergeInto(other any) error {
	if c == nil || other == nil {
		return nil
	}

	dst, ok := other.(*Linux)
	if !ok {
		return fmt.Errorf(
			"%w: expected a '%T' but was '%T'",
			config.ErrInvalidInput,
			new(Linux),
			other,
		)
	}

	return config.Merge(dst, *c)
}
