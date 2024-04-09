package export

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	doublestar "github.com/bmatcuk/doublestar/v4"
	"golang.org/x/exp/maps"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
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
	// Name is the name of the pack file. If omitted a name will be chosen based
	// on the pack files index within the target configuration.
	Name osutil.Path `toml:"name"`
	// PackFilePartition is a ruleset for how to split the files matched by
	// 'glob' into one or more '.pck' files.
	Partition PackFilePartition `toml:"partition"`
	// Visuals sets whether the resources in this pack file will have visuals
	// stripped. Only usable when exporting as a server and defaults to 'true'.
	Visuals *bool `toml:"visuals"`
	// Zip defines whether to compress the matching game files. The pack files
	// will use the '.zip' extension instead of '.pck'.
	Zip *bool `toml:"zip"`
}

/* ----------------------------- Method: Preset ----------------------------- */

// Preset constructs a 'Preset' for the pack file based on the current context.
func (c *PackFile) Preset(rc *run.Context, xp *Export, index int) (Preset, error) {
	var preset Preset

	preset.Options = map[string]string{}

	for key, value := range xp.Options {
		if value, ok := value.(string); ok {
			preset.Options[key] = value
		}
	}

	preset.Arch = xp.Arch
	preset.Embed = config.Dereference(c.Embed)
	preset.Features = slices.Clone(rc.Features)
	preset.Name = c.Filename(rc.Platform, rc.Target, index)
	preset.PathTemplate = xp.PathTemplate
	preset.Platform = rc.Platform
	preset.Runnable = xp.Runnable
	preset.Server = xp.Server

	if config.Dereference(c.Encrypt) {
		preset.Encrypt = true
		preset.EncryptIndex = true
		preset.Encrypted = slices.Clone(c.Glob)
		preset.EncryptionKey = xp.EncryptionKey
	}

	ff, err := c.Files(rc.PathWorkspace)
	if err != nil {
		return Preset{}, err
	}

	shouldStripVisuals := c.StripVisuals()

	if !xp.Server && !shouldStripVisuals {
		preset.ExportMode = ModeResources

		for _, f := range ff {
			if err := preset.AddFile(rc, f); err != nil {
				return Preset{}, err
			}
		}
	} else {
		preset.ExportMode = ModeCustomized

		for _, f := range ff {
			visuals := FileVisualMode(FileVisualModeKeep)
			if shouldStripVisuals {
				visuals = FileVisualModeStrip
			}

			if err := preset.AddServerFile(rc, f, visuals); err != nil {
				return Preset{}, err
			}
		}
	}

	return preset, nil
}

/* ---------------------------- Method: Filename ---------------------------- */

func (c *PackFile) Filename(pl platform.OS, targetName string, index int) string {
	if config.Dereference(c.Embed) {
		return targetName + c.Extension(pl)
	}

	ext := c.Extension(pl)

	if c.Name != "" {
		return strings.TrimSuffix(c.Name.String(), ext) + ext
	}

	return targetName + "." + strconv.Itoa(index) + ext
}

/* ---------------------------- Method: Extension --------------------------- */

func (c *PackFile) Extension(pl platform.OS) string {
	if config.Dereference(c.Embed) {
		switch pl {
		case platform.OSMacOS:
			return ".app/"
		case platform.OSWindows:
			return ".exe"
		default:
			return ""
		}
	}

	if config.Dereference(c.Zip) {
		return ".zip"
	}

	return ".pck"
}

/* -------------------------- Method: ResourceMode -------------------------- */

// StripVisuals determines whether the included resources should have visuals
// stripped as part of a server-side optimization.
func (c *PackFile) StripVisuals() bool {
	return c.Visuals != nil && !*c.Visuals
}

/* ------------------------------ Method: Files ----------------------------- */

func (c *PackFile) Files(path osutil.Path) ([]osutil.Path, error) { //nolint:cyclop,funlen
	pathRoot, err := filepath.Abs(path.String())
	if err != nil {
		return nil, err
	}

	files := make(map[osutil.Path]struct{})

	for _, pattern := range c.Glob {
		// NOTE: See https://github.com/bmatcuk/doublestar?tab=readme-ov-file#glob.
		pattern = filepath.Clean(pattern)
		pattern = filepath.ToSlash(pattern)

		matches, err := doublestar.Glob(
			os.DirFS(pathRoot),
			pattern,
			doublestar.WithNoFollow(),
			doublestar.WithFailOnIOErrors(),
			doublestar.WithFailOnPatternNotExist(),
		)
		if err != nil {
			return nil, err
		}

		var mm []string

		for _, pathMatch := range matches {
			pathMatch = filepath.Join(pathRoot, pathMatch)

			info, err := os.Stat(pathMatch)
			if err != nil {
				return nil, err
			}

			if !info.IsDir() {
				mm = append(mm, pathMatch)

				continue
			}

			if err := fs.WalkDir(
				os.DirFS(pathMatch),
				".",
				func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}

					if !d.IsDir() {
						mm = append(mm, filepath.Join(pathMatch, path))
					}

					return nil
				},
			); err != nil {
				return nil, err
			}
		}

		baseGlob := filepath.Base(pattern)

		for _, m := range mm {
			baseMatch := filepath.Base(m)

			// Ignore hidden files unless explicitly searched for.
			if !strings.HasPrefix(baseGlob, ".") &&
				strings.HasPrefix(baseMatch, ".") {
				continue
			}

			files[osutil.Path(m)] = struct{}{}
		}
	}

	return maps.Keys(files), nil
}

/* ------------------------- Impl: config.Configurer ------------------------ */

func (c *PackFile) Configure(rc *run.Context) error {
	if err := c.Name.RelTo(rc.PathWorkspace); err != nil {
		return err
	}

	c.Name = osutil.Path(strings.TrimPrefix(
		c.Name.String(),
		rc.PathWorkspace.String()+string(os.PathSeparator)),
	)

	return nil
}

/* ------------------------- Impl: config.Validator ------------------------- */

func (c *PackFile) Validate(rc *run.Context) error {
	if config.Dereference(c.Embed) && c.Name != "" {
		return fmt.Errorf(
			"%w: cannot set the name of an embedded pack file: %s",
			ErrInvalidInput,
			c.Name,
		)
	}

	if c.Name != "" {
		got := c.Extension(rc.Platform)
		if want := filepath.Ext(c.Name.String()); got != want {
			return fmt.Errorf(
				"%w: incorrect extension for pack file: %s: was %s but wanted %s",
				ErrInvalidInput,
				c.Name,
				got,
				want,
			)
		}
	}

	if len(c.Glob) == 0 {
		return fmt.Errorf(
			"%w: missing required 'glob' property for pack file",
			ErrInvalidInput,
		)
	}

	ff, err := c.Files(rc.PathWorkspace)
	if err != nil {
		return err
	}

	if len(ff) == 0 {
		return fmt.Errorf(
			"%w: no files match pack file: %s",
			ErrInvalidInput,
			strings.Join(c.Glob, ","),
		)
	}

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
