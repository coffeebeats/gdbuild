package config

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/config/platform"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                             Function: Template                             */
/* -------------------------------------------------------------------------- */

// Template creates a `Template` instance which contains an action for
// compiling Godot based on the specified configuration.
func Template(rc *run.Context, m *Manifest) (*template.Template, error) { //nolint:cyclop,funlen
	var merged struct {
		godot    Godot
		template platform.Templater
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

		// Build 'Template' properties.
		t, err := cfg.manifest.Template.Combine(&rc)
		if err != nil {
			return nil, err
		}

		// Configure 'Template' properties.
		if err := t.Configure(&rc); err != nil {
			return nil, err
		}

		if merged.template == nil {
			merged.template = t

			continue
		}

		// Merge 'Template' properties.
		if err := t.MergeInto(merged.template); err != nil {
			return nil, err
		}
	}

	if merged.template == nil {
		return nil, fmt.Errorf("%w: failed to build template", ErrMissingInput)
	}

	// Validate 'Template' properties.
	if err := merged.godot.Validate(rc); err != nil {
		return nil, err
	}

	if err := merged.template.Validate(rc); err != nil {
		return nil, err
	}

	return merged.template.Collect(*merged.godot.Source, rc), nil
}
