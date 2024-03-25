package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"

	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/config"
	"github.com/coffeebeats/gdbuild/pkg/run"
)

var ErrTargetUsageProfiles = errors.New("cannot specify both '--release' and '--release_debug'")

// A 'urfave/cli' command to compile and export a Godot project target.
func NewTarget() *cli.Command { //nolint:cyclop,funlen
	return &cli.Command{
		Name:     "target",
		Category: "Build",

		Usage:     "compile required Godot export template(s) and then export the specified 'TARGET'",
		UsageText: "gdbuild target [OPTIONS] <TARGET>",

		Flags: []cli.Flag{
			newVerboseFlag(),

			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "log the build command without running it",
			},
			&cli.BoolFlag{
				Name:  "force",
				Usage: "export the target even if it was cached in the store (does not rebuild the export template)",
			},
			&cli.BoolFlag{
				Name:  "print-hash",
				Usage: "log the unique hash of the game binary (skips exporting)",
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
				Usage:   "write generated artifacts to 'PATH'",
			},
			&cli.StringSliceFlag{
				Name:     "feature",
				Aliases:  []string{"f"},
				Category: "Export",
				Usage:    "enable the provided feature tag 'FEATURE' (can be specified more than once)",
			},
			&cli.StringFlag{
				Name:     "platform",
				Aliases:  []string{"p"},
				Category: "Template",
				Usage:    "build for the specified Godot platform 'PLATFORM'",
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
			target := c.Args().First()
			if target == "" {
				return UsageError{ctx: c, err: fmt.Errorf("%w: target", ErrMissingInput)}
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

			dryRun := c.Bool("dry-run")
			// force := c.Bool("force")
			// printHash := c.Bool("print-hash")

			// Determine path to store.
			storePath, err := touchStore()
			if err != nil {
				return err
			}

			log.Debugf("using store at path: %s", storePath)

			// Determine paths for export context.

			pathOut, err := parseWorkDir(c.Path("out"), dryRun)
			if err != nil {
				return err
			}

			log.Debugf("placing template artifacts at path: %s", pathOut)

			pathBuild, err := parseBuildDir(c.Path("build-dir"), dryRun)
			if err != nil {
				return err
			}

			log.Debugf("using build directory: %s", pathBuild)

			pathManifest, err := parseManifestPath(c.Path("config"))
			if err != nil {
				return err
			}

			m, err := config.ParseFile(pathManifest)
			if err != nil {
				return err
			}

			log.Debugf("using manifest at path: %s", pathManifest)

			// Evaluate export context.

			features := c.StringSlice("feature")

			log.Infof("features: %s", strings.Join(features, ","))

			pr := parseProfile(c.Bool("release"), c.Bool("release_debug"))

			log.Infof("profile: %s", pr)

			pl, err := parsePlatform(c.String("platform"))
			if err != nil {
				return err
			}

			log.Infof("platform: %s", pl)

			rc := run.Context{
				Features:     features,
				PathBuild:    osutil.Path(pathBuild),
				PathManifest: osutil.Path(pathManifest),
				PathOut:      osutil.Path(pathOut),
				Platform:     pl,
				Profile:      pr,
				Verbose:      log.GetLevel() == log.DebugLevel,
			}

			templateAction, err := exportTemplate(
				c.Context,
				storePath,
				m,
				&rc,
				/* force= */ false,
				/* print-hash= */ false,
			)
			if err != nil {
				return err
			}

			if dryRun {
				log.Print(templateAction.Sprint())

				return nil
			}

			return templateAction.Run(c.Context)
		},
	}
}
