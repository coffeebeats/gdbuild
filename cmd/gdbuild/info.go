package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"

	"github.com/coffeebeats/gdbuild/pkg/manifest"
)

// A 'urfave/cli' command to inspect a GDBuild manifest.
func NewInfo() *cli.Command { //nolint:funlen
	return &cli.Command{
		Name:     "info",
		Category: "Inspect",

		Usage:     "inspect various properties of the GDBuild manifest",
		UsageText: "gdbuild info [OPTIONS] <PROPERTY>",

		Flags: []cli.Flag{
			newVerboseFlag(),

			&cli.PathFlag{
				Name:     "path",
				FilePath: ".",
				Usage:    "use the Godot project found at 'PATH'",
			},
			&cli.BoolFlag{
				Name:  "json",
				Usage: "print the property values in JSON format",
			},
		},

		Action: func(c *cli.Context) error {
			if c.Args().Len() > 1 {
				return UsageError{
					ctx: c,
					err: fmt.Errorf("%w: %s", ErrTooManyArguments, strings.Join(c.Args().Slice()[1:], " "))}
			}

			// Parse manifest.
			pathManifest, err := parseManifestPath(c.Path("path"))
			if err != nil {
				return err
			}

			m, err := manifest.ParseFile(pathManifest)
			if err != nil {
				return err
			}

			log.Debugf("using manifest at path: %s", pathManifest)

			isJSON := c.Bool("json")

			var output any
			switch a := c.Args().First(); a {
			case "target", "targets":
				targets := make([]string, 0)
				for target := range m.Target {
					targets = append(targets, target)
				}

				if !isJSON {
					output = strings.Join(targets, "\n")
				} else {
					output = targets
				}

			default:
				return UsageError{ctx: c, err: fmt.Errorf("%w: unsupported property: %s", ErrInvalidInput, a)}
			}

			if isJSON {
				bb, err := json.Marshal(output)
				if err != nil {
					return err
				}

				output = string(bb)
			}

			log.Print(output)

			return nil
		},
	}
}
