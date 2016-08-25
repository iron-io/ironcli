package cmd

import (
	"fmt"

	"github.com/iron-io/iron_go3/config"
	"github.com/urfave/cli"
)

type Register struct {
	cli.Command
}

func NewRegister(settings *config.Settings) *Register {
	register := &Register{
		Command: cli.Command{
			Name:      "register",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[test]",
			Action: func(c *cli.Context) error {
				fmt.Println(settings)
				fmt.Println("added task: ", c.Args().First())
				return nil
			},
		},
	}

	return register
}

func (r Register) GetCmd() cli.Command {
	return r.Command
}
