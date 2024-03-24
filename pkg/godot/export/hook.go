package export

import (
	"fmt"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/exec"
	"github.com/coffeebeats/gdbuild/pkg/godot/build"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: Hook                                */
/* -------------------------------------------------------------------------- */

// Hook contains commands to execute before and after a build step.
//
// TODO: Allow per-hook shell settings.
type Hook struct {
	// Pre contains a command to run *before* an export step.
	Pre []action.Command `toml:"preexport"`
	// Post contains a command to run *after* an export step.
	Post []action.Command `toml:"postexport"`
	// Shell defines which shell process to run these commands in.
	Shell exec.Shell `toml:"shell"`
}

/* --------------------------- Method: PreActions --------------------------- */

// PreActions is a utility function to convert pre-build commands into a slice
// of 'Action' types.
func (h Hook) PreActions(bc *build.Context) action.Action { //nolint:ireturn
	actions := make([]action.Action, 0, len(h.Pre))

	for _, a := range h.Pre {
		p := a.Process()
		p.Directory = bc.PathBuild.String()
		p.Shell = h.Shell
		p.Verbose = bc.Verbose

		actions = append(actions, action.Action(p))
	}

	return action.InOrder(actions...)
}

/* --------------------------- Method: PostActions -------------------------- */

// PostActions is a utility function to convert post-build commands into a slice
// of 'Action' types.
func (h Hook) PostActions(bc *build.Context) action.Action { //nolint:ireturn
	actions := make([]action.Action, 0, len(h.Post))

	for _, a := range h.Post {
		p := a.Process()
		p.Directory = bc.PathBuild.String()
		p.Shell = h.Shell
		p.Verbose = bc.Verbose

		actions = append(actions, action.Action(p))
	}

	return action.InOrder(actions...)
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (h Hook) Validate(_ *build.Context) error {
	if h.Shell != exec.ShellUnknown {
		if _, err := exec.ParseShell(h.Shell.String()); err != nil {
			return fmt.Errorf("%w: unsupported shell: %s", ErrInvalidInput, h.Shell)
		}
	}

	return nil
}
