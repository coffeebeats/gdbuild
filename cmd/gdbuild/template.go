package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"

	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/config"
	"github.com/coffeebeats/gdbuild/pkg/godot/build"
	"github.com/coffeebeats/gdbuild/pkg/template"
)

var ErrPrintHashUsage = errors.New("cannot set option with '--print-hash'")

// A 'urfave/cli' command to compile a Godot export template.
func NewTemplate() *cli.Command { //nolint:cyclop,funlen,gocognit
	return &cli.Command{
		Name:     "template",
		Category: "Build",

		Usage:     "compile an export template for the specified Godot platform 'PLATFORM'",
		UsageText: "gdbuild template [OPTIONS] <PLATFORM>",

		Flags: []cli.Flag{
			newVerboseFlag(),

			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "log the build command without running it",
			},
			&cli.BoolFlag{
				Name:  "print-hash",
				Usage: "log the unique hash of the export template (skips compilation)",
			},
			&cli.PathFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "gdbuild.toml",
				Usage:   "use the 'gdbuild' configuration file found at 'PATH'",
			},
			&cli.PathFlag{
				Name:  "build-dir",
				Usage: "build the template within 'PATH' (defaults to a temporary directory)",
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
			// Validate arguments.
			platformInput := c.Args().First()
			if platformInput == "" {
				return UsageError{ctx: c, err: fmt.Errorf("%w: 'platform'", ErrMissingInput)}
			}

			if c.Args().Len() > 1 {
				return UsageError{
					ctx: c,
					err: fmt.Errorf("%w: %s", ErrTooManyArguments, strings.Join(c.Args().Slice()[1:], " "))}
			}

			// Validate flag options.
			if c.IsSet("release") && c.IsSet("release_debug") {
				return UsageError{ctx: c, err: ErrTargetUsageProfiles}
			}

			if c.IsSet("print-hash") {
				for _, opt := range []string{"build-dir", "dry-run", "out"} {
					if c.IsSet(opt) {
						return UsageError{
							ctx: c,
							err: fmt.Errorf("%w: --%s", ErrPrintHashUsage, opt),
						}
					}
				}

				// Don't log anything lower than error since that will obstruct
				// parsing of the hash.
				log.SetLevel(log.ErrorLevel)
			}

			// Determine path to store.
			storePath, err := touchStore()
			if err != nil {
				return err
			}

			log.Debugf("using store at path: %s", storePath)

			// Determine paths for build context.

			pathOut, err := parseWorkDir(c.Path("out"), c.Bool("dry-run"))
			if err != nil {
				return err
			}

			log.Debugf("placing template artifacts at path: %s", pathOut)

			pathBuild := c.Path("build-dir")
			if pathBuild == "" {
				p, err := os.MkdirTemp("", "gdbuild-*")
				if err != nil {
					return err
				}

				defer os.RemoveAll(p)

				pathBuild = p
			}

			log.Debugf("using build directory: %s", pathBuild)

			// Parse manifest.
			pathManifest, err := parseManifestPath(c.Path("config"))
			if err != nil {
				return err
			}

			m, err := config.ParseFile(pathManifest)
			if err != nil {
				return err
			}

			log.Debugf("using manifest at path: %s", pathManifest)

			// Evaluate build context.

			features := c.StringSlice("feature")

			log.Infof("features: %s", strings.Join(features, ","))

			pr := parseProfile(c.Bool("release"), c.Bool("release_debug"))

			log.Infof("profile: %s", pr)

			pl, err := parsePlatform(platformInput)
			if err != nil {
				return err
			}

			log.Infof("platform: %s", pl)

			bc := build.Context{
				Features:     features,
				PathBuild:    osutil.Path(pathBuild),
				PathManifest: osutil.Path(pathManifest),
				PathOut:      osutil.Path(pathOut),
				Platform:     pl,
				Profile:      pr,
				Verbose:      log.GetLevel() == log.DebugLevel,
			}

			t, err := template.Build(m, &bc)
			if err != nil {
				return err
			}

			cs, err := t.Checksum()
			if err != nil {
				return err
			}

			if c.Bool("print-hash") {
				log.Print(cs)

				return nil
			}

			action, err := template.Action(t, &bc)
			if err != nil {
				return err
			}

			if c.Bool("dry-run") {
				log.Print(action.Sprint())

				log.Debug("template hash: " + cs)

				return nil
			}

			return action.Run(c.Context)
		},
	}
}

/* ------------------------- Function: parseWorkDir ------------------------- */

func parseWorkDir(path string, dryRun bool) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		if !dryRun {
			if err := os.MkdirAll(path, osutil.ModeUserRWXGroupRX); err != nil {
				return "", err
			}
		}
	}

	if info != nil && !info.IsDir() {
		path = filepath.Dir(path)
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return path, nil
}
