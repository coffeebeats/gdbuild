package config

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/pathutil"
	cfgtemplate "github.com/coffeebeats/gdbuild/pkg/config/template"
	"github.com/coffeebeats/gdbuild/pkg/godot/build"
	"github.com/coffeebeats/gdbuild/pkg/template"
)

var (
	ErrInvalidInput = config.ErrInvalidInput
	ErrMissingInput = config.ErrMissingInput
)

/* -------------------------------------------------------------------------- */
/*                              Struct: Manifest                              */
/* -------------------------------------------------------------------------- */

// Manifest defines the supported structure of the GDBuild manifest file.
type Manifest struct {
	// Config contains GDBuild configuration-related settings.
	Config Config `toml:"config"`
	// Godot contains settings on which Godot version/source code to use.
	Godot build.Source `toml:"godot"`
	// Template includes settings for building custom export templates.
	Template cfgtemplate.Templates `toml:"template"`
}

/* -------------------------- Method: BuildTemplate ------------------------- */

type configuration struct {
	manifest   *Manifest
	invocation *build.Context
}

// BuildTemplate creates a `Template` instance which contains an action for
// compiling Godot based on the specified configuration.
func (m *Manifest) BuildTemplate(bc build.Context) (template.Template, error) { //nolint:cyclop,funlen
	var merged struct {
		source   build.Source
		template cfgtemplate.Template
	}

	toBuild := []configuration{{invocation: &bc, manifest: m}}
	visited := map[pathutil.Path]struct{}{}

	for len(toBuild) > 0 {
		// Remove the next manifest from the queue.
		cfg := toBuild[0]
		toBuild = toBuild[1:]

		inv := *cfg.invocation

		// First, determine whether this manifest extends another one.

		if err := cfg.manifest.Config.Extends.RelTo(bc.Invoke.PathManifest); err != nil {
			return template.Template{}, fmt.Errorf(
				"%w: cannot find inherited manifest: %w",
				config.ErrInvalidInput,
				err,
			)
		}

		extends := cfg.manifest.Config.Extends

		// Skip block below if this manifest has already been "visited".
		if _, ok := visited[extends]; !ok && extends != "" {
			baseManifest, err := ParseFile(extends.String())
			if err != nil {
				return template.Template{}, fmt.Errorf("cannot parse inherited manifest: %w", err)
			}

			bc.Invoke.PathManifest = extends

			base := configuration{invocation: &inv, manifest: baseManifest}
			toBuild = append(toBuild, base, cfg)

			visited[extends] = struct{}{}

			continue
		}

		// Configure 'Godot' properties.
		if err := cfg.manifest.Godot.Configure(bc.Invoke); err != nil {
			return template.Template{}, err
		}

		// Merge 'Godot' properties.
		if err := config.Merge(&merged.source, cfg.manifest.Godot); err != nil {
			return template.Template{}, err
		}

		// Build 'Template' properties.
		t, err := cfg.manifest.Template.Build(bc)
		if err != nil {
			return template.Template{}, err
		}

		// Configure 'Template' properties.
		if err := t.Configure(bc.Invoke); err != nil {
			return template.Template{}, err
		}

		if merged.template == nil {
			merged.template = t

			continue
		}

		// Merge 'Template' properties.
		if err := t.MergeInto(merged.template); err != nil {
			return template.Template{}, err
		}
	}

	if merged.template == nil {
		return template.Template{}, fmt.Errorf("%w: failed to build template", ErrMissingInput)
	}

	// Validate 'Template' properties.
	if err := merged.template.Validate(bc.Invoke); err != nil {
		return template.Template{}, err
	}

	return merged.template.ToTemplate(merged.source, bc), nil
}

/* -------------------------------------------------------------------------- */
/*                               Struct: Config                               */
/* -------------------------------------------------------------------------- */

// Configs specifies GDBuild manifest-related settings.
type Config struct {
	// Extends is a path to another GDBuild manifest to extend. Note that value
	// override rules work the same as within a manifest; any primitive values
	// will override those defined in the base configuration, while arrays will
	// be appended to the base configuration's arrays.
	Extends pathutil.Path `toml:"extends"`
}
