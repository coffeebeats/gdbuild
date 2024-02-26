package manifest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/build"
	"github.com/coffeebeats/gdbuild/pkg/build/template"
)

var ErrInvalidInput = errors.New("invalid input")

/* -------------------------------------------------------------------------- */
/*                              Struct: Manifest                              */
/* -------------------------------------------------------------------------- */

// Manifest defines the supported structure of the GDBuild manifest file.
type Manifest struct {
	// Project contains project-wide settings, like the Godot version.
	Project Project `toml:"project"`
	// Target contains all exportable artifact specifications.
	Target map[string]Target `toml:"target"`
	// Template includes settings for building custom export templates.
	Template Template `toml:"template"`
}

/* --------------------------- Function: Filename --------------------------- */

// Filename returns the name of the GDBuild manifest file.
func Filename() string {
	return "gdbuild.toml"
}

/* -------------------------- Method: BuildTemplate ------------------------- */

func (m *Manifest) BuildTemplate( //nolint:cyclop,funlen,gocognit,gocyclo,ireturn
	pathManifest,
	pathBuild string,
	pl build.OS,
	pr build.Profile,
	ff ...string,
) (action.Action, error) {
	base := m.Template.Base
	if base == nil {
		base = &template.Base{} //nolint:exhaustruct
	}

	pathManifest, err := filepath.Abs(pathManifest)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(pathManifest)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		if err := os.MkdirAll(pathManifest, osutil.ModeUserRWXGroupRX); err != nil {
			return nil, err
		}
	}

	if !info.IsDir() {
		pathManifest = filepath.Dir(pathManifest)
	}

	pathBuild, err = filepath.Abs(pathBuild)
	if err != nil {
		return nil, err
	}

	base.Invocation = build.Invocation{
		Features:     ff,
		PathBuild:    build.Path(pathBuild),
		PathManifest: build.Path(pathManifest),
		Platform:     pl,
		Profile:      pr,
	}

	base.Godot = m.Project.Godot
	if err := base.Godot.Validate(); err != nil {
		return nil, err
	}

	// Merge template base.
	if err := base.Configure(&base.Invocation); err != nil {
		return nil, err
	}

	var cmd action.Actioner

	// Merge platform-specific template.
	switch pl {
	case build.OSAndroid:
		base := template.Android{Base: base} //nolint:exhaustruct

		template := m.Template.Platform.Android
		if err := template.Configure(&base.Invocation); err != nil {
			return nil, err
		}

		if err := base.Merge(template.Android); err != nil {
			return nil, err
		}

		cmd = &base
	case build.OSIOS:
		base := template.IOS{Base: base} //nolint:exhaustruct

		template := m.Template.Platform.IOS
		if err := template.Configure(&base.Invocation); err != nil {
			return nil, err
		}

		if err := base.Merge(template.IOS); err != nil {
			return nil, err
		}

		cmd = &base
	case build.OSLinux:
		base := template.Linux{Base: base}

		template := m.Template.Platform.Linux
		if err := template.Configure(&base.Invocation); err != nil {
			return nil, err
		}

		if err := base.Merge(template.Linux); err != nil {
			return nil, err
		}

		cmd = &base
	case build.OSMacOS:
		base := template.MacOS{Base: base} //nolint:exhaustruct

		template := m.Template.Platform.MacOS
		if err := template.Configure(&base.Invocation); err != nil {
			return nil, err
		}

		if err := base.Merge(template.MacOS); err != nil {
			return nil, err
		}

		cmd = &base
	case build.OSWeb:
		base := template.Web{Base: base} //nolint:exhaustruct

		template := m.Template.Platform.Web
		if err := template.Configure(&base.Invocation); err != nil {
			return nil, err
		}

		if err := base.Merge(template.Web); err != nil {
			return nil, err
		}

		cmd = &base
	case build.OSWindows:
		base := template.Windows{Base: base} //nolint:exhaustruct

		template := m.Template.Platform.Windows
		if err := template.Configure(&base.Invocation); err != nil {
			return nil, err
		}

		if err := base.Merge(template.Windows); err != nil {
			return nil, err
		}

		cmd = &base
	default:
		return nil, fmt.Errorf("%w: unsupported platform", ErrInvalidInput)
	}

	if value, ok := cmd.(build.Configurer); ok {
		if err := value.Configure(&base.Invocation); err != nil {
			return nil, err
		}
	}

	if value, ok := cmd.(build.Validater); ok {
		if err := value.Validate(); err != nil {
			return nil, err
		}
	}

	return cmd.Action()
}

/* ------------------------- Function: getOrDefault ------------------------- */

// getOrDefault is a convenience method to safely access a value from a
// potentially nil map.
func getOrDefault[K comparable, V any](m map[K]V, key K) V { //nolint:ireturn
	if m == nil {
		return *new(V)
	}

	return m[key]
}
