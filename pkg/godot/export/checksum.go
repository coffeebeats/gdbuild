package export

import (
	"fmt"
	"hash/crc64"
	"io"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/mitchellh/hashstructure/v2"

	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

/* -------------------------------------------------------------------------- */
/*                             Function: Checksum                             */
/* -------------------------------------------------------------------------- */

// Checksum produces a checksum hash of the export specification. When the
// checksums of two 'Export' definitions matches, the resulting exported
// artifacts will be equivalent.
func Checksum(rc *run.Context, x *Export) (string, error) { //nolint:funlen
	if x == nil {
		return "", fmt.Errorf("%w: export configuration", ErrMissingInput)
	}

	xp := *x

	// If 'PathTemplateArchive' is set then don't include the cached template
	// archive specification.
	if xp.PathTemplateArchive != "" {
		xp.Template = nil
	}

	hash, err := hashstructure.Hash(
		xp,
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

	log.Debugf(
		"target export hash before file dependencies: %s",
		strconv.FormatUint(cs.Sum64(), 16),
	)

	files := make([]osutil.Path, 0)
	pathRoot := osutil.Path(filepath.Dir(rc.PathManifest.String()))

	for _, pck := range xp.PackFiles {
		ff, err := pck.Files(pathRoot)
		if err != nil {
			return "", err
		}

		files = append(files, ff...)
	}

	// Make the path list unique and sorted.
	slices.Sort(files)
	files = slices.Compact(files)

	for _, path := range files {
		if err := osutil.HashFiles(cs, path.String()); err != nil {
			return "", err
		}

		log.Debugf(
			"target export hash after file dependency: %s: %s",
			path.String(),
			strconv.FormatUint(cs.Sum64(), 16),
		)
	}

	// Include the optional 'PathTemplateArchive' in the checksum.
	if xp.PathTemplateArchive != "" {
		if err := osutil.HashFile(cs, xp.PathTemplateArchive.String()); err != nil {
			return "", err
		}

		log.Debugf(
			"target export hash after file dependency: %s: %s",
			xp.PathTemplateArchive.String(),
			strconv.FormatUint(cs.Sum64(), 16),
		)
	}

	return strconv.FormatUint(cs.Sum64(), 16), nil
}
