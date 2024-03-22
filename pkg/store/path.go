package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdbuild/internal/archive"
	"github.com/coffeebeats/gdbuild/pkg/template"
)

const envStore = "GDBUILD_HOME"

var (
	ErrInvalidPath   = errors.New("invalid file path")
	ErrMissingEnvVar = errors.New("missing environment variable")
)

/* -------------------------------------------------------------------------- */
/*                          Function: TemplateArchive                         */
/* -------------------------------------------------------------------------- */

// TemplateArchive returns the full path (starting with the store path) to the
// export template artifact archive within the store.
//
// NOTE: This does *not* mean the template archive exists.
func TemplateArchive(storePath string, t template.Template) (string, error) {
	if storePath == "" {
		return "", ErrMissingStore
	}

	cs, err := t.Checksum()
	if err != nil {
		return "", err
	}

	return filepath.Join(storePath, storeDirTemplate, cs+archive.FileExtension), nil
}

/* -------------------------------------------------------------------------- */
/*                               Function: Path                               */
/* -------------------------------------------------------------------------- */

// Returns the user-configured path to the 'gdbuild' store.
func Path() (string, error) {
	path := os.Getenv(envStore)
	if path == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingEnvVar, envStore)
	}

	if !filepath.IsAbs(path) {
		return "", fmt.Errorf("%w; expected absolute path: %s", ErrInvalidPath, path)
	}

	return filepath.Clean(path), nil
}
