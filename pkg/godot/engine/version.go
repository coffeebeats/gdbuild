package engine

import (
	"errors"

	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Version                              */
/* -------------------------------------------------------------------------- */

// Version is a type wrapper around a 'version.Version' struct from 'gdenv'.
type Version version.Version

/* ----------------------------- Method: IsZero ----------------------------- */

// IsZero returns whether the 'Version' struct is equal to the zero value.
func (v Version) IsZero() bool {
	return v == Version{} //nolint:exhaustruct
}

/* --------------------------- Impl: fmt.Stringer --------------------------- */

func (v Version) String() string {
	return version.Version(v).String()
}

/* ---------------------- Impl: encoding.UnmarshalText ---------------------- */

func (v *Version) UnmarshalText(text []byte) error {
	parsed, err := version.Parse(string(text))
	if err != nil {
		if !errors.Is(err, version.ErrMissing) {
			return err
		}

		*v = Version{} //nolint:exhaustruct

		return nil
	}

	*v = Version(parsed)

	return nil
}
