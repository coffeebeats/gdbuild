package build

/* -------------------------------------------------------------------------- */
/*                               Struct: Target                               */
/* -------------------------------------------------------------------------- */

// Target specifies a single exportable artifact within the Godot project.
type Target struct {
	// Name is the display name of the target. Not used by Godot.
	Name string `json:"name"`

	// Runnable is whether the export artifact should be executable. This should
	// be true for client and server targets and false for artifacts like DLC.
	Runnable bool `json:"runnable" toml:"runnable"`
	// Server configures the target as a server-only executable, enabling some
	// optimizations like disabling graphics.
	Server bool `json:"server" toml:"server"`
	// DefaultFeatures contains the slice of Godot project feature tags to build
	// with.
	DefaultFeatures []string `json:"default_features" toml:"default_features"` //nolint:tagliatelle
	// PackFiles defines the game files exported as part of this artifact.
	PackFiles []*PackFile `json:"pack_files" toml:"pack_files"` //nolint:tagliatelle

	// Hook defines commands to be run before or after the target artifact is
	// generated.
	Hook Hook `json:"hook" toml:"hook"`
	// Options are 'export_presets.cfg' overrides, specifically the preset
	// 'options' table, for the exported artifact.
	Options map[string]any `json:"options" toml:"options"`
}

/* -------------------------------------------------------------------------- */
/*                              Struct: PackFile                              */
/* -------------------------------------------------------------------------- */

// PackFile defines instructions for assembling one or more '.pck' files
// containing exported game files.
type PackFile struct {
	// Embed defines whether the associated '.pck' file should be embedded in
	// the binary. If true, then the target this 'PackFile' is associated with
	// must be runnable.
	Embed bool `json:"embed" toml:"embed"`
	// Encrypt determines whether or not to encrypt the game files contained in
	// the resulting '.pck' files.
	Encrypt bool `json:"encrypt" toml:"encrypt"`
	// Glob is a slice of glob expressions to match game files against. These
	// will be evaluated from the directory containing the GDBuild manifest.
	Glob []string `json:"glob" toml:"glob"`
	// PackFilePartition is a ruleset for how to split the files matched by
	// 'glob' into one or more '.pck' files.
	Partition PackFilePartition `json:"partition" toml:"partition"`
	// Zip defines whether to compress the matching game files. The pack files
	// will use the '.zip' extension instead of '.pck'.
	Zip bool `json:"zip" toml:"zip"`
}

/* ------------------------ Struct: PackFilePartition ----------------------- */

// PackFilePartition describes how to automatically partition a collection of
// files into multiple '.pck' files.
//
// NOTE: This struct contains multiple different expressions of limits, multiple
// of which may be true at a time. If any of the contained rules would trigger a
// new '.pck' to be formed within a partition, then that rule will be respected.
type PackFilePartition struct {
	// Depth is the maximum folder depth from the project directory containing
	// the GDBuild manifest to split files between. Any folders past this depth
	// limit will all be included within the same '.pck' file.
	Depth uint `json:"depth" toml:"depth"`
	// Limit describes limits on the files within individual '.pck' files in the
	// partition.
	Limit PackFilePartitionLimit `json:"limit" toml:"limit"`
}

/* --------------------- Struct: PackFilePartitionLimit --------------------- */

// PackFilePartitionLimit describes limits used to determine when a new '.pck'
// file within a partition should be started.
type PackFilePartitionLimit struct {
	// Size is a human-readable file size limit that all '.pck' files within the
	// partition must adhere to.
	Size string `json:"size" toml:"size"`
	// Files is the maximum count of files within a single '.pck' file within a
	// partition.
	Files uint `json:"files" toml:"files"`
}
