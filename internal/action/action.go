package action

import (
	"context"
	"fmt"
)

/* -------------------------------------------------------------------------- */
/*                              Interface: Runner                             */
/* -------------------------------------------------------------------------- */

// Action defines a type which can be executed, both in live mode and in dry-run
// mode, as well as combine with other Action implementers.
type Action interface {
	Combinable
	Runner
	Printer

	fmt.Stringer // Required for checksums.
}

/* ---------------------------- Interface: Runner --------------------------- */

// Runner is a type which can execute an arbitrary action or command.
type Runner interface {
	Run(ctx context.Context) error
}

/* -------------------------- Interface: Combinable ------------------------- */

type Combinable interface {
	After(a Action) Action
	AndThen(a Action) Action
}

/* -------------------------------------------------------------------------- */
/*                             Interface: Printer                             */
/* -------------------------------------------------------------------------- */

type Printer interface {
	Sprint() string
}

/* -------------------------------------------------------------------------- */
/*                             Interface: Actioner                            */
/* -------------------------------------------------------------------------- */

type Actioner interface {
	Action() (Action, error)
}

/* -------------------------------------------------------------------------- */
/*                               Function: NoOp                               */
/* -------------------------------------------------------------------------- */

type NoOp struct{}

/* ------------------------------ Impl: Runner ------------------------------ */

// Run executes a no-op function.
func (n NoOp) Run(_ context.Context) error {
	return nil
}

/* -------------------------- Interface: Combinable ------------------------- */

// After creates a new action which executes the provided action and then the
// no-op function.
func (n NoOp) After(a Action) Action { //nolint:ireturn
	if a == nil {
		return n
	}

	return Sequence{Action: n, Pre: a} //nolint:exhaustruct
}

// AndThen creates a new action which executes the no-op function and then the
// provided action.
func (n NoOp) AndThen(a Action) Action { //nolint:ireturn
	if a == nil {
		return n
	}

	return Sequence{Action: n, Post: a} //nolint:exhaustruct
}

/* ------------------------------ Impl: Printer ----------------------------- */

// Sprint displays the action without actually executing it.
func (n NoOp) Sprint() string {
	return ""
}

/* --------------------------- Impl: fmt.Stringer --------------------------- */

func (n NoOp) String() string {
	return n.Sprint()
}
