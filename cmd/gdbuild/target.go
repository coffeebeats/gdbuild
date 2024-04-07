package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"

	"github.com/coffeebeats/gdbuild/internal/action"
	"github.com/coffeebeats/gdbuild/internal/archive"
	"github.com/coffeebeats/gdbuild/internal/osutil"
	"github.com/coffeebeats/gdbuild/pkg/config"
	"github.com/coffeebeats/gdbuild/pkg/godot/export"
	"github.com/coffeebeats/gdbuild/pkg/godot/template"
	"github.com/coffeebeats/gdbuild/pkg/run"
	"github.com/coffeebeats/gdbuild/pkg/store"
	"github.com/coffeebeats/gdbuild/pkg/target"
)

var ErrTargetUsageProfiles = errors.New("cannot specify more than one of '--debug', '--release_debug', and '--release'")

// A 'urfave/cli' command to compile and export a Godot project target.
func NewTarget() *cli.Command { //nolint:cyclop,funlen,gocognit,maintidx
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
				Usage:   "use the 'gdbuild' configuration file found at 'PATH'",
			},
			&cli.PathFlag{
				Name:  "project",
				Usage: "use the Godot project found at 'PATH'",
			},
			&cli.PathFlag{
				Name:    "out",
				Aliases: []string{"o"},
				Value:   ".",
				Usage:   "write generated artifacts to 'PATH'",
			},
			&cli.PathFlag{
				Name:  "template-archive",
				Usage: "extract the template from the archive found at 'PATH' (skips template build)",
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
				Usage:    "use a release export template (cannot be used with '--release_debug' or '--debug')",
			},
			&cli.BoolFlag{
				Name:     "release_debug",
				Category: "Profile",
				Usage:    "use a release export template with debug symbols (cannot be used with '--release' or '--debug')",
			},
			&cli.BoolFlag{
				Name:     "debug",
				Category: "Profile",
				Usage:    "use a debug export template (cannot be used with '--release' or '--release_debug')",
			},
		},

		Action: func(c *cli.Context) error {
			// Validate arguments.
			targetName := c.Args().First()
			if targetName == "" {
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
			pathOut, err := parseOutDir(c.Path("out"), dryRun)
			if err != nil {
				return err
			}

			pathConfig := c.Path("config")
			pathProject := c.Path("project")

			switch {
			case pathConfig == "" && pathProject != "":
				pathConfig = filepath.Join(pathProject, config.DefaultFilename())
			case pathProject == "" && pathConfig != "":
				pathProject = filepath.Dir(pathConfig)
			case pathProject == "" && pathConfig == "":
				pathProject = "."
				pathConfig = config.DefaultFilename()
			}

			// Parse manifest.
			pathManifest, err := parseManifestPath(pathConfig)
			if err != nil {
				return err
			}

			m, err := config.ParseFile(pathManifest)
			if err != nil {
				return err
			}

			// Evaluate build context.
			rc, err := buildTemplateContext(c, pathManifest, "", c.String("platform"), dryRun)
			if err != nil {
				return err
			}

			defer cleanTemporaryDirectory(&rc)

			tl, err := config.Template(&rc, m)
			if err != nil {
				return err
			}

			ec, err := buildExportContext(rc, targetName, pathProject, pathOut)
			if err != nil {
				return err
			}

			defer cleanTemporaryDirectory(&ec)

			xp, err := config.Export(&ec, m, tl, targetName)
			if err != nil {
				return err
			}

			pathTemplateArchive, err := templateArchivePath(c, storePath, tl)
			if err != nil {
				return err
			}

			if c.IsSet("template-archive") {
				xp.PathTemplateArchive = pathTemplateArchive
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
				tl,
				xp,
				force,
			)
			if err != nil {
				return err
			}

			extractTemplateAction, err := target.NewExtractTemplateAction(&ec, pathTemplateArchive)
			if err != nil {
				return err
			}

			exportAction = action.InOrder(
				templateAction,
				extractTemplateAction,
				export.NewInstallEditorGodotAction(&ec, xp.Version, ec.GodotPath()),
				exportAction,
			)

			if dryRun {
				log.Print(exportAction.Sprint())

				return nil
			}

			return exportAction.Run(c.Context)
		},
	}
}

/* ---------------------- Function: buildExportContext ---------------------- */

func buildExportContext(rc run.Context, targetName, pathProject, pathOut string) (run.Context, error) {
	pathWorkspace := osutil.Path(filepath.Dir(rc.PathManifest.String()))
	if pathProject != "" {
		pathWorkspace = osutil.Path(pathProject)

		wd, err := os.Getwd()
		if err != nil {
			return run.Context{}, err
		}

		if err := pathWorkspace.RelTo(osutil.Path(wd)); err != nil {
			return run.Context{}, err
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return run.Context{}, err
	}

	// Update the workspace path to the project directory.
	rc.PathWorkspace = pathWorkspace
	if err := rc.PathWorkspace.RelTo(osutil.Path(wd)); err != nil {
		return run.Context{}, err
	}

	// Update output directory to option value.
	rc.PathOut = osutil.Path(pathOut)
	if err := rc.PathWorkspace.RelTo(osutil.Path(pathOut)); err != nil {
		return run.Context{}, err
	}

	rc.Target = targetName

	if err := rc.Validate(); err != nil {
		return run.Context{}, err
	}

	pathGodotManifest := rc.GodotProjectManifestPath()
	if err := pathGodotManifest.CheckIsFile(); err != nil {
		return run.Context{}, fmt.Errorf(
			"%w: Godot project configuration: %s",
			ErrMissingInput,
			pathGodotManifest.String(),
		)
	}

	return rc, nil
}

/* ------------------------- Function: exportProject ------------------------ */

func exportProject( //nolint:ireturn
	_ context.Context,
	rc *run.Context,
	storePath string,
	tl *template.Template,
	xp *export.Export,
	force bool,
) (action.Action, error) {
	cs, err := export.Checksum(rc, xp)
	if err != nil {
		return nil, err
	}

	pathTmp, err := rc.TempDir()
	if err != nil {
		return nil, err
	}

	templateName := template.Name(rc.Platform, tl.Arch, rc.Profile)
	xp.PathTemplate = osutil.Path(filepath.Join(pathTmp, templateName))

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

	log.Debugf("using project directory: %s", rc.PathWorkspace)

	// Target was not cached; create build action.
	return target.Action(rc, xp)
}

/* ---------------------- Function: pathTemplateArchive --------------------- */

func templateArchivePath(c *cli.Context, storePath string, tl *template.Template) (osutil.Path, error) {
	if c.IsSet("template-archive") {
		path := osutil.Path(c.Path("template-archive"))

		if err := path.CheckIsFile(); err != nil {
			return "", err
		}

		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		if err := path.RelTo(osutil.Path(wd)); err != nil {
			return "", err
		}

		return path, nil
	}

	cs, err := template.Checksum(tl)
	if err != nil {
		return "", err
	}

	pathArchive, err := store.TemplateArchive(storePath, cs)
	if err != nil {
		return "", err
	}

	return osutil.Path(pathArchive), nil
}

/* ------------------------ Function: printTargetHash ----------------------- */

func printTargetHash(rc *run.Context, xp *export.Export) error {
	cs, err := export.Checksum(rc, xp)
	if err != nil {
		return err
	}

	fmt.Println(cs) //nolint:forbidigo

	return nil
}
