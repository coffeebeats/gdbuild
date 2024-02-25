package build

/* -------------------------------------------------------------------------- */
/*                                Struct: SCons                               */
/* -------------------------------------------------------------------------- */

// SCons defines options and settings for use with the _Godot_ build system.
type SCons struct {
	// CFlags are additional 'CFLAGS' to append to the SCons build command. Note
	// that 'CFLAGS=...' will be appended *before* 'ExtraArgs'.
	CFlags []string `toml:"cflags"`
	// CCFlags are additional 'CFLAGS' to append to the SCons build command.
	// Note that 'CCFLAGS=...' will be appended *before* 'ExtraArgs'.
	CCFlags []string `toml:"ccflags"`
	// CXXFlags are additional 'CXXFLAGS' to append to the SCons build command.
	// Note that 'CXXFLAGS=...' will be appended *before* 'ExtraArgs'.
	CXXFlags []string `toml:"cxxflags"`
	// ExtraArgs are additional arguments to append to the SCons build command.
	ExtraArgs []string `toml:"extra_args"`
	// PathCache is the path to the SCons cache, relative to the manifest.
	PathCache string `toml:"cache_path"`
}
