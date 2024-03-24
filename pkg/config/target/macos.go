package target

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: MacOS                               */
/* -------------------------------------------------------------------------- */

type MacOS struct {
	*Base
}

/* ----------------------------- Impl: Exporter ----------------------------- */

func (b *MacOS) Export(_ engine.Source, _ *run.Context) *export.Export {
	return nil
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (b *MacOS) Configure(_ *run.Context) error {
	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (b *MacOS) Validate(_ *run.Context) error {
	return nil
}

/* --------------------------- Impl: config.Merger -------------------------- */

func (b *MacOS) MergeInto(other any) error {
	if b == nil || other == nil {
		return nil
	}

	dst, ok := other.(*MacOS)
	if !ok {
		return fmt.Errorf(
			"%w: expected a '%T' but was '%T'",
			config.ErrInvalidInput,
			new(MacOS),
			other,
		)
	}

	return config.Merge(dst, *b)
}
