package build

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/exec"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: Hook                                */
/* -------------------------------------------------------------------------- */

// Hook contains commands to execute before and after a build step.
//
// TODO: Allow per-hook shell settings.
type Hook struct {
	// Pre contains a command to run *before* a build step.
	Pre []action.Command `toml:"prebuild"`
	// Post contains a command to run *after* a build step.
	Post []action.Command `toml:"postbuild"`
	// Shell defines which shell process to run these commands in.
	Shell exec.Shell `toml:"shell"`
}

/* --------------------------- Method: PreActions --------------------------- */

// PreActions is a utility function to convert pre-build commands into a slice
// of 'Action' types.
func (h Hook) PreActions(inv Invocation) []action.Action {
	out := make([]action.Action, 0, len(h.Pre))

	for _, a := range h.Pre {
		p := a.Process()
		p.Directory = inv.PathBuild.String()
		p.Shell = h.Shell
		p.Verbose = inv.Verbose

		out = append(out, action.Action(p))
	}

	return out
}

/* --------------------------- Method: PostActions -------------------------- */

// PostActions is a utility function to convert post-build commands into a slice
// of 'Action' types.
func (h Hook) PostActions(inv Invocation) []action.Action {
	out := make([]action.Action, 0, len(h.Post))

	for _, a := range h.Post {
		p := a.Process()
		p.Directory = inv.PathBuild.String()
		p.Shell = h.Shell
		p.Verbose = inv.Verbose

		out = append(out, action.Action(p))
	}

	return out
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (h Hook) Validate() error {
	if h.Shell != exec.ShellUnknown {
		if _, err := exec.ParseShell(h.Shell.String()); err != nil {
			return fmt.Errorf("%w: unsupported shell: %s", ErrInvalidInput, h.Shell)
		}
	}

	return nil
}
