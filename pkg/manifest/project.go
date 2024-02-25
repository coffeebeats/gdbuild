package manifest

import (
	"github.com/coffeebeats/gdenv/pkg/godot/version"

	"github.com/coffeebeats/gdbuild/pkg/build"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Project                              */
/* -------------------------------------------------------------------------- */

// Project defines the project-wide configuration which affects all exportable
// artifacts.
type Project struct {
	// Godot contains a specification for which Godot version to use.
	Godot build.Godot `toml:"godot"`
}

/* -------------------------------------------------------------------------- */
/*                                Type: Version                               */
/* -------------------------------------------------------------------------- */

// Version wraps 'version.Version' and implements TOML unmarshaling for it.
type Version version.Version

/* ---------------------- Impl: encoding.UnmarshalText ---------------------- */

func (v *Version) UnmarshalText(bb []byte) error {
	value, err := version.Parse(string(bb))
	if err != nil {
		return err
	}

	*v = Version(value)

	return nil
}
