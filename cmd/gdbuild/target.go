package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"

	"github.com/coffeebeats/gdbuild/pkg/build"
	"github.com/coffeebeats/gdbuild/pkg/build/platform"
	"github.com/coffeebeats/gdbuild/pkg/manifest"
)

var ErrTargetUsageProfiles = errors.New("cannot specify both '--release' and '--release_debug'")

// A 'urfave/cli' command to compile and export a Godot project target.
func NewTarget() *cli.Command { //nolint:funlen
	return &cli.Command{
		Name:     "target",
		Category: "Build",

		Usage:     "compile required Godot export template(s) and then export the specified 'TARGET'",
		UsageText: "gdbuild target [OPTIONS] <TARGET>",

		Flags: []cli.Flag{
			newVerboseFlag(),

			&cli.PathFlag{
				Name:  "path",
				Value: ".",
				Usage: "use the Godot project found at 'PATH'",
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
			target := c.Args().First()
			if target == "" {
				return UsageError{ctx: c, err: fmt.Errorf("%w: target", ErrMissingInput)}
			}

			if c.Args().Len() > 1 {
				return UsageError{
					ctx: c,
					err: fmt.Errorf("%w: %s", ErrTooManyArguments, strings.Join(c.Args().Slice()[1:], " "))}
			}

			if c.IsSet("release") && c.IsSet("release_debug") {
				return UsageError{ctx: c, err: ErrTargetUsageProfiles}
			}

			pathOut, err := parseWorkDir(c.Path("out"))
			if err != nil {
				return err
			}

			log.Debugf("placing template artifacts at path: %s", pathOut)

			// Parse manifest.
			pathManifest, err := parseManifestPath(c.Path("path"))
			if err != nil {
				return err
			}

			_, err = manifest.ParseFile(pathManifest)
			if err != nil {
				return err
			}

			log.Debugf("using manifest at path: %s", pathManifest)

			pathManifest = filepath.Dir(pathManifest) //nolint:ineffassign,staticcheck,wastedassign

			// Collect build modifiers.

			features := c.StringSlice("feature")

			log.Infof("features: %s", strings.Join(features, ","))

			pr := parseProfile(c.Bool("release"), c.Bool("release_debug"))

			log.Infof("profile: %s", pr)

			pl, err := parsePlatform(c.String("platform"))
			if err != nil {
				return err
			}

			log.Infof("platform: %s", pl)

			return nil
		},
	}
}

func parsePlatform(platformInput string) (platform.OS, error) {
	if platformInput == "" {
		platformInput = runtime.GOOS
	}

	godotPlatform, err := platform.ParseOS(platformInput)
	if err != nil {
		return platform.OS(0), err
	}

	return godotPlatform, nil
}

func parseProfile(releaseInput, releaseDebugInput bool) build.Profile {
	var profile build.Profile

	switch {
	case releaseInput:
		profile = build.ProfileRelease
	case releaseDebugInput:
		profile = build.ProfileReleaseDebug
	default:
		profile = build.ProfileDebug
	}

	return profile
}
