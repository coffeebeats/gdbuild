package run

import (
	"errors"
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/exec"
)

var ErrInvalidInput = errors.New("invalid input")

/* -------------------------------------------------------------------------- */
/*                                Struct: Hook                                */
/* -------------------------------------------------------------------------- */

// Hook contains commands to execute before and after a build step.
//
// TODO: Allow per-hook shell settings.
type Hook struct {
	// Pre contains a command to run *before* an export step.
	Pre []action.Command `toml:"run_before"`
	// Post contains a command to run *after* an export step.
	Post []action.Command `toml:"run_after"`
	// Shell defines which shell process to run these commands in.
	Shell exec.Shell `toml:"shell"`
}

/* --------------------------- Method: PreActions --------------------------- */

// PreActions is a utility function to convert pre-build commands into a slice
// of 'Action' types.
func (h Hook) PreActions(rc *Context) action.Action { //nolint:ireturn
	actions := make([]action.Action, 0, len(h.Pre))

	for _, a := range h.Pre {
		p := a.Process()
		p.Directory = rc.PathWorkspace.String()
		p.Shell = h.Shell
		p.Verbose = rc.Verbose

		actions = append(actions, action.Action(p))
	}

	return action.InOrder(actions...)
}

/* --------------------------- Method: PostActions -------------------------- */

// PostActions is a utility function to convert post-build commands into a slice
// of 'Action' types.
func (h Hook) PostActions(rc *Context) action.Action { //nolint:ireturn
	actions := make([]action.Action, 0, len(h.Post))

	for _, a := range h.Post {
		p := a.Process()
		p.Directory = rc.PathWorkspace.String()
		p.Shell = h.Shell
		p.Verbose = rc.Verbose

		actions = append(actions, action.Action(p))
	}

	return action.InOrder(actions...)
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (h Hook) Validate(_ *Context) error {
	if h.Shell != exec.ShellUnknown {
		if _, err := exec.ParseShell(h.Shell.String()); err != nil {
			return fmt.Errorf("%w: unsupported shell: %s", ErrInvalidInput, h.Shell)
		}
	}

	return nil
}
