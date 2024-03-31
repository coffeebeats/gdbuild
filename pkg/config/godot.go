package config

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

var ErrConflictingValue = errors.New("conflicting setting")

type Godot struct {
	*engine.Source
}

/* ---------------------------- config.Configurer --------------------------- */

func (g *Godot) Configure(rc *run.Context) error {
	if g.Source == nil {
		return nil
	}

	if err := g.PathSource.RelTo(rc.PathManifest); err != nil {
		return err
	}

	if err := g.VersionFile.RelTo(rc.PathManifest); err != nil {
		return err
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (g *Godot) Validate(_ *run.Context) error {
	if g.Source == nil || g.IsEmpty() {
		return fmt.Errorf("%w: no Godot version specified in manifest", ErrMissingInput)
	}

	if g.PathSource != "" {
		if !g.Version.IsZero() || g.VersionFile != "" {
			return fmt.Errorf(
				"%w: cannot specify 'version' or 'version_file' with 'src_path'",
				ErrConflictingValue,
			)
		}

		if err := g.PathSource.CheckIsDir(); err != nil {
			// NOTE: A hook might generate this file, so don't return an error.
			if !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		}

		return nil
	}

	if _, err := g.ParseVersion(); err != nil {
		// NOTE: A hook might generate this file, so don't return an error.
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		return err
	}

	return nil
}

/* --------------------------- Impl: config.Merger -------------------------- */

func (g *Godot) MergeInto(other any) error {
	if g == nil || other == nil {
		return nil
	}

	dst, ok := other.(*Godot)
	if !ok {
		return fmt.Errorf(
			"%w: expected a '%T' but was '%T'",
			config.ErrInvalidInput,
			new(Godot),
			other,
		)
	}

	return config.Merge(dst, *g)
}
