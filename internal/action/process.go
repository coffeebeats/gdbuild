package action

import (
	"context"
	"errors"
	"fmt"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdbuild/internal/exec"
)

var ErrMissingInput = errors.New("missing input")

/* -------------------------------------------------------------------------- */
/*                               Struct: Process                              */
/* -------------------------------------------------------------------------- */

// Process is a child process.
type Process exec.Process

// Compile-time check that 'Action' is implemented.
var _ Action = (*Process)(nil)

/* ------------------------------ Impl: Runner ------------------------------ */

// Run executes the underlying function.
func (p Process) Run(ctx context.Context) error {
	process := exec.Process(p)

	log.Infof("running command: %s", process.String())

	return process.Run(ctx)
}

/* -------------------------- Interface: Combinable ------------------------- */

// After creates a new action which executes the provided action and then the
// wrapped function.
func (p Process) After(a Action) Action { //nolint:ireturn
	return Sequence{Action: p, Pre: a} //nolint:exhaustruct
}

// AndThen creates a new action which executes the wrapped function and then the
// provided action.
func (p Process) AndThen(a Action) Action { //nolint:ireturn
	return Sequence{Action: p, Post: a} //nolint:exhaustruct
}

/* ------------------------------ Impl: Printer ----------------------------- */

// Sprint displays the action without actually executing it.
func (p Process) Sprint() string {
	dir := fmt.Sprintf("in directory '%s':\n  ", p.Directory)

	return dir + exec.Process(p).String()
}
