package build

import (
	"errors"
	"fmt"
	"path/filepath"
)

var (
	// ErrConflictingValue is returned when two settings conflict with eachother.
	ErrConflictingValue = errors.New("conflicting setting")
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

/* -------------------------------------------------------------------------- */
/*                                 Type: Path                                 */
/* -------------------------------------------------------------------------- */

// Path is a string type that's expected to be a path.
type Path string

/* ------------------------------ Method: RelTo ----------------------------- */

// RelTo converts the underlying path into a cleaned, absolute path. If the path
// is relative it's first resolved against the provided base path (or current
// working directory if 'base' is empty).
func (p *Path) RelTo(base Path) error {
	if base == "" {
		return fmt.Errorf("%w: base path", ErrMissingInput)
	}

	if p == nil {
		return nil
	}

	path := string(*p)
	if path == "" || filepath.IsAbs(path) {
		return nil
	}

	path, err := filepath.Abs(filepath.Join(string(base), path))
	if err != nil {
		return err
	}

	*p = Path(path)

	return nil
}

/* --------------------------- Impl: fmt.Stringer --------------------------- */

func (p Path) String() string {
	return string(p)
}
