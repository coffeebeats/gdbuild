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

	Feature  map[string]*TargetWithoutFeature        `toml:"feature"`
	Platform map[platform.OS]*TargetWithoutPlatform  `toml:"platform"`
	Profile  map[build.Profile]*TargetWithoutProfile `toml:"profile"`
}

/* ---------------------- Struct: TargetWithoutFeature ---------------------- */

type TargetWithoutFeature struct {
	*build.Target

	Platform map[platform.OS]*TargetWithoutFeatureAndPlatform  `toml:"platform"`
	Profile  map[build.Profile]*TargetWithoutFeatureAndProfile `toml:"profile"`
}

/* ---------------------- Struct: TargetWithoutPlatform --------------------- */

type TargetWithoutPlatform struct {
	*build.Target

	Feature map[string]*TargetWithoutFeatureAndPlatform        `toml:"feature"`
	Profile map[build.Profile]*TargetWithoutPlatformAndProfile `toml:"profile"`
}

/* ---------------------- Struct: TargetWithoutProfile ---------------------- */

type TargetWithoutProfile struct {
	*build.Target

	Feature  map[string]*TargetWithoutFeatureAndProfile       `toml:"feature"`
	Platform map[platform.OS]*TargetWithoutPlatformAndProfile `toml:"platform"`
}

/* ----------------- Struct: TargetWithoutFeatureAndPlatform ---------------- */

type TargetWithoutFeatureAndPlatform struct {
	*build.Target

	Profile map[build.Profile]*build.Target `toml:"profile"`
}

/* ----------------- Struct: TargetWithoutFeatureAndProfile ----------------- */

type TargetWithoutFeatureAndProfile struct {
	*build.Target

	Platform map[platform.OS]*build.Target `toml:"platform"`
}

/* ----------------- Struct: TargetWithoutPlatformAndProfile ---------------- */

type TargetWithoutPlatformAndProfile struct {
	*build.Target

	Feature map[string]*build.Target `toml:"feature"`
}
