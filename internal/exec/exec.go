package exec

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var ErrMissingInput = errors.New("missing input")

/* -------------------------------------------------------------------------- */
/*                               Struct: Process                              */
/* -------------------------------------------------------------------------- */

// Process is a child process.
type Process struct {
	Directory   string
	Environment []string

	Shell Shell

	Verbose bool

	Args []string
}

/* ------------------------------ Method: args ------------------------------ */

func (p Process) args() ([]string, error) {
	if len(p.Args) == 0 {
		return nil, fmt.Errorf("%w: args", ErrMissingInput)
	}

	shell := p.Shell
	if shell == ShellUnknown {
		shell = DefaultShell()
	}

	var flag string

	switch shell {
	case ShellBash, ShellSh, ShellZsh:
		flag = "-c"
	case ShellPowerShell, ShellPwsh:
		flag = "-Command"
	case ShellCmd:
		flag = "/K"
	default:
		return nil, fmt.Errorf("%w: unsupported shell: %s", ErrUnrecognizedShell, p.Shell)
	}

	return []string{shell.String(), flag, strings.Join(p.Args, " ")}, nil
}

/* ------------------------------ Impl: Runner ------------------------------ */

// Run executes the underlying function.
func (p Process) Run(ctx context.Context) error {
	args, err := p.args()
	if err != nil {
		return err
	}

	if len(args) < 3 { //nolint:gomnd
		return fmt.Errorf("%w: missing arguments: %s", ErrMissingInput, args)
	}

	program, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, program, args[1:]...)

	cmd.Dir = p.Directory
	cmd.Env = p.Environment

	if p.Verbose {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	return cmd.Run()
}

/* --------------------------- Impl: fmt.Stringer --------------------------- */

func (p Process) String() string {
	args, err := p.args()
	if err != nil {
		return ""
	}

	return strings.Join(args, " ")
}
