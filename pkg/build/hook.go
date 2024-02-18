package build

/* -------------------------------------------------------------------------- */
/*                                Struct: Hook                                */
/* -------------------------------------------------------------------------- */

// Hook contains commands to execute before and after a build step.
type Hook struct {
	// Pre contains a command to run *before* a build step.
	Pre string `json:"prebuild" toml:"prebuild"`
	// Post contains a command to run *after* a build step.
	Post string `json:"postbuild" toml:"postbuild"`
}
