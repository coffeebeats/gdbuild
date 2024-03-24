package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/config/platform"
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
	Godot Godot `toml:"godot"`
	// Target includes settings for exporting Godot game executables and packs.
	Target map[string]platform.Targets `toml:"target"`
	// Template includes settings for building custom export templates.
	Template platform.Templates `toml:"template"`
}

/* ----------------------------- Function: Init ----------------------------- */

// Init initializes a GDBuild manifest at the specified path. Note that 'path'
// can be a directory or a '.toml' file.
func Init(path string) error { //nolint:cyclop
	if path == "" {
		return fmt.Errorf("%w: 'path'", ErrMissingInput)
	}

	info, err := os.Stat(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		// Assume the path is a directory if it doesn't end with '.toml'.
		if !strings.HasSuffix(path, ".toml") {
			if strings.Contains(filepath.Base(path), ".") {
				return fmt.Errorf(
					"%w: path must be a directory or a '.toml' file: %s",
					ErrInvalidInput,
					path,
				)
			}

			path = filepath.Join(path, DefaultFilename())
		}
	}

	if info != nil {
		if !info.IsDir() {
			return fmt.Errorf("%w: %s", fs.ErrExist, path)
		}

		path = filepath.Join(path, DefaultFilename())
	}

	// Check again if the file exists.
	info, err = os.Stat(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	if info != nil {
		return fmt.Errorf("%w: %s", fs.ErrExist, path)
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := io.Copy(f, bytes.NewReader([]byte(defaultContents()))); err != nil {
		return err
	}

	return nil
}

/* ------------------------ Function: defaultContents ----------------------- */

// defaultContents contains the default GDBuild manifest contents.
func defaultContents() string {
	return `[config]
  # Inherit from the specified manifest file, merging the configuration in
  # this file on top of the settings in the specified file. 
  extends = ""

[godot]
  # The version of Godot to use for compiling and exporting.
  version = "4.2.1-stable"

[target.client]
  runnable = true
  server   = false

  default_features = ["steam"]

  pack_files = [
    {glob = ["*"], embed = true, encrypt = true},
  ]

[template]
  # A path to a 'custom.py' file which defines export template build options.
  custom_py_path = "$PWD/custom.py"

[template.scons]
  cache_path = "$PWD/.scons"
  command    = ["python3", "-m", "SCons"]

[template.profile.release]
  # EncryptionKey is the encryption key to embed in the export template.
  encryption_key = "$SCRIPT_AES256_ENCRYPTION_KEY"
`
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
