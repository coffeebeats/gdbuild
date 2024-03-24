package platform //nolint:dupl

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/config/platform/common"
	"github.com/coffeebeats/gdbuild/pkg/config/platform/linux"
	"github.com/coffeebeats/gdbuild/pkg/config/platform/macos"
	"github.com/coffeebeats/gdbuild/pkg/config/platform/windows"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

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

func (t *Templates) Combine(rc *run.Context) (Templater, error) { //nolint:ireturn
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
