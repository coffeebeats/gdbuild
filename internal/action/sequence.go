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
	Action Action
	Post   Runner
}

// Compile-time check that 'Runner' is implemented.
var _ Runner = (*Sequence)(nil)

/* ------------------------------ Impl: Action ------------------------------ */

// Run executes all actions in the sequence.
func (s Sequence) Run(ctx context.Context) error {
	if err := s.Pre.Run(ctx); err != nil {
		return err
	}

	if s.Action != nil {
		if err := s.Action.Run(ctx); err != nil {
			return err
		}
	}

	if err := s.Post.Run(ctx); err != nil {
		return err
	}

	return nil
}

/* -------------------------- Interface: Combinable ------------------------- */

// After creates a new action which executes the provided action and then all of
// the actions in this sequence.
func (s Sequence) After(r Runner) Runner { //nolint:ireturn
	return Sequence{Action: s, Pre: r} //nolint:exhaustruct
}

// AndThen creates a new action which executes all actions in this sequence and
// then the provided action.
func (s Sequence) AndThen(r Runner) Runner { //nolint:ireturn
	return Sequence{Action: s, Post: r} //nolint:exhaustruct
}

/* ------------------------------ Impl: Printer ----------------------------- */

// Print displays the action without actually executing it.
func (s Sequence) Print() string {
	runners := make([]string, 0)

	if p, ok := s.Pre.(Printer); ok {
		runners = append(runners, p.Print())
	}

	if p, ok := s.Action.(Printer); ok {
		runners = append(runners, p.Print())
	}

	if p, ok := s.Post.(Printer); ok {
		runners = append(runners, p.Print())
	}

	return strings.Join(runners, "\n")
}
