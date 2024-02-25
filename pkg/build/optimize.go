package build

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMissingOptimize      = errors.New("missing optimize")
	ErrUnrecognizedOptimize = errors.New("unrecognized optimize")
)

/* -------------------------------------------------------------------------- */
/*                               Enum: Optimize                               */
/* -------------------------------------------------------------------------- */

// Optimize is the level of optimization for a _Godot_ export template.
type Optimize uint

const (
	OptimizeUnknown Optimize = iota
	OptimizeCustom
	OptimizeDebug
	OptimizeNone
	OptimizeSize
	OptimizeSpeed
	OptimizeSpeedTrace
)

/* ----------------------------- Impl: Stringer ----------------------------- */

// String implements 'fmt.Stringer' for 'Optimize' according to the optimization
// levels passed to SCons during compilation.
func (o Optimize) String() string {
	switch o {
	case OptimizeCustom:
		return "custom"
	case OptimizeDebug:
		return "debug" //nolint:goconst
	case OptimizeNone:
		return "none"
	case OptimizeSize:
		return "size"
	case OptimizeSpeed:
		return "speed"
	case OptimizeSpeedTrace:
		return "speed_trace"
	default:
		return ""
	}
}

/* ------------------------- Function: ParseOptimize ------------------------ */

// Parses an input string as an operating system specification.
//
// NOTE: This is a best effort implementation. Please open an issue on GitHub
// if some values are missing:
// github.com/coffeebeats/gdbuild/issues/new?labels=bug&template=%F0%9F%90%9B-bug-report.md.
func ParseOptimize(input string) (Optimize, error) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "":
		return 0, ErrMissingOptimize

	case "custom":
		return OptimizeCustom, nil

	case "debug":
		return OptimizeDebug, nil

	case "none":
		return OptimizeNone, nil

	case "size":
		return OptimizeSize, nil

	case "speed":
		return OptimizeSpeed, nil

	case "speed_trace":
		return OptimizeSpeedTrace, nil

	default:
		return 0, fmt.Errorf("%w: '%s'", ErrUnrecognizedOptimize, input)
	}
}

/* -------------------------- Function: MustParseOptimize ------------------------- */

// Parses an input string as an operating system specification but panics if it
// would fail.
func MustParseOptimize(input string) Optimize {
	opt, err := ParseOptimize(input)
	if err != nil {
		panic(err)
	}

	return opt
}

/* ---------------------- Impl: encoding.UnmarshalText ---------------------- */

func (o *Optimize) UnmarshalText(bb []byte) error {
	value, err := ParseOptimize(string(bb))
	if err != nil {
		return err
	}

	*o = value

	return nil
}
