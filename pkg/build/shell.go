package build

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMissingShell      = errors.New("missing shell")
	ErrUnrecognizedShell = errors.New("unrecognized shell")
	ErrUnsupportedShell  = errors.New("unsupported shell")
)

/* -------------------------------------------------------------------------- */
/*                                 Enum: Shell                                */
/* -------------------------------------------------------------------------- */

// Shell is the name of a shell program to use to execute a command in.
type Shell uint

const (
	ShellUnknown Shell = iota
	ShellBash
	ShellCmd
	ShellFish
	ShellPowershell
	ShellSh
	ShellZsh
)

/* ----------------------------- Impl: Stringer ----------------------------- */

// String implements 'fmt.Stringer' for 'Shell'.
func (s Shell) String() string {
	switch s {
	case ShellBash:
		return "bash"
	case ShellCmd:
		return "cmd"
	case ShellFish:
		return "fish"
	case ShellPowershell:
		return "ps"
	case ShellSh:
		return "sh"
	case ShellZsh:
		return "zsh"
	default:
		return ""
	}
}

/* -------------------------- Function: ParseShell -------------------------- */

// Parses an input string as a shell name.
//
// NOTE: This is a best effort implementation. Please open an issue on GitHub
// if some values are missing:
// github.com/coffeebeats/gdbuild/issues/new?labels=bug&template=%F0%9F%90%9B-bug-report.md.
func ParseShell(input string) (Shell, error) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "":
		return 0, ErrMissingShell

	case "bash":
		return ShellBash, nil
	case "cmd":
		return ShellCmd, nil
	case "fish":
		return ShellFish, nil
	case "ps":
		return ShellPowershell, nil
	case "sh":
		return ShellSh, nil
	case "zsh":
		return ShellZsh, nil

	default:
		return 0, fmt.Errorf("%w: '%s'", ErrUnrecognizedShell, input)
	}
}

/* ------------------------ Function: MustParseShell ------------------------ */

// Parses an input string as a shell name but panics if it would fail.
func MustParseShell(input string) Shell {
	shell, err := ParseShell(input)
	if err != nil {
		panic(err)
	}

	return shell
}

/* ---------------------- Impl: encoding.UnmarshalText ---------------------- */

func (s *Shell) UnmarshalText(bb []byte) error {
	value, err := ParseShell(string(bb))
	if err != nil {
		return err
	}

	*s = value

	return nil
}
