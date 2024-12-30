package cmd

import "github.com/urfave/cli/v2"

// Command is the interface that all commands must implement
type Command interface {
	// Name returns the name of the command
	Name() string
	// Usage returns the usage description of the command
	Usage() string
	// Flags returns the flags for the command
	Flags() []cli.Flag
	// Action returns the action function for the command
	Action() func(c *cli.Context) error
}

// BaseCommand provides a base implementation of the Command interface
type BaseCommand struct {
	name   string
	usage  string
	flags  []cli.Flag
	action func(c *cli.Context) error
}

func NewBaseCommand(name, usage string, flags []cli.Flag, action func(c *cli.Context) error) *BaseCommand {
	return &BaseCommand{
		name:   name,
		usage:  usage,
		flags:  flags,
		action: action,
	}
}

func (b *BaseCommand) Name() string {
	return b.name
}

func (b *BaseCommand) Usage() string {
	return b.usage
}

func (b *BaseCommand) Flags() []cli.Flag {
	return b.flags
}

func (b *BaseCommand) Action() func(c *cli.Context) error {
	return b.action
}

// ToCLICommand converts a Command to a cli.Command
func ToCLICommand(cmd Command) *cli.Command {
	return &cli.Command{
		Name:   cmd.Name(),
		Usage:  cmd.Usage(),
		Flags:  cmd.Flags(),
		Action: cmd.Action(),
	}
}
