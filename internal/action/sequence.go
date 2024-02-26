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
	Runners []Runner
}

// Compile-time check that 'Runner' is implemented.
var _ Runner = (*Sequence)(nil)

/* ------------------------------ Impl: Action ------------------------------ */

// Run executes all actions in the sequence.
func (s Sequence) Run(ctx context.Context) error {
	for _, r := range s.Runners {
		if err := r.Run(ctx); err != nil {
			return err
		}
	}

	return nil
}

/* -------------------------- Interface: Combinable ------------------------- */

// After creates a new action which executes the provided action and then all of
// the actions in this sequence.
func (s Sequence) After(r Runner) Runner { //nolint:ireturn
	return Sequence{Runners: append([]Runner{r}, s.Runners...)}
}

// AndThen creates a new action which executes all actions in this sequence and
// then the provided action.
func (s Sequence) AndThen(r Runner) Runner { //nolint:ireturn
	return Sequence{Runners: append(s.Runners, r)}
}

/* ------------------------------ Impl: Printer ----------------------------- */

// Print displays the action without actually executing it.
func (s Sequence) Print() string {
	runners := make([]string, 0, len(s.Runners))

	for _, runner := range s.Runners {
		if p, ok := runner.(Printer); ok {
			runners = append(runners, p.Print())
		}
	}

	return strings.Join(runners, "\n")
}
