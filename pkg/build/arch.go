package build

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	ErrMissingArch      = errors.New("missing architecture")
	ErrUnrecognizedArch = errors.New("unrecognized architecture")
)

/* -------------------------------------------------------------------------- */
/*                                 Enum: Arch                                 */
/* -------------------------------------------------------------------------- */

// CPU architectures which Godot supports targeting.
type Arch uint

const (
	ArchUnknown Arch = iota
	ArchAmd64
	ArchArm32
	ArchArm64
	ArchI386
	ArchUniversal
)

/* ----------------------------- Method: IsOneOf ---------------------------- */

// IsOneOf returns whether the 'Arch' is in the allowed set.
func (a Arch) IsOneOf(aa ...Arch) bool {
	return slices.Contains(aa, a)
}

/* ----------------------------- Impl: Stringer ----------------------------- */

// String implements 'fmt.Stringer' for 'Arch' according to the architecture
// names passed to SCons during compilation.
func (a Arch) String() string {
	switch a {
	case ArchAmd64:
		return "x86_64"
	case ArchArm32:
		return "arm32"
	case ArchArm64:
		return "arm64"
	case ArchI386:
		return "x86_32"
	case ArchUniversal:
		return "universal"
	default:
		return ""
	}
}

/* --------------------------- Function: ParseArch -------------------------- */

// Parses an input string as a CPU architecture specification.
//
// NOTE: This is a best effort implementation. Please open an issue on GitHub
// if some values are missing:
// github.com/coffeebeats/gdbuild/issues/new?labels=bug&template=%F0%9F%90%9B-bug-report.md.
func ParseArch(input string) (Arch, error) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "":
		return 0, ErrMissingArch

	case "amd64", "x86_64", "x86-64":
		return ArchAmd64, nil

	case "arm32":
		return ArchArm32, nil

	case "arm64", "arm64be":
		return ArchArm64, nil

	case "386", "i386", "x86", "x86_32":
		return ArchI386, nil

	case "fat", "universal":
		return ArchUniversal, nil

	default:
		return 0, fmt.Errorf("%w: '%s'", ErrUnrecognizedArch, input)
	}
}

/* ------------------------- Function: MustParseArch ------------------------ */

// Parses an input string as a CPU architecture specification but panics if it
// would fail.
func MustParseArch(input string) Arch {
	arch, err := ParseArch(input)
	if err != nil {
		panic(err)
	}

	return arch
}

/* ---------------------- Impl: encoding.UnmarshalText ---------------------- */

func (a *Arch) UnmarshalText(bb []byte) error {
	value, err := ParseArch(string(bb))
	if err != nil {
		return err
	}

	*a = value

	return nil
}
