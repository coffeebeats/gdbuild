package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"

	"github.com/coffeebeats/gdbuild/pkg/manifest"
)

var ErrInvalidPath = fmt.Errorf(
	"%w: expected 'path' to be a gdbuild.toml manifest file",
	ErrInvalidInput,
)

// A 'urfave/cli' command to compile and export a Godot project.
func NewExport() *cli.Command {
	return &cli.Command{
		Name:     "export",
		Category: "Export",

		Aliases: []string{"x"},

		Usage:     "Build and export the specified target.",
		UsageText: "gdbuild export [OPTIONS] <TARGET>",

		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Value:   ".",
				Usage:   "use the Godot project found at 'PATH'",
			},
		},

		Action: func(c *cli.Context) error {
			target := c.Args().First()
			if target == "" {
				return UsageError{ctx: c, err: fmt.Errorf("%w: target", ErrMissingInput)}
			}

			path := c.Path("path")

			log.Debugf("using manifest at path: %s", path)

			info, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("%w: %s: %w", ErrInvalidPath, path, err)
			}

			if info.IsDir() {
				path = filepath.Join(path, manifest.Filename())
			}

			_, err = manifest.ParseFile(path)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
