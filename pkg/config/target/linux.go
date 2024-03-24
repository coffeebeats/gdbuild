package target

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/godot/build"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: Linux                               */
/* -------------------------------------------------------------------------- */

type Linux struct {
	*Base
}

/* ----------------------------- Impl: Exporter ----------------------------- */

func (b *Linux) Export(_ build.Source, _ *build.Context) *export.Export {
	return nil
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (b *Linux) Configure(_ *build.Context) error {
	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (b *Linux) Validate(_ *build.Context) error {
	return nil
}

/* --------------------------- Impl: config.Merger -------------------------- */

func (b *Linux) MergeInto(other any) error {
	if b == nil || other == nil {
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

	return config.Merge(dst, *b)
}
