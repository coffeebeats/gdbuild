package export_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
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
				Platform: platform.OSLinux,
				Runnable: true,
			},

			want: `[preset.0]
platform                   = "Linux/X11"
encrypt_pck                = false
encrypt_directory          = false
encryption_include_filters = ""
exclude_filter             = ""
export_files               = ""
export_filter              = ""
custom_features            = ""
include_filter             = ""
name                       = "Export name"
runnable                   = true
dedicated_server           = false

[preset.0.options]
`,
		},
		{
			name: "lists are encoded correctly",

			index: 1,
			preset: export.Preset{
				Platform: platform.OSLinux,
				Features: []string{"feature1", "feature2"},
				ExportedFiles: []string{
					"res://A/B/C.gd",
					"res://B/C/D.tscn",
				},
			},

			want: `[preset.1]
platform                   = "Linux/X11"
encrypt_pck                = false
encrypt_directory          = false
encryption_include_filters = ""
exclude_filter             = ""
export_files               = PackedStringArray("res://A/B/C.gd","res://B/C/D.tscn")
export_filter              = ""
custom_features            = "feature1,feature2"
include_filter             = ""
name                       = ""
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
