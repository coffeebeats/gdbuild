package platform

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
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                             Interface: Exporter                            */
/* -------------------------------------------------------------------------- */

type Exporter interface {
	config.Configurable[*run.Context]

	Collect(rc *run.Context, tl *template.Template, ev engine.Version) *export.Export
}

/* -------------------------------------------------------------------------- */
/*                              Struct: Targets                             */
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
	*common.TargetWithFeaturesAndProfile

	Platform TargetPlatforms `toml:"platform"`
}

/* -------------------- Struct: BaseTargetWithoutFeature -------------------- */

type BaseTargetWithoutFeature struct {
	*common.Target

	Profile map[engine.Profile]common.Target `toml:"profile"`
}

/* ---------------------------- Struct: Platforms --------------------------- */

type TargetPlatforms struct {
	Linux   linux.TargetWithFeaturesAndProfile   `toml:"linux"`
	MacOS   macos.TargetWithFeaturesAndProfile   `toml:"macos"`
	Windows windows.TargetWithFeaturesAndProfile `toml:"windows"`
}

/* ------------------------ Interface: TargetBuilder ------------------------ */

type TargetBuilder[T Exporter] interface {
	Build(rc *run.Context, dst T) error
}

// Compile-time check that 'Builder' is implemented.
var _ TargetBuilder[*common.Target] = (*common.TargetWithFeaturesAndProfile)(nil)
var _ TargetBuilder[*linux.Target] = (*linux.TargetWithFeaturesAndProfile)(nil)
var _ TargetBuilder[*macos.Target] = (*macos.TargetWithFeaturesAndProfile)(nil)
var _ TargetBuilder[*windows.Target] = (*windows.TargetWithFeaturesAndProfile)(nil)

/* ----------------------------- Method: Combine ---------------------------- */

func (t Targets) Combine(rc *run.Context) (Exporter, error) { //nolint:dupl,ireturn
	// Root params.
	base := new(common.Target)

	if err := t.TargetWithFeaturesAndProfile.Build(rc, base); err != nil {
		return nil, err
	}

	switch p := rc.Platform; p {
	case platform.OSLinux:
		out := &linux.Target{Target: base}

		if err := t.Platform.Linux.Build(rc, out); err != nil {
			return nil, err
		}

		return out, nil
	case platform.OSMacOS:
		out := &macos.Target{Target: base}

		if err := t.Platform.MacOS.Build(rc, out); err != nil {
			return nil, err
		}

		return out, nil
	case platform.OSWindows:
		out := &windows.Target{Target: base}

		if err := t.Platform.Windows.Build(rc, out); err != nil {
			return nil, err
		}

		return out, nil
	default:
		return nil, fmt.Errorf("%w: unsupported platform: %s", config.ErrInvalidInput, p)
	}
}
