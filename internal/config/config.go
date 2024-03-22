package config

import (
	"errors"
	"reflect"

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
		mergo.WithTransformers(configTransformer{}),
	)
}

/* --------------------------- Struct: transformer -------------------------- */

type configTransformer struct{}

/* ---------------------------- Impl: Transformer --------------------------- */

func (ct configTransformer) Transformer(ty reflect.Type) func(dst, src reflect.Value) error {
	// Handle pointers to 'bool' types.
	if ty == reflect.TypeOf((*bool)(nil)) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() && !src.IsNil() {
				dst.Set(src)
			}

			return nil
		}
	}

	// Handle pointers to 'uint' types.
	if ty == reflect.TypeOf((*uint)(nil)) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() && !src.IsNil() {
				dst.Set(src)
			}

			return nil
		}
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*                            Function: Dereference                           */
/* -------------------------------------------------------------------------- */

// Dereference safely dereferences a pointer into the underlying value or the
// empty value for the type.
func Dereference[T any](ptr *T) T { //nolint:ireturn
	if ptr == nil {
		return *new(T)
	}

	return *ptr
}
