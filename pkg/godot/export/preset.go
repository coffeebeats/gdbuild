package export

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
	"github.com/coffeebeats/gdbuild/pkg/run"

	ini "gopkg.in/ini.v1"
)

const resourcePathPrefix = "res://"

/* -------------------------------------------------------------------------- */
/*                               Struct: Preset                               */
/* -------------------------------------------------------------------------- */

// Preset defines the parameters used in a Godot export preset.
type Preset struct {
	CustomizedFiles map[string]string `ini:"customized_files,omitempty"`
	Encrypt         bool              `ini:"encrypt_pck"`
	Encrypted       []string          `ini:"encryption_include_filters"`
	EncryptIndex    bool              `ini:"encrypt_directory"`
	Exclude         string            `ini:"exclude_filter"`
	ExportedFiles   []string          `ini:"export_files,omitempty"`
	// ExportMode sets the type of export to use. Should be 'resources' for a
	// standard pack file and 'customized' for dedicated server pack files.
	ExportMode string   `ini:"export_filter"`
	Features   []string `ini:"custom_features"`
	Include    string   `ini:"include_filter"`
	Name       string   `ini:"name"`
	PathExport string   `ini:"export_path"`
	Platform   string   `ini:"platform"`
	Runnable   bool     `ini:"runnable"`
	Server     bool     `ini:"dedicated_server"`

	Options map[string]string `ini:"-"`
}

/* ----------------------------- Method: AddFile ---------------------------- */

func (p *Preset) AddFile(rc *run.Context, path osutil.Path) error {
	if err := path.CheckIsFile(); err != nil {
		return err
	}

	if err := path.RelTo(rc.PathWorkspace); err != nil {
		return err
	}

	pathString := strings.TrimPrefix(path.String(), rc.PathWorkspace.String())
	pathString = strings.TrimPrefix(filepath.Clean("./"+pathString), "./")
	pathString = resourcePathPrefix + strings.TrimPrefix(pathString, resourcePathPrefix)

	p.ExportedFiles = append(p.ExportedFiles, pathString)

	slices.Sort(p.ExportedFiles)
	p.ExportedFiles = slices.Compact(p.ExportedFiles)

	return nil
}

/* ------------------------- Method: SetArchitecture ------------------------ */

func (p *Preset) SetArchitecture(arch platform.Arch) {
	if p.Options == nil {
		p.Options = map[string]string{}
	}

	p.Options["binary_format/architecture"] = arch.String()
}

/* --------------------------- Method: SetPlatform -------------------------- */

func (p *Preset) SetPlatform(pl platform.OS) error {
	switch pl {
	case platform.OSLinux:
		p.Platform = "Linux/X11"
	case platform.OSMacOS:
		p.Platform = "macOS"
	case platform.OSWindows:
		p.Platform = "Windows Desktop"
	default:
		return fmt.Errorf("%w: unsupported platform: %s", ErrInvalidInput, pl)
	}

	return nil
}

/* --------------------------- Method: SetTemplate -------------------------- */

func (p *Preset) SetTemplate(path string) {
	if p.Options == nil {
		p.Options = map[string]string{}
	}

	p.Options["custom_template/debug"] = path
	p.Options["custom_template/release"] = path
}

/* ----------------------------- Method: Marshal ---------------------------- */

func (p *Preset) Marshal(w io.Writer, index int) error {
	lo := ini.LoadOptions{PreserveSurroundedQuote: true} //nolint:exhaustruct

	cfg := ini.Empty(lo)
	cfg.ValueMapper = valueMapper

	presetName := "preset." + strconv.Itoa(index)
	section := cfg.Section(presetName)

	// Initially hydrate the 'ini' file.
	if err := section.ReflectFrom(p); err != nil {
		return err
	}

	// Unmarshal the file back into the struct to trigger the value mapper.
	if err := section.MapTo(p); err != nil {
		return err
	}

	// Finally, re-populate the file with the updated values.
	if err := section.ReflectFrom(p); err != nil {
		return err
	}

	section = cfg.Section(presetName + ".options")

	for key, value := range p.Options {
		section.Key(key).SetValue(valueMapper(value))
	}

	if _, err := cfg.WriteTo(w); err != nil {
		return err
	}

	return nil
}

/* -------------------------- Function: valueMapper ------------------------- */

func valueMapper(s string) string {
	s = os.ExpandEnv(s)

	// NOTE: There doesn't seem to be a way to apply this value mapper to
	// specific fields, so this is the best way to match the 'export_files'
	// field. This may need to be updated in the future.
	if strings.Contains(s, resourcePathPrefix) && strings.Contains(s, ",") {
		paths := strings.Split(s, ",")

		elements := make([]string, len(paths))
		for i, path := range paths {
			elements[i] = fmt.Sprintf(`"%s"`, path)
		}

		s = fmt.Sprintf("PackedStringArray(%s)", strings.Join(elements, ", "))
	}

	return `"` + s + `"`
}
