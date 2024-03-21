package exec

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

var (
	ErrMissingShell      = errors.New("missing shell")
	ErrUnrecognizedShell = errors.New("unrecognized shell")
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
	ShellPowerShell
	ShellPwsh
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
		return "cmd.exe"
	case ShellPowerShell:
		return "powershell.exe"
	case ShellPwsh:
		return "pwsh.exe"
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
	switch strings.TrimSuffix(strings.ToLower(strings.TrimSpace(input)), ".exe") {
	case "":
		return 0, ErrMissingShell

	case "bash":
		return ShellBash, nil
	case "cmd":
		return ShellCmd, nil
	case "powershell":
		return ShellPowerShell, nil
	case "pwsh":
		return ShellPwsh, nil
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

/* ------------------------- Function: DefaultShell ------------------------- */

// DefaultShell returns the default shell, based on the platform.
func DefaultShell() Shell {
	switch runtime.GOOS {
	case "windows":
		return ShellCmd
	default:
		return ShellSh
	}
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
