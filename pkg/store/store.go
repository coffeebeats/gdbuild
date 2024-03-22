package store

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
)

const (
	storeDirTemplate = "templates"
	storeFileLayout  = "layout.v0" // simplify migrating in the future
)

var (
	ErrMissingInput = errors.New("missing input")
	ErrMissingStore = errors.New("missing store")
)

/* -------------------------------------------------------------------------- */
/*                               Function: Clear                              */
/* -------------------------------------------------------------------------- */

// Removes all cached artifacts in the store.
func Clear(storePath string) error {
	if storePath == "" {
		return ErrMissingStore
	}

	// Clear the entire export template cache directory.
	if err := os.RemoveAll(filepath.Join(storePath, storeDirTemplate)); err != nil {
		return err
	}

	// Remake the deleted directories.
	return Touch(storePath)
}

/* -------------------------------------------------------------------------- */
/*                                Function: Has                               */
/* -------------------------------------------------------------------------- */

// Return whether the store has the specified version cached.
func Has(storePath string, t template.Template) (bool, error) {
	if storePath == "" {
		return false, ErrMissingStore
	}

	path, err := TemplateArchive(storePath, t)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

/* -------------------------------------------------------------------------- */
/*                              Function: Remove                              */
/* -------------------------------------------------------------------------- */

// Removes the specified version from the store.
func Remove(storePath string, t template.Template) error {
	if storePath == "" {
		return ErrMissingStore
	}

	path, err := TemplateArchive(storePath, t)
	if err != nil {
		return err
	}

	// Remove the specific template archive from the store.
	if err := os.Remove(path); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	log.Debugf("removed template archive from store: %s", filepath.Base(path))

	return removeUnusedCacheDirectories(storePath, path)
}

/* ----------------- Function: removeUnusedCacheDirectories ----------------- */

// A utility method which cleans up unused directories from the specified path
// up to the store's cache directories.
func removeUnusedCacheDirectories(storePath, path string) error {
	if path == "" {
		return fmt.Errorf("%w: 'path'", ErrMissingInput)
	}

	for {
		path = filepath.Dir(path)

		// Add a safeguard to not escape the store cache directories.
		if !strings.HasPrefix(path, filepath.Join(storePath, storeDirTemplate)) {
			return nil
		}

		files, err := os.ReadDir(path)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		if len(files) > 0 {
			return nil
		}

		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                               Function: Touch                              */
/* -------------------------------------------------------------------------- */

// Touch ensures a store is initialized at the specified path; no effect if it
// exists already.
func Touch(storePath string) error {
	if storePath == "" {
		return ErrMissingStore
	}

	// Create the 'Store' directory, if needed.
	if err := os.MkdirAll(storePath, osutil.ModeUserRWXGroupRX); err != nil {
		return err
	}

	// Create the required subdirectories, if needed.
	for _, d := range []string{storeDirTemplate} {
		path := filepath.Join(storePath, d)
		if err := os.MkdirAll(path, osutil.ModeUserRWXGroupRX); err != nil {
			return err
		}
	}

	// Create the required files, if needed.
	for _, f := range []string{storeFileLayout} {
		path := filepath.Join(storePath, f)
		if err := os.WriteFile(path, nil, osutil.ModeUserRW); err != nil {
			return err
		}
	}

	return nil
}
