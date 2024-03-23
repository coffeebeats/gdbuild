package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"

	"github.com/coffeebeats/gdbuild/pkg/store"
)

const (
	envLogLevel = "GDBUILD_LOG"

	lenLevelLabel = 5

	colorCyanBright    = 14
	colorGreenBright   = 10
	colorMagentaBright = 13
	colorRedBright     = 9
	colorWhiteBright   = 15
	colorYellowBright  = 11
)

var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrInvalidManifestPath = fmt.Errorf("%w: expected 'path' to be a gdbuild.toml manifest file", ErrInvalidInput)
	ErrMissingInput        = errors.New("missing required argument")
	ErrTooManyArguments    = errors.New("too many arguments (were options passed after args?)")
	ErrUnrecognizedLevel   = errors.New("unrecognized level")
)

func main() { //nolint:funlen
	cli.VersionPrinter = versionPrinter
	cli.VersionFlag = &cli.BoolFlag{
		Name:               "version",
		Aliases:            []string{"V"},
		Usage:              "print the version",
		DisableDefaultText: true,
	}

	app := &cli.App{
		Name:    "gdbuild",
		Version: "v0.1.3", // x-release-please-version

		Suggest:                true,
		UseShortOptionHandling: true,

		Flags: []cli.Flag{
			newVerboseFlag(),
		},

		Commands: []*cli.Command{
			/* ----------------------------- Build/Export ---------------------------- */

			NewTarget(),
			NewTemplate(),

			/* ------------------------------- Inspect ------------------------------- */

			NewInfo(),
		},
	}

	// Call 'os.Exit' as the first-in/last-out defer; ensures an exit code is
	// returned to the caller.
	var exitCode int
	defer func() {
		if err := recover(); err != nil {
			exitCode = 1

			log.Error(err)
		}

		os.Exit(exitCode)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Ensure that the signal handler is removed after first interrupt.
	go func() {
		<-ctx.Done()
		stop()
	}()

	if err := setUpLogger(); err != nil {
		panic(err)
	}

	if err := app.RunContext(ctx, os.Args); err != nil {
		var usageErr UsageError
		if errors.As(err, &usageErr) {
			usageErr.PrintUsage()
		}

		panic(err)
	}
}

/* -------------------------------------------------------------------------- */
/*                              Type: UsageError                              */
/* -------------------------------------------------------------------------- */

// UsageError is any error returned from a subcommand implementation that should
// have subcommand usage instructions printed.
type UsageError struct {
	ctx *cli.Context
	err error
}

/* -------------------------- Function: PrintUsage -------------------------- */

// PrintUsage prints the usage associated with the subcommand that failed.
func (e UsageError) PrintUsage() {
	// NOTE: This never returns a meaningful error so ignore it.
	cli.ShowSubcommandHelp(e.ctx) //nolint:errcheck
}

/* ------------------------------- Impl: Error ------------------------------ */

func (e UsageError) Error() string {
	return e.err.Error()
}

/* -------------------------------------------------------------------------- */
/*                            Function: setUpLogger                           */
/* -------------------------------------------------------------------------- */

// setUpLogger configures the package-level charm.sh 'log' logger.
func setUpLogger() error {
	// Configure timestamp reporting.
	log.SetReportTimestamp(false)

	// Configure styles for each log level.
	s := log.DefaultStyles()
	s.Levels[log.DebugLevel] = newStyleWithColor("debug", colorCyanBright)
	s.Levels[log.InfoLevel] = newStyleWithColor("info", colorGreenBright)
	s.Levels[log.WarnLevel] = newStyleWithColor("warn", colorYellowBright)
	s.Levels[log.ErrorLevel] = newStyleWithColor("error", colorRedBright)
	s.Levels[log.FatalLevel] = newStyleWithColor("fatal", colorMagentaBright)

	log.SetStyles(s)

	// Try to parse a log level override.
	if envLevel := os.Getenv(envLogLevel); envLevel != "" {
		level, err := log.ParseLevel(envLevel)
		if err != nil {
			return err
		}

		// Configure the default logging level.
		log.SetLevel(level)
	}

	return nil
}

/* ----------------------- Function: newStyleWithColor ---------------------- */

// newStyleWithColor creates a new 'lipgloss.Style' for the given log level and
// ANSI escape color.
//
// NOTE: This function assumes that the width of the level strings is '5'.
func newStyleWithColor(name string, ansiColor int) lipgloss.Style {
	if name == "" {
		panic("missing style name")
	}

	return lipgloss.NewStyle().
		SetString(name + ":").
		PaddingRight(int(math.Max(float64(lenLevelLabel-len(name)), 0))).
		Bold(true).
		Foreground(lipgloss.ANSIColor(ansiColor))
}

/* -------------------------------------------------------------------------- */
/*                          Function: newVerboseFlag                          */
/* -------------------------------------------------------------------------- */

// newVerboseFlag creates a new standardize verbosity flag which handles
// updating the log level.
func newVerboseFlag() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:               "verbose",
		Usage:              "increase log verbosity",
		Aliases:            []string{"v"},
		DisableDefaultText: true,

		Action: func(_ *cli.Context, isVerbose bool) error {
			if !isVerbose || log.GetLevel() == log.DebugLevel {
				return nil
			}

			if l := log.GetLevel(); isVerbose {
				log.SetLevel(l - (log.InfoLevel - log.DebugLevel))
			}

			return nil
		},
	}
}

/* -------------------------------------------------------------------------- */
/*                         Function: parseManifestPath                        */
/* -------------------------------------------------------------------------- */

func parseManifestPath(path string) (string, error) {
	path = filepath.Clean(path)

	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("%w: %s: %w", ErrInvalidManifestPath, path, err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("%w: %s", ErrInvalidManifestPath, path)
	}

	return path, nil
}

/* -------------------------------------------------------------------------- */
/*                            Function: touchStore                            */
/* -------------------------------------------------------------------------- */

// touchStore determines the store path and ensures it has the expected layout.
func touchStore() (string, error) {
	// Determine the store path.
	storePath, err := store.Path()
	if err != nil {
		return "", err
	}

	// Ensure the store exists.
	if err := store.Touch(storePath); err != nil {
		return "", err
	}

	return storePath, nil
}

/* -------------------------------------------------------------------------- */
/*                          Function: versionPrinter                          */
/* -------------------------------------------------------------------------- */

// versionPrinter prints a 'gdbuild' version string to the terminal.
func versionPrinter(cCtx *cli.Context) {
	log.Printf("gdbuild %s", cCtx.App.Version)
}
