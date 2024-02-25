package merge

import (
	"errors"
	"fmt"
)

var (
	// ErrConflictingValue is returned when mergeable types have conflicting
	// values for the same property.
	ErrConflictingValue = errors.New("conflicting value")
)

// Merger is a type which can merge with other instances of itself.
type Merger[T any] interface {
	Merge(other T) error
}

/* -------------------------------------------------------------------------- */
/*                               Function: Bool                               */
/* -------------------------------------------------------------------------- */

func Bool(others ...bool) bool {
	var out bool

	for _, o := range others {
		out = out || o
	}

	return out
}

/* -------------------------------------------------------------------------- */
/*                              Function: Number                              */
/* -------------------------------------------------------------------------- */

func Number[T number](others ...T) T { //nolint:ireturn
	var out T

	for _, o := range others {
		if o == *new(T) {
			continue
		}

		out = o
	}

	return out
}

/* ---------------------------- Interface: number --------------------------- */

type number interface {
	~float32 | ~float64 | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

/* -------------------------------------------------------------------------- */
/*                                Function: Map                               */
/* -------------------------------------------------------------------------- */

func Map[K comparable, V comparable](base *map[K]V, other map[K]V) error {
	if other == nil {
		return nil
	}

	if base == nil {
		*base = other
	}

	b := *base
	for k, v := range other {
		if _, ok := b[k]; ok {
			return fmt.Errorf("%w: map key: %v", ErrConflictingValue, k)
		}

		(*base)[k] = v
	}

	return nil
}

func Primitive[T comparable](base *T, other T) error {
	var t T

	if other == t {
		return nil
	}

	if base != nil && *base != t {
		return fmt.Errorf("%w: %v,%v", ErrConflictingValue, *base, other)
	}

	*base = other

	return nil
}

func Pointer[T comparable](base *T, other *T) error {
	if other == nil {
		return nil
	}

	if base == nil {
		*base = *other
	}

	if base != nil && base != other {
		return fmt.Errorf("%w: %v,%v", ErrConflictingValue, *base, other)
	}

	*base = *other

	return nil
}
