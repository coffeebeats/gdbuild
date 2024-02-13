package manifest

/* -------------------------------------------------------------------------- */
/*                              Struct: Manifest                              */
/* -------------------------------------------------------------------------- */

// Manifest defines the supported structure of the GDBuild manifest file.
type Manifest struct {
	// Project contains project-wide settings, like the Godot version.
	Project *Project `toml:"project"`
	// Target contains all exportable artifact specifications.
	Target map[string]*Target `toml:"target"`
	// Template includes settings for building custom export templates.
	Template *Template `toml:"template"`
}

/* ---------------------------- Method: Validate ---------------------------- */

// Validate checks that the 'Manifest' contents are valid.
//
// TODO: Implement this method, as well as for all contained types.
func (m *Manifest) Validate() error {
	return nil
}

/* --------------------------- Function: Filename --------------------------- */

// Filename returns the name of the GDBuil manifest file.
func Filename() string {
	return "gdbuild.toml"
}
