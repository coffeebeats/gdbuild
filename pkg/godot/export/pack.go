package export

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                              Struct: PackFile                              */
/* -------------------------------------------------------------------------- */

// PackFile defines instructions for assembling one or more '.pck' files
// containing exported game files.
type PackFile struct {
	// Embed defines whether the associated '.pck' file should be embedded in
	// the binary. If true, then the target this 'PackFile' is associated with
	// must be runnable.
	Embed *bool `toml:"embed"`
	// Encrypt determines whether or not to encrypt the game files contained in
	// the resulting '.pck' files.
	Encrypt *bool `toml:"encrypt"`
	// Glob is a slice of glob expressions to match game files against. These
	// will be evaluated from the directory containing the GDBuild manifest.
	Glob []string `toml:"glob"`
	// PackFilePartition is a ruleset for how to split the files matched by
	// 'glob' into one or more '.pck' files.
	Partition PackFilePartition `toml:"partition"`
	// Zip defines whether to compress the matching game files. The pack files
	// will use the '.zip' extension instead of '.pck'.
	Zip *bool `toml:"zip"`
}

/* ------------------------------ Method: Files ----------------------------- */

func (c *PackFile) Files(path osutil.Path) ([]osutil.Path, error) {
	pathRoot, err := filepath.Abs(path.String())
	if err != nil {
		return nil, err
	}

	files := make(map[osutil.Path]struct{})

	for _, g := range c.Glob {
		if !strings.HasPrefix(g, "/") &&
			!strings.HasPrefix(g, string(os.PathSeparator)) {
			g = filepath.Join(pathRoot, g)
		}

		matches, err := filepath.Glob(g)
		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			baseGlob := filepath.Base(g)
			baseMatch := filepath.Base(match)

			// Ignore hidden files unless explicitly searched for.
			if !strings.HasPrefix(baseGlob, ".") &&
				strings.HasPrefix(baseMatch, ".") {
				continue
			}

			files[osutil.Path(match)] = struct{}{}
		}
	}

	return maps.Keys(files), nil
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (c *PackFile) Configure(_ *run.Context) error {
	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *PackFile) Validate(_ *run.Context) error {
	return nil
}

/* -------------------------------------------------------------------------- */
/*                          Struct: PackFilePartition                         */
/* -------------------------------------------------------------------------- */

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
	Depth uint `toml:"depth"`
	// Limit describes limits on the files within individual '.pck' files in the
	// partition.
	Limit PackFilePartitionLimit `toml:"limit"`
}

/* --------------------- Struct: PackFilePartitionLimit --------------------- */

// PackFilePartitionLimit describes limits used to determine when a new '.pck'
// file within a partition should be started.
type PackFilePartitionLimit struct {
	// Size is a human-readable file size limit that all '.pck' files within the
	// partition must adhere to.
	Size string `toml:"size"`
	// Files is the maximum count of files within a single '.pck' file within a
	// partition.
	Files uint `toml:"files"`
}
