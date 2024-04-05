package common

import (
	"fmt"
	"os"

	"github.com/coffeebeats/gdbuild/pkg/godot/template"
)

/* -------------------------------------------------------------------------- */
/*                       Function: resolveEncryptionKey                       */
/* -------------------------------------------------------------------------- */

// resolveEncryptionKey determines the resource encryption key to use, first
// checking for an environment variable and then using a configuration value.
func resolveEncryptionKey(value string) (string, error) {
	if ek := template.EncryptionKeyFromEnv(); ek != "" {
		return ek, nil
	}

	if value == "" {
		return "", nil
	}

	ek := os.ExpandEnv(value)
	if ek != "" {
		return ek, nil
	}

	return "", fmt.Errorf(
		"%w: encryption key set in manifest, but value was empty: %s",
		ErrMissingInput,
		value,
	)
}
