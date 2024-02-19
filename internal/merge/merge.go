package merge

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

func Map[K comparable, V any](others ...map[K]V) map[K]V {
	out := make(map[K]V)

	for _, other := range others {
		for k, v := range other {
			out[k] = v
		}
	}

	return out
}

/* -------------------------------------------------------------------------- */
/*                              Function: String                              */
/* -------------------------------------------------------------------------- */

func String(others ...string) string {
	var out string

	for _, o := range others {
		if o == "" {
			continue
		}

		out = o
	}

	return out
}
