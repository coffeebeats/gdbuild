package main

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/archive"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/config"
	"github.com/coffeebeats/gdbuild/pkg/export"
	"github.com/coffeebeats/gdbuild/pkg/run"
	"github.com/coffeebeats/gdbuild/pkg/store"
)

var ErrTargetUsageProfiles = errors.New("cannot specify both '--release' and '--release_debug'")

// A 'urfave/cli' command to compile and export a Godot project target.
func NewTarget() *cli.Command { //nolint:cyclop,funlen,gocognit
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

			dryRun := c.Bool("dry-run")
			force := c.Bool("force")
			printHash := c.Bool("print-hash")

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
				PathOut:      "", // No need to copy export templates anywhere.
				Platform:     pl,
				Profile:      pr,
				Verbose:      log.GetLevel() == log.DebugLevel,
			}

			if printHash {
				return printTargetHash(&rc, m, target)
			}

			templateAction, err := exportTemplate(
				c.Context,
				storePath,
				m,
				&rc,
				/* force= */ false,
			)
			if err != nil {
				return err
			}

			ec := rc

			// Update the build path to the project directory.
			ec.PathBuild = osutil.Path(filepath.Dir(ec.PathManifest.String()))

			// Update output directory to option value.
			ec.PathOut = osutil.Path(pathOut)

			exportAction, err := exportProject(
				c.Context,
				storePath,
				&ec,
				m,
				target,
				force,
			)
			if err != nil {
				return err
			}

			exportAction = templateAction.AndThen(exportAction)

			if dryRun {
				log.Print(exportAction.Sprint())

				return nil
			}

			return exportAction.Run(c.Context)
		},
	}
}

/* ------------------------- Function: exportProject ------------------------ */

func exportProject( //nolint:ireturn
	_ context.Context,
	storePath string,
	rc *run.Context,
	m *config.Manifest,
	target string,
	force bool,
) (action.Action, error) {
	x, err := config.Export(rc, m, target)
	if err != nil {
		return nil, err
	}

	hasTarget, err := store.HasTarget(storePath, rc, x)
	if err != nil {
		return nil, err
	}

	// Target is cached; create cache extraction action.
	if hasTarget && !force {
		log.Info("found target in cache; skipping build.")

		pathOut := rc.PathOut.String()

		if rc.PathOut == "" {
			log.Info("no output path set; exiting without changes")

			return action.NoOp{}, nil
		}

		pathArchive, err := store.TargetArchive(storePath, rc, x)
		if err != nil {
			return nil, err
		}

		fn := func(ctx context.Context) error {
			log.Infof("extracting artifacts from cached archive: %s", pathArchive)

			return archive.Extract(ctx, pathArchive, pathOut)
		}

		return action.WithDescription[action.Function]{
			Action:      fn,
			Description: "extract cached artifacts to path: " + pathOut,
		}, nil
	}

	// Target was not cached; create build action.
	return export.Action(rc, x)
}

/* ------------------------ Function: printTargetHash ----------------------- */

func printTargetHash(rc *run.Context, m *config.Manifest, target string) error {
	x, err := config.Export(rc, m, target)
	if err != nil {
		return err
	}

	cs, err := x.Checksum(rc)
	if err != nil {
		return err
	}

	log.Print(cs)

	return nil
}
