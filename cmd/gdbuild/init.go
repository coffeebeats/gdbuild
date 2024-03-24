package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/coffeebeats/gdbuild/pkg/config"
)

// A 'urfave/cli' command to initialize a GDBuild manifest.
func NewInit() *cli.Command {
	return &cli.Command{
		Name:     "init",
		Category: "Configuration",

		Usage:     "initialize a project with a GDBuild manifest",
		UsageText: "gdbuild init [OPTIONS]",

		Flags: []cli.Flag{
			newVerboseFlag(),

			&cli.PathFlag{
				Name:  "path",
				Value: ".",
				Usage: "use the Godot project found at 'PATH'",
			},
		},

		Action: func(c *cli.Context) error {
			if c.Args().Len() > 0 {
				return UsageError{
					ctx: c,
					err: fmt.Errorf(
						"%w: %s", ErrTooManyArguments,
						strings.Join(c.Args().Slice()[1:], " "),
					),
				}
			}

			path := c.Path("path")
			if path == "" {
				return UsageError{
					ctx: c,
					err: fmt.Errorf("%w: '--path'", ErrMissingInput),
				}
			}

			return config.Init(path)
		},
	}
}
