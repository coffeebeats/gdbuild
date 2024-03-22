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

// Compile-time check that 'Action' is implemented.
var _ Action = (*Function)(nil)

/* ------------------------------ Impl: Runner ------------------------------ */

// Run executes the underlying function.
func (f Function) Run(ctx context.Context) error {
	return f(ctx)
}

/* -------------------------- Interface: Combinable ------------------------- */

// After creates a new action which executes the provided action and then the
// wrapped function.
func (f Function) After(a Action) Action { //nolint:ireturn
	if a == nil {
		return f
	}

	if f == nil {
		return a
	}

	return Sequence{Action: f, Pre: a} //nolint:exhaustruct
}

// AndThen creates a new action which executes the wrapped function and then the
// provided action.
func (f Function) AndThen(a Action) Action { //nolint:ireturn
	if a == nil {
		return f
	}

	if f == nil {
		return a
	}

	return Sequence{Action: f, Post: a} //nolint:exhaustruct
}

/* ------------------------------ Impl: Printer ----------------------------- */

// Sprint displays the action without actually executing it.
func (f Function) Sprint() string {
	return fmt.Sprintf("%#v", f)
}

/* --------------------------- Impl: fmt.Stringer --------------------------- */

func (f Function) String() string {
	return f.Sprint()
}
