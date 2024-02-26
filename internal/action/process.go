package action

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
	Environment map[string]string

	Args []string
}

// Compile-time check that 'Runner' is implemented.
var _ Runner = (*Process)(nil)

/* ------------------------------ Impl: Runner ------------------------------ */

// Run executes the underlying function.
func (p Process) Run(ctx context.Context) error {
	if len(p.Args) < 1 {
		return fmt.Errorf("%w: args", ErrMissingInput)
	}

	cmd := exec.CommandContext(ctx, p.Args[0], p.Args[1:]...) //nolint:gosec

	cmd.Dir = p.Directory

	for k, v := range p.Environment {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	cmd.Stdout = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

/* -------------------------- Interface: Combinable ------------------------- */

// After creates a new action which executes the provided action and then the
// wrapped function.
func (p Process) After(a Runner) Runner { //nolint:ireturn
	return Sequence{Runners: append([]Runner{p}, a)}
}

// AndThen creates a new action which executes the wrapped function and then the
// provided action.
func (p Process) AndThen(a Runner) Runner { //nolint:ireturn
	return Sequence{Runners: append([]Runner{p}, a)}
}

/* ------------------------------ Impl: Printer ----------------------------- */

func (p Process) Print() string {
	return strings.Join(p.Args, " ")
}

/* -------------------------------------------------------------------------- */
/*                                Type: Command                               */
/* -------------------------------------------------------------------------- */

// Command is a string which is interpreted as a shell script.
type Command string

// Compile-time check that 'Runner' is implemented.
var _ Runner = (*Command)(nil)

/* ------------------------------ Impl: Runner ------------------------------ */

// Run executes the underlying shell command as a 'Process'.
func (c Command) Run(ctx context.Context) error {
	p := Process{Args: strings.Split(string(c), " ")} //nolint:exhaustruct

	return p.Run(ctx)
}

/* -------------------------- Interface: Combinable ------------------------- */

// After creates a new action which executes the provided action and then the
// wrapped command.
func (c Command) After(a Runner) Runner { //nolint:ireturn
	return Sequence{Runners: append([]Runner{c}, a)}
}

// AndThen creates a new action which executes the wrapped command and then the
// provided action.
func (c Command) AndThen(a Runner) Runner { //nolint:ireturn
	return Sequence{Runners: append([]Runner{c}, a)}
}

/* ------------------------------ Impl: Printer ----------------------------- */

// Print displays the action without actually executing it.
func (c Command) Print() string {
	return string(c)
}
