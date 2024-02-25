package manifest

import (
	"github.com/coffeebeats/gdbuild/pkg/build"
	"github.com/coffeebeats/gdbuild/pkg/build/template"
)

/* -------------------------------------------------------------------------- */
/*                              Struct: Template                              */
/* -------------------------------------------------------------------------- */

// Template defines the parameters for building a Godot export template for a
// specified platform. A 'Template' definition can be customized based on
// 'feature', 'platform', and 'profile' labels used in the property names. Note
// that each specifier label can only be used once per property name (i.e.
// 'target.profile.release.profile.debug' is not allowed). Additionally, the
// order of specifies is strict: 'platform' > 'feature' > 'profile'.
//
// For example, the following are all valid table names:
//
//	[template]
//	[template.profile.release]
//	[template.platform.macos.feature.client]
//	[template.platform.linux.feature.server.profile.release_debug]
type Template struct {
	*template.Base

	Platform TemplatePlatform                      `toml:"platform"`
	Feature  map[string]TemplateBaseWithoutFeature `toml:"feature"`
	Profile  map[build.Profile]template.Base       `toml:"profile"`
}

func (t *Template) Features() map[string]TemplateBaseWithoutFeature {
	if t == nil || t.Feature == nil {
		return map[string]TemplateBaseWithoutFeature{}
	}

	return t.Feature
}

func (t *Template) Profiles() map[build.Profile]template.Base {
	if t == nil || t.Profile == nil {
		return map[build.Profile]template.Base{}
	}

	return t.Profile
}

/* ------------------- Struct: TemplateBaseWithoutFeature ------------------- */

type TemplateBaseWithoutFeature struct {
	*template.Base

	Profile map[build.Profile]template.Base `toml:"profile"`
}

func (t *TemplateBaseWithoutFeature) Profiles() map[build.Profile]template.Base {
	if t == nil || t.Profile == nil {
		return map[build.Profile]template.Base{}
	}

	return t.Profile
}

/* -------------------------------------------------------------------------- */
/*                          Struct: TemplatePlatform                          */
/* -------------------------------------------------------------------------- */

type TemplatePlatform struct {
	Android TemplateAndroid `toml:"android"`
	IOS     TemplateIOS     `toml:"ios"`
	Linux   TemplateLinux   `toml:"linux"`
	MacOS   TemplateMacOS   `toml:"macos"`
	Web     TemplateWeb     `toml:"web"`
	Windows TemplateWindows `toml:"windows"`
}

/* ---------------------------- Platform: Android --------------------------- */

type TemplateAndroid struct {
	*template.Android

	Feature map[string]TemplateAndroidWithoutFeature `toml:"feature"`
	Profile map[build.Profile]template.Android       `toml:"profile"`
}

func (t *TemplateAndroid) Features() map[string]TemplateAndroidWithoutFeature {
	if t == nil || t.Feature == nil {
		return map[string]TemplateAndroidWithoutFeature{}
	}

	return t.Feature
}

func (t *TemplateAndroid) Profiles() map[build.Profile]template.Android {
	if t == nil || t.Profile == nil {
		return map[build.Profile]template.Android{}
	}

	return t.Profile
}

type TemplateAndroidWithoutFeature struct {
	*template.Android

	Profile map[build.Profile]template.Android `toml:"profile"`
}

func (t *TemplateAndroidWithoutFeature) Profiles() map[build.Profile]template.Android {
	if t == nil || t.Profile == nil {
		return map[build.Profile]template.Android{}
	}

	return t.Profile
}

/* ------------------------------ Platform: IOS ----------------------------- */

type TemplateIOS struct {
	*template.IOS

	Feature map[string]TemplateIOSWithoutFeature `toml:"feature"`
	Profile map[build.Profile]template.IOS       `toml:"profile"`
}

func (t *TemplateIOS) Features() map[string]TemplateIOSWithoutFeature {
	if t == nil || t.Feature == nil {
		return map[string]TemplateIOSWithoutFeature{}
	}

	return t.Feature
}

func (t *TemplateIOS) Profiles() map[build.Profile]template.IOS {
	if t == nil || t.Profile == nil {
		return map[build.Profile]template.IOS{}
	}

	return t.Profile
}

type TemplateIOSWithoutFeature struct {
	*template.IOS

	Profile map[build.Profile]template.IOS `toml:"profile"`
}

