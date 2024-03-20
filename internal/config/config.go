package config

import (
	"errors"

	"dario.cat/mergo"

	"github.com/coffeebeats/gdbuild/pkg/build"
)

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrMissingInput  = errors.New("missing input")
	ErrUnimplemented = errors.New("unimplemented")
)

/* -------------------------------------------------------------------------- */
/*                           Interface: Configurable                          */
/* -------------------------------------------------------------------------- */

type Configurable interface {
	Configurer
	Validator
	Merger
}

/* -------------------------- Interface: Configurer ------------------------- */

// Configurer is a type which can configure itself based on the current
// invocation. Note that 'Configure' must *not* set default values as this
// method might be called prior to complete resolution of user inputs.
type Configurer interface {
	Configure(inv build.Invocation) error
}

/* -------------------------- Interface: Validator -------------------------- */

// Validator is a type which can validate itself.
type Validator interface {
	Validate(inv build.Invocation) error
}

/* ---------------------------- Interface: Merger --------------------------- */

// Merger is a type which can merge from or into other objects (typically, but
// not always, of the same type).
type Merger interface {
	MergeInto(other any) error
}

/* -------------------------------------------------------------------------- */
/*                           Function: GetOrDefault                           */
/* -------------------------------------------------------------------------- */

// GetOrDefault is a convenience method to safely access a value from a
// potentially nil map.
func GetOrDefault[K comparable, V any](m map[K]V, key K) V { //nolint:ireturn
	if m == nil {
		return *new(V)
	}

	return m[key]
}

/* -------------------------------------------------------------------------- */
/*                               Function: Merge                              */
/* -------------------------------------------------------------------------- */

// Merge is a helper function for invoking 'mergo.Merge' with consistent
// settings.
func Merge[T any](dst *T, src T) error {
	return mergo.Merge(
		dst,
		src,
		mergo.WithAppendSlice,
		mergo.WithTypeCheck,
		mergo.WithOverride,
	)
}
