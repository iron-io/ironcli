package helpers

import (
	"flag"
	"fmt"
	"os"

	"github.com/iron-io/ironcli/commands"
	"github.com/spf13/cobra"
)

func OldCommands(cobra *cobra.Command, args []string) {
	var c commands.Commander
	var ok bool

	if cobra.Parent().Name() != "iron" {
		if c, ok = commands.Commands[cobra.Parent().Name()]; !ok {
			fmt.Println("Command not found")
			os.Exit(-1)
		}
		args = append([]string{cobra.Name()}, args...)
	} else {
		if c, ok = commands.Commands[cobra.Name()]; !ok {
			fmt.Println("Command not found")
			os.Exit(-1)
		}
	}

	cmd, err := c.Command(args...)

	if err != nil {
		if err == flag.ErrHelp && cmd != nil {
			cmd.Usage()
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = cmd.Config()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	cmd.Run()
}
