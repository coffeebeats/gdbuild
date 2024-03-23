package template

import (
	"errors"
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/config"
	"github.com/coffeebeats/gdbuild/pkg/config/template"
	"github.com/coffeebeats/gdbuild/pkg/godot/build"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrMissingInput = errors.New("missing input")
)

/* -------------------------------------------------------------------------- */
/*                               Function: Build                              */
/* -------------------------------------------------------------------------- */

// Build creates a `Template` instance which contains an action for compiling
// Godot based on the specified configuration.
func Build(m *config.Manifest, bc *build.Context) (*build.Template, error) { //nolint:cyclop,funlen
	var merged struct {
		source   build.Source
		template template.Template
	}

	toBuild := []configuration{{context: bc, manifest: m}}
	visited := map[osutil.Path]struct{}{}

	for len(toBuild) > 0 {
		// Remove the next manifest from the queue.
		cfg := toBuild[0]
		toBuild = toBuild[1:]

		// Copy build context so it can be modified.
		bc := *cfg.context

		// First, determine whether this manifest extends another one.

		if err := cfg.manifest.Config.Extends.RelTo(bc.PathManifest); err != nil {
			return nil, fmt.Errorf(
				"%w: cannot find inherited manifest: %w",
				ErrInvalidInput,
				err,
			)
		}

		extends := cfg.manifest.Config.Extends

		// Skip block below if this manifest has already been "visited".
		if _, ok := visited[extends]; !ok && extends != "" {
			baseManifest, err := config.ParseFile(extends.String())
			if err != nil {
				return nil, fmt.Errorf("cannot parse inherited manifest: %w", err)
			}

			bc.PathManifest = extends

			base := configuration{context: &bc, manifest: baseManifest}
			toBuild = append(toBuild, base, cfg)

			visited[extends] = struct{}{}

			continue
		}

		// Configure 'Godot' properties.
		if err := cfg.manifest.Godot.Configure(&bc); err != nil {
			return nil, err
		}

		// Merge 'Godot' properties.
		if err := cfg.manifest.Godot.MergeInto(&merged.source); err != nil {
			return nil, err
		}

		// Build 'Template' properties.
		t, err := cfg.manifest.Template.Build(&bc)
		if err != nil {
			return nil, err
		}

		// Configure 'Template' properties.
		if err := t.Configure(&bc); err != nil {
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
	if err := merged.source.Validate(bc); err != nil {
		return nil, err
	}

	if err := merged.template.Validate(bc); err != nil {
		return nil, err
	}

	return merged.template.Template(merged.source, bc), nil
}

/* -------------------------- Struct: configuration ------------------------- */

type configuration struct {
	manifest *config.Manifest
	context  *build.Context
}
