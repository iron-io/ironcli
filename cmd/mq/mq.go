package mq

import (
	"github.com/iron-io/ironcli/commands"
	"github.com/iron-io/ironcli/helpers"
	"github.com/spf13/cobra"
)

var commandName = "mq"

var RootCmd = &cobra.Command{
	Use: commandName,
}

// TODO: Convert old commands to cobra and put it here

func init() {
	commands := commands.Commands[commandName].(commands.Mapper)
	for name := range commands {
		RootCmd.AddCommand(&cobra.Command{Use: name, Run: helpers.OldCommands})
	}
}
