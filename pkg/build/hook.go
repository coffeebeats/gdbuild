package build

/* -------------------------------------------------------------------------- */
/*                                Struct: Hook                                */
/* -------------------------------------------------------------------------- */

// Hook contains commands to execute before and after a build step.
type Hook struct {
	// Pre contains a command to run *before* a build step.
	Pre []string `json:"prebuild" toml:"prebuild"`
	// Post contains a command to run *after* a build step.
	Post []string `json:"postbuild" toml:"postbuild"`
}

/* ------------------------------ Impl: Merger ------------------------------ */

func (h *Hook) CombineWith(hooks ...*Hook) *Hook {
	base := h
	if h == nil {
		base = &Hook{} //nolint:exhaustruct
	}

	for _, other := range hooks {
		if other == nil {
			continue
		}

		base.Pre = append(append(base.Pre, ";"), other.Pre...)
		base.Post = append(append(base.Post, ";"), other.Post...)
	}

	return base
}
