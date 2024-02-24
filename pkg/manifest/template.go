package manifest //nolint:dupl

import "github.com/coffeebeats/gdbuild/pkg/build"

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

	Feature  map[string]*TemplateWithoutFeature        `json:"feature"  toml:"feature"`
	Platform map[build.OS]*TemplateWithoutPlatform     `json:"platform" toml:"platform"`
	Profile  map[build.Profile]*TemplateWithoutProfile `json:"profile"  toml:"profile"`
}

// TODO: Improve merging logic to detect conflicts instead of silently, and
// unpredictably, overriding.
func (t *Template) merge(pl build.OS, pr build.Profile, ff ...string) *build.Template {
	out := t.Template

	if cfg, ok := t.Profile[pr]; ok {
		out = out.CombineWith(cfg.merge(pl, pr, ff...))
	}

	if cfg, ok := t.Platform[pl]; ok {
		out = out.CombineWith(cfg.merge(pl, pr, ff...))
	}

	for _, f := range ff {
		if cfg, ok := t.Feature[f]; ok {
			out = out.CombineWith(cfg.merge(pl, pr, ff...))
		}
	}

	return out
}

/* --------------------- Struct: TemplateWithoutFeature --------------------- */

type TemplateWithoutFeature struct {
	*build.Template

	Platform map[build.OS]*TemplateWithoutFeatureAndPlatform     `json:"platform" toml:"platform"`
	Profile  map[build.Profile]*TemplateWithoutFeatureAndProfile `json:"profile"  toml:"profile"`
}

func (t *TemplateWithoutFeature) merge(pl build.OS, pr build.Profile, ff ...string) *build.Template {
	out := t.Template

	if cfg, ok := t.Profile[pr]; ok {
		out = out.CombineWith(cfg.merge(pl, pr, ff...))
	}

	if cfg, ok := t.Platform[pl]; ok {
		out = out.CombineWith(cfg.merge(pl, pr, ff...))
	}

	return out
}

/* --------------------- Struct: TemplateWithoutPlatform -------------------- */

type TemplateWithoutPlatform struct {
	*build.Template

	Feature map[string]*TemplateWithoutFeatureAndPlatform        `json:"feature" toml:"feature"`
	Profile map[build.Profile]*TemplateWithoutPlatformAndProfile `json:"profile" toml:"profile"`
}

func (t *TemplateWithoutPlatform) merge(pl build.OS, pr build.Profile, ff ...string) *build.Template {
	out := t.Template

	if cfg, ok := t.Profile[pr]; ok {
		out = out.CombineWith(cfg.merge(pl, pr, ff...))
	}

	for _, f := range ff {
		if cfg, ok := t.Feature[f]; ok {
			out = out.CombineWith(cfg.merge(pl, pr, ff...))
		}
	}

	return out
}

/* --------------------- Struct: TemplateWithoutProfile --------------------- */

type TemplateWithoutProfile struct {
	*build.Template

	Feature  map[string]*TemplateWithoutFeatureAndProfile    `json:"feature"  toml:"feature"`
	Platform map[build.OS]*TemplateWithoutPlatformAndProfile `json:"platform" toml:"platform"`
}

func (t *TemplateWithoutProfile) merge(pl build.OS, pr build.Profile, ff ...string) *build.Template {
	out := t.Template

	if cfg, ok := t.Platform[pl]; ok {
		out = out.CombineWith(cfg.merge(pl, pr, ff...))
	}

	for _, f := range ff {
		if cfg, ok := t.Feature[f]; ok {
			out = out.CombineWith(cfg.merge(pl, pr, ff...))
		}
	}

	return out
}

/* ---------------- Struct: TemplateWithoutFeatureAndPlatform --------------- */

type TemplateWithoutFeatureAndPlatform struct {
	*build.Template

	Profile map[build.Profile]*build.Template `json:"profile" toml:"profile"`
}

func (t *TemplateWithoutFeatureAndPlatform) merge(_ build.OS, pr build.Profile, _ ...string) *build.Template {
	out := t.Template

	if cfg, ok := t.Profile[pr]; ok {
		out = out.CombineWith(cfg)
	}

	return out
}

/* ---------------- Struct: TemplateWithoutFeatureAndProfile ---------------- */

type TemplateWithoutFeatureAndProfile struct {
	*build.Template

	Platform map[build.OS]*build.Template `json:"platform" toml:"platform"`
}

func (t *TemplateWithoutFeatureAndProfile) merge(pl build.OS, _ build.Profile, _ ...string) *build.Template {
	out := t.Template

	if cfg, ok := t.Platform[pl]; ok {
		out = out.CombineWith(cfg)
	}

	return out
}

/* ---------------- Struct: TemplateWithoutPlatformAndProfile --------------- */

type TemplateWithoutPlatformAndProfile struct {
	*build.Template

	Feature map[string]*build.Template `json:"feature" toml:"feature"`
}

func (t *TemplateWithoutPlatformAndProfile) merge(_ build.OS, _ build.Profile, ff ...string) *build.Template {
	out := t.Template

	for _, f := range ff {
		if cfg, ok := t.Feature[f]; ok {
			out = out.CombineWith(cfg)
		}
	}

	return out
}
