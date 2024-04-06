package common

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/pkg/godot/engine"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrMissingInput = errors.New("missing input")
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Target                               */
/* -------------------------------------------------------------------------- */

// Target specifies a single, platform-agnostic exportable artifact within the
// Godot project.
type Target struct {
	// DefaultFeatures contains the slice of Godot project feature tags to build
	// with.
	DefaultFeatures []string `toml:"default_features"`
	// Encrypt sets whether the exported artifacts will be encrypted or not.
	Encrypt *bool `toml:"encrypt"`
	// EncryptionKey is the encryption key to encrypt game assets with.
	EncryptionKey string `toml:"encryption_key"`
	// Hook defines commands to be run before or after the target artifact is
	// generated.
	Hook run.Hook `toml:"hook"`
	// Options are 'export_presets.cfg' overrides, specifically the preset
	// 'options' table, for the exported artifact.
	Options map[string]any `toml:"options"`
	// PackFiles defines the game files exported as part of this artifact.
	PackFiles []export.PackFile `toml:"pack_files"`
	// Runnable is whether the export artifact should be executable. This should
	// be true for client and server targets and false for artifacts like DLC.
	Runnable *bool `toml:"runnable"`
	// Server configures the target as a server-only executable, enabling some
	// optimizations like disabling graphics.
	Server *bool `toml:"server"`
}

/* ----------------------------- Impl: Exporter ----------------------------- */

func (t *Target) Collect(rc *run.Context, tl *template.Template, ev engine.Version) *export.Export {
	// Set the encryption key environment variable; see
	// https://docs.godotengine.org/en/stable/contributing/development/compiling/compiling_with_script_encryption_key.html.
	encryptionKey, _ := resolveEncryptionKey(t.EncryptionKey)

	ff := make([]string, 0, len(t.DefaultFeatures)+len(rc.Features))
	ff = append(ff, t.DefaultFeatures...)
	ff = append(ff, rc.Features...)

	return &export.Export{
		EncryptionKey: encryptionKey,
		Features:      ff,
		Options:       t.Options,
		PackFiles:     t.PackFiles,
		PathTemplate:  "",
		RunBefore:     t.Hook.PreActions(rc),
		RunAfter:      t.Hook.PostActions(rc),
		Runnable:      config.Dereference(t.Runnable),
		Server:        config.Dereference(t.Server),
		Template:      tl,
		Version:       ev,
	}
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (t *Target) Configure(rc *run.Context) error {
	hasEncrypt := false
	isEncrypted := config.Dereference(t.Encrypt)

	for _, pf := range t.PackFiles {
		if err := pf.Configure(rc); err != nil {
			return err
		}

		isPackFileEncrypted := config.Dereference(pf.Encrypt)
		hasEncrypt = hasEncrypt || isPackFileEncrypted

		// Disable pack file encryption if the top level encrypt flag is off.
		encrypt := isEncrypted && isPackFileEncrypted
		pf.Encrypt = &encrypt
	}

	if t.EncryptionKey != "" && !isEncrypted {
		log.Warn("ignoring encryption key because encryption is disabled.")

		t.EncryptionKey = ""
	}

	if isEncrypted && !hasEncrypt {
		log.Warn("encryption was enabled but no encrypted pack files included.")

		disable := false
		t.Encrypt = &disable
		t.EncryptionKey = ""
	}

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (t *Target) Validate(rc *run.Context) error { //nolint:cyclop,funlen
	if err := t.Hook.Validate(rc); err != nil {
		return err
	}

	encryptionKey, err := resolveEncryptionKey(t.EncryptionKey)
	if err != nil {
		return err
	}

	hasEmbed := false
	hasVisualsStripped := false
	packNames := make(map[string]struct{})

	isRunnable := config.Dereference(t.Runnable)
	isServer := config.Dereference(t.Server)

	if isServer && !isRunnable {
		return fmt.Errorf(
			"%w: cannot specify server optimizations for a non-runnable target",
			ErrInvalidInput,
		)
	}

	projectFiles := map[string]struct{}{}

	for i, pf := range t.PackFiles {
		if err := pf.Validate(rc); err != nil {
			return err
		}

		name := pf.Filename(rc.Platform, rc.Target, i)
		if _, ok := packNames[name]; ok {
			return fmt.Errorf(
				"%w: duplicate pack filename found: %s",
				ErrInvalidInput,
				name,
			)
		}

		hasEmbed = hasEmbed || config.Dereference(pf.Embed)
		hasVisualsStripped = hasVisualsStripped || pf.StripVisuals()

		ff, err := pf.Files(rc.PathWorkspace)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}

			return err
		}

		for _, f := range ff {
			path := f.String()

			if _, ok := projectFiles[path]; ok {
				log.Warn("Same file found in multiple packs.", "file", path)

				continue
			}

			projectFiles[path] = struct{}{}
		}
	}

	if !isRunnable && hasEmbed {
		return fmt.Errorf(
			"%w: cannot embed a pack file into a non-runnable target",
			ErrInvalidInput,
		)
	}

	if isRunnable && !hasEmbed {
		return fmt.Errorf(
			"%w: missing embedded pack file for runnable target",
			ErrMissingInput,
		)
	}

	if !isServer && hasVisualsStripped {
		return fmt.Errorf(
			"%w: cannot strip visuals from a pack file for a non-server target",
			ErrInvalidInput,
		)
	}

	if config.Dereference(t.Encrypt) && encryptionKey == "" {
		return fmt.Errorf(
			"%w: encryption is enabled but no encryption key is set",
			ErrInvalidInput,
		)
	}

	return nil
}

/* --------------------------- Impl: config.Merger -------------------------- */

func (t *Target) MergeInto(other any) error {
	if t == nil || other == nil {
		return nil
	}

	dst, ok := other.(*Target)
	if !ok {
		return fmt.Errorf(
			"%w: expected a '%T' but was '%T'",
			config.ErrInvalidInput,
			new(Target),
			other,
		)
	}

	return config.Merge(dst, *t)
}

/* -------------------------------------------------------------------------- */
/*                    Struct: TargetWithFeaturesAndProfile                    */
/* -------------------------------------------------------------------------- */

type TargetWithFeaturesAndProfile struct {
	*Target

	Feature map[string]TargetWithProfile `toml:"feature"`
	Profile map[engine.Profile]Target    `toml:"profile"`
}

/* ------------------------ Struct: TargetWithProfile ----------------------- */

type TargetWithProfile struct {
	*Target

	Profile map[engine.Profile]Target `toml:"profile"`
}

/* ---------------------- Impl: platform.targetBuilder ---------------------- */

func (t *TargetWithFeaturesAndProfile) Build(rc *run.Context, dst *Target) error {
	if t == nil {
		return nil
	}

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
