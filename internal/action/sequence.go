package action

import (
	"context"
	"strings"
)

/* -------------------------------------------------------------------------- */
/*                              Struct: Sequence                              */
/* -------------------------------------------------------------------------- */

// Sequence is an action which is solely comprised of a sequence of other action
// types.
type Sequence struct {
	Pre    Runner
	Action Runner
	Post   Runner
}

// Compile-time check that 'Action' is implemented.
var _ Action = (*Sequence)(nil)

/* ----------------------------- Method: Unwrap ----------------------------- */

// Unwrap recursively unwraps a sequence to get the central 'Runner'.
func (s Sequence) Unwrap() Runner { //nolint:ireturn
	inner, ok := s.Action.(Sequence)
	if ok {
		return inner.Unwrap()
	}

	return s.Action
}

/* ------------------------------ Impl: Action ------------------------------ */

// Run executes all actions in the sequence.
func (s Sequence) Run(ctx context.Context) error {
	if s.Pre != nil {
		if err := s.Pre.Run(ctx); err != nil {
			return err
		}
	}

	if s.Action != nil {
		if err := s.Action.Run(ctx); err != nil {
			return err
		}
	}

	if s.Post != nil {
		if err := s.Post.Run(ctx); err != nil {
			return err
		}
	}

	return nil
}

/* -------------------------- Interface: Combinable ------------------------- */

// After creates a new action which executes the provided action and then all of
// the actions in this sequence.
func (s Sequence) After(r Action) Action { //nolint:ireturn
	return Sequence{Action: s, Pre: r} //nolint:exhaustruct
}

// AndThen creates a new action which executes all actions in this sequence and
// then the provided action.
func (s Sequence) AndThen(r Action) Action { //nolint:ireturn
	return Sequence{Action: s, Post: r} //nolint:exhaustruct
}

/* ------------------------------ Impl: Printer ----------------------------- */

// Sprint displays the action without actually executing it.
func (s Sequence) Sprint() string {
	cmds := make([]string, 0)

	if p, ok := s.Pre.(Printer); ok {
		cmds = append(cmds, p.Sprint())
	}

	if p, ok := s.Action.(Printer); ok {
		cmds = append(cmds, p.Sprint())
	}

	if p, ok := s.Post.(Printer); ok {
		cmds = append(cmds, p.Sprint())
	}

	return strings.Join(cmds, "\n")
}
