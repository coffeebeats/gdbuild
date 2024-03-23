package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"

	"github.com/coffeebeats/gdbuild/internal/config"
	"github.com/coffeebeats/gdbuild/internal/osutil"
)

/* -------------------------------------------------------------------------- */
/*                          Function: DefaultFilename                         */
/* -------------------------------------------------------------------------- */

// DefaultFilename returns the default name of the 'gdbuild' manifest file.
func DefaultFilename() string {
	return "gdbuild.toml"
}

/* -------------------------------------------------------------------------- */
/*                               Function: Parse                              */
/* -------------------------------------------------------------------------- */

// Parse parses a 'Manifest' struct from a 'toml' document.
func Parse(bb []byte) (*Manifest, error) {
	d := toml.NewDecoder(bytes.NewReader(bb))
	d.DisallowUnknownFields()

	var m Manifest
	if err := d.Decode(&m); err != nil {
		var detailedErr *toml.StrictMissingError
		if errors.As(err, &detailedErr) {
			return nil, fmt.Errorf("%w\n%s", err, detailedErr.String())
		}

		return nil, err
	}

	return &m, nil
}

/* -------------------------------------------------------------------------- */
/*                             Function: ParseFile                            */
/* -------------------------------------------------------------------------- */

// Parse parses a 'Manifest' struct from a 'toml' file.
func ParseFile(path string) (*Manifest, error) {
	if err := osutil.Path(path).CheckIsFile(); err != nil {
		return nil, fmt.Errorf(
			"%w: expected a path to a 'gdbuild.toml' manifest file",
			config.ErrInvalidInput,
		)
	}

	bb, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Parse(bb)
}
