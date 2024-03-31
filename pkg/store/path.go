package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdbuild/internal/archive"
)

const envStore = "GDBUILD_HOME"

var (
	ErrInvalidInput  = errors.New("invalid input")
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
func TemplateArchive(storePath string, checksum string) (string, error) {
	if storePath == "" {
		return "", ErrMissingStore
	}

	if checksum == "" {
		return "", fmt.Errorf("%w: checksum: %s", ErrInvalidInput, checksum)
	}

	return filepath.Join(storePath, storeDirTemplate, checksum+archive.FileExtension), nil
}

/* -------------------------------------------------------------------------- */
/*                           Function: TargetArchive                          */
/* -------------------------------------------------------------------------- */

// TargetArchive returns the full path (starting with the store path) to the
// exported project archive within the store.
//
// NOTE: This does *not* mean the export archive exists.
func TargetArchive(storePath string, checksum string) (string, error) {
	if storePath == "" {
		return "", ErrMissingStore
	}

	return filepath.Join(storePath, storeDirTemplate, checksum+archive.FileExtension), nil
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
