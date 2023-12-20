package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/gavsidhu/miflo/internal/helpers"
)

type Command struct {
	Name        string
	Description string
	Usage       string
	Flags       *flag.FlagSet
	SubCommands map[string]*Command
	Run         func(cmd *Command, args []string)
}

func (c *Command) AddCommand(subCommand *Command) {
	c.SubCommands[subCommand.Name] = subCommand
}

func (c *Command) Execute() {
	args := os.Args[1:]

	err := c.Flags.Parse(args)
	if err != nil {
		helpers.ErrAndExit(fmt.Sprintf("Error: %v", err))
	}

	nonFlagArgs := c.Flags.Args()

	if len(nonFlagArgs) > 0 {
		commandArg := nonFlagArgs[0]

		command, ok := c.SubCommands[commandArg]
		if !ok {
			helpers.ErrAndExit("Command does not exist")
			return
		}

		command.Run(command, nonFlagArgs[1:])
	}
	c.Run(c, nonFlagArgs)
}
