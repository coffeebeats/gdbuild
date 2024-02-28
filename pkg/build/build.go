package build

import "errors"

var (
	// ErrInvalidInput is returned when a function is provided invalid input.
	ErrInvalidInput = errors.New("invalid input")
	// ErrMissingInput is returned when a function is missing required input.
	ErrMissingInput = errors.New("missing input")
)

/* -------------------------------------------------------------------------- */
/*                            Interface: Configure                            */
/* -------------------------------------------------------------------------- */

// Configurer is a type which can configure itself based on the current
// invocation.
type Configurer interface {
	Configure(i *Invocation) error
}

/* -------------------------------------------------------------------------- */
/*                            Interface: Validater                            */
/* -------------------------------------------------------------------------- */

// Validater is a type which can validate itself.
type Validater interface {
	Validate() error
}
