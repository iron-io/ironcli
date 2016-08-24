package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

type Register struct {
	cli.Command
}

func NewRegister() *Register {
	register := &Register{
		Command: cli.Command{
			Name:      "register",
			Usage:     "do the doo",
			UsageText: "doo - does the dooing",
			ArgsUsage: "[test]",
			Action: func(c *cli.Context) error {
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
