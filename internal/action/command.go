package action

import (
	"context"
	"strings"

	"github.com/coffeebeats/gdbuild/internal/exec"
)

/* -------------------------------------------------------------------------- */
/*                                Type: Command                               */
/* -------------------------------------------------------------------------- */

// Command is a string which is interpreted as a shell script.
type Command string

// Compile-time check that 'Action' is implemented.
var _ Action = (*Command)(nil)

/* ----------------------------- Method: Process ---------------------------- */

func (c Command) Process() *Process {
	return &Process{Args: strings.Split(string(c), " ")} //nolint:exhaustruct
}

/* ------------------------------ Impl: Runner ------------------------------ */

// Run executes the underlying shell command as a 'Process'.
func (c Command) Run(ctx context.Context) error {
	return c.Process().Run(ctx)
}

/* -------------------------- Interface: Combinable ------------------------- */

// After creates a new action which executes the provided action and then the
// wrapped command.
func (c Command) After(a Action) Action { //nolint:ireturn
	if c == "" {
		return a
	}

	return Sequence{Action: c, Pre: a} //nolint:exhaustruct
}

// AndThen creates a new action which executes the wrapped command and then the
// provided action.
func (c Command) AndThen(a Action) Action { //nolint:ireturn
	if c == "" {
		return a
	}

	return Sequence{Action: c, Post: a} //nolint:exhaustruct
}

/* ------------------------------ Impl: Printer ----------------------------- */

// Sprint displays the action without actually executing it.
func (c Command) Sprint() string {
	return string(c)
}

/* -------------------------------------------------------------------------- */
/*                               Type: Commands                               */
/* -------------------------------------------------------------------------- */

// Commands is a group of string commands which are executed in order.
type Commands struct {
	Commands []Command
	Shell    exec.Shell
	Verbose  bool
}

// Compile-time check that 'Action' is implemented.
var _ Action = (*Commands)(nil)

/* ------------------------------ Impl: Runner ------------------------------ */

// Run executes the underlying shell commands in order.
func (c Commands) Run(ctx context.Context) error {
	for _, cmd := range c.Commands {
		p := cmd.Process()

		p.Shell = c.Shell
		p.Verbose = c.Verbose

		if err := p.Run(ctx); err != nil {
			return err
		}
	}

	return nil
}

/* -------------------------- Interface: Combinable ------------------------- */

// After creates a new action which executes the provided action and then the
// wrapped list of commands.
func (c Commands) After(a Action) Action { //nolint:ireturn
	return Sequence{Action: c, Pre: a} //nolint:exhaustruct
}

// AndThen creates a new action which executes the wrapped list of commands and
// then the provided action.
func (c Commands) AndThen(a Action) Action { //nolint:ireturn
	return Sequence{Action: c, Post: a} //nolint:exhaustruct
}

/* ------------------------------ Impl: Printer ----------------------------- */

// Sprint displays the action without actually executing it.
func (c Commands) Sprint() string {
	cmds := make([]string, 0, len(c.Commands))

	for _, cmd := range c.Commands {
		cmds = append(cmds, string(cmd))
	}

	return strings.Join(cmds, "\n")
}
