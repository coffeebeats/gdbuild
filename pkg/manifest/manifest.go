package manifest

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/coffeebeats/gdbuild/pkg/build"
	"github.com/coffeebeats/gdbuild/pkg/build/platform"
	"github.com/coffeebeats/gdbuild/pkg/build/template"
)

var ErrInvalidInput = errors.New("invalid input")

/* -------------------------------------------------------------------------- */
/*                              Struct: Manifest                              */
/* -------------------------------------------------------------------------- */

// Manifest defines the supported structure of the GDBuild manifest file.
type Manifest struct {
	// Config contains GDBuild configuration-related settings.
	Config Config `toml:"config"`
	// Godot contains settings on which Godot version/source code to use.
	Godot build.Godot `toml:"godot"`
	// Target contains all exportable artifact specifications.
	Target map[string]Target `toml:"target"`
	// Template includes settings for building custom export templates.
	Template Template `toml:"template"`

	// Parent is a reference to another 'Manifest' from which this one inherits
	// properties from. This must be set manually.
	Parent *Manifest
}

/* --------------------------- Function: Filename --------------------------- */

// Filename returns the name of the GDBuild manifest file.
func Filename() string {
	return "gdbuild.toml"
}

/* -------------------------- Method: BuildTemplate ------------------------- */

func (m *Manifest) BuildTemplate( //nolint:cyclop,funlen,gocognit,ireturn
	inv build.Invocation,
) (template.Template, error) {
	// First, determine whether this manifest extends another one.
	if err := m.Config.Extends.RelTo(inv.PathManifest); err != nil {
		return nil, fmt.Errorf(
			"%w: cannot find inherited manifest: %w",
			ErrInvalidInput,
			err,
		)
	}

	// If it doesn't, simply build the template from this manifest alone.
	if m.Config.Extends == "" {
		t, err := m.mergeTemplateForInvocation(&inv)
		if err != nil {
			return nil, err
		}

		if err := t.Configure(&inv); err != nil {
			return nil, err
		}

		if err := t.Validate(); err != nil {
			return nil, err
		}

		return t, nil
	}

	baseManifest, err := ParseFile(string(m.Config.Extends))
	if err != nil {
		return nil, fmt.Errorf(
			"%w: cannot parse inherited manifest: %w",
			ErrInvalidInput,
			err,
		)
	}

	m.Parent = baseManifest

	baseInv := inv
	baseInv.PathManifest = build.Path(filepath.Dir(string(m.Config.Extends)))

	baseTemplate, err := baseManifest.BuildTemplate(baseInv)
	if err != nil {
		return nil, err
	}

	childTemplate, err := m.mergeTemplateForInvocation(&inv)
	if err != nil {
		return nil, err
	}

	var out template.Template

	switch base := baseTemplate.(type) {
	case *template.Android:
		child, ok := childTemplate.(*template.Android)
		if !ok {
			return nil, fmt.Errorf("%w: incompatible template type", ErrInvalidInput)
		}

		if err := base.Merge(child); err != nil {
			return nil, err
		}

		out = base
	case *template.IOS:
		child, ok := childTemplate.(*template.IOS)
		if !ok {
			return nil, fmt.Errorf("%w: incompatible template type", ErrInvalidInput)
		}

		if err := base.Merge(child); err != nil {
			return nil, err
		}

		out = base
	case *template.Linux:
		child, ok := childTemplate.(*template.Linux)
		if !ok {
			return nil, fmt.Errorf("%w: incompatible template type", ErrInvalidInput)
		}

		if err := base.Merge(child); err != nil {
			return nil, err
		}

		out = base
	case *template.MacOS:
		child, ok := childTemplate.(*template.MacOS)
		if !ok {
			return nil, fmt.Errorf("%w: incompatible template type", ErrInvalidInput)
		}

		if err := base.Merge(child); err != nil {
			return nil, err
		}

		out = base
	case *template.Web:
		child, ok := childTemplate.(*template.Web)
		if !ok {
			return nil, fmt.Errorf("%w: incompatible template type", ErrInvalidInput)
		}

		if err := base.Merge(child); err != nil {
			return nil, err
		}

		out = base
	case *template.Windows:
		child, ok := childTemplate.(*template.Windows)
		if !ok {
			return nil, fmt.Errorf("%w: incompatible template type", ErrInvalidInput)
		}

		if err := base.Merge(child); err != nil {
			return nil, err
		}

		out = base
	default:
		return nil, fmt.Errorf("%w: unknown platform type: %T", ErrInvalidInput, base)
	}

	if err := out.Configure(&inv); err != nil {
		return nil, err
	}

	if err := out.Validate(); err != nil {
		return nil, err
	}

	return out, nil
}

/* ------------------- Method: mergeTemplateForInvocation ------------------- */

func (m *Manifest) mergeTemplateForInvocation( //nolint:cyclop,funlen,ireturn
	inv *build.Invocation,
) (template.Template, error) {
	base := m.Template.Base
	if base == nil {
		base = &template.Base{} //nolint:exhaustruct
	}

	if err := inv.Validate(); err != nil {
		return nil, err
	}

	base.Invocation = *inv
	base.Godot = m.Godot

	var out template.Template

	// Merge platform-specific template.
	switch inv.Platform {
	case platform.OSAndroid:
		base := template.Android{Base: base} //nolint:exhaustruct

		t := m.Template.Platform.Android
		if err := t.Configure(&base.Invocation); err != nil {
			return nil, err
		}

		if err := base.Merge(t.Android); err != nil {
			return nil, err
		}

		out = &base
	case platform.OSIOS:
		base := template.IOS{Base: base} //nolint:exhaustruct

		template := m.Template.Platform.IOS
		if err := template.Configure(&base.Invocation); err != nil {
			return nil, err
		}

		if err := base.Merge(template.IOS); err != nil {
			return nil, err
		}

		out = &base
	case platform.OSLinux:
		base := template.Linux{Base: base}

		template := m.Template.Platform.Linux
		if err := template.Configure(&base.Invocation); err != nil {
			return nil, err
		}

		if err := base.Merge(template.Linux); err != nil {
			return nil, err
		}

		out = &base
	case platform.OSMacOS:
		base := template.MacOS{Base: base} //nolint:exhaustruct

		template := m.Template.Platform.MacOS
		if err := template.Configure(&base.Invocation); err != nil {
			return nil, err
		}

		if err := base.Merge(template.MacOS); err != nil {
			return nil, err
		}

		out = &base
	case platform.OSWeb:
		base := template.Web{Base: base} //nolint:exhaustruct

		template := m.Template.Platform.Web
		if err := template.Configure(&base.Invocation); err != nil {
			return nil, err
		}

		if err := base.Merge(template.Web); err != nil {
			return nil, err
		}

		out = &base
	case platform.OSWindows:
		base := template.Windows{Base: base} //nolint:exhaustruct

		template := m.Template.Platform.Windows
		if err := template.Configure(&base.Invocation); err != nil {
			return nil, err
		}

		if err := base.Merge(template.Windows); err != nil {
			return nil, err
		}

		out = &base
	default:
		return nil, fmt.Errorf("%w: unsupported platform", ErrInvalidInput)
	}

	return out, nil
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
