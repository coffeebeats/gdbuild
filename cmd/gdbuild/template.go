package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to compile a Godot export template.
func NewTemplate() *cli.Command { //nolint:funlen
	return &cli.Command{
		Name:     "template",
		Category: "Build",

		Usage:     "compile an export template for the specified Godot platform 'PLATFORM'",
		UsageText: "gdbuild template [OPTIONS] <PLATFORM>",

		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:  "path",
				Value: ".",
				Usage: "use the Godot project found at 'PATH'",
			},
			&cli.PathFlag{
				Name:    "out",
				Aliases: []string{"o"},
				Value:   ".",
				Usage:   "place the compiled artifacts at 'PATH'",
			},
			&cli.StringSliceFlag{
				Name:     "feature",
				Aliases:  []string{"f"},
				Category: "Export",
				Usage:    "enable the provided feature tag 'FEATURE' (can be specified more than once)",
			},
			&cli.BoolFlag{
				Name:     "release",
				Category: "Profile",
				Usage:    "use a release export template (cannot be used with '--release_debug')",
			},
			&cli.BoolFlag{
				Name:     "release_debug",
				Category: "Profile",
				Usage:    "use a release export template with debug symbols (cannot be used with '--release')",
			},
		},

		Action: func(c *cli.Context) error {
			platformInput := c.Args().First()
			if platformInput == "" {
				return UsageError{ctx: c, err: fmt.Errorf("%w: 'platform'", ErrMissingInput)}
			}

			if c.Args().Len() > 1 {
				return UsageError{
					ctx: c,
					err: fmt.Errorf("%w: %s", ErrTooManyArguments, strings.Join(c.Args().Slice()[1:], " "))}
			}

			if c.IsSet("release") && c.IsSet("release_debug") {
				return UsageError{ctx: c, err: ErrTargetUsageProfiles}
			}

			pathOut, err := parseDirectory(c.Path("out"))
			if err != nil {
				return err
			}

			log.Debugf("moving template artifacts to path: %s", pathOut)

			// Collect build modifiers.

			features := c.StringSlice("feature")

			log.Debugf("enabling features for target: %s", strings.Join(features, ","))

			pr := parseProfile(c.Bool("release"), c.Bool("release_debug"))

			log.Debugf("using template profile: %s", pr)

			pl, err := parsePlatform(platformInput)
			if err != nil {
				return err
			}

			log.Debugf("building for platform: %s", pl)

			// Parse manifest.
			pathManifest := c.Path("path")
			m, err := parseManifest(pathManifest)
			if err != nil {
				return err
			}

			log.Debugf("using manifest at path: %s", pathManifest)

			log.Print(m.BuildTemplate(pl, pr, features...))

			return nil
		},
	}
}

/* ------------------------ Function: parseDirectory ------------------------ */

func parseDirectory(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if !info.IsDir() {
		return "", fmt.Errorf("%w: expected a directory: %s", ErrInvalidInput, path)
	}

	return path, nil
}
