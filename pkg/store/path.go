package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdbuild/internal/archive"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
	"github.com/coffeebeats/gdbuild/pkg/run"
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
func TemplateArchive(storePath string, t *template.Template) (string, error) {
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
/*                           Function: TargetArchive                          */
/* -------------------------------------------------------------------------- */

// TargetArchive returns the full path (starting with the store path) to the
// exported project archive within the store.
//
// NOTE: This does *not* mean the export archive exists.
func TargetArchive(storePath string, rc *run.Context, x *export.Export) (string, error) {
	if storePath == "" {
		return "", ErrMissingStore
	}

	cs, err := x.Checksum(rc)
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
