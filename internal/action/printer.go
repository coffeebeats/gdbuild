package action

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
)

/* -------------------------------------------------------------------------- */
/*                           Struct: WithDescription                          */
/* -------------------------------------------------------------------------- */

// WithDescription is a utility type for wrapping an action with a description
// that replaces a 'Print' call.
type WithDescription[T Action] struct {
	Action      T
	Description string
}

// Compile-time check that 'Action' is implemented.
var _ Action = (*Function)(nil)

/* ------------------------------ Impl: Runner ------------------------------ */

// Run executes the underlying action.
func (d WithDescription[T]) Run(ctx context.Context) error {
	log.Infof("calling function: %s", d.Description)

	return d.Action.Run(ctx)
}

/* -------------------------- Interface: Combinable ------------------------- */

// After creates a new action which executes the provided action and then the
// wrapped action.
func (d WithDescription[T]) After(a Action) Action { //nolint:ireturn
	return Sequence{Action: d, Pre: a} //nolint:exhaustruct
}

// AndThen creates a new action which executes the wrapped action and then the
// provided action.
func (d WithDescription[T]) AndThen(a Action) Action { //nolint:ireturn
	return Sequence{Action: d, Post: a} //nolint:exhaustruct
}

/* ------------------------------ Impl: Printer ----------------------------- */

// Sprint displays the action without actually executing it.
func (d WithDescription[T]) Sprint() string {
	return fmt.Sprintf("%T:\n  %s", d.Action, d.Description)
}
