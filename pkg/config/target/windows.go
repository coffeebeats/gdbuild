package target

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Windows                              */
/* -------------------------------------------------------------------------- */

type Windows struct {
	*Base
}

/* ----------------------------- Impl: Exporter ----------------------------- */

func (b *Windows) Export(_ engine.Source, _ *run.Context) *export.Export {
	return nil
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (b *Windows) Configure(_ *run.Context) error {
	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (b *Windows) Validate(_ *run.Context) error {
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
