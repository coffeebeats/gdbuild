package manifest //nolint:dupl

import (
	"github.com/coffeebeats/gdbuild/pkg/build"
	"github.com/coffeebeats/gdbuild/pkg/platform"
)

/* -------------------------------------------------------------------------- */
/*                              Struct: Template                              */
/* -------------------------------------------------------------------------- */

// Template defines the parameters for building a Godot export template for a
// specified platform. A 'Template' definition can be customized based on
// 'feature', 'platform', and 'profile' labels used in the property names. Note
// that each specifier label can only be used once per property name (i.e.
// 'target.profile.release.profile.debug' is not allowed).
//
// For example, the following are all valid table names:
//
//	[template]
//	[template.feature.client]
//	[template.platform.windows]
//	[template.profile.release]
//	[template.profile.release.platform.macos.feature.client]
type Template struct {
	*build.Template

	Feature  map[string]*TemplateWithoutFeature        `toml:"feature"`
	Platform map[platform.OS]*TemplateWithoutPlatform  `toml:"platform"`
	Profile  map[build.Profile]*TemplateWithoutProfile `toml:"profile"`
}

/* --------------------- Struct: TemplateWithoutFeature --------------------- */

type TemplateWithoutFeature struct {
	*build.Template

	Platform map[platform.OS]*TemplateWithoutFeatureAndPlatform  `toml:"platform"`
	Profile  map[build.Profile]*TemplateWithoutFeatureAndProfile `toml:"profile"`
}

/* --------------------- Struct: TemplateWithoutPlatform -------------------- */

type TemplateWithoutPlatform struct {
	*build.Template

	Feature map[string]*TemplateWithoutFeatureAndPlatform        `toml:"feature"`
	Profile map[build.Profile]*TemplateWithoutPlatformAndProfile `toml:"profile"`
}

/* --------------------- Struct: TemplateWithoutProfile --------------------- */

type TemplateWithoutProfile struct {
	*build.Template

	Feature  map[string]*TemplateWithoutFeatureAndProfile       `toml:"feature"`
	Platform map[platform.OS]*TemplateWithoutPlatformAndProfile `toml:"platform"`
}

/* ---------------- Struct: TemplateWithoutFeatureAndPlatform --------------- */

type TemplateWithoutFeatureAndPlatform struct {
	*build.Template

	Profile map[build.Profile]*build.Template `toml:"profile"`
}

/* ---------------- Struct: TemplateWithoutFeatureAndProfile ---------------- */

type TemplateWithoutFeatureAndProfile struct {
	*build.Template

	Platform map[platform.OS]*build.Template `toml:"platform"`
}

/* ---------------- Struct: TemplateWithoutPlatformAndProfile --------------- */

type TemplateWithoutPlatformAndProfile struct {
	*build.Template

	Feature map[string]*build.Template `toml:"feature"`
}
