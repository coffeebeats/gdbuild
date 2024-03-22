package osutil

import (
	"io/fs"
	"os"
)

const ModeUserRW = 0600         // rw-------
const ModeUserRWX = 0700        // rwx------
const ModeUserRWXGroupRX = 0750 // rwxr-x---

/* -------------------------------------------------------------------------- */
/*                              Function: ModeOf                              */
/* -------------------------------------------------------------------------- */

// ModeOf returns the file mode of the specified file.
func ModeOf(path string) (fs.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return info.Mode(), nil
}
