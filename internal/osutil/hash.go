package osutil

import (
	"errors"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var ErrInvalidInput = errors.New("invalid input")

/* -------------------------------------------------------------------------- */
/*                             Function: HashFiles                            */
/* -------------------------------------------------------------------------- */

// HashFiles updates the provided 'hash.Hash' with the file structure contents
// rooted at 'root'.
func HashFiles(h hash.Hash, root string) error {
	info, err := os.Stat(root)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return HashFile(h, root)
	}

	if err := fs.WalkDir(os.DirFS(root), ".", func(path string, d fs.DirEntry, err error) error {
		// Propagate an error walking the directory.
		if err != nil {
			return err
		}

		// Only hash files.
		if d.IsDir() {
			return nil
		}

		// Hash the file.
		return HashFileWithName(h, root, path)
	}); err != nil {
		return err
	}

	return nil
}

/* --------------------------- Function: HashFile --------------------------- */

func HashFile(h hash.Hash, path string) error {
	if !filepath.IsAbs(path) {
		return fmt.Errorf("%w: expected an absolute path: %s", ErrInvalidInput, path)
	}

	// Hash the file contents.
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	return nil
}

/* ----------------------- Function: HashFileWithName ----------------------- */

func HashFileWithName(h hash.Hash, root, path string) error {
	// Hash the filename.
	if _, err := io.Copy(h, strings.NewReader(path)); err != nil {
		return err
	}

	return HashFile(h, filepath.Join(root, path))
}
