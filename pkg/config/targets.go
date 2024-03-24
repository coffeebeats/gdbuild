package config

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/config/platform/common"
	"github.com/coffeebeats/gdbuild/pkg/config/platform/linux"
	"github.com/coffeebeats/gdbuild/pkg/config/platform/macos"
	"github.com/coffeebeats/gdbuild/pkg/config/platform/windows"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                             Interface: Exporter                            */
/* -------------------------------------------------------------------------- */

type Exporter interface {
	config.Configurable[*run.Context]

	Collect(src engine.Source, rc *run.Context) *export.Export
}

/* -------------------------------------------------------------------------- */
/*                               Struct: Targets                              */
/* -------------------------------------------------------------------------- */

// Targets defines the parameters for exporting a game binary or pack file for
// a specified platform. A 'Target' definition can be customized  based on
// 'feature', 'platform', and 'profile' labels used in the property names. Note
// that each specifier label can only be used once per property name (i.e.
// 'target.profile.release.profile.debug' is not allowed). Additionally, the
// order of specifiers is strict: 'platform' < 'feature' < 'profile'.
//
// For example, the following are all valid table names:
//
//	[target]
//	[target.profile.release]
//	[target.platform.macos.feature.client]
//	[target.platform.linux.feature.server.profile.release_debug]
type Targets struct {
	Target *common.Target

	Platform TargetPlatforms                  `toml:"platform"`
	Feature  map[string]TargetWithoutFeature  `toml:"feature"`
	Profile  map[engine.Profile]common.Target `toml:"profile"`
}

/* ----------------------- Struct: TargetWithoutFeature ----------------------- */

type TargetWithoutFeature struct {
	Target *common.Target

	Profile map[engine.Profile]common.Target `toml:"profile"`
}

/* ---------------------------- Struct: Platforms --------------------------- */

type TargetPlatforms struct {
	Linux   LinuxTargetWithFeaturesAndProfile   `toml:"linux"`
	MacOS   MacOSTargetWithFeaturesAndProfile   `toml:"macos"`
	Windows WindowsTargetWithFeaturesAndProfile `toml:"windows"`
}

/* ------------------------------ Method: Build ----------------------------- */

func (t *Targets) Build(rc *run.Context) (Exporter, error) { //nolint:cyclop,ireturn
	// Target params (root)
	var out Exporter = new(common.Target)

	// Target params (root)
	if err := t.Target.MergeInto(out); err != nil {
		return nil, err
	}

	// Target params (feature-constrained)
	for _, f := range rc.Features {
		bwof := t.Feature[f].Target
		if err := bwof.MergeInto(out); err != nil {
			return nil, err
		}
	}

	// Target params (profile-constrained)
	b := t.Profile[rc.Profile]
	if err := b.MergeInto(out); err != nil {
		return nil, err
	}

	// Feature-and-profile-constrained params
	for _, f := range rc.Features {
		bwof := t.Feature[f].Profile[rc.Profile]
		if err := bwof.MergeInto(out); err != nil {
			return nil, err
		}
	}

	switch p := rc.Platform; p {
	case platform.OSLinux:
		out = &linux.Target{Target: out.(*common.Target)} //nolint:forcetypeassert

		if err := t.Platform.Linux.build(rc, out); err != nil {
			return nil, err
		}
	case platform.OSMacOS:
		out = &macos.Target{Target: out.(*common.Target)} //nolint:forcetypeassert

		if err := t.Platform.MacOS.build(rc, out); err != nil {
			return nil, err
		}
	case platform.OSWindows:
		out = &windows.Target{Target: out.(*common.Target)} //nolint:forcetypeassert

		if err := t.Platform.Windows.build(rc, out); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%w: unsupported platform: %s", ErrInvalidInput, p)
	}

	return out, nil
}

/* ------------------------ Interface: exportBuilder ------------------------ */

type exportBuilder interface {
	build(rc *run.Context, dst Exporter) error
}

/* -------------------------------------------------------------------------- */
/*                               Platform: Linux                              */
/* -------------------------------------------------------------------------- */

/* ---------------- Struct: LinuxTargetWithFeaturesAndProfile --------------- */

type LinuxTargetWithFeaturesAndProfile struct {
	Target *linux.Target

	Feature map[string]LinuxTargetWithProfile `toml:"feature"`
	Profile map[engine.Profile]linux.Target   `toml:"profile"`
}

