package osutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrUnsupportedFileType is returned when attempting to copy a file type that's
// not supported.
var ErrUnsupportedFileType = errors.New("unsupported file type")

// CopyDir recursively copies a directory from 'srcDir' to 'dstDir', preserving
// soft links. All regular files will be hard copied. Note that file attributes
// are not preserved, so this should only be used when the folder contents are
// required in the original structure. This implementation is based on [1].
//
// [1] https://github.com/moby/moby/blob/master/daemon/graphdriver/copy/copy.go
func CopyDir(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(srcPath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Rebase path
		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, relPath)

		switch mode := f.Mode(); {
		case mode.IsRegular():
			if err2 := os.Link(srcPath, dstPath); err2 != nil {
				return err2
			}

		case mode.IsDir():
			if err := os.Mkdir(dstPath, f.Mode()); err != nil && !os.IsExist(err) {
				return err
			}

		case mode&os.ModeSymlink != 0:
			link, err := os.Readlink(srcPath)
			if err != nil {
				return err
			}

			if err := os.Symlink(link, dstPath); err != nil {
				return err
			}

		default:
			return fmt.Errorf("%w: (%d / %s) for %s", ErrUnsupportedFileType, f.Mode(), f.Mode().String(), srcPath)
		}

		return nil
	})
}
