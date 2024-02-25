package command

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/coffeebeats/gdbuild/internal/merge"
)

/* -------------------------------------------------------------------------- */
/*                            Interface: Commander                            */
/* -------------------------------------------------------------------------- */

// Commander is a type which can be translated into a shell command.
type Commander interface {
	Command() (*Command, error)
}

/* -------------------------------------------------------------------------- */
/*                               Struct: Command                              */
/* -------------------------------------------------------------------------- */

type Command struct {
	Pre  []string
	Post []string

	Shell Shell

	Directory   string
	Environment map[string]string

	Args []string
}

/* ----------------------------- Method: Execute ---------------------------- */

func (c *Command) Execute(ctx context.Context) error {
	args, err := c.args()
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, c.Shell.String(), args...) //nolint:gosec

	cmd.Dir = c.Directory

	for k, v := range c.Environment {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	cmd.Stdout = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

/* --------------------------- Impl: fmt.Stringer --------------------------- */

func (c *Command) String() string {
	args, err := c.args()
	if err != nil {
		return ""
	}

	return c.Shell.String() + " " + strings.Join(args, " ")
}

/* ------------------------------ Method: args ------------------------------ */

func (c *Command) args() ([]string, error) {
	var args []string
	args = append(args, c.Pre...)
	args = append(args, strings.Join(c.Args, " "))
	args = append(args, c.Post...)

	command := strings.Join(args, " && ")

	switch c.Shell {
	case ShellBash, ShellSh, ShellZsh:
		command = fmt.Sprintf(`-C "%s"`, command)
	default:
		return nil, fmt.Errorf("%w: unsupport shell: %s", ErrUnsupportedShell, c.Shell)
	}

	return strings.Split(command, " "), nil
}

/* --------------------------- Impl: merge.Merger --------------------------- */

func (c *Command) Merge(other *Command) error {
	if c == nil || other == nil {
		return nil
	}

	c.Pre = append(c.Pre, other.Pre...)
	c.Post = append(c.Post, other.Post...)
	c.Args = append(c.Args, other.Args...)

	if err := merge.Primitive(&c.Shell, other.Shell); err != nil {
		return fmt.Errorf("%w: shell", err)
	}

	if err := merge.Primitive(&c.Directory, other.Directory); err != nil {
		return fmt.Errorf("%w: directory", err)
	}

	if err := merge.Map(&c.Environment, other.Environment); err != nil {
		return fmt.Errorf("%w: environment", err)
	}

	return nil
}
