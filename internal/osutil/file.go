package osutil

import (
	"context"
	"io"
	"io/fs"
	"os"

	"github.com/coffeebeats/gdbuild/internal/ioutil"
)

// Only write to 'out'; create a new file/overwrite an existing.
const copyFileWriteFlag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC

/* -------------------------------------------------------------------------- */
/*                             Function: CopyFile                             */
/* -------------------------------------------------------------------------- */

// CopyFile is a utility function for copying an 'io.Reader' to a new file
// created with the specified 'os.FileMode'.
func CopyFile(ctx context.Context, src, out string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}

	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	dst, err := os.OpenFile(out, copyFileWriteFlag, info.Mode())
	if err != nil {
		return err
	}

	defer dst.Close()

	if _, err := io.Copy(dst, ioutil.NewReaderWithContext(ctx, f.Read)); err != nil {
		return err
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*                        Function: CopyReaderWithMode                        */
/* -------------------------------------------------------------------------- */

// CopyReaderWithMode is a utility function for copying an 'io.Reader' to a new
// file created with the specified 'os.FileMode'.
func CopyReaderWithMode(ctx context.Context, src io.Reader, mode fs.FileMode, out string) error {
	dst, err := os.OpenFile(out, copyFileWriteFlag, mode)
	if err != nil {
		return err
	}

	defer dst.Close()

	if _, err := io.Copy(dst, ioutil.NewReaderWithContext(ctx, src.Read)); err != nil {
		return err
	}

	return nil
}
