package manifest

import (
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Project                              */
/* -------------------------------------------------------------------------- */

// Project defines the project-wide configuration which affects all exportable
// artifacts.
type Project struct {
	Icon        string  `json:"icon"        toml:"icon"`
	Version     Version `json:"version"     toml:"version"`
	VersionFile string  `json:"versionFile" toml:"version_file"`
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
