package target

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/godot/build"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Windows                              */
/* -------------------------------------------------------------------------- */

type Windows struct {
	*Base
}

/* ----------------------------- Impl: Exporter ----------------------------- */

func (b *Windows) Export(_ build.Source, _ *build.Context) *export.Export {
	return nil
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (b *Windows) Configure(_ *build.Context) error {
	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (b *Windows) Validate(_ *build.Context) error {
	return nil
}

/* --------------------------- Impl: config.Merger -------------------------- */

func (b *Windows) MergeInto(other any) error {
	if b == nil || other == nil {
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

	return config.Merge(dst, *b)
}
