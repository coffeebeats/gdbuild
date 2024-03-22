package build

import "errors"

var (
	// ErrInvalidInput is returned when a function is provided invalid input.
	ErrInvalidInput = errors.New("invalid input")
	// ErrMissingInput is returned when a function is missing required input.
	ErrMissingInput = errors.New("missing input")
)
