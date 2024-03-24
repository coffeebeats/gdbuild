package target

import (
	"errors"
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/godot/build"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrMissingInput = errors.New("missing input")
)

/* -------------------------------------------------------------------------- */
/*                             Interface: Exporter                            */
/* -------------------------------------------------------------------------- */

type Exporter interface {
	config.Configurable[*build.Context]

	Export(src build.Source, bc *build.Context) *export.Export
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
	*Base

	Platform Platforms                     `toml:"platform"`
	Feature  map[string]BaseWithoutFeature `toml:"feature"`
	Profile  map[build.Profile]Base        `toml:"profile"`
}

/* ----------------------- Struct: BaseWithoutFeature ----------------------- */

type BaseWithoutFeature struct {
	*Base

	Profile map[build.Profile]Base `toml:"profile"`
}

/* ---------------------------- Struct: Platforms --------------------------- */

type Platforms struct {
	Linux   LinuxWithFeaturesAndProfile   `toml:"linux"`
	MacOS   MacOSWithFeaturesAndProfile   `toml:"macos"`
	Windows WindowsWithFeaturesAndProfile `toml:"windows"`
}

/* ------------------------------ Method: Build ----------------------------- */

func (t *Targets) Build(bc *build.Context) (Exporter, error) { //nolint:cyclop,ireturn
	// Target params (root)
	var out Exporter = new(Base)

	// Target params (root)
	if err := t.Base.MergeInto(out); err != nil {
		return nil, err
	}

	// Target params (feature-constrained)
	for _, f := range bc.Features {
		bwof := t.Feature[f].Base
		if err := bwof.MergeInto(out); err != nil {
			return nil, err
		}
	}

	// Target params (profile-constrained)
	b := t.Profile[bc.Profile]
	if err := b.MergeInto(out); err != nil {
		return nil, err
	}

	// Feature-and-profile-constrained params
	for _, f := range bc.Features {
		bwof := t.Feature[f].Profile[bc.Profile]
		if err := bwof.MergeInto(out); err != nil {
			return nil, err
		}
	}

	switch p := bc.Platform; p {
	case platform.OSLinux:
		out = &Linux{Base: out.(*Base)} //nolint:forcetypeassert

		if err := t.Platform.Linux.build(bc, out); err != nil {
			return nil, err
		}
	case platform.OSMacOS:
		out = &MacOS{Base: out.(*Base)} //nolint:forcetypeassert

		if err := t.Platform.MacOS.build(bc, out); err != nil {
			return nil, err
		}
	case platform.OSWindows:
		out = &Windows{Base: out.(*Base)} //nolint:forcetypeassert

		if err := t.Platform.Windows.build(bc, out); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%w: unsupported platform: %s", ErrInvalidInput, p)
	}

	return out, nil
}

/* ----------------------- Interface: templateBuilder ----------------------- */

type templateBuilder interface {
	build(bc *build.Context, dst Exporter) error
}

/* -------------------------------------------------------------------------- */
/*                               Platform: Linux                              */
/* -------------------------------------------------------------------------- */

/* ------------------ Struct: LinuxWithFeaturesAndProfile ----------------- */

type LinuxWithFeaturesAndProfile struct {
	*Linux

	Feature map[string]LinuxWithProfile `toml:"feature"`
	Profile map[build.Profile]Linux     `toml:"profile"`
}

/* ----------------------- Struct: LinuxWithProfile ----------------------- */

type LinuxWithProfile struct {
	*Linux

	Profile map[build.Profile]Linux `toml:"profile"`
}

/* -------------------------- Impl: templateBuilder ------------------------- */

// Compile-time check that 'Builder' is implemented.
var _ templateBuilder = (*LinuxWithFeaturesAndProfile)(nil)

func (t *LinuxWithFeaturesAndProfile) build(bc *build.Context, dst Exporter) error {
	// Root-level params
	if err := t.Linux.MergeInto(dst); err != nil {
		return err
	}

	// Feature-constrained params
	for _, f := range bc.Features {
		if err := t.Feature[f].Linux.MergeInto(dst); err != nil {
			return err
		}
	}

	// Profile-constrained params
	l := t.Profile[bc.Profile]
	if err := l.MergeInto(dst); err != nil {
		return err
	}

	// Feature-and-profile-constrained params
	for _, f := range bc.Features {
		l := t.Feature[f].Profile[bc.Profile]
		if err := l.MergeInto(dst); err != nil {
			return err
		}
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*                               Platform: MacOS                              */
/* -------------------------------------------------------------------------- */

/* ------------------ Struct: MacOSWithFeaturesAndProfile ----------------- */

type MacOSWithFeaturesAndProfile struct {
	*MacOS

	Feature map[string]MacOSWithProfile `toml:"feature"`
	Profile map[build.Profile]MacOS     `toml:"profile"`
}

/* ----------------------- Struct: MacOSWithProfile ----------------------- */

type MacOSWithProfile struct {
	*MacOS

	Profile map[build.Profile]MacOS `toml:"profile"`
}

/* -------------------------- Impl: templateBuilder ------------------------- */

// Compile-time check that 'Builder' is implemented.
var _ templateBuilder = (*MacOSWithFeaturesAndProfile)(nil)

func (t *MacOSWithFeaturesAndProfile) build(bc *build.Context, dst Exporter) error {
	// Root-level params
	if err := t.MacOS.MergeInto(dst); err != nil {
		return err
	}

	// Feature-constrained params
	for _, f := range bc.Features {
		if err := t.Feature[f].MacOS.MergeInto(dst); err != nil {
			return err
		}
	}

	// Profile-constrained params
	l := t.Profile[bc.Profile]
	if err := l.MergeInto(dst); err != nil {
		return err
	}

	// Feature-and-profile-constrained params
	for _, f := range bc.Features {
		l := t.Feature[f].Profile[bc.Profile]
		if err := l.MergeInto(dst); err != nil {
			return err
		}
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*                              Platform: Windows                             */
/* -------------------------------------------------------------------------- */

/* ------------------ Struct: WindowsWithFeaturesAndProfile ----------------- */

type WindowsWithFeaturesAndProfile struct {
	*Windows

	Feature map[string]WindowsWithProfile `toml:"feature"`
	Profile map[build.Profile]Windows     `toml:"profile"`
}

/* ----------------------- Struct: WindowsWithProfile ----------------------- */

type WindowsWithProfile struct {
	*Windows

	Profile map[build.Profile]Windows `toml:"profile"`
}

/* -------------------------- Impl: templateBuilder ------------------------- */

// Compile-time check that 'Builder' is implemented.
var _ templateBuilder = (*WindowsWithFeaturesAndProfile)(nil)

func (t *WindowsWithFeaturesAndProfile) build(bc *build.Context, dst Exporter) error {
	// Root-level params
	if err := t.Windows.MergeInto(dst); err != nil {
		return err
	}

	// Feature-constrained params
	for _, f := range bc.Features {
		if err := t.Feature[f].Windows.MergeInto(dst); err != nil {
			return err
		}
	}

	// Profile-constrained params
	l := t.Profile[bc.Profile]
	if err := l.MergeInto(dst); err != nil {
		return err
	}

	// Feature-and-profile-constrained params
	for _, f := range bc.Features {
		l := t.Feature[f].Profile[bc.Profile]
		if err := l.MergeInto(dst); err != nil {
			return err
		}
	}

	return nil
}
