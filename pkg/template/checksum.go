package template

import (
	"hash/crc64"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/mitchellh/hashstructure/v2"

	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/build"
)

/* -------------------------------------------------------------------------- */
/*                             Function: Checksum                             */
/* -------------------------------------------------------------------------- */

// Checksum produces a checksum hash of the export template specification. When
// the checksums of two 'Template' definitions matches, the resulting export
// templates will be equivalent.
//
// NOTE: This implementation relies on producers of 'Template' to correctly
// register all file system dependencies within 'Paths'.
func Checksum(t *build.Template) (string, error) {
	hash, err := hashstructure.Hash(
		t,
		hashstructure.FormatV2,
		&hashstructure.HashOptions{ //nolint:exhaustruct
			IgnoreZeroValue: true,
			SlicesAsSets:    true,
			ZeroNil:         true,
		},
	)
	if err != nil {
		return "", err
	}

	cs := crc64.New(crc64.MakeTable(crc64.ECMA))

	// Update the 'crc64' hash with the struct hash.
	if _, err := io.Copy(cs, strings.NewReader(strconv.FormatUint(hash, 16))); err != nil {
		return "", err
	}

	for _, p := range uniquePaths(t) {
		root := p.String()

		log.Debugf("hashing files rooted at path: %s", root)

		if err := osutil.HashFiles(cs, root); err != nil {
			return "", err
		}
	}

	return strconv.FormatUint(cs.Sum64(), 16), nil
}

/* -------------------------- Function: uniquePaths ------------------------- */

// uniquePaths returns the unique list of expanded path dependencies.
func uniquePaths(t *build.Template) []osutil.Path {
	paths := t.Paths

	for _, b := range t.Builds {
		paths = append(paths, b.CustomModules...)

		if b.CustomPy != "" {
			paths = append(paths, b.CustomPy)
		}

		switch g := b.Source; {
		case g.PathSource != "":
			paths = append(paths, g.PathSource)
		case g.VersionFile != "":
			paths = append(paths, g.VersionFile)
		}
	}

	slices.Sort(paths)

	return slices.Compact(paths)
}
