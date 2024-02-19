package manifest //nolint:dupl

import (
	"github.com/coffeebeats/gdbuild/pkg/build"
	"github.com/coffeebeats/gdbuild/pkg/platform"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Target                               */
/* -------------------------------------------------------------------------- */

// Target specifies a single exportable artifact within the Godot project. A
// 'Target' can be customized based on 'feature', 'platform', and 'profile'
// labels used in the property names. Note that each specifier label can only
// be used once per property name (i.e. 'target.profile.release.profile.debug'
// is not allowed).
//
// For example, the following are all valid table names:
//
//	[target]
//	[target.feature.client]
//	[target.platform.windows]
//	[target.profile.release]
//	[target.profile.release.platform.macos.feature.client]
type Target struct {
	*build.Target

	Feature  map[string]*TargetWithoutFeature        `json:"feature"  toml:"feature"`
	Platform map[platform.OS]*TargetWithoutPlatform  `json:"platform" toml:"platform"`
	Profile  map[build.Profile]*TargetWithoutProfile `json:"profile"  toml:"profile"`
}

func (t *Target) merge(pl platform.OS, pr build.Profile, ff ...string) *build.Target {
	out := t.Target
	if out == nil {
		out = &build.Target{} //nolint:exhaustruct
	}

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

/* ---------------------- Struct: TargetWithoutFeature ---------------------- */

type TargetWithoutFeature struct {
	*build.Target

	Platform map[platform.OS]*TargetWithoutFeatureAndPlatform  `json:"platform" toml:"platform"`
	Profile  map[build.Profile]*TargetWithoutFeatureAndProfile `json:"profile"  toml:"profile"`
}

func (t *TargetWithoutFeature) merge(pl platform.OS, pr build.Profile, ff ...string) *build.Target {
	out := t.Target
	if out == nil {
		out = &build.Target{} //nolint:exhaustruct
	}

	if cfg, ok := t.Profile[pr]; ok {
		out = out.CombineWith(cfg.merge(pl, pr, ff...))
	}

	if cfg, ok := t.Platform[pl]; ok {
		out = out.CombineWith(cfg.merge(pl, pr, ff...))
	}

	return out
}

/* ---------------------- Struct: TargetWithoutPlatform --------------------- */

type TargetWithoutPlatform struct {
	*build.Target

	Feature map[string]*TargetWithoutFeatureAndPlatform        `json:"feature" toml:"feature"`
	Profile map[build.Profile]*TargetWithoutPlatformAndProfile `json:"profile" toml:"profile"`
}

func (t *TargetWithoutPlatform) merge(pl platform.OS, pr build.Profile, ff ...string) *build.Target {
	out := t.Target
	if out == nil {
		out = &build.Target{} //nolint:exhaustruct
	}

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

/* ---------------------- Struct: TargetWithoutProfile ---------------------- */

type TargetWithoutProfile struct {
	*build.Target

	Feature  map[string]*TargetWithoutFeatureAndProfile       `json:"feature"  toml:"feature"`
	Platform map[platform.OS]*TargetWithoutPlatformAndProfile `json:"platform" toml:"platform"`
}

func (t *TargetWithoutProfile) merge(pl platform.OS, pr build.Profile, ff ...string) *build.Target {
	out := t.Target
	if out == nil {
		out = &build.Target{} //nolint:exhaustruct
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

/* ----------------- Struct: TargetWithoutFeatureAndPlatform ---------------- */

type TargetWithoutFeatureAndPlatform struct {
	*build.Target

	Profile map[build.Profile]*build.Target `json:"profile" toml:"profile"`
}

func (t *TargetWithoutFeatureAndPlatform) merge(_ platform.OS, pr build.Profile, _ ...string) *build.Target {
	out := t.Target
	if out == nil {
		out = &build.Target{} //nolint:exhaustruct
	}

	if cfg, ok := t.Profile[pr]; ok {
		out = out.CombineWith(cfg)
	}

	return out
}

/* ----------------- Struct: TargetWithoutFeatureAndProfile ----------------- */

type TargetWithoutFeatureAndProfile struct {
	*build.Target

	Platform map[platform.OS]*build.Target `json:"platform" toml:"platform"`
}

func (t *TargetWithoutFeatureAndProfile) merge(pl platform.OS, _ build.Profile, _ ...string) *build.Target {
	out := t.Target
	if out == nil {
		out = &build.Target{} //nolint:exhaustruct
	}

	if cfg, ok := t.Platform[pl]; ok {
		out = out.CombineWith(cfg)
	}

	return out
}

/* ----------------- Struct: TargetWithoutPlatformAndProfile ---------------- */

type TargetWithoutPlatformAndProfile struct {
	*build.Target

	Feature map[string]*build.Target `json:"feature" toml:"feature"`
}

func (t *TargetWithoutPlatformAndProfile) merge(_ platform.OS, _ build.Profile, ff ...string) *build.Target {
	out := t.Target
	if out == nil {
		out = &build.Target{} //nolint:exhaustruct
	}

	for _, f := range ff {
		if cfg, ok := t.Feature[f]; ok {
			out = out.CombineWith(cfg)
		}
	}

	return out
}
