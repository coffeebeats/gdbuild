package config

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/config/common"
	"github.com/coffeebeats/gdbuild/pkg/config/linux"
	"github.com/coffeebeats/gdbuild/pkg/config/macos"
	"github.com/coffeebeats/gdbuild/pkg/config/windows"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
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
		template Templater
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

/* -------------------------------------------------------------------------- */
/*                            Interface: Templater                            */
/* -------------------------------------------------------------------------- */

type Templater interface {
	config.Configurable[*run.Context]

	Collect(src engine.Source, rc *run.Context) *template.Template
}

/* -------------------------------------------------------------------------- */
/*                              Struct: Templates                             */
/* -------------------------------------------------------------------------- */

// Templates defines the parameters for building a Godot export template for a
// specified platform. A 'Template' definition can be customized based on
// 'feature', 'platform', and 'profile' labels used in the property names. Note
// that each specifier label can only be used once per property name (i.e.
// 'target.profile.release.profile.debug' is not allowed). Additionally, the
// order of specifiers is strict: 'platform' < 'feature' < 'profile'.
//
// For example, the following are all valid table names:
//
//	[template]
//	[template.profile.release]
//	[template.platform.macos.feature.client]
//	[template.platform.linux.feature.server.profile.release_debug]
type Templates struct {
	*common.TemplateWithFeaturesAndProfile

	Platform TemplatePlatforms `toml:"platform"`
}

/* ------------------- Struct: BaseTemplateWithoutFeature ------------------- */

type BaseTemplateWithoutFeature struct {
	*common.Template

	Profile map[engine.Profile]common.Template `toml:"profile"`
}

/* ---------------------------- Struct: Platforms --------------------------- */

type TemplatePlatforms struct {
	Linux   linux.TemplateWithFeaturesAndProfile   `toml:"linux"`
	MacOS   macos.TemplateWithFeaturesAndProfile   `toml:"macos"`
	Windows windows.TemplateWithFeaturesAndProfile `toml:"windows"`
}

/* ----------------------- Interface: TemplateBuilder ----------------------- */

type TemplateBuilder[T Templater] interface {
	Build(rc *run.Context, dst T) error
}

// Compile-time check that 'Builder' is implemented.
var _ TemplateBuilder[*common.Template] = (*common.TemplateWithFeaturesAndProfile)(nil)
var _ TemplateBuilder[*linux.Template] = (*linux.TemplateWithFeaturesAndProfile)(nil)
var _ TemplateBuilder[*macos.Template] = (*macos.TemplateWithFeaturesAndProfile)(nil)
var _ TemplateBuilder[*windows.Template] = (*windows.TemplateWithFeaturesAndProfile)(nil)

/* ----------------------------- Method: Combine ---------------------------- */

func (t *Templates) Combine(rc *run.Context) (Templater, error) { //nolint:dupl,ireturn
	// Root params.
	base := new(common.Template)

	if err := t.TemplateWithFeaturesAndProfile.Build(rc, base); err != nil {
		return nil, err
	}

	switch p := rc.Platform; p {
	case platform.OSLinux:
		out := &linux.Template{Template: base} //nolint:exhaustruct

		if err := t.Platform.Linux.Build(rc, out); err != nil {
			return nil, err
		}

		return out, nil
	case platform.OSMacOS:
		out := &macos.Template{Template: base} //nolint:exhaustruct

		if err := t.Platform.MacOS.Build(rc, out); err != nil {
			return nil, err
		}

		return out, nil
	case platform.OSWindows:
		out := &windows.Template{Template: base} //nolint:exhaustruct

		if err := t.Platform.Windows.Build(rc, out); err != nil {
			return nil, err
		}

		return out, nil
	default:
		return nil, fmt.Errorf("%w: unsupported platform: %s", config.ErrInvalidInput, p)
	}
}
