package config

import (
	"errors"
	"reflect"

	"dario.cat/mergo"
)

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrMissingInput  = errors.New("missing input")
	ErrUnimplemented = errors.New("unimplemented")
)

/* -------------------------------------------------------------------------- */
/*                           Interface: Configurable                          */
/* -------------------------------------------------------------------------- */

type Configurable[T any] interface {
	Configurer[T]
	Validator[T]
	Merger
}

/* -------------------------- Interface: Configurer ------------------------- */

// Configurer is a type which can configure itself based on the current
// invocation. Note that 'Configure' must *not* set default values as this
// method might be called prior to complete resolution of user inputs.
type Configurer[T any] interface {
	Configure(c T) error
}

/* -------------------------- Interface: Validator -------------------------- */

// Validator is a type which can validate itself.
type Validator[T any] interface {
	Validate(c T) error
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
func Dereference[T any](ptr *T) T { //nolint:ireturn,nolintlint
	if ptr == nil {
		return *new(T)
	}

	return *ptr
}