func (t *TemplateIOSWithoutFeature) Profiles() map[build.Profile]template.IOS {
	if t == nil || t.Profile == nil {
		return map[build.Profile]template.IOS{}
	}

	return t.Profile
}

/* ----------------------------- Platform: Linux ---------------------------- */

type TemplateLinux struct {
	*template.Linux

	Feature map[string]TemplateLinuxWithoutFeature `toml:"feature"`
	Profile map[build.Profile]template.Linux       `toml:"profile"`
}

func (t *TemplateLinux) Features() map[string]TemplateLinuxWithoutFeature {
	if t == nil || t.Feature == nil {
		return map[string]TemplateLinuxWithoutFeature{}
	}

	return t.Feature
}

func (t *TemplateLinux) Profiles() map[build.Profile]template.Linux {
	if t == nil || t.Profile == nil {
		return map[build.Profile]template.Linux{}
	}

	return t.Profile
}

type TemplateLinuxWithoutFeature struct {
	*template.Linux

	Profile map[build.Profile]template.Linux `toml:"profile"`
}

func (t *TemplateLinuxWithoutFeature) Profiles() map[build.Profile]template.Linux {
	if t == nil || t.Profile == nil {
		return map[build.Profile]template.Linux{}
	}

	return t.Profile
}

/* ----------------------------- Platform: MacOS ---------------------------- */

type TemplateMacOS struct {
	*template.MacOS

	Feature map[string]TemplateMacOSWithoutFeature `toml:"feature"`
	Profile map[build.Profile]template.MacOS       `toml:"profile"`
}

func (t *TemplateMacOS) Features() map[string]TemplateMacOSWithoutFeature {
	if t == nil || t.Feature == nil {
		return map[string]TemplateMacOSWithoutFeature{}
	}

	return t.Feature
}

func (t *TemplateMacOS) Profiles() map[build.Profile]template.MacOS {
	if t == nil || t.Profile == nil {
		return map[build.Profile]template.MacOS{}
	}

	return t.Profile
}

type TemplateMacOSWithoutFeature struct {
	*template.MacOS

	Profile map[build.Profile]template.MacOS `toml:"profile"`
}

func (t *TemplateMacOSWithoutFeature) Profiles() map[build.Profile]template.MacOS {
	if t == nil || t.Profile == nil {
		return map[build.Profile]template.MacOS{}
	}

	return t.Profile
}

/* ------------------------------ Platform: Web ----------------------------- */

type TemplateWeb struct {
	*template.Web

	Feature map[string]TemplateWebWithoutFeature `toml:"feature"`
	Profile map[build.Profile]template.Web       `toml:"profile"`
}

func (t *TemplateWeb) Features() map[string]TemplateWebWithoutFeature {
	if t == nil || t.Feature == nil {
		return map[string]TemplateWebWithoutFeature{}
	}

	return t.Feature
}

func (t *TemplateWeb) Profiles() map[build.Profile]template.Web {
	if t == nil || t.Profile == nil {
		return map[build.Profile]template.Web{}
	}

	return t.Profile
}

type TemplateWebWithoutFeature struct {
	*template.Web

	Profile map[build.Profile]template.Web `toml:"profile"`
}

func (t *TemplateWebWithoutFeature) Profiles() map[build.Profile]template.Web {
	if t == nil || t.Profile == nil {
		return map[build.Profile]template.Web{}
	}

	return t.Profile
}

/* ---------------------------- Platform: Windows --------------------------- */

type TemplateWindows struct {
	*template.Windows

	Feature map[string]TemplateWindowsWithoutFeature `toml:"feature"`
	Profile map[build.Profile]template.Windows       `toml:"profile"`
}

func (t *TemplateWindows) Features() map[string]TemplateWindowsWithoutFeature {
	if t == nil || t.Feature == nil {
		return map[string]TemplateWindowsWithoutFeature{}
	}

	return t.Feature
}

func (t *TemplateWindows) Profiles() map[build.Profile]template.Windows {
	if t == nil || t.Profile == nil {
		return map[build.Profile]template.Windows{}
	}

	return t.Profile
}

type TemplateWindowsWithoutFeature struct {
	*template.Windows

	Profile map[build.Profile]template.Windows `toml:"profile"`
}

func (t *TemplateWindowsWithoutFeature) Profiles() map[build.Profile]template.Windows {
	if t == nil || t.Profile == nil {
		return map[build.Profile]template.Windows{}
	}

	return t.Profile
}
