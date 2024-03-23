package config

import (
	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/config/template"
	"github.com/coffeebeats/gdbuild/pkg/godot/build"
)

var (
	ErrInvalidInput = config.ErrInvalidInput
	ErrMissingInput = config.ErrMissingInput
)

/* -------------------------------------------------------------------------- */
/*                              Struct: Manifest                              */
/* -------------------------------------------------------------------------- */

// Manifest defines the supported structure of the GDBuild manifest file.
type Manifest struct {
	// Config contains GDBuild configuration-related settings.
	Config Config `toml:"config"`
	// Godot contains settings on which Godot version/source code to use.
	Godot build.Source `toml:"godot"`
	// Template includes settings for building custom export templates.
	Template template.Templates `toml:"template"`
}

/* -------------------------------------------------------------------------- */
/*                               Struct: Config                               */
/* -------------------------------------------------------------------------- */

// Configs specifies GDBuild manifest-related settings.
type Config struct {
	// Extends is a path to another GDBuild manifest to extend. Note that value
	// override rules work the same as within a manifest; any primitive values
	// will override those defined in the base configuration, while arrays will
	// be appended to the base configuration's arrays.
	Extends osutil.Path `toml:"extends"`
}
