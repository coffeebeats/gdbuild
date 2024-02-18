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

/* ---------------------- Struct: TargetWithoutFeature ---------------------- */

type TargetWithoutFeature struct {
	*build.Target

	Platform map[platform.OS]*TargetWithoutFeatureAndPlatform  `json:"platform" toml:"platform"`
	Profile  map[build.Profile]*TargetWithoutFeatureAndProfile `json:"profile"  toml:"profile"`
}

/* ---------------------- Struct: TargetWithoutPlatform --------------------- */

type TargetWithoutPlatform struct {
	*build.Target

	Feature map[string]*TargetWithoutFeatureAndPlatform        `json:"feature" toml:"feature"`
	Profile map[build.Profile]*TargetWithoutPlatformAndProfile `json:"profile" toml:"profile"`
}

/* ---------------------- Struct: TargetWithoutProfile ---------------------- */

type TargetWithoutProfile struct {
	*build.Target

	Feature  map[string]*TargetWithoutFeatureAndProfile       `json:"feature"  toml:"feature"`
	Platform map[platform.OS]*TargetWithoutPlatformAndProfile `json:"platform" toml:"platform"`
}

/* ----------------- Struct: TargetWithoutFeatureAndPlatform ---------------- */

type TargetWithoutFeatureAndPlatform struct {
	*build.Target

	Profile map[build.Profile]*build.Target `json:"profile" toml:"profile"`
}

/* ----------------- Struct: TargetWithoutFeatureAndProfile ----------------- */

type TargetWithoutFeatureAndProfile struct {
	*build.Target

	Platform map[platform.OS]*build.Target `json:"platform" toml:"platform"`
}

/* ----------------- Struct: TargetWithoutPlatformAndProfile ---------------- */

type TargetWithoutPlatformAndProfile struct {
	*build.Target

	Feature map[string]*build.Target `json:"feature" toml:"feature"`
}
