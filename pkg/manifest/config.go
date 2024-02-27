package manifest

import "github.com/coffeebeats/gdbuild/pkg/build"

/* -------------------------------------------------------------------------- */
/*                               Struct: Config                               */
/* -------------------------------------------------------------------------- */

// Configs specifies GDBuild manifest-related settings.
type Config struct {
	// Extends is a path to another GDBuild manifest to extend. Note that value
	// override rules work the same as within a manifest; any primitive values
	// will override those defined in the base configuration, while arrays will
	// be appended to the base configuration's arrays.
	Extends build.Path `toml:"extends"`
}
