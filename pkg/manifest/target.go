package manifest

import (
	"github.com/coffeebeats/gdbuild/pkg/build"
	"github.com/coffeebeats/gdbuild/pkg/target"
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
	*target.Base

	Platform TargetPlatform                      `toml:"platform"`
	Feature  map[string]TargetBaseWithoutFeature `toml:"feature"`
	Profile  map[build.Profile]target.Base       `toml:"profile"`
}

/* -------------------- Struct: TargetBaseWithoutFeature -------------------- */

type TargetBaseWithoutFeature struct {
	*target.Base

	Profile map[build.Profile]target.Base `toml:"profile"`
}

/* -------------------------------------------------------------------------- */
/*                           Struct: TargetPlatform                           */
/* -------------------------------------------------------------------------- */

type TargetPlatform struct {
	Android TargetMacOS   `toml:"android"`
	IOS     TargetIOS     `toml:"ios"`
	MacOS   TargetMacOS   `toml:"macos"`
	Linux   TargetLinux   `toml:"linux"`
	Web     TargetWeb     `toml:"web"`
	Windows TargetWindows `toml:"windows"`
}

/* ---------------------------- Platform: Android --------------------------- */

type TargetAndroid struct {
	*target.Android

	Feature map[string]TargetAndroidWithoutFeature `toml:"feature"`
	Profile map[build.Profile]target.Android       `toml:"profile"`
}

type TargetAndroidWithoutFeature struct {
	*target.Android

	Profile map[build.Profile]target.Android `toml:"profile"`
}

/* ------------------------------ Platform: IOS ----------------------------- */

type TargetIOS struct {
	*target.IOS

	Feature map[string]TargetIOSWithoutFeature `toml:"feature"`
	Profile map[build.Profile]target.IOS       `toml:"profile"`
}

type TargetIOSWithoutFeature struct {
	*target.IOS

	Profile map[build.Profile]target.IOS `toml:"profile"`
}

/* ----------------------------- Platform: MacOS ---------------------------- */

type TargetMacOS struct {
	*target.MacOS

	Feature map[string]TargetMacOSWithoutFeature `toml:"feature"`
	Profile map[build.Profile]target.MacOS       `toml:"profile"`
}

type TargetMacOSWithoutFeature struct {
	*target.MacOS

	Profile map[build.Profile]target.MacOS `toml:"profile"`
}

/* ----------------------------- Platform: Linux ---------------------------- */

type TargetLinux struct {
	*target.Linux

	Feature map[string]TargetLinuxWithoutFeature `toml:"feature"`
	Profile map[build.Profile]target.Linux       `toml:"profile"`
}

type TargetLinuxWithoutFeature struct {
	*target.Linux

	Profile map[build.Profile]target.Linux `toml:"profile"`
}

/* ------------------------------ Platform: Web ----------------------------- */

type TargetWeb struct {
	*target.Web

	Feature map[string]TargetWebWithoutFeature `toml:"feature"`
	Profile map[build.Profile]target.Web       `toml:"profile"`
}

type TargetWebWithoutFeature struct {
	*target.Web

	Profile map[build.Profile]target.Web `toml:"profile"`
}

/* ---------------------------- Platform: Windows --------------------------- */

type TargetWindows struct {
	*target.Windows

	Feature map[string]TargetWindowsWithoutFeature `toml:"feature"`
	Profile map[build.Profile]target.Windows       `toml:"profile"`
}

type TargetWindowsWithoutFeature struct {
	*target.Windows

	Profile map[build.Profile]target.Windows `toml:"profile"`
}
