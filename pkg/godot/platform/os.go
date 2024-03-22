package platform

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMissingOS      = errors.New("missing OS")
	ErrUnrecognizedOS = errors.New("unrecognized OS")
)

/* -------------------------------------------------------------------------- */
/*                                  Enum: OS                                  */
/* -------------------------------------------------------------------------- */

// Operating systems which the Godot project supports building for.
type OS uint

const (
	OSUnknown OS = iota
	OSAndroid
	OSIOS
	OSLinux
	OSMacOS
	OSWeb
	OSWindows
)

/* ----------------------------- Impl: Stringer ----------------------------- */

// String implements 'fmt.Stringer' for 'OS' according to the platform names
// passed to SCons during compilation.
func (o OS) String() string {
	switch o {
	case OSAndroid:
		return "android"
	case OSIOS:
		return "ios"
	case OSLinux:
		return "linuxbsd"
	case OSMacOS:
		return "macos"
	case OSWeb:
		return "web"
	case OSWindows:
		return "windows"
	default:
		return ""
	}
}

/* ---------------------------- Function: ParseOS --------------------------- */

// Parses an input string as an operating system specification.
//
// NOTE: This is a best effort implementation. Please open an issue on GitHub
// if some values are missing:
// github.com/coffeebeats/gdbuild/issues/new?labels=bug&template=%F0%9F%90%9B-bug-report.md.
func ParseOS(input string) (OS, error) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "":
		return 0, ErrMissingOS

	case "android":
		return OSAndroid, nil

	case "ios":
		return OSIOS, nil

	case "linux", "linuxbsd", "x11":
		return OSLinux, nil

	case "darwin", "macos", "osx":
		return OSMacOS, nil

	case "web":
		return OSWeb, nil

	case "win", "windows":
		return OSWindows, nil

	default:
		return 0, fmt.Errorf("%w: '%s'", ErrUnrecognizedOS, input)
	}
}

/* -------------------------- Function: MustParseOS ------------------------- */

// Parses an input string as an operating system specification but panics if it
// would fail.
func MustParseOS(input string) OS {
	os, err := ParseOS(input)
	if err != nil {
		panic(err)
	}

	return os
}

/* ---------------------- Impl: encoding.UnmarshalText ---------------------- */

func (o *OS) UnmarshalText(bb []byte) error {
	value, err := ParseOS(string(bb))
	if err != nil {
		return err
	}

	*o = value

	return nil
}