/* --------------------- Struct: LinuxTargetWithProfile --------------------- */

type LinuxTargetWithProfile struct {
	Target *linux.Target

	Profile map[engine.Profile]linux.Target `toml:"profile"`
}

/* --------------------------- Impl: exportBuilder -------------------------- */

// Compile-time check that 'Builder' is implemented.
var _ exportBuilder = (*LinuxTargetWithFeaturesAndProfile)(nil)

func (t *LinuxTargetWithFeaturesAndProfile) build(rc *run.Context, dst Exporter) error {
	// Root-level params
	if err := t.Target.MergeInto(dst); err != nil {
		return err
	}

	// Feature-constrained params
	for _, f := range rc.Features {
		if err := t.Feature[f].Target.MergeInto(dst); err != nil {
			return err
		}
	}

	// Profile-constrained params
	l := t.Profile[rc.Profile]
	if err := l.MergeInto(dst); err != nil {
		return err
	}

	// Feature-and-profile-constrained params
	for _, f := range rc.Features {
		l := t.Feature[f].Profile[rc.Profile]
		if err := l.MergeInto(dst); err != nil {
			return err
		}
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*                               Platform: MacOS                              */
/* -------------------------------------------------------------------------- */

/* ---------------- Struct: MacOSTargetWithFeaturesAndProfile --------------- */

type MacOSTargetWithFeaturesAndProfile struct {
	Target *macos.Target

	Feature map[string]MacOSTargetWithProfile `toml:"feature"`
	Profile map[engine.Profile]macos.Target   `toml:"profile"`
}

/* --------------------- Struct: MacOSTargetWithProfile --------------------- */

type MacOSTargetWithProfile struct {
	Target *macos.Target

	Profile map[engine.Profile]macos.Target `toml:"profile"`
}

/* --------------------------- Impl: exportBuilder -------------------------- */

// Compile-time check that 'Builder' is implemented.
var _ exportBuilder = (*MacOSTargetWithFeaturesAndProfile)(nil)

func (t *MacOSTargetWithFeaturesAndProfile) build(rc *run.Context, dst Exporter) error {
	// Root-level params
	if err := t.Target.MergeInto(dst); err != nil {
		return err
	}

	// Feature-constrained params
	for _, f := range rc.Features {
		if err := t.Feature[f].Target.MergeInto(dst); err != nil {
			return err
		}
	}

	// Profile-constrained params
	l := t.Profile[rc.Profile]
	if err := l.MergeInto(dst); err != nil {
		return err
	}

	// Feature-and-profile-constrained params
	for _, f := range rc.Features {
		l := t.Feature[f].Profile[rc.Profile]
		if err := l.MergeInto(dst); err != nil {
			return err
		}
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*                              Platform: Windows                             */
/* -------------------------------------------------------------------------- */

/* ---------------- Struct: WindowsTargetWithFeaturesAndProfile --------------- */

type WindowsTargetWithFeaturesAndProfile struct {
	Target *windows.Target

	Feature map[string]WindowsTargetWithProfile `toml:"feature"`
	Profile map[engine.Profile]windows.Target   `toml:"profile"`
}

/* --------------------- Struct: WindowsTargetWithProfile --------------------- */

type WindowsTargetWithProfile struct {
	Target *windows.Target

	Profile map[engine.Profile]windows.Target `toml:"profile"`
}

/* --------------------------- Impl: exportBuilder -------------------------- */

// Compile-time check that 'Builder' is implemented.
var _ exportBuilder = (*WindowsTargetWithFeaturesAndProfile)(nil)

func (t *WindowsTargetWithFeaturesAndProfile) build(rc *run.Context, dst Exporter) error {
	// Root-level params
	if err := t.Target.MergeInto(dst); err != nil {
		return err
	}

	// Feature-constrained params
	for _, f := range rc.Features {
		if err := t.Feature[f].Target.MergeInto(dst); err != nil {
			return err
		}
	}

	// Profile-constrained params
	l := t.Profile[rc.Profile]
	if err := l.MergeInto(dst); err != nil {
		return err
	}

	// Feature-and-profile-constrained params
	for _, f := range rc.Features {
		l := t.Feature[f].Profile[rc.Profile]
		if err := l.MergeInto(dst); err != nil {
			return err
		}
	}

	return nil
}
