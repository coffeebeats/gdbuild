package manifest

import "github.com/coffeebeats/gdbuild/pkg/build"

/* -------------------------------------------------------------------------- */
/*                              Struct: Manifest                              */
/* -------------------------------------------------------------------------- */

// Manifest defines the supported structure of the GDBuild manifest file.
type Manifest struct {
	// Project contains project-wide settings, like the Godot version.
	Project *Project `json:"project" toml:"project"`
	// Target contains all exportable artifact specifications.
	Target map[string]*Target `json:"target" toml:"target"`
	// Template includes settings for building custom export templates.
	Template *Template `json:"template" toml:"template"`
}

/* --------------------------- Function: Filename --------------------------- */

// Filename returns the name of the GDBuild manifest file.
func Filename() string {
	return "gdbuild.toml"
}

/* --------------------------- Method: BuildTarget -------------------------- */

func (m *Manifest) BuildTarget(name string, pl build.OS, pr build.Profile, ff ...string) *build.Target {
	target := m.Target[name]
	if target != nil {
		ff = append(target.DefaultFeatures, ff...)
	}

	return target.merge(pl, pr, ff...)
}

/* -------------------------- Method: BuildTemplate ------------------------- */

func (m *Manifest) BuildTemplate(pl build.OS, pr build.Profile, ff ...string) *build.Template {
	return m.Template.merge(pl, pr, ff...)
}

/* ---------------------------- Method: Validate ---------------------------- */

// Validate checks that the 'Manifest' contents are valid.
//
// TODO: Implement this method, as well as for all contained types.
func (m *Manifest) Validate() error {
	return nil
}
