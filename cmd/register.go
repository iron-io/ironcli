package cmd

import (
	"github.com/iron-io/ironcli/helpers"
	"github.com/spf13/cobra"
)

var RegisterCmd = &cobra.Command{
	Use: "register",
	Run: helpers.OldCommands,
}

// TODO: Convert old commands to cobra and put it here
