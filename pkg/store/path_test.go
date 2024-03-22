package store

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coffeebeats/gdbuild/internal/archive"
	"github.com/coffeebeats/gdbuild/pkg/template"
)

/* -------------------------- Test: TemplateArchive ------------------------- */

func TestTemplateArchive(t *testing.T) {
	storeName := "store"

	tests := []struct {
		name string

		store    string
		template template.Template

		want string
		err  error
	}{
		{
			name: "missing store returns an error",

			store: "",

			err: ErrMissingStore,
		},

		{
			name: "template archive successfully found",

			store: storeName,

			want: func() string {
				cs, err := (&template.Template{}).Checksum()
				assert.NoError(t, err)

				return filepath.Join(
					storeName,
					storeDirTemplate,
					cs+archive.FileExtension,
				)
			}(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// When: The path to the cached executable is determined.
			got, err := TemplateArchive(tc.store, tc.template)

			// Then: The expected error value is returned.
			assert.ErrorIs(t, err, tc.err)

			// Then: The expected filepath is returned.
			assert.Equal(t, tc.want, got)
		})
	}
}

/* ------------------------------- Test: Path ------------------------------- */

func TestPath(t *testing.T) {
	storeName := "store"

	tests := []struct {
		env  string
		want string
		err  error
	}{
		// Invalid inputs
		{env: "", err: ErrMissingEnvVar},
		{env: "a", err: ErrInvalidPath},
		{env: "a/b/c", err: ErrInvalidPath},

		// Valid inputs
		{env: "/" + storeName, want: "/" + storeName},
		{env: "/." + storeName, want: "/." + storeName},
		{env: "/a/b/" + storeName, want: "/a/b/" + storeName},
		{env: "/a/b/." + storeName, want: "/a/b/." + storeName},
	}

	for _, tc := range tests {
		t.Run(tc.env, func(t *testing.T) {
			// Given: The store environment variable is set.
			t.Setenv(envStore, tc.env)

			// When: The store path is determined.
			got, err := Path()

			// Then: The returned error matches expectations.
			assert.ErrorIs(t, err, tc.err)

			// Then: The resulting path matches expectations.
			assert.Equal(t, tc.want, got)
		})
	}
}
