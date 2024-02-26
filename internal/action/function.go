package action

import (
	"context"
	"fmt"
)

/* -------------------------------------------------------------------------- */
/*                               Type: Function                               */
/* -------------------------------------------------------------------------- */

// Function is a wrapper for a Go-based action.
type Function func(context.Context) error

// Compile-time check that 'Runner' is implemented.
var _ Runner = (*Function)(nil)

/* ------------------------------ Impl: Runner ------------------------------ */

// Run executes the underlying function.
func (f Function) Run(ctx context.Context) error {
	return f(ctx)
}

/* -------------------------- Interface: Combinable ------------------------- */

// After creates a new action which executes the provided action and then the
// wrapped function.
func (f Function) After(a Runner) Runner { //nolint:ireturn
	return Sequence{Action: f, Pre: a} //nolint:exhaustruct
}

// AndThen creates a new action which executes the wrapped function and then the
// provided action.
func (f Function) AndThen(a Runner) Runner { //nolint:ireturn
	return Sequence{Action: f, Post: a} //nolint:exhaustruct
}

/* ------------------------------ Impl: Printer ----------------------------- */

// Print displays the action without actually executing it.
func (f Function) Print() string {
	return fmt.Sprintf("%##v", f)
}
