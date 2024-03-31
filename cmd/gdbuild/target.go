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
	godotexport "github.com/coffeebeats/gdbuild/pkg/godot/export"
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

			// Determine output path.
			pathOut, err := parseWorkDir(c.Path("out"), dryRun)
			if err != nil {
				return err
			}

			log.Debugf("placing template artifacts at path: %s", pathOut)

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
			rc, err := buildTemplateContext(c, pathManifest, "", c.String("platform"), dryRun)
			if err != nil {
				return err
			}

			tl, err := config.Template(&rc, m)
			if err != nil {
				return err
			}

			ec, err := buildExportContext(rc, pathOut)
			if err != nil {
				return err
			}

			xp, err := config.Export(&ec, m, tl, target)
			if err != nil {
				return err
			}

			if printHash {
				return printTargetHash(&ec, xp)
			}

			templateAction, err := exportTemplate(
				c.Context,
				&rc,
				storePath,
				tl,
				/* force= */ false,
			)
			if err != nil {
				return err
			}

			exportAction, err := exportProject(
				c.Context,
				&ec,
				storePath,
				xp,
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

/* ---------------------- Function: buildExportContext ---------------------- */

func buildExportContext(rc run.Context, pathOut string) (run.Context, error) {
	// Update the workspace path to the project directory.
	rc.PathWorkspace = osutil.Path(filepath.Dir(rc.PathManifest.String()))

	// Update output directory to option value.
	rc.PathOut = osutil.Path(pathOut)

	if err := rc.Validate(); err != nil {
		return run.Context{}, err
	}

	if err := rc.ProjectManifest().CheckIsFile(); err != nil {
		return run.Context{}, fmt.Errorf(
			"%w: Godot project configuration: %s",
			ErrMissingInput,
			rc.ProjectPath(),
		)
	}

	return rc, nil
}

/* ------------------------- Function: exportProject ------------------------ */

func exportProject( //nolint:ireturn
	_ context.Context,
	rc *run.Context,
	storePath string,
	xp *godotexport.Export,
	force bool,
) (action.Action, error) {
	cs, err := xp.Checksum(rc)
	if err != nil {
		return nil, err
	}

	hasTarget, err := store.HasTarget(storePath, cs)
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

		pathArchive, err := store.TargetArchive(storePath, cs)
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
	return export.Action(rc, xp)
}

/* ------------------------ Function: printTargetHash ----------------------- */

func printTargetHash(rc *run.Context, xp *godotexport.Export) error {
	cs, err := xp.Checksum(rc)
	if err != nil {
		return err
	}

	log.Print(cs)

	return nil
}
