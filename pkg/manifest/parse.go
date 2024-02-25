package manifest

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

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

	if err := m.Validate(); err != nil {
		return nil, err
	}

	return &m, nil
}

// Parse parses a 'Manifest' struct from a 'toml' file.
func ParseFile(path string) (*Manifest, error) {
	if !strings.HasSuffix(path, Filename()) {
		return nil, fmt.Errorf(
			"%w: expected a path to a 'gdbuild.toml' manifest file",
			ErrInvalidInput,
		)
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, fmt.Errorf(
			"%w: expected a path to a 'gdbuild.toml' manifest file",
			ErrInvalidInput,
		)
	}

	bb, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Parse(bb)
}
