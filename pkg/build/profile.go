package build

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMissingProfile      = errors.New("missing profile")
	ErrUnrecognizedProfile = errors.New("unrecognized profile")
)

/* -------------------------------------------------------------------------- */
/*                                Enum: Profile                               */
/* -------------------------------------------------------------------------- */

// Optimization levels supported by Godot.
type Profile uint

const (
	ProfileDebug Profile = iota + 1
	ProfileReleaseDebug
	ProfileRelease
)

/* ----------------------------- Impl: Stringer ----------------------------- */

// String implements 'fmt.Stringer' for 'Profile' according to the profile names
// passed to SCons during compilation.
func (p Profile) String() string {
	switch p {
	case ProfileDebug:
		return "debug"
	case ProfileReleaseDebug:
		return "release_debug"
	case ProfileRelease:
		return "release"
	default:
		return ""
	}
}

/* ------------------------- Function: ParseProfile ------------------------- */

// Parses an input string as a build 'Profile' optimization level.
func ParseProfile(input string) (Profile, error) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "":
		return 0, ErrMissingProfile

	case "dbg", "debug":
		return ProfileDebug, nil

	case "release_debug", "releasedebug", "release_dbg", "releasedbg":
		return ProfileReleaseDebug, nil

	case "release":
		return ProfileRelease, nil

	default:
		return 0, fmt.Errorf("%w: '%s'", ErrUnrecognizedProfile, input)
	}
}

/* ----------------------- Function: MustParseProfile ----------------------- */

// Parses an input string as a build profile specification but panics if it
// would fail.
func MustParseProfile(input string) Profile {
	arch, err := ParseProfile(input)
	if err != nil {
		panic(err)
	}

	return arch
}

/* ---------------------- Impl: encoding.UnmarshalText ---------------------- */

func (p *Profile) UnmarshalText(bb []byte) error {
	value, err := ParseProfile(string(bb))
	if err != nil {
		return err
	}

	*p = value

	return nil
}
