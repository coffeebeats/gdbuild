package config

import (
	"errors"
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/config/common"
	"github.com/coffeebeats/gdbuild/pkg/config/linux"
	"github.com/coffeebeats/gdbuild/pkg/config/macos"
	"github.com/coffeebeats/gdbuild/pkg/config/windows"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                              Function: Export                              */
/* -------------------------------------------------------------------------- */

// Export creates an `Export` instance which contains an action for exporting
// the specified target.
func Export( //nolint:cyclop,funlen,gocognit
	rc *run.Context,
	m *Manifest,
	tl *template.Template,
	target string,
) (*export.Export, error) {
	var mr merged

	toBuild := []configuration{{context: rc, manifest: m}}
	visited := map[osutil.Path]struct{}{}

	for len(toBuild) > 0 {
		// Remove the next manifest from the queue.
		cfg := toBuild[0]
		toBuild = toBuild[1:]

		// Copy build context so it can be modified.
		rc := *cfg.context

		// First, determine whether this manifest extends another one.

		if err := cfg.manifest.Config.Extends.RelTo(rc.PathManifest); err != nil {
			return nil, fmt.Errorf(
				"%w: cannot find inherited manifest: %w",
				ErrInvalidInput,
				err,
			)
		}

		extends := cfg.manifest.Config.Extends

		// Skip block below if this manifest has already been "visited".
		if _, ok := visited[extends]; !ok && extends != "" {
			baseManifest, err := ParseFile(extends.String())
			if err != nil {
				return nil, fmt.Errorf("cannot parse inherited manifest: %w", err)
			}

			rc.PathManifest = extends

			base := configuration{context: &rc, manifest: baseManifest}
			toBuild = append(toBuild, base, cfg)

			visited[extends] = struct{}{}

			continue
		}

		// Configure 'Godot' properties.
		if err := cfg.manifest.Godot.Configure(&rc); err != nil {
			return nil, err
		}

		// Merge 'Godot' properties.
		if err := cfg.manifest.Godot.MergeInto(&mr.godot); err != nil {
			return nil, err
		}

		tr, ok := cfg.manifest.Target[target]
		if !ok {
			continue
		}

		// Build 'Target' properties.
		t, err := tr.Combine(&rc)
		if err != nil {
			return nil, err
		}

		// Configure 'Target' properties.
		if err := t.Configure(&rc); err != nil {
			return nil, err
		}

		if mr.target == nil {
			mr.target = t

			continue
		}

		// Merge 'Target' properties.
		if err := t.MergeInto(mr.target); err != nil {
			return nil, err
		}
	}

	if mr.target == nil {
		return nil, fmt.Errorf("%w: no target found: %s", ErrInvalidInput, target)
	}

	if err := mr.Validate(rc); err != nil {
		return nil, err
	}

	ev, err := mr.godot.ParseVersion()
	if err != nil {
		if errors.Is(err, ErrConflictingValue) {
			return nil, fmt.Errorf("%w: 'src_path' is unsupported at this time", err)
		}

		return nil, err
	}

	xp := mr.target.Collect(rc, tl, ev)

	// Set the encryption key on the template builds in the event that the key
	// was just set on the target. This is the only property that needs to be
	// synchronized between the target/template builds, so do it here.
	for i, tb := range tl.Builds {
		if tb.EncryptionKey != "" && xp.EncryptionKey == "" {
			return nil, fmt.Errorf(
				"%w: template has encryption key set but target does not",
				ErrInvalidInput,
			)
		}

		if tb.EncryptionKey != xp.EncryptionKey {
			tb.EncryptionKey = xp.EncryptionKey
			tl.Builds[i] = tb // Update the slice since 'tb' is a copy.
		}
	}

	return xp, nil
}

/* ----------------------------- Struct: merged ----------------------------- */

type merged struct {
	godot  Godot
	target Exporter
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (m *merged) Validate(rc *run.Context) error {
	if m.target == nil {
		return fmt.Errorf("%w: failed to build target", ErrMissingInput)
	}

	// Validate 'Target' properties.
	if err := m.godot.Validate(rc); err != nil {
		return err
	}

	if err := m.target.Validate(rc); err != nil {
		return err
	}

	return nil
}

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
		out := &macos.Target{Target: base} //nolint:exhaustruct

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
