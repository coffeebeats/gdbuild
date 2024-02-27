package action

import "context"

/* -------------------------------------------------------------------------- */
/*                              Interface: Runner                             */
/* -------------------------------------------------------------------------- */

// Action defines a type which can be executed, both in live mode and in dry-run
// mode, as well as combine with other Action implementers.
type Action interface {
	Combinable
	Runner
	Printer
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
