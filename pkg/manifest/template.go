package manifest

import (
	"github.com/coffeebeats/gdbuild/pkg/build"
	"github.com/coffeebeats/gdbuild/pkg/template"
)

/* -------------------------------------------------------------------------- */
/*                              Struct: Template                              */
/* -------------------------------------------------------------------------- */

// Template defines the parameters for building a Godot export template for a
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
type Template struct {
	*template.Base

	Platform TemplatePlatform                      `toml:"platform"`
	Feature  map[string]TemplateBaseWithoutFeature `toml:"feature"`
	Profile  map[build.Profile]template.Base       `toml:"profile"`
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (t *Template) Configure(inv *build.Invocation) error {
	if t == nil {
		return nil
	}

	if t.Base == nil {
		t.Base = &template.Base{} //nolint:exhaustruct
	}

	p := getOrDefault(t.Profile, inv.Profile)
	if err := t.Base.Merge(&p); err != nil {
		return err
	}

	for _, f := range inv.Features {
		f := getOrDefault(t.Feature, f)
		if err := t.Base.Merge(f.Base); err != nil {
			return err
		}
	}

	for _, f := range inv.Features {
		f := getOrDefault(t.Feature, f)
		p := getOrDefault(f.Profile, inv.Profile)

		if err := t.Base.Merge(&p); err != nil {
			return err
		}
	}

	return nil
}

/* ------------------- Struct: TemplateBaseWithoutFeature ------------------- */

type TemplateBaseWithoutFeature struct {
	*template.Base

	Profile map[build.Profile]template.Base `toml:"profile"`
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

/* -------------------------------------------------------------------------- */
/*                              Platform: Android                             */
/* -------------------------------------------------------------------------- */

type TemplateAndroid struct {
	*template.Android

	Feature map[string]TemplateAndroidWithoutFeature `toml:"feature"`
	Profile map[build.Profile]template.Android       `toml:"profile"`
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (t *TemplateAndroid) Configure(inv *build.Invocation) error {
	if t == nil {
		return nil
	}

	if t.Android == nil {
		t.Android = &template.Android{} //nolint:exhaustruct
	}

	p := getOrDefault(t.Profile, inv.Profile)
	if err := t.Android.Merge(&p); err != nil {
		return err
	}

	for _, f := range inv.Features {
		f := getOrDefault(t.Feature, f)
		if err := t.Android.Merge(f.Android); err != nil {
			return err
		}
	}

	for _, f := range inv.Features {
		f := getOrDefault(t.Feature, f)
		p := getOrDefault(f.Profile, inv.Profile)

		if err := t.Android.Merge(&p); err != nil {
			return err
		}
	}

	return nil
}

/* --------------------- Struct: TemplateWithoutFeature --------------------- */

type TemplateAndroidWithoutFeature struct {
	*template.Android

	Profile map[build.Profile]template.Android `toml:"profile"`
}

/* -------------------------------------------------------------------------- */
/*                                Platform: IOS                               */
/* -------------------------------------------------------------------------- */

type TemplateIOS struct {
	*template.IOS

	Feature map[string]TemplateIOSWithoutFeature `toml:"feature"`
	Profile map[build.Profile]template.IOS       `toml:"profile"`
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (t *TemplateIOS) Configure(inv *build.Invocation) error {
	if t == nil {
		return nil
	}

	if t.IOS == nil {
		t.IOS = &template.IOS{} //nolint:exhaustruct
	}

	p := getOrDefault(t.Profile, inv.Profile)
	if err := t.IOS.Merge(&p); err != nil {
		return err
	}

	for _, f := range inv.Features {
		f := getOrDefault(t.Feature, f)
		if err := t.IOS.Merge(f.IOS); err != nil {
			return err
		}
	}

	for _, f := range inv.Features {
		f := getOrDefault(t.Feature, f)
		p := getOrDefault(f.Profile, inv.Profile)

		if err := t.IOS.Merge(&p); err != nil {
			return err
		}
	}

	return nil
}

/* --------------------- Struct: TemplateWithoutFeature --------------------- */

type TemplateIOSWithoutFeature struct {
	*template.IOS

	Profile map[build.Profile]template.IOS `toml:"profile"`
}

/* -------------------------------------------------------------------------- */
/*                               Platform: Linux                              */
/* -------------------------------------------------------------------------- */

type TemplateLinux struct {
	*template.Linux

	Feature map[string]TemplateLinuxWithoutFeature `toml:"feature"`
	Profile map[build.Profile]template.Linux       `toml:"profile"`
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (t *TemplateLinux) Configure(inv *build.Invocation) error {
	if t == nil {
		return nil
	}

	if t.Linux == nil {
		t.Linux = &template.Linux{} //nolint:exhaustruct
	}

	p := getOrDefault(t.Profile, inv.Profile)
	if err := t.Linux.Merge(&p); err != nil {
		return err
	}

	for _, f := range inv.Features {
		f := getOrDefault(t.Feature, f)
		if err := t.Linux.Merge(f.Linux); err != nil {
			return err
		}
	}

	for _, f := range inv.Features {
		f := getOrDefault(t.Feature, f)
		p := getOrDefault(f.Profile, inv.Profile)

		if err := t.Linux.Merge(&p); err != nil {
			return err
		}
	}

	return nil
}

/* --------------------- Struct: TemplateWithoutFeature --------------------- */

type TemplateLinuxWithoutFeature struct {
	*template.Linux

	Profile map[build.Profile]template.Linux `toml:"profile"`
}

/* -------------------------------------------------------------------------- */
/*                               Platform: MacOS                              */
/* -------------------------------------------------------------------------- */

type TemplateMacOS struct {
	*template.MacOS

	Feature map[string]TemplateMacOSWithoutFeature `toml:"feature"`
	Profile map[build.Profile]template.MacOS       `toml:"profile"`
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (t *TemplateMacOS) Configure(inv *build.Invocation) error {
	if t == nil {
		return nil
	}

	if t.MacOS == nil {
		t.MacOS = &template.MacOS{} //nolint:exhaustruct
	}

	p := getOrDefault(t.Profile, inv.Profile)
	if err := t.MacOS.Merge(&p); err != nil {
		return err
	}

	for _, f := range inv.Features {
		f := getOrDefault(t.Feature, f)
		if err := t.MacOS.Merge(f.MacOS); err != nil {
			return err
		}
	}

	for _, f := range inv.Features {
		f := getOrDefault(t.Feature, f)
		p := getOrDefault(f.Profile, inv.Profile)

		if err := t.MacOS.Merge(&p); err != nil {
			return err
		}
	}

	return nil
}

/* --------------------- Struct: TemplateWithoutFeature --------------------- */

type TemplateMacOSWithoutFeature struct {
	*template.MacOS

	Profile map[build.Profile]template.MacOS `toml:"profile"`
}

/* -------------------------------------------------------------------------- */
/*                                Platform: Web                               */
/* -------------------------------------------------------------------------- */

type TemplateWeb struct {
	*template.Web

	Feature map[string]TemplateWebWithoutFeature `toml:"feature"`
	Profile map[build.Profile]template.Web       `toml:"profile"`
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (t *TemplateWeb) Configure(inv *build.Invocation) error {
	if t == nil {
		return nil
	}

	if t.Web == nil {
		t.Web = &template.Web{} //nolint:exhaustruct
	}

	p := getOrDefault(t.Profile, inv.Profile)
	if err := t.Web.Merge(&p); err != nil {
		return err
	}

	for _, f := range inv.Features {
		f := getOrDefault(t.Feature, f)
		if err := t.Web.Merge(f.Web); err != nil {
			return err
		}
	}

	for _, f := range inv.Features {
		f := getOrDefault(t.Feature, f)
		p := getOrDefault(f.Profile, inv.Profile)

		if err := t.Web.Merge(&p); err != nil {
			return err
		}
	}

	return nil
}

/* --------------------- Struct: TemplateWithoutFeature --------------------- */

type TemplateWebWithoutFeature struct {
	*template.Web

	Profile map[build.Profile]template.Web `toml:"profile"`
}

/* -------------------------------------------------------------------------- */
/*                              Platform: Windows                             */
/* -------------------------------------------------------------------------- */

type TemplateWindows struct {
	*template.Windows

	Feature map[string]TemplateWindowsWithoutFeature `toml:"feature"`
	Profile map[build.Profile]template.Windows       `toml:"profile"`
}

/* ------------------------- Impl: build.Configurer ------------------------- */

func (t *TemplateWindows) Configure(inv *build.Invocation) error {
	if t == nil {
		return nil
	}

	if t.Windows == nil {
		t.Windows = &template.Windows{} //nolint:exhaustruct
	}

	p := getOrDefault(t.Profile, inv.Profile)
	if err := t.Windows.Merge(&p); err != nil {
		return err
	}

	for _, f := range inv.Features {
		f := getOrDefault(t.Feature, f)
		if err := t.Windows.Merge(f.Windows); err != nil {
			return err
		}
	}

	for _, f := range inv.Features {
		f := getOrDefault(t.Feature, f)
		p := getOrDefault(f.Profile, inv.Profile)

		if err := t.Windows.Merge(&p); err != nil {
			return err
		}
	}

	return nil
}

/* --------------------- Struct: TemplateWithoutFeature --------------------- */

type TemplateWindowsWithoutFeature struct {
	*template.Windows

	Profile map[build.Profile]template.Windows `toml:"profile"`
}
