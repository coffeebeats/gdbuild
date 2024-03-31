package export_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coffeebeats/gdbuild/pkg/godot/export"
)

func TestPresetMarshal(t *testing.T) {
	tests := []struct {
		name string

		index  int
		preset export.Preset

		want string
		err  error
	}{
		{
			name: "empty preset returns empty options",

			index:  0,
			preset: export.Preset{},

			want: "",
		},
		{
			name: "simple preset returns correct string",

			preset: export.Preset{
				Name:     "Export name",
				Runnable: true,
			},

			want: `[preset.0]
encrypt_pck                = false
encryption_include_filters = ""
encrypt_directory          = false
exclude_filter             = ""
export_filter              = ""
custom_features            = ""
include_filter             = ""
name                       = "Export name"
export_path                = ""
platform                   = ""
runnable                   = true
dedicated_server           = false

[preset.0.options]
`,
		},
		{
			name: "lists are encoded correctly",

			index: 1,
			preset: export.Preset{
				Features: []string{"feature1", "feature2"},
				ExportedFiles: []string{
					"res://A/B/C.gd",
					"res://B/C/D.tscn",
				},
			},

			want: `[preset.1]
encrypt_pck                = false
encryption_include_filters = ""
encrypt_directory          = false
exclude_filter             = ""
export_files               = PackedStringArray("res://A/B/C.gd","res://B/C/D.tscn")
export_filter              = ""
custom_features            = "feature1,feature2"
include_filter             = ""
name                       = ""
export_path                = ""
platform                   = ""
runnable                   = false
dedicated_server           = false

[preset.1.options]
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given: A string builder to write to.
			var sb strings.Builder

			// When: The preset is serialized.
			err := test.preset.Marshal(&sb, test.index)

			// Then: The returned error matches expectations.
			assert.ErrorIs(t, err, test.err)

			// Then: The resulting string matches expectations.
			assert.Equal(t, test.want, sb.String())
		})
	}
}
