package export

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/godot/platform"
	"github.com/coffeebeats/gdbuild/pkg/run"

	ini "gopkg.in/ini.v1"
)

const resourcePathPrefix = "res://"

/* -------------------------------------------------------------------------- */
/*                                 Enum: Mode                                 */
/* -------------------------------------------------------------------------- */

type Mode string

const (
	ModeUnknown    = ""
	ModeResources  = "resources"
	ModeCustomized = "customized"
)

/* -------------------------- Enum: FileVisualMode -------------------------- */

type FileVisualMode string

const (
	FileVisualModeUnknown = ""
	FileVisualModeKeep    = "keep"
	FileVisualModeRemove  = "remove"
	FileVisualModeStrip   = "strip"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Preset                               */
/* -------------------------------------------------------------------------- */

// Preset defines the parameters used in a Godot export preset.
type Preset struct {
	Arch            platform.Arch     `ini:"-"`
	CustomizedFiles map[string]string `ini:"-"`
	Embed           bool              `ini:"-"`
	Encrypt         bool              `ini:"encrypt_pck"`
	EncryptIndex    bool              `ini:"encrypt_directory"`
	Encrypted       []string          `ini:"encryption_include_filters"`
	EncryptionKey   string            `ini:"-"`
	Exclude         string            `ini:"exclude_filter"`
	ExportedFiles   []string          `ini:"export_files"`
	// ExportMode sets the type of export to use. Should be 'resources' for a
	// standard pack file and 'customized' for dedicated server pack files.
	ExportMode   Mode        `ini:"export_filter"`
	Features     []string    `ini:"custom_features"`
	Include      string      `ini:"include_filter"`
	Name         string      `ini:"name"`
	PathTemplate osutil.Path `ini:"-"`
	Platform     platform.OS `ini:"-"`
	Runnable     bool        `ini:"runnable"`
	Server       bool        `ini:"dedicated_server"`

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

	log.Debugf("adding file to preset: %s", pathString)

	p.ExportedFiles = append(p.ExportedFiles, pathString)

	slices.Sort(p.ExportedFiles)
	p.ExportedFiles = slices.Compact(p.ExportedFiles)

	return nil
}

/* ----------------------------- Method: AddFile ---------------------------- */

func (p *Preset) AddServerFile(rc *run.Context, path osutil.Path, visuals FileVisualMode) error {
	if visuals == FileVisualModeUnknown {
		return fmt.Errorf("%w: visuals", ErrMissingInput)
	}

	if err := path.CheckIsFile(); err != nil {
		return err
	}

	if err := path.RelTo(rc.PathWorkspace); err != nil {
		return err
	}

	pathString := strings.TrimPrefix(path.String(), rc.PathWorkspace.String())
	pathString = strings.TrimPrefix(filepath.Clean("./"+pathString), "./")
	pathString = resourcePathPrefix + strings.TrimPrefix(pathString, resourcePathPrefix)

	log.Debugf("adding server file to preset: %s", pathString)

	if p.CustomizedFiles == nil {
		p.CustomizedFiles = map[string]string{}
	}

	if _, ok := p.CustomizedFiles[pathString]; ok {
		return fmt.Errorf("%w: path already present: %s", ErrConflictingValue, pathString)
	}

	p.CustomizedFiles[pathString] = string(visuals)

	return nil
}

/* ------------------------- Method: exportPlatform ------------------------- */

func (p *Preset) exportPlatform() string {
	switch pl := p.Platform; pl {
	case platform.OSLinux:
		return "Linux/X11"
	case platform.OSMacOS:
		return "macOS"
	case platform.OSWindows:
		return "Windows Desktop"
	default:
		return ""
	}
}

/* ----------------------------- Method: Marshal ---------------------------- */

func (p *Preset) Marshal(w io.Writer, index int) error { //nolint:funlen
	// Don't write anything for an empty 'Preset'.
	if reflect.DeepEqual(p, new(Preset)) {
		return nil
	}

	preset := p

	lo := ini.LoadOptions{PreserveSurroundedQuote: true} //nolint:exhaustruct

	cfg := ini.Empty(lo)
	cfg.ValueMapper = valueMapper

	presetName := "preset." + strconv.Itoa(index)
	section := cfg.Section(presetName)

	section.Key("platform").SetValue(valueMapper(p.exportPlatform()))

	// Initially hydrate the 'ini' file.
	if err := section.ReflectFrom(preset); err != nil {
		return err
	}

	// Unmarshal the file back into the struct to trigger the value mapper.
	if err := section.MapTo(preset); err != nil {
		return err
	}

	// Finally, re-populate the file with the updated values.
	if err := section.ReflectFrom(preset); err != nil {
		return err
	}

	if len(preset.CustomizedFiles) > 0 {
		customizedFiles, err := json.Marshal(preset.CustomizedFiles)
		if err != nil {
			return err
		}

		section.Key("customized_files").SetValue(string(customizedFiles))
	}

	section = cfg.Section(presetName + ".options")

	options := maps.Clone(preset.Options)
	if options == nil {
		options = map[string]string{}
	}

	if p.Embed {
		options["binary_format/embed_pck"] = strconv.FormatBool(preset.Embed)
	}

	options["binary_format/architecture"] = preset.Arch.String()
	options["custom_template/debug"] = preset.PathTemplate.String()
	options["custom_template/release"] = preset.PathTemplate.String()

	for key, value := range options {
		if value == "" {
			delete(options, key)

			continue
		}

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
	if strings.Contains(s, resourcePathPrefix) {
		paths := strings.Split(s, ",")

		elements := make([]string, len(paths))
		for i, path := range paths {
			elements[i] = fmt.Sprintf(`"%s"`, path)
		}

		return fmt.Sprintf("PackedStringArray(%s)", strings.Join(elements, ", "))
	}

	return `"` + s + `"`
}

/* -------------------------------------------------------------------------- */
/*                    Function: NewWriteExportPresetsAction                   */
/* -------------------------------------------------------------------------- */

// NewWriteExportPresetsAction creates a new 'action.Action' which constructs an
// 'export_presets.cfg' file based on the target. It will be written to the
// workspace directory and overwrite any existing files.
func NewWriteExportPresetsAction(
	rc *run.Context,
	x *Export,
) action.WithDescription[action.Function] {
	path := filepath.Join(rc.PathWorkspace.String(), "export_presets.cfg")

	fn := func(_ context.Context) error {
		presets, err := x.Presets(rc)
		if err != nil {
			return err
		}

		var cfg strings.Builder

		for i, preset := range presets {
			if err := preset.Marshal(&cfg, i); err != nil {
				return err
			}
		}

		f, err := os.Create(filepath.Join(rc.PathWorkspace.String(), "export_presets.cfg"))
		if err != nil {
			return err
		}

		defer f.Close()

		if _, err := io.Copy(f, strings.NewReader(cfg.String()+"\n")); err != nil {
			return err
		}

		return nil
	}

	return action.WithDescription[action.Function]{
		Action:      fn,
		Description: "generate export presets file: " + path,
	}
}
