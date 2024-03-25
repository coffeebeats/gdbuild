package config

import (
	"errors"
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/config/platform"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                              Function: Export                              */
/* -------------------------------------------------------------------------- */

// Export creates an `Export` instance which contains an action for exporting
// the specified target.
func Export(rc *run.Context, m *Manifest, target string) (*export.Export, error) { //nolint:cyclop,funlen
	var merged struct {
		godot  Godot
		target platform.Exporter
	}

	tl, err := Template(rc, m)
	if err != nil {
		return nil, err
	}

	toBuild := []configuration{{context: rc, manifest: m}}
	visited := map[osutil.Path]struct{}{}

	for len(toBuild) > 0 {
		// Remove the next manifest from the queue.
		cfg := toBuild[0]
		toBuild = toBuild[1:]

		// Copy build context so it can be modified.
		rc := *cfg.context

		// First, determine whether this manifest extends another one.

		if err := cfg.manifest.Config.Extends.RelTo(rc.PathManifest); err != nil {
			return nil, fmt.Errorf(
				"%w: cannot find inherited manifest: %w",
				ErrInvalidInput,
				err,
			)
		}

		extends := cfg.manifest.Config.Extends

		// Skip block below if this manifest has already been "visited".
		if _, ok := visited[extends]; !ok && extends != "" {
			baseManifest, err := ParseFile(extends.String())
			if err != nil {
				return nil, fmt.Errorf("cannot parse inherited manifest: %w", err)
			}

			rc.PathManifest = extends

			base := configuration{context: &rc, manifest: baseManifest}
			toBuild = append(toBuild, base, cfg)

			visited[extends] = struct{}{}

			continue
		}

		// Configure 'Godot' properties.
		if err := cfg.manifest.Godot.Configure(&rc); err != nil {
			return nil, err
		}

		// Merge 'Godot' properties.
		if err := cfg.manifest.Godot.MergeInto(&merged.godot); err != nil {
			return nil, err
		}

		// Build 'Target' properties.
		t, err := cfg.manifest.Target[target].Combine(&rc)
		if err != nil {
			return nil, err
		}

		// Configure 'Target' properties.
		if err := t.Configure(&rc); err != nil {
			return nil, err
		}

		if merged.target == nil {
			merged.target = t

			continue
		}

		// Merge 'Target' properties.
		if err := t.MergeInto(merged.target); err != nil {
			return nil, err
		}
	}

	if merged.target == nil {
		return nil, fmt.Errorf("%w: failed to build target", ErrMissingInput)
	}

	// Validate 'Target' properties.
	if err := merged.godot.Validate(rc); err != nil {
		return nil, err
	}

	if err := merged.target.Validate(rc); err != nil {
		return nil, err
	}

	ev, err := merged.godot.ParseVersion()
	if err != nil {
		if errors.Is(err, ErrConflictingValue) {
			return nil, fmt.Errorf("%w: 'src_path' is unsupported at this time", err)
		}

		return nil, err
	}

	return merged.target.Collect(rc, tl, ev), nil
}
