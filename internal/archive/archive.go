package archive

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdbuild/internal/osutil"
)

const FileExtension = ".tar.gz"

var (
	ErrExtractFailed = errors.New("extract failed")
	ErrMissingInput  = errors.New("missing input")
)

/* -------------------------------------------------------------------------- */
/*                              Function: Create                              */
/* -------------------------------------------------------------------------- */

// Create writes the provided files to a compressed archive at 'out'. The
// implementation follows from https://www.arthurkoziel.com/writing-tar-gz-files-in-go/.
func Create(root string, files []string, out string) error {
	if len(files) == 0 {
		return fmt.Errorf("%w: 'files'", ErrMissingInput)
	}

	if out == "" {
		return fmt.Errorf("%w: 'out'", ErrMissingInput)
	}

	if !strings.HasSuffix(out, FileExtension) {
		out += FileExtension
	}

	log.Debugf("creating archive at path: %s", out)

	f, err := os.Create(out)
	if err != nil {
		return err
	}

	defer f.Close()

	return addFilesToArchive(root, files, f)
}

/* ----------------------- Function: addFilesToArchive ---------------------- */

func addFilesToArchive(root string, files []string, w io.Writer) error {
	gw := gzip.NewWriter(w)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Iterate over files and add each of them to the tar archive.
	for _, path := range files {
		path := filepath.Clean(path)

		info, err := os.Stat(filepath.Join(root, path))
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if err := addFileToArchive(tw, root, path, info); err != nil {
				return err
			}
		}

		if err := fs.WalkDir(os.DirFS(root), path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			info, err := d.Info()
			if err != nil {
				return err
			}

			if err := addFileToArchive(tw, root, path, info); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

/* ----------------------- Function: addFileToArchive ----------------------- */

func addFileToArchive(tw *tar.Writer, root, path string, info fs.FileInfo) error {
	log.Debugf("adding path to archive: %s", path)

	header, err := tar.FileInfoHeader(info, path)
	if err != nil {
		return err
	}

	header.Name = path
	if info.IsDir() {
		header.Name += "/"
	}

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	f, err := os.Open(filepath.Join(root, path))
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.Copy(tw, f)

	return err
}

/* -------------------------------------------------------------------------- */
/*                              Function: Extract                             */
/* -------------------------------------------------------------------------- */

// Extract uncompresses the files within the archive at 'archive' and copies
// them to the directory 'out'.
func Extract(ctx context.Context, archive, out string) error { //nolint:cyclop,funlen
	if archive == "" {
		return fmt.Errorf("%w: 'archive'", ErrMissingInput)
	}

	if out == "" {
		return fmt.Errorf("%w: 'out'", ErrMissingInput)
	}

	baseDirMode, err := osutil.ModeOf(out)
	if err != nil {
		return err
	}

	prefix := strings.TrimSuffix(filepath.Base(archive), FileExtension)

	a, err := os.Open(archive)
	if err != nil {
		return err
	}

	defer a.Close()

	gr, err := gzip.NewReader(a)
	if err != nil {
		return err
	}

	tr := tar.NewReader(gr)

	// Extract all files within the archive.
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}

			break
		}

		name := hdr.Name

		// See https://cs.opensource.google/go/go/+/refs/tags/go1.21.3:src/archive/tar/reader.go;l=60-67.
		if !filepath.IsLocal(name) || strings.Contains(name, `\`) || strings.Contains(name, "..") {
			return fmt.Errorf("%w: %s", tar.ErrInsecurePath, name)
		}

		// Remove the name of the tar-file from the filepath; this is to
		// facilitate extracting contents directly into the 'out' path.
		name = strings.TrimPrefix(name, prefix+string(os.PathSeparator))
		if strings.HasPrefix(name, prefix) {
			return fmt.Errorf(
				"%w: couldn't trim prefix: %s from %s",
				ErrExtractFailed,
				prefix, name,
			)
		}

		out := filepath.Join(out, name) //nolint:gosec

		if err := extractTarFile(ctx, tr, hdr, out, baseDirMode); err != nil {
			return err
		}
	}

	return nil
}

/* ------------------------ Function: extractTarFile ------------------------ */

// extractFile handles the extraction logic for each file in the Tar archive.
func extractTarFile(
	ctx context.Context,
	archive *tar.Reader,
	hdr *tar.Header,
	out string,
	baseDirMode fs.FileMode,
) error {
	// Ensure the parent directory exists with best-effort permissions. If
	// the zip archive already contains the directory as an entry then this
	// will have no effect.
	if err := os.MkdirAll(filepath.Dir(out), baseDirMode); err != nil {
		return err
	}

	mode := hdr.FileInfo().Mode()

	switch hdr.Typeflag {
	case tar.TypeDir:
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := os.MkdirAll(out, mode); err != nil {
			return err
		}

	case tar.TypeReg:
		if err := osutil.CopyReaderWithMode(ctx, archive, mode, out); err != nil {
			return err
		}
	}

	return nil
}
